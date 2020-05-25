package game

import (
	"math"
	"math/rand"
	"time"
)

type Inputs struct {
	Up   bool
	Down bool
}

const TICKRATE = 128

type Engine struct {
}

func (e *Engine) Process(s State, p1in Inputs, p2in Inputs) State {
	if s.Finished {
		return s
	}

	if s.Paused {
		s.BallSpeed = Position{}
		if time.Now().After(s.Restart) {
			s.Paused = false
			s.BallSpeed = Position{
				X: 0.005,
			}
			s.BallSpeed = rotate(s.BallSpeed, rand.Float64()*math.Pi)
		}
	}

	s.P1 = e.processPlayerInput(s.P1, p1in)
	s.P2 = e.processPlayerInput(s.P2, p2in)
	e.processBall(&s)
	return s
}

func (e *Engine) processPlayerInput(p Player, in Inputs) Player {
	dy := 0.0
	if in.Up {
		dy += 0.01
		if p.Top.Y+dy > 1 {
			return p
		}
	}
	if in.Down {
		dy -= 0.01
		if p.Bottom.Y+dy < 0 {
			return p
		}
	}

	p.Top.Y += dy
	p.Bottom.Y += dy

	// Now verify
	//fmt.Println(p.top.y)
	return p
}

func (e *Engine) processBall(s *State) {
	b := s.Ball

	b.P.Y += s.BallSpeed.Y
	b.P.X += s.BallSpeed.X

	// Top and top wall collisions
	if b.P.Y+b.R > 1 {
		s.BallSpeed.Y = -s.BallSpeed.Y
		b.P.Y = 1 - b.R
		return
	}
	if b.P.Y-b.R < 0 {
		s.BallSpeed.Y = -s.BallSpeed.Y
		b.P.Y = b.R
		return
	}

	// Player collisions?

	s.Ball = b
	e.CollisionsWithPlayer(s, s.P1)
	e.CollisionsWithPlayer(s, s.P2)

	// end?

	if s.Ball.P.X > 1 {
		s.P1.Score++
		s.Countdown()
	}

	if s.Ball.P.X < 0 {
		s.P2.Score++
		s.Countdown()
	}
}

func (e *Engine) CollisionsWithPlayer(s *State, p Player) {
	b := s.Ball
	tx := b.P.X
	ty := b.P.Y
	if b.P.X > p.Top.X {
		tx = p.Top.X + p.Width/2
	} else if b.P.X < p.Top.X {
		tx = p.Top.X - p.Width/2
	}

	if b.P.Y > p.Top.Y {
		ty = p.Top.Y
	} else if b.P.Y < p.Bottom.Y {
		ty = p.Bottom.Y
	}
	dx := tx - b.P.X
	dy := ty - b.P.Y

	d := math.Sqrt(dx*dx + dy*dy)
	//fmt.Println("dx dy", dx, dy)
	//fmt.Println("DISTANCE", d)
	//fmt.Println("BALL", b)
	//fmt.Println("BOTOOM", p.top)
	//fmt.Println("TOP", p.bottom)
	if d < b.R {
		//fmt.Println("COLLISION", b.p.y, p.bottom.y, p.offsetRatio(b.p.y))

		offsetAngle := p.OffsetRatio(b.P.Y) * math.Pi / 3
		//fmt.Println("added angle", offsetAngle*(180/math.Pi))

		speed := s.BallSpeed.L2() + 0.02*s.BallSpeed.L2()
		s.BallSpeed.Y = 0
		if s.BallSpeed.X < 0 {
			s.BallSpeed.X = speed
			s.BallSpeed = rotate(s.BallSpeed, offsetAngle)

		} else {
			s.BallSpeed.X = -speed
			s.BallSpeed = rotate(s.BallSpeed, -offsetAngle)
		}

		//inAngle := math.Atan(s.BallSpeed.y / s.BallSpeed.x)
		//fmt.Println("angle", inAngle*(180/math.Pi))
		//s.BallSpeed.y = p.offsetRatio(b.p.y) * 0.005
		if b.P.X+b.R > p.Bottom.X+p.Width/2 {
			b.P.X = p.Bottom.X - p.Width/2 - b.R
		}
		if b.P.X-b.R < p.Bottom.X-p.Width/2 {
			b.P.X = p.Bottom.X + p.Width/2 + b.R
		}
	}
}

func rotate(v Position, angle float64) Position {
	x2 := v.X*math.Cos(angle) - v.Y*math.Sin(angle)
	y2 := v.X*math.Sin(angle) + v.Y*math.Cos(angle)
	return Position{
		X: x2,
		Y: y2,
	}
}
