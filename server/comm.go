package server

type ConnectPayload struct {
	GameID   string
	PlayerID int
}

type ActionPayload struct {
	message string
}
