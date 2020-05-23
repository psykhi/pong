package game

import (
	"math"
)

type Inputs struct {
	Up   bool
	Down bool
}

type Engine struct {
}

func (e *Engine) Process(s State, p1in Inputs, p2in Inputs) State {
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

	b.P.Y += s.ballSpeed.Y
	b.P.X += s.ballSpeed.X

	// Top and top wall collisions
	if b.P.Y+b.R > 1 {
		s.ballSpeed.Y = -s.ballSpeed.Y
		b.P.Y = 1 - b.R
		return
	}
	if b.P.Y-b.R < 0 {
		s.ballSpeed.Y = -s.ballSpeed.Y
		b.P.Y = b.R
		return
	}

	// Player collisions?

	s.Ball = b
	e.CollisionsWithPlayer(s, s.P1)
	e.CollisionsWithPlayer(s, s.P2)

	// end?

	if s.Ball.P.X > 1 || s.Ball.P.X < 0 {
		s.Reset()
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

		speed := s.ballSpeed.L2() + 0.01*s.ballSpeed.L2()
		s.ballSpeed.Y = 0
		if s.ballSpeed.X < 0 {
			s.ballSpeed.X = speed
			s.ballSpeed = rotate(s.ballSpeed, offsetAngle)

		} else {
			s.ballSpeed.X = -speed
			s.ballSpeed = rotate(s.ballSpeed, -offsetAngle)
		}

		//inAngle := math.Atan(s.ballSpeed.y / s.ballSpeed.x)
		//fmt.Println("angle", inAngle*(180/math.Pi))
		//s.ballSpeed.y = p.offsetRatio(b.p.y) * 0.005
		if b.P.X+b.R > p.Bottom.X+p.Width/2 {
			b.P.X = p.Bottom.X - p.Width/2 - b.R
		}
		if b.P.X-b.R < p.Bottom.X-p.Width/2 {
			b.P.X = p.Bottom.X + p.Width/2 + b.R
		}
	}
}

func rotate(v position, angle float64) position {
	x2 := v.X*math.Cos(angle) - v.Y*math.Sin(angle)
	y2 := v.X*math.Sin(angle) + v.Y*math.Cos(angle)
	return position{
		X: x2,
		Y: y2,
	}
}
