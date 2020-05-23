package main

import (
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

func (c *Canvas) render(s GameState) {
	c.context().Call("clearRect", 0, 0, c.w, c.h)
	c.drawLine(s.p1)
	c.drawLine(s.p2)
	c.drawBall(s.ball)
}

func (c *Canvas) context() js.Value {
	canvas := js.
		Global().
		Get("document").
		Call("getElementById", c.id)

	return canvas.Call("getContext", "2d")
}

func (c *Canvas) drawLine(l player) {
	c.context().Set("lineWidth", l.width*c.w)
	c.context().Call("beginPath")
	c.context().Call("moveTo", l.top.x*c.w, l.top.y*c.h)
	c.context().Call("lineTo", l.bottom.x*c.w, l.bottom.y*c.h)
	c.context().Call("stroke")
}

func (c *Canvas) drawBall(b ball) {
	c.context().Call("beginPath")
	c.context().Call("arc", b.p.x*c.w, b.p.y*c.h, b.r*c.w, 0, 3*math.Pi)
	c.context().Set("fillStyle", "red")
	c.context().Call("fill")
}
