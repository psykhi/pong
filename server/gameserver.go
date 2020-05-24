package server

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/psykhi/pong/game"
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
}

func NewGameInstance() *GameInstance {
	s := game.NewState()
	s.WaitingForPlayer = true
	return &GameInstance{
		State:      s,
		e:          &game.Engine{},
		spectators: []*websocket.Conn{},
		updates:    make(chan InputUpdate, 2),
	}
}

func (g *GameInstance) ConnectPlayer(playerID int, c *websocket.Conn) {
	pc := &PlayerConn{
		id:       playerID,
		Conn:     c,
		updateCh: g.updates,
	}
	if playerID == 0 {
		fmt.Println("player 1 connected!!")
		g.p1Conn = pc
		g.State.P1.Connected = true
	} else if playerID == 1 {
		fmt.Println("player 2 connected!!")
		g.p2Conn = pc
		g.State.P2.Connected = true
	}
	// Start the game if players are there
	if g.p1Conn != nil && g.p2Conn != nil {
		g.State.WaitingForPlayer = false
		g.State.Countdown()
		// Start the game engine
		go g.loop()
		// Start managing the render connections and send the start signal
		go g.p1Conn.Start()
		go g.p2Conn.Start()
	}
}

func (g *GameInstance) addSpectator(c *websocket.Conn) {
	g.spectators = append(g.spectators, c)
}

func (g *GameInstance) loop() {
	fmt.Println("starting game!")
	tc := time.Tick(time.Second / game.TICKRATE)
	p1In := game.Inputs{}
	p2In := game.Inputs{}
	for {
		select {
		case <-tc:
			*g.State = g.e.Process(*g.State, p1In, p2In)
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
			}
			if inUpdate.playerID == 0 {
				p1In = inUpdate.inputs
			}
			if inUpdate.playerID == 1 {
				p2In = inUpdate.inputs
			}
		}
	}
}
