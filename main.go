package main

import (
	"container/list"
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
		//fmt.Println("render", time.Now().UnixNano())
		return nil
	})
	defer render.Release()

	sCh := make(chan game.State)

	client := NewClient(sCh)
	client.Connect()

	go func() {
		//lag := 2
		tempSeverState := &game.State{}
		*tempSeverState = *s
		t := time.Tick(time.Second / game.TICKRATE)
		tTick := time.Time{}
		for {
			select {
			case <-t:
				ts := time.Second / game.TICKRATE
				if !tTick.IsZero() {
					ts = time.Since(tTick)
					tTick = time.Now()
				}
				if client.playerID == 0 {
					*s = e.Process(*s, inputs, s.P2.Inputs, ts)
				}
				if client.playerID == 1 {
					*s = e.Process(*s, s.P1.Inputs, inputs, ts)
				}
				//fmt.Println(s.Ball.P.X)
				//fmt.Println(s.Ball.P.Y)
				//fmt.Println(s.Ball.P.Y)
				//fmt.Println("local render")
				js.Global().Call("requestAnimationFrame", render)
				inputs.SequenceID++
				client.sendInputs(inputs)
				//case <-sCh:
			case serverState := <-sCh:
				limit := 8
				reject := false
				if client.playerID == 0 {
					if inputs.SequenceID-serverState.P1.Inputs.SequenceID > limit {
						fmt.Println(serverState.P1.Inputs.SequenceID, s.P1.Inputs.SequenceID)
						reject = true
					} else {
						//fmt.Println("ok")
					}
				}
				if client.playerID == 1 {
					if inputs.SequenceID-serverState.P2.Inputs.SequenceID > limit {
						fmt.Println(serverState.P2.Inputs.SequenceID, s.P2.Inputs.SequenceID)
						reject = true
					} else {
						//fmt.Println("ok")
					}
				}

				if !reject {
					//*s = *tempSeverState
					//*tempSeverState = serverState
					*s = serverState
					js.Global().Call("requestAnimationFrame", render)
				}

			}
		}
	}()

	js.Global().Call("requestAnimationFrame", render)

	<-done
}

type ServerStates struct {
	states list.List
	size   int
}

func (ss *ServerStates) Add(s game.State) {
	ss.states.PushBack(s)
}
