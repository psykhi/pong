package server

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/psykhi/pong/game"
	"nhooyr.io/websocket"
)

type InputUpdate struct {
	playerID   int
	inputs     game.Inputs
	disconnect bool
}

type PlayerConn struct {
	id int
	*websocket.Conn
	in       game.Inputs
	updateCh chan InputUpdate
}

func (pc *PlayerConn) Start() {
	for {
		_, b, err := pc.Read(context.Background())
		if err != nil {
			// player disconnected!
			fmt.Printf("Player disconnected from game\n")
			pc.updateCh <- InputUpdate{
				playerID:   pc.id,
				inputs:     pc.in,
				disconnect: true,
			}
			return
		}
		err = json.Unmarshal(b, &pc.in)
		if err != nil {
			panic(err)
		}
		//fmt.Println("Got player inputs")
		pc.updateCh <- InputUpdate{
			playerID: pc.id,
			inputs:   pc.in,
		}
	}
}
