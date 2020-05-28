package server

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/psykhi/pong/game"
	uuid "github.com/satori/go.uuid"
	"nhooyr.io/websocket"
	"time"
)

const MaxPlayerCount = 2

type GameInstance struct {
	*game.State
	playerConnections []*PlayerConn
	e                 *game.Engine
	spectators        []*websocket.Conn
	updates           chan InputUpdate
	endCh             chan string
	id                string
}

func NewGameInstance(endCh chan string) *GameInstance {
	s := game.NewState()
	s.WaitingForPlayer = true
	g := &GameInstance{
		playerConnections: make([]*PlayerConn, MaxPlayerCount),
		State:             s,
		e:                 &game.Engine{},
		spectators:        []*websocket.Conn{},
		updates:           make(chan InputUpdate, MaxPlayerCount),
		endCh:             endCh,
		id:                uuid.NewV4().String(),
	}
	go g.loop()
	return g
}

func (g *GameInstance) AllConnected() bool {
	for _, pc := range g.playerConnections {
		if pc == nil {
			return false
		}
	}
	return true
}

func (g *GameInstance) ConnectPlayer(playerID int, c *websocket.Conn) {
	pc := &PlayerConn{
		id:       playerID,
		Conn:     c,
		updateCh: g.updates,
	}
	g.playerConnections[playerID] = pc
	g.State.Players[playerID].Connected = true
	go g.playerConnections[playerID].Start()

	if g.AllConnected() {
		g.State.WaitingForPlayer = false
		g.State.Countdown()
	}
}

func (g *GameInstance) addSpectator(c *websocket.Conn) {
	g.spectators = append(g.spectators, c)
}

func (g *GameInstance) closeConnections() {
	for _, c := range g.playerConnections {
		if c != nil {
			c.Close(200, "Game ended")
		}
	}
}

func (g *GameInstance) loop() {
	fmt.Printf("Starting game %s\n", g.id)
	tc := time.Tick(time.Second / game.Tickrate)
	tcupdate := time.Tick(time.Second / 64)
	playerInputs := make([]game.Inputs, MaxPlayerCount)

	tTick := time.Time{}
	for {
		select {
		case <-tc:
			ts := time.Second / game.Tickrate
			if !tTick.IsZero() {
				ts = time.Since(tTick)
				tTick = time.Now()
			}
			*g.State = g.e.Process(*g.State, playerInputs, ts)
		case <-tcupdate:
			for _, s := range g.spectators {
				b, _ := json.Marshal(g.State)
				err := s.Write(context.Background(), websocket.MessageText, b)
				if err != nil {
					panic(err)
				}
			}
		case inUpdate := <-g.updates:
			if inUpdate.disconnect {
				g.State.Finished = true
				for _, s := range g.spectators {
					b, _ := json.Marshal(g.State)
					err := s.Write(context.Background(), websocket.MessageText, b)
					if err != nil {
						panic(err)
					}
				}
				g.closeConnections()
				g.endCh <- g.id
				fmt.Println("Game finished")
			}
			pid := inUpdate.playerID
			playerInputs[pid] = inUpdate.inputs
			g.State.Players[pid].Inputs.SequenceID = playerInputs[pid].SequenceID
			g.State.Players[pid].Ping = inUpdate.ping
		}
	}
}
