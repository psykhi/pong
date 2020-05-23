package main

import (
	"math"
	"syscall/js"
)

type Canvas struct {
	w  int
	h  int
	id string
}

func NewCanvas(id string, width, height int) *Canvas {
	canvas := js.
		Global().
		Get("document").
		Call("getElementById", id)

	// reset
	canvas.Set("height", height)
	canvas.Set("width", width)
	//context.Call("clearRect", 0, 0, width, height)
	return &Canvas{w: width, h: height, id: id}
}

func (c *Canvas) render(objects []player, ball ball) {
	c.context().Call("clearRect", 0, 0, c.w, c.h)
	for _, o := range objects {
		c.drawLine(o)
	}
	c.drawBall(ball)
}

func (c *Canvas) context() js.Value {
	canvas := js.
		Global().
		Get("document").
		Call("getElementById", c.id)

	return canvas.Call("getContext", "2d")
}

func (c *Canvas) drawLine(l player) {
	c.context().Set("lineWidth", 10)
	c.context().Call("beginPath")
	c.context().Call("moveTo", l.p2.x*float64(c.w), l.p2.y*float64(c.h))
	c.context().Call("lineTo", l.p1.x*float64(c.w), l.p1.y*float64(c.h))
	c.context().Call("stroke")
}

func (c *Canvas) drawBall(b ball) {
	c.context().Call("beginPath")
	c.context().Call("arc", b.p.x*float64(c.w), b.p.y*float64(c.h), b.r*float64(c.w), 0, 3*math.Pi)
	c.context().Set("fillStyle", "red")
	c.context().Call("fill")

}
