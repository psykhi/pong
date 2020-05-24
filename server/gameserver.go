package server

import (
	"github.com/gorilla/websocket"
	"github.com/psykhi/pong/game"
	"time"
)

type GameInstance struct {
	*game.State
	p1Conn     *PlayerConn
	p2Conn     *PlayerConn
	e          *game.Engine
	spectators []*websocket.Conn
}

func NewGameInstance() *GameInstance {
	return &GameInstance{
		State:      game.NewState(),
		e:          &game.Engine{},
		spectators: []*websocket.Conn{}}
}

func (g *GameInstance) ConnectPlayer(playerID int, c *websocket.Conn) {
	pc := &PlayerConn{
		Conn: c,
	}
	if playerID == 0 {
		g.p1Conn = pc
	} else if playerID == 1 {
		g.p2Conn = pc
	}
	// Start the game if players are there
	if g.p1Conn != nil && g.p2Conn != nil {
		// Start the game engine
		go g.loop()
		// Start managing the client connections and send the start signal
	}
}

func (g *GameInstance) addSpectator(c *websocket.Conn) {
	g.spectators = append(g.spectators, c)
}

func (g *GameInstance) loop() {
	tc := time.Tick(time.Second / 128)
	for {
		select {
		case <-tc:
			*g.State = g.e.Process(*g.State, g.p1Conn.in, g.p2Conn.in)
			for _, s := range g.spectators {
				err := s.WriteJSON(g.State)
				if err != nil {
					panic(err)
				}
			}
		}
	}
}
