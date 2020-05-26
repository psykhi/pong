package server

import (
	"github.com/psykhi/pong/game"
	"time"
)

type ConnectPayload struct {
	GameID   string
	PlayerID int
}

type InputPayload struct {
	Inputs game.Inputs
	Ping   time.Duration
}
