package client

import (
	"github.com/psykhi/pong/game"
	"math"
	"syscall/js"
)

type Canvas struct {
	w  float64
	h  float64
	id string
}

func NewCanvas(id string, width, height float64) *Canvas {
	canvas := js.
		Global().
		Get("document").
		Call("getElementById", id)

	canvas.Set("height", height)
	canvas.Set("width", width)
	return &Canvas{w: width, h: height, id: id}
}

func (c *Canvas) Render(s game.State) {
	c.context().Call("clearRect", 0, 0, c.w, c.h)
	c.drawLine(s.P1)
	c.drawLine(s.P2)
	c.drawBall(s.Ball)
}

func (c *Canvas) context() js.Value {
	canvas := js.
		Global().
		Get("document").
		Call("getElementById", c.id)

	return canvas.Call("getContext", "2d")
}

func (c *Canvas) drawLine(l game.Player) {
	c.context().Set("lineWidth", l.Width*c.w)
	c.context().Call("beginPath")
	c.context().Call("moveTo", l.Top.X*c.w, l.Top.Y*c.h)
	c.context().Call("lineTo", l.Bottom.X*c.w, l.Bottom.Y*c.h)
	c.context().Call("stroke")
}

func (c *Canvas) drawBall(b game.Ball) {
	c.context().Call("beginPath")
	c.context().Call("arc", b.P.X*c.w, b.P.Y*c.h, b.R*c.w, 0, 3*math.Pi)
	c.context().Set("fillStyle", "red")
	c.context().Call("fill")
}
