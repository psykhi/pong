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

type connection struct {
	updateConn   *websocket.Conn
	spectateConn *websocket.Conn
	pingConn     *websocket.Conn
	gameID       string
	playerID     int
	StateCh      chan game.State
	PingCh       chan time.Duration
}

type Config struct {
	Port       string `default:"8080"`
	Address    string `default:"localhost"`
	HTTPScheme string `envconfig:"http" default:"http"`
	WsScheme   string `envconfig:"ws" default:"ws"`
}

func NewClientConnection(ch chan game.State, pingCh chan time.Duration) *connection {
	resp, err := http.Get(Conf.HTTPScheme + Conf.Address + Conf.Port + "/play")
	if err != nil {
		return nil
	}
	r := server.Response{}

	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		panic(err)
	}
	return &connection{
		updateConn:   nil,
		spectateConn: nil,
		gameID:       r.GameID,
		playerID:     r.PlayerID,
		StateCh:      ch,
		PingCh:       pingCh,
	}
}

func (c *connection) close() {
	c.spectateConn.Close(200, "client disconnecting")
	c.updateConn.Close(200, "client disconnecting")
	c.pingConn.Close(200, "client disconnecting")
}

func (c *connection) Connect() {
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
				c.close()
				s.Finished = true
				c.StateCh <- s
				return
			}
			err = json.Unmarshal(b, &s)
			//fmt.Println("got state from server", state)
			if err != nil {
				panic(err)
			}
			c.StateCh <- s
		}
	}()
	c.updateConn = inputConn
	c.spectateConn = spectateConn
	c.pingConn = pingConn
}

func (c *connection) sendInputs(in *game.Inputs, ping time.Duration) {
	b, _ := json.Marshal(server.InputPayload{
		Inputs: *in,
		Ping:   ping,
	})
	err := c.updateConn.Write(context.Background(), websocket.MessageText, b)
	if err != nil {
		c.close()
		c.StateCh <- game.State{Finished: true}
	}
}
