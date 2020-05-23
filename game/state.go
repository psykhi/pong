package game

import (
	"math"
)

type position struct {
	X float64
	Y float64
}

func (p position) L2() float64 {
	return math.Sqrt(p.X*p.X + p.Y*p.Y)
}

type Ball struct {
	P position
	R float64
}

type State struct {
	Ball      Ball
	P1        Player
	P2        Player
	ballSpeed position
}

func (s *State) Reset() {
	p1 := Player{
		Bottom: position{
			X: 0.1,
			Y: 0.4,
		}, Top: position{
			X: 0.1,
			Y: 0.6,
		},
		Width: 0.01,
	}
	p2 := Player{
		Bottom: position{
			X: 0.9,
			Y: 0.4,
		}, Top: position{
			X: 0.9,
			Y: 0.6,
		},
		Width: 0.01,
	}
	ball := Ball{
		P: position{
			X: 0.5,
			Y: 0.5,
		},
		R: 0.01,
	}
	speed := position{
		X: 0.005,
		Y: 0.001,
	}

	s.Ball = ball
	s.P1 = p1
	s.P2 = p2
	s.ballSpeed = speed
}

func NewState() *State {
	s := State{}
	s.Reset()
	return &s
}
