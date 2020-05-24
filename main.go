package main

import (
	"fmt"
	"github.com/psykhi/pong/game"
	"github.com/psykhi/pong/render"
	"syscall/js"
	"time"
)

var done = make(chan struct{})

func main() {
	fmt.Println("BLAB Start")
	doc := js.Global().Get("document")
	width := doc.Get("body").Get("clientWidth").Float()
	height := doc.Get("body").Get("clientHeight").Float()
	c := render.NewCanvas("canvas", width, height)
	e := game.Engine{}

	s := game.NewState()
	s.WaitingForPlayer = true
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
		//fmt.Println("render")
		return nil
	})
	defer render.Release()

	sCh := make(chan game.State)

	client := NewClient(sCh)
	client.Connect()

	go func() {
		t := time.Tick(time.Second / game.TICKRATE)
		for {
			select {
			case <-t:
				if client.playerID == 0 {
					*s = e.Process(*s, inputs, game.Inputs{})
				}
				if client.playerID == 1 {
					*s = e.Process(*s, game.Inputs{}, inputs)
				}
				//fmt.Println(s.Ball.P.X)
				//fmt.Println(s.Ball.P.Y)
				//fmt.Println(s.Ball.P.Y)
				//fmt.Println("local render")
				js.Global().Call("requestAnimationFrame", render)
				client.sendInputs(inputs)
			case serverState := <-sCh:
				//fmt.Println("Got server state")
				//fmt.Println(serverState)
				*s = serverState
				js.Global().Call("requestAnimationFrame", render)

			}
		}
	}()

	js.Global().Call("requestAnimationFrame", render)

	<-done
}
