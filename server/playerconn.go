package server

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/psykhi/pong/game"
	"nhooyr.io/websocket"
	"time"
)

type InputUpdate struct {
	playerID   int
	inputs     game.Inputs
	ping       time.Duration
	disconnect bool
}

type PlayerConn struct {
	id int
	*websocket.Conn
	ip       InputPayload
	updateCh chan InputUpdate
}

func (pc *PlayerConn) Start() {
	for {
		_, b, err := pc.Read(context.Background())
		if err != nil {
			pc.OnDisconnect(err)
			pc.Close(400, err.Error())
			return
		}
		ip := InputPayload{}
		err = json.Unmarshal(b, &ip)
		if err != nil {
			panic(err)
		}
		pc.ip = ip
		pc.updateCh <- InputUpdate{
			playerID: pc.id,
			inputs:   pc.ip.Inputs,
			ping:     ip.Ping,
		}
	}
}

func (pc *PlayerConn) OnDisconnect(err error) {
	fmt.Printf("Player disconnected from game: %s\n", err)
	pc.updateCh <- InputUpdate{
		playerID:   pc.id,
		inputs:     pc.ip.Inputs,
		disconnect: true,
	}
}
