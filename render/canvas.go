package render

import (
	"fmt"
	"github.com/psykhi/pong/game"
	"math"
	"syscall/js"
	"time"
)

type Canvas struct {
	W  float64
	H  float64
	id string
}

const TextColor = "white"

//const BackgroundColor = "black"
const BallColor = "white"

func NewCanvas(id string, width, height float64) *Canvas {
	canvas := js.
		Global().
		Get("document").
		Call("getElementById", id)

	canvas.Set("height", height)
	canvas.Set("width", width)
	return &Canvas{W: width, H: height, id: id}
}

func (c *Canvas) Resize(w, h float64) {
	canvas := js.
		Global().
		Get("document").
		Call("getElementById", c.id)

	canvas.Set("height", w)
	canvas.Set("width", h)
	c.W = w
	c.H = h

}

func (c *Canvas) Render(s game.State) {
	c.context().Call("clearRect", 0, 0, c.W, c.H)
	if s.Finished {
		c.drawFinished()
		return
	}
	if s.WaitingForPlayer {
		c.drawWaiting()
		return
	}
	if s.Players[0].Connected && s.Players[1].Connected {
		c.drawScore(s.Players[0].Score, s.Players[1].Score)
		c.drawPing(s.Players[0].Ping, s.Players[1].Ping)
	}
	if s.Paused {
		c.drawCountdown(s.Restart)
	}

	c.drawLine(s.Players[0])
	c.drawLine(s.Players[1])
	c.drawBall(s.Ball)
}

func (c *Canvas) drawPing(p1 time.Duration, p2 time.Duration) {
	c.context().Set("font", "20px Arial")
	c.context().Set("textBaseline", "top")
	c.context().Set("fillStyle", TextColor)
	c.context().Set("textAlign", "left")
	c.context().Call("fillText", formatPing(p1), 0.0, 0)
	c.context().Set("textAlign", "right")
	c.context().Call("fillText", formatPing(p2), c.W, 0)
}

func formatPing(p time.Duration) string {
	return fmt.Sprintf("ping: %dms", p.Milliseconds())
}

func (c *Canvas) drawWaiting() {
	c.context().Set("font", "50px Arial")
	c.context().Set("fillStyle", TextColor)
	c.context().Set("textAlign", "center")
	c.context().Call("fillText", "Waiting for opponent...", 0.5*c.W, c.H/2)
}

func (c *Canvas) drawFinished() {
	c.context().Set("font", "30px Arial")
	c.context().Set("fillStyle", TextColor)
	c.context().Set("textAlign", "center")
	c.context().Call("fillText", "Opponent left or connection was lost. Refresh to look for game", 0.5*c.W, c.H/2)
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
	c.context().Set("fillStyle", TextColor)
	c.context().Set("textAlign", "center")
	c.context().Set("textBaseline", "top")
	c.context().Call("fillText", p1, c.W/4, 0.05*c.H)
	c.context().Call("fillText", p2, 3*c.W/4, 0.05*c.H)
}

func (c *Canvas) drawLine(l game.Player) {
	c.context().Set("lineWidth", l.Width*c.W)
	c.context().Set("strokeStyle", BallColor)
	c.context().Call("beginPath")
	c.context().Call("moveTo", l.Top.X*c.W, l.Top.Y*c.H)
	c.context().Call("lineTo", l.Bottom.X*c.W, l.Bottom.Y*c.H)
	c.context().Call("stroke")
}

func (c *Canvas) drawBall(b game.Ball) {
	c.context().Call("beginPath")
	c.context().Call("arc", b.P.X*c.W, b.P.Y*c.H, b.R*c.W, 0, 3*math.Pi)
	c.context().Set("fillStyle", BallColor)
	c.context().Call("fill")
}

func (c *Canvas) drawCountdown(restart time.Time) {
	c.context().Set("font", "30px Arial")
	c.context().Set("fillStyle", TextColor)
	c.context().Set("textAlign", "center")
	c.context().Set("textBaseline", "top")
	delta := restart.Sub(time.Now()).Seconds()
	if delta < 0 {
		delta = 0
	}
	c.context().Call("fillText", fmt.Sprintf("%.0fs", delta), c.W/2, 0.05*c.H)
}
