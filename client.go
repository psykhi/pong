package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/psykhi/pong/game"
	"github.com/psykhi/pong/server"
	"net/http"
	"net/url"
	"nhooyr.io/websocket"
)

type client struct {
	inputConn *websocket.Conn
	spectate  *websocket.Conn
	gameID    string
	playerID  int
	StateCh   chan game.State
}

func NewClient(ch chan game.State) *client {
	resp, err := http.Get("https://" + server.ADDRESS + "/play")
	if err != nil {
		panic(err)
	}
	r := server.Response{}
	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		panic(err)
	}
	return &client{
		inputConn: nil,
		spectate:  nil,
		gameID:    r.GameID,
		playerID:  r.PlayerID,
		StateCh:   ch,
	}
}

func (c *client) Connect() {
	cp := server.ConnectPayload{
		GameID:   c.gameID,
		PlayerID: c.playerID,
	}
	b, _ := json.Marshal(cp)
	// Connect
	urlSpectate := url.URL{Scheme: server.WS_SCHEME, Host: server.ADDRESS, Path: "/watch"}
	urlMove := url.URL{Scheme: server.WS_SCHEME, Host: server.ADDRESS, Path: "/game"}
	fmt.Println("Connecting to server..", urlSpectate.String())
	spectateConn, _, err := websocket.Dial(context.Background(), urlSpectate.String(), nil)
	if err != nil {
		panic(err)
	}
	err = spectateConn.Write(context.Background(), websocket.MessageText, b)
	if err != nil {
		panic(err)
	}
	defer spectateConn.Close(200, "end")
	fmt.Println("Connected to server!")

	// input conn
	inputConn, _, err := websocket.Dial(context.Background(), urlMove.String(), nil)
	if err != nil {
		panic(err)
	}

	err = inputConn.Write(context.Background(), websocket.MessageText, b)
	if err != nil {
		panic(err)
	}
	defer inputConn.Close(200, "end")
	fmt.Println("Connected to server for inputs!")

	go func() {
		for {
			s := game.State{}
			_, b, err := spectateConn.Read(context.Background())
			if err != nil {
				panic(err)
			}
			err = json.Unmarshal(b, &s)
			//fmt.Println("got state from server", s)
			if err != nil {
				panic(err)
			}
			c.StateCh <- s
		}

	}()
	c.inputConn = inputConn
	c.spectate = spectateConn
}

func (c *client) sendInputs(in game.Inputs) {
	//Send inputs to server as well
	b, _ := json.Marshal(in)
	err := c.inputConn.Write(context.Background(), websocket.MessageText, b)
	if err != nil {
		panic(err)
	}
	//fmt.Println("Sent inputs!", in)
}
