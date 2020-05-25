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

type GameInstance struct {
	*game.State
	p1Conn     *PlayerConn
	p2Conn     *PlayerConn
	e          *game.Engine
	spectators []*websocket.Conn
	updates    chan InputUpdate
	endCh      chan string
	id         string
}

func NewGameInstance(endCh chan string) *GameInstance {
	s := game.NewState()
	s.WaitingForPlayer = true
	return &GameInstance{
		State:      s,
		e:          &game.Engine{},
		spectators: []*websocket.Conn{},
		updates:    make(chan InputUpdate, 2),
		endCh:      endCh,
		id:         uuid.NewV4().String(),
	}
}

func (g *GameInstance) ConnectPlayer(playerID int, c *websocket.Conn) {
	pc := &PlayerConn{
		id:       playerID,
		Conn:     c,
		updateCh: g.updates,
	}
	if playerID == 0 {
		fmt.Printf("Player 1 connected on game %s\n", g.id)
		g.p1Conn = pc
		g.State.P1.Connected = true
		go g.p1Conn.Start()
	} else if playerID == 1 {
		fmt.Printf("Player 2 connected on game %s\n", g.id)
		g.p2Conn = pc
		g.State.P2.Connected = true
		go g.p2Conn.Start()
	}
	// Start the game if players are there
	if g.p1Conn != nil && g.p2Conn != nil {
		g.State.WaitingForPlayer = false
		g.State.Countdown()
		// Start the game engine
		go g.loop()
		// Start managing the render connections and send the start signal

	}
}

func (g *GameInstance) addSpectator(c *websocket.Conn) {
	g.spectators = append(g.spectators, c)
}

func (g *GameInstance) loop() {
	fmt.Printf("Starting game %s\n", g.id)
	tc := time.Tick(time.Second / game.TICKRATE)
	tcupdate := time.Tick(time.Second / 64)
	p1In := game.Inputs{}
	p2In := game.Inputs{}
	tTick := time.Time{}
	for {
		select {
		case <-tc:
			ts := time.Second / game.TICKRATE
			if !tTick.IsZero() {
				ts = time.Since(tTick)
				tTick = time.Now()
			}
			*g.State = g.e.Process(*g.State, p1In, p2In, ts)
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
				g.endCh <- g.id
				fmt.Println("Game finished")
			}
			if inUpdate.playerID == 0 {
				p1In = inUpdate.inputs
				g.State.P1.Inputs.SequenceID = p1In.SequenceID
			}
			if inUpdate.playerID == 1 {
				p2In = inUpdate.inputs
				g.State.P2.Inputs.SequenceID = p2In.SequenceID
			}
		}
	}
}
