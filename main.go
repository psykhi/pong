package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/psykhi/pong/client"
	"github.com/psykhi/pong/game"
	"net/url"
	"syscall/js"
	"time"
)

var done = make(chan struct{})

func main() {
	fmt.Println("BLAB Start")
	doc := js.Global().Get("document")
	width := doc.Get("body").Get("clientWidth").Float()
	height := doc.Get("body").Get("clientHeight").Float()
	c := client.NewCanvas("canvas", width, height)
	e := game.Engine{}

	s := game.NewState()
	inputs := game.Inputs{}
	keyDown := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		e := args[0]
		code := e.Get("keyCode")
		switch code.Int() {
		case 38:
			inputs.Down = true
		case 40:
			inputs.Up = true
		}
		return nil
	})
	defer keyDown.Release()

	keyUp := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		e := args[0]
		code := e.Get("keyCode")
		switch code.Int() {
		case 38:
			inputs.Down = false
		case 40:
			inputs.Up = false
		}
		return nil
	})
	defer keyUp.Release()

	js.Global().Get("document").Call("addEventListener", "keydown", keyDown)
	js.Global().Get("document").Call("addEventListener", "keyup", keyUp)
	render := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		c.Render(*s)
		return nil
	})
	defer render.Release()

	// Connect
	u := url.URL{Scheme: "ws", Host: "localhost:3010", Path: "/game"}

	fmt.Println("Connecting to server..", u.String())
	spectateConn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		panic(err)
	}
	defer spectateConn.Close()
	fmt.Println("Connected to server!")

	sCh := make(chan game.State)

	go func() {
		s := game.State{}
		err := spectateConn.ReadJSON(&s)
		fmt.Println("got state from server", s)
		if err != nil {
			panic(err)
		}
		sCh <- s
	}()

	go func() {
		t := time.Tick(time.Second / game.TICKRATE)
		for {
			select {
			case <-t:
				*s = e.Process(*s, inputs, inputs)
				js.Global().Call("requestAnimationFrame", render)
			case serverState := <-sCh:
				fmt.Println("Got server state")
				*s = serverState
			}
		}
	}()

	js.Global().Call("requestAnimationFrame", render)

	<-done
}
