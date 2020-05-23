package main

import (
	"math"
)

type Engine struct {
}

func (e *Engine) Process(s GameState, p1in inputs, p2in inputs) GameState {
	s.p1 = e.processPlayerInput(s.p1, p1in)
	s.p2 = e.processPlayerInput(s.p2, p2in)
	e.processBall(&s)
	return s
}

func (e *Engine) processPlayerInput(p player, in inputs) player {
	dy := 0.0
	if in.up {
		dy += 0.01
		if p.top.y+dy > 1 {
			return p
		}
	}
	if in.down {
		dy -= 0.01
		if p.bottom.y+dy < 0 {
			return p
		}
	}

	p.top.y += dy
	p.bottom.y += dy

	// Now verify
	//fmt.Println(p.top.y)
	return p
}

func (e *Engine) processBall(s *GameState) {
	b := s.ball

	b.p.y += s.ballSpeed.y
	b.p.x += s.ballSpeed.x

	// Top and top wall collisions
	if b.p.y+b.r > 1 {
		s.ballSpeed.y = -s.ballSpeed.y
		b.p.y = 1 - b.r
		return
	}
	if b.p.y-b.r < 0 {
		s.ballSpeed.y = -s.ballSpeed.y
		b.p.y = b.r
		return
	}

	// Player collisions?

	s.ball = b
	e.CollisionsWithPlayer(s, s.p1)
	e.CollisionsWithPlayer(s, s.p2)

	// end?

	if s.ball.p.x > 1 || s.ball.p.x < 0 {
		s.Reset()
	}
}

func (e *Engine) CollisionsWithPlayer(s *GameState, p player) {
	b := s.ball
	tx := b.p.x
	ty := b.p.y
	if b.p.x > p.top.x {
		tx = p.top.x + p.width/2
	} else if b.p.x < p.top.x {
		tx = p.top.x - p.width/2
	}

	if b.p.y > p.top.y {
		ty = p.top.y
	} else if b.p.y < p.bottom.y {
		ty = p.bottom.y
	}
	dx := tx - b.p.x
	dy := ty - b.p.y

	d := math.Sqrt(dx*dx + dy*dy)
	//fmt.Println("dx dy", dx, dy)
	//fmt.Println("DISTANCE", d)
	//fmt.Println("BALL", b)
	//fmt.Println("BOTOOM", p.top)
	//fmt.Println("TOP", p.bottom)
	if d < b.r {
		//fmt.Println("COLLISION", b.p.x, p.bottom.x)
		s.ballSpeed.x = -s.ballSpeed.x
		if b.p.x+b.r > p.bottom.x+p.width/2 {
			b.p.x = p.bottom.x - p.width/2 - b.r
		}
		if b.p.x-b.r < p.bottom.x-p.width/2 {
			b.p.x = p.bottom.x + p.width/2 + b.r
		}
	}
}
