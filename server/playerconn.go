package server

import (
	"github.com/gorilla/websocket"
	"github.com/psykhi/pong/game"
)

type InputUpdate struct {
	playerID int
	inputs   game.Inputs
}

type PlayerConn struct {
	id int
	*websocket.Conn
	in       game.Inputs
	updateCh chan InputUpdate
}

func (pc *PlayerConn) Start() {
	err := pc.WriteJSON(ActionPayload{message: "start"})
	if err != nil {
		panic(err)
	}

	for {
		err := pc.ReadJSON(&pc.in)
		if err != nil {
			panic(err)
		}
		pc.updateCh <- InputUpdate{
			playerID: pc.id,
			inputs:   pc.in,
		}
	}
}
