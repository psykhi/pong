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
	//fmt.Println(s.Ball.P)
	c.context().Call("clearRect", 0, 0, c.w, c.h)
	if s.Finished {
		c.drawFinished()
		return
	}
	if s.WaitingForPlayer {
		c.drawWaiting()
		return
	}
	if s.P1.Connected && s.P2.Connected {
		c.drawScore(s.P1.Score, s.P2.Score)
		c.drawPing(s.P1.Ping, s.P2.Ping)
	}
	if s.Paused {
		c.drawCountdown(s.Restart)
	}

	c.drawLine(s.P1)
	c.drawLine(s.P2)
	c.drawBall(s.Ball)
}

func (c *Canvas) drawPing(p1 time.Duration, p2 time.Duration) {
	c.context().Set("font", "20px Arial")
	c.context().Set("fillStyle", "black")
	c.context().Set("textAlign", "left")
	c.context().Call("fillText", formatPing(p1), 0.01, 0.05*c.h)
	c.context().Set("textAlign", "right")
	c.context().Call("fillText", formatPing(p2), c.w, 0.05*c.h)
}

func formatPing(p time.Duration) string {
	return fmt.Sprintf("ping: %dms", p.Milliseconds())
}

func (c *Canvas) drawWaiting() {
	c.context().Set("font", "50px Arial")
	c.context().Set("fillStyle", "black")
	c.context().Set("textAlign", "center")
	c.context().Call("fillText", "Waiting for opponent...", 0.5*c.w, c.h/2)
}

func (c *Canvas) drawFinished() {
	fmt.Println("finished")
	c.context().Set("font", "30px Arial")
	c.context().Set("fillStyle", "black")
	c.context().Set("textAlign", "center")
	c.context().Call("fillText", "Opponent left or connection was lost. Refresh to look for game", 0.5*c.w, c.h/2)
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
	c.context().Set("textAlign", "center")
	c.context().Set("textBaseline", "top")
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
	c.context().Set("textAlign", "center")
	c.context().Set("textBaseline", "top")
	delta := restart.Sub(time.Now()).Seconds()
	if delta < 0 {
		delta = 0
	}
	c.context().Call("fillText", fmt.Sprintf("%.0fs", delta), c.w/2, 0.05*c.h)
}
