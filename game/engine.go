package game

import (
	"math"
	"time"
)

type Inputs struct {
	Up         bool
	Down       bool
	SequenceID int
}

const Tickrate = 128
const StartBallSpeed = 0.45
const PlayerMaxSpeed = 0.8
const BallSpeedIncrease = 0.023

type Engine struct {
}

func (e *Engine) Process(s State, inputs []Inputs, interval time.Duration) State {
	for i, in := range inputs {
		s.Players[i].Inputs = in
	}
	if s.Finished {
		return s
	}
	if s.Paused {
		if time.Now().After(s.Restart) {
			s.Paused = false
		}
	}
	for i, in := range inputs {
		s.Players[i] = e.processPlayerInput(s.Players[i], in, interval)
	}

	e.processBall(&s, interval)
	return s
}

func (e *Engine) processPlayerInput(p Player, in Inputs, interval time.Duration) Player {
	p.Inputs = in
	dy := 0.0
	if in.Up {
		dy += PlayerMaxSpeed * interval.Seconds()
		if p.Top.Y+dy > 1 {
			return p
		}
	}
	if in.Down {
		dy -= PlayerMaxSpeed * interval.Seconds()
		if p.Bottom.Y+dy < 0 {
			return p
		}
	}

	p.Top.Y += dy
	p.Bottom.Y += dy
	return p
}

func (e *Engine) processBall(s *State, interval time.Duration) {
	if s.Paused {
		return
	}
	b := s.Ball

	b.P.Y += s.BallSpeed.Y * interval.Seconds()
	b.P.X += s.BallSpeed.X * interval.Seconds()

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
	for _, p := range s.Players {
		e.collisionsWithPlayer(s, p)
	}

	// end?
	if s.Ball.P.X > 1 {
		s.Players[0].Score++
		s.Countdown()
	}

	if s.Ball.P.X < 0 {
		s.Players[1].Score++
		s.Countdown()

	}
}

func (e *Engine) collisionsWithPlayer(s *State, p Player) {
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

	if d < b.R {
		offsetAngle := p.OffsetRatio(b.P.Y) * math.Pi / 3

		speed := s.BallSpeed.L2() + BallSpeedIncrease*s.BallSpeed.L2()
		s.BallSpeed.Y = 0
		if s.BallSpeed.X < 0 {
			s.BallSpeed.X = speed
			s.BallSpeed = rotate(s.BallSpeed, offsetAngle)

		} else {
			s.BallSpeed.X = -speed
			s.BallSpeed = rotate(s.BallSpeed, -offsetAngle)
		}

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
