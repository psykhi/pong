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
	"time"
)

type client struct {
	updateConn *websocket.Conn
	spectate   *websocket.Conn
	gameID     string
	playerID   int
	StateCh    chan game.State
	PingCh     chan time.Duration
}

type Config struct {
	Port       string `default:"8080"`
	Address    string `default:"localhost"`
	HTTPScheme string `envconfig:"http" default:"http"`
	WsScheme   string `envconfig:"ws" default:"ws"`
}

func NewClient(ch chan game.State, pingCh chan time.Duration) *client {
	resp, err := http.Get(Conf.HTTPScheme + Conf.Address + Conf.Port + "/play")
	if err != nil {
		return nil
	}
	r := server.Response{}

	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		panic(err)
	}
	return &client{
		updateConn: nil,
		spectate:   nil,
		gameID:     r.GameID,
		playerID:   r.PlayerID,
		StateCh:    ch,
		PingCh:     pingCh,
	}
}

func (c *client) Connect() {
	cp := server.ConnectPayload{
		GameID:   c.gameID,
		PlayerID: c.playerID,
	}
	b, _ := json.Marshal(cp)
	// Connect
	urlSpectate := url.URL{Scheme: Conf.WsScheme, Host: Conf.Address + Conf.Port, Path: "/watch"}
	urlMove := url.URL{Scheme: Conf.WsScheme, Host: Conf.Address + Conf.Port, Path: "/game"}
	urlPing := url.URL{Scheme: Conf.WsScheme, Host: Conf.Address + Conf.Port, Path: "/ping"}
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
	// ping
	pingConn, _, err := websocket.Dial(context.Background(), urlPing.String(), nil)
	if err != nil {
		panic(err)
	}
	defer pingConn.Close(200, "end")
	go func() {
		tick := time.Tick(100 * time.Millisecond)
		for {
			select {
			case <-tick:
				start := time.Now()
				ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
				err := pingConn.Write(ctx, websocket.MessageBinary, []byte(""))
				if err != nil {
					c.StateCh <- game.State{Finished: true}
					pingConn.Close(200, "ping timeout")
					return
				}
				_, _, err = pingConn.Read(ctx)
				if err != nil {
					c.StateCh <- game.State{Finished: true}
					pingConn.Close(200, "ping timeout")
					return
				}
				c.PingCh <- time.Since(start)
			}

		}
	}()

	err = pingConn.Write(context.Background(), websocket.MessageText, b)
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			s := game.State{}
			_, b, err := spectateConn.Read(context.Background())
			if err != nil {
				s.Finished = true
				c.StateCh <- s
				return
			}
			err = json.Unmarshal(b, &s)
			//fmt.Println("got state from server", s)
			if err != nil {
				panic(err)
			}
			c.StateCh <- s
		}
	}()
	c.updateConn = inputConn
	c.spectate = spectateConn
}

func (c *client) sendInputs(in game.Inputs, ping time.Duration) {
	//Send inputs to server as well
	b, _ := json.Marshal(server.InputPayload{
		Inputs: in,
		Ping:   ping,
	})
	err := c.updateConn.Write(context.Background(), websocket.MessageText, b)
	if err != nil {
		c.StateCh <- game.State{Finished: true}
	}
	//fmt.Println("Sent inputs!", in)
}
