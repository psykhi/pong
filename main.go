package main

import (
	"fmt"
	"github.com/psykhi/pong/client"
	"github.com/psykhi/pong/game"
	"syscall/js"
	"time"
)

const ticks = 120

var done = make(chan struct{})

func main() {
	fmt.Println("Start")
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

	go func() {
		t := time.Tick(time.Second / ticks)
		for {
			select {
			case <-t:
				*s = e.Process(*s, inputs, inputs)
				js.Global().Call("requestAnimationFrame", render)
			}
		}
	}()

	js.Global().Call("requestAnimationFrame", render)

	<-done
}
