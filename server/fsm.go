package server

type GameState int

const (
	WAITING_FOR_PLAYERS GameState = iota
	IN_GAME
	DONE
)

type GameFSM struct {
	state GameState
}
