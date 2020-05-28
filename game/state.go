package game

import (
	"math"
	"math/rand"
	"time"
)

type Position struct {
	X float64
	Y float64
}

func (p Position) L2() float64 {
	return math.Sqrt(p.X*p.X + p.Y*p.Y)
}

type Ball struct {
	P Position
	R float64
}

type State struct {
	Ball    Ball
	Players []Player
	//P1               Player
	//P2               Player
	BallSpeed        Position
	Finished         bool
	WaitingForPlayer bool
	Paused           bool
	Restart          time.Time
}

func (s *State) Inputs() []Inputs {
	ret := make([]Inputs, 0)
	for _, p := range s.Players {
		ret = append(ret, p.Inputs)
	}
	return ret
}

func (s *State) Countdown() {
	s.Reset()
	s.Paused = true
	s.Restart = time.Now().Add(2 * time.Second)
	s.BallSpeed = Position{}
	s.BallSpeed = Position{
		X: StartBallSpeed,
	}
	if rand.Float64() > 0.5 {
		s.BallSpeed.X = -s.BallSpeed.X
	}
	s.BallSpeed = rotate(s.BallSpeed, rand.Float64()*math.Pi/3)

}

func (s *State) Reset() {
	s.Players[0].Bottom.X = 0.1
	s.Players[0].Bottom.Y = 0.4
	s.Players[0].Top.X = 0.1
	s.Players[0].Top.Y = 0.6
	s.Players[0].Width = 0.01

	s.Players[1].Bottom.X = 0.9
	s.Players[1].Bottom.Y = 0.4
	s.Players[1].Top.X = 0.9
	s.Players[1].Top.Y = 0.6
	s.Players[1].Width = 0.01

	ball := Ball{
		P: Position{
			X: 0.5,
			Y: 0.5,
		},
		R: 0.01,
	}

	s.Ball = ball
	s.Finished = false
}

func NewState() *State {
	s := State{
		Players: make([]Player, 2),
	}
	s.Reset()
	return &s
}
