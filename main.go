package main

import (
	"fmt"
	"math"
	"syscall/js"
	"time"
)

type position struct {
	x float64
	y float64
}

func (p position) L2() float64 {
	return math.Sqrt(p.x*p.x + p.y*p.y)
}

type ball struct {
	p position
	r float64
}

const ticks = 120

var done = make(chan struct{})

type inputs struct {
	up   bool
	down bool
}

func main() {
	fmt.Println("Start")
	doc := js.Global().Get("document")
	width := doc.Get("body").Get("clientWidth").Float()
	height := doc.Get("body").Get("clientHeight").Float()
	c := NewCanvas("canvas", width, height)
	e := Engine{}

	s := NewState()
	inputs := inputs{}
	keyDown := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		e := args[0]
		code := e.Get("keyCode")
		switch code.Int() {
		case 38:
			inputs.down = true
		case 40:
			inputs.up = true
		}
		return nil
	})
	defer keyDown.Release()

	keyUp := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		e := args[0]
		code := e.Get("keyCode")
		switch code.Int() {
		case 38:
			inputs.down = false
		case 40:
			inputs.up = false
		}
		return nil
	})
	defer keyUp.Release()

	js.Global().Get("document").Call("addEventListener", "keydown", keyDown)
	js.Global().Get("document").Call("addEventListener", "keyup", keyUp)
	render := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		c.render(s)
		return nil
	})
	defer render.Release()

	go func() {
		t := time.Tick(time.Second / ticks)
		for {
			select {
			case <-t:
				s = e.Process(s, inputs, inputs)
				//c.render(s)
				js.Global().Call("requestAnimationFrame", render)
			}
		}
	}()

	js.Global().Call("requestAnimationFrame", render)

	<-done
}
