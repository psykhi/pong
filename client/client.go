package main

import (
	"fmt"
	"github.com/psykhi/pong/game"
	"github.com/psykhi/pong/render"
	"math"
	"syscall/js"
	"time"
)

type touch struct {
	touch bool
	x     float64
	y     float64
}

type client struct {
	e            *game.Engine
	c            *render.Canvas
	s            *game.State
	inputs       *game.Inputs
	doc          js.Value
	funcs        []js.Func
	render       js.Func
	width        float64
	height       float64
	conn         *connection
	touch        touch
	sCh          chan game.State
	pingCh       chan time.Duration
	inputManager *InputManager
}

func NewClient() *client {
	doc := js.Global().Get("document")
	width := doc.Get("body").Get("clientWidth").Float()
	height := doc.Get("body").Get("clientHeight").Float()
	c := render.NewCanvas("canvas", width, height)
	e := &game.Engine{}

	s := game.NewState()
	s.WaitingForPlayer = true
	inputs := &game.Inputs{}
	cl := &client{
		e:            e,
		c:            c,
		s:            s,
		inputs:       inputs,
		doc:          doc,
		funcs:        []js.Func{},
		sCh:          make(chan game.State),
		pingCh:       make(chan time.Duration),
		inputManager: NewInputManager(100),
	}

	w, h := cl.Dimensions()
	cl.width = w
	cl.height = h
	conn := NewClientConnection(cl.sCh, cl.pingCh)
	if conn == nil {
		cl.s.Finished = true
		js.Global().Call("requestAnimationFrame", cl.render)
		select {}
	} else {
		conn.Connect()
	}
	cl.conn = conn

	doRender := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		cl.c.Render(*cl.s)
		return nil
	})
	cl.render = doRender
	cl.registerKeyboard()
	cl.registerTouch()
	cl.registerResize()

	return cl
}
func (cl *client) Stop() {
	cl.render.Release()
	for _, f := range cl.funcs {
		f.Release()
	}
}

func (cl *client) UpdateTouchInput(x, y float64) {
	xp := x / cl.width
	yp := y / cl.height
	cl.touch.x = xp
	cl.touch.y = yp
}

func (cl *client) registerKeyboard() {
	keyUpDown := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		e := args[0]
		code := e.Get("keyCode")
		typ := e.Get("type").String()
		switch code.Int() {
		case 38:
			cl.inputs.Down = typ == "keydown"
		case 40:
			cl.inputs.Up = typ == "keydown"
		}
		return nil
	})

	js.Global().Get("document").Call("addEventListener", "keydown", keyUpDown)
	js.Global().Get("document").Call("addEventListener", "keyup", keyUpDown)
	cl.funcs = append(cl.funcs, keyUpDown)
}

func (cl *client) registerTouch() {
	onTouchMove := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		e := args[0]
		x := e.Get("touches").Index(0).Get("clientX").Float()
		y := e.Get("touches").Index(0).Get("clientY").Float()
		cl.UpdateTouchInput(x, y)
		cl.touch.touch = true
		return nil
	})
	onTouchStart := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		e := args[0]
		x := e.Get("touches").Index(0).Get("clientX").Float()
		y := e.Get("touches").Index(0).Get("clientY").Float()
		cl.UpdateTouchInput(x, y)
		cl.touch.touch = true
		return nil
	})
	onTouchEnd := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		cl.inputs.Up = false
		cl.inputs.Down = false
		cl.touch.touch = false
		return nil
	})

	js.Global().Get("document").Call("addEventListener", "touchmove", onTouchMove)
	js.Global().Get("document").Call("addEventListener", "touchstart", onTouchStart)
	js.Global().Get("document").Call("addEventListener", "touchend", onTouchEnd)

	cl.funcs = append(cl.funcs, onTouchStart)
	cl.funcs = append(cl.funcs, onTouchMove)
	cl.funcs = append(cl.funcs, onTouchEnd)
}

func (cl *client) getInputs() *game.Inputs {
	if !cl.touch.touch {
		return cl.inputs
	}
	p := cl.s.Players[cl.conn.playerID]

	if math.Abs(cl.touch.y-p.Center()) < 0.01 {
		cl.inputs.Down = false
		cl.inputs.Up = false
		return cl.inputs
	}
	if cl.touch.y > p.Center() {
		cl.inputs.Up = true
		cl.inputs.Down = false
	} else {
		cl.inputs.Down = true
		cl.inputs.Up = false
	}
	return cl.inputs
}
func (cl *client) Dimensions() (float64, float64) {
	w := js.Global().Get("innerWidth").Float()
	h := js.Global().Get("innerHeight").Float()
	return w, h
}

func (cl *client) registerResize() {
	onResize := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		w, h := cl.Dimensions()
		cl.c.Resize(w, h)
		fmt.Println(w, h)
		cl.width = w
		cl.height = h
		js.Global().Call("requestAnimationFrame", cl.render)
		return nil
	})
	js.Global().Call("addEventListener", "resize", onResize)
	cl.funcs = append(cl.funcs, onResize)
}

func (cl *client) replay(s game.State, from int) game.State {
	state := s
	delta := time.Second / game.Tickrate
	for sid := from; sid < cl.inputs.SequenceID; sid++ {
		in, err := cl.inputManager.Get(sid)
		if err != nil {
			// give up, we're replaying from too far
			return s
		}
		inputs := cl.s.Inputs()
		inputs[cl.conn.playerID] = in
		state = cl.e.Process(state, inputs, delta)
	}
	return state
}

func (cl *client) Start() {
	go func() {
		gameTick := time.Tick(time.Second / game.Tickrate)
		frameTs := time.Time{}
		ping := 0 * time.Millisecond
		//tsLastValidState := time.Now().Add(24 * time.Hour)
		for {
			select {
			case <-gameTick:
				tSinceLastFrame := time.Second / game.Tickrate
				if !frameTs.IsZero() {
					tSinceLastFrame = time.Since(frameTs)
					frameTs = time.Now()
				}
				// Process movement based on the local inputs, not what the server sees
				inputs := cl.s.Inputs()
				inputs[cl.conn.playerID] = *cl.getInputs()

				*cl.s = cl.e.Process(*cl.s, inputs, tSinceLastFrame)

				js.Global().Call("requestAnimationFrame", cl.render)

				in := cl.getInputs()
				cl.inputManager.Set(*in)
				cl.conn.sendInputs(in, ping)
				in.SequenceID++
			case p := <-cl.pingCh:
				ping = p
			case serverState := <-cl.sCh:
				//Reconcile server input with local state
				*cl.s = cl.replay(serverState, serverState.Players[cl.conn.playerID].Inputs.SequenceID)
				js.Global().Call("requestAnimationFrame", cl.render)
			}
		}
	}()

	js.Global().Call("requestAnimationFrame", cl.render)
}
