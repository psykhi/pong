package main

import (
	"fmt"
	"syscall/js"
	"time"
)

type position struct {
	x float64
	y float64
}

type ball struct {
	p position
	r float64
}

type player struct {
	top    position
	bottom position
}

const ticks = 60

var done = make(chan struct{})

func initState() GameState {
	p1 := player{
		top: position{
			x: 0,
			y: 0.4,
		}, bottom: position{
			x: 0,
			y: 0.6,
		},
	}
	p2 := player{
		top: position{
			x: 1,
			y: 0.4,
		}, bottom: position{
			x: 1,
			y: 0.6,
		},
	}
	ball := ball{
		p: position{
			x: 0.5,
			y: 0.5,
		},
		r: 0.01,
	}
	speed := position{
		x: 0.01,
		y: 0,
	}

	return GameState{
		ball:      ball,
		p1:        p1,
		p2:        p2,
		ballSpeed: speed,
	}

}

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

	s := initState()
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

	go func() {
		t := time.Tick(time.Second / ticks)
		for {
			select {
			case <-t:
				s = e.Process(s, inputs, inputs)
				c.render(s)
			}
		}
	}()

	<-done
}

func input() {

}
