package render

import (
	"fmt"
	"github.com/psykhi/pong/game"
	"math"
	"syscall/js"
	"time"
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
	if s.WaitingForPlayer {
		c.drawWaiting()
		return
	}
	if s.P1.Connected && s.P2.Connected {
		c.drawScore(s.P1.Score, s.P2.Score)
	}
	if s.Paused {
		c.drawCountdown(s.Restart)
	}
	if s.Finished {
		c.drawFinished()
		return
	}

	c.drawLine(s.P1)
	c.drawLine(s.P2)
	c.drawBall(s.Ball)
}

func (c *Canvas) drawWaiting() {
	c.context().Set("font", "50px Arial")
	c.context().Set("fillStyle", "black")
	c.context().Call("fillText", "Waiting for opponent...", 0.4*c.w, c.h/2)
}

func (c *Canvas) drawFinished() {
	c.context().Set("font", "50px Arial")
	c.context().Set("fillStyle", "black")
	c.context().Call("fillText", "Opponent left. Refresh to look for game", 0.5*c.w/2, c.h/2)
}

func (c *Canvas) context() js.Value {
	canvas := js.
		Global().
		Get("document").
		Call("getElementById", c.id)

	return canvas.Call("getContext", "2d")
}

func (c *Canvas) drawScore(p1 int, p2 int) {
	c.context().Set("font", "50px Arial")
	c.context().Set("fillStyle", "black")
	c.context().Call("fillText", p1, c.w/4, 0.05*c.h)
	c.context().Call("fillText", p2, 3*c.w/4, 0.05*c.h)
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

func (c *Canvas) drawCountdown(restart time.Time) {
	c.context().Set("font", "30px Arial")
	c.context().Set("fillStyle", "black")
	c.context().Call("fillText", fmt.Sprintf("%.1f", restart.Sub(time.Now()).Seconds())+"s", c.w/2, 0.05*c.h)
}
