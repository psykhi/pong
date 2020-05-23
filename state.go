package main

type GameState struct {
	ball      ball
	p1        player
	p2        player
	ballSpeed position
}

func (s *GameState) Reset() {
	p1 := player{
		bottom: position{
			x: 0,
			y: 0.4,
		}, top: position{
			x: 0,
			y: 0.6,
		},
		width: 0.01,
	}
	p2 := player{
		bottom: position{
			x: 1,
			y: 0.4,
		}, top: position{
			x: 1,
			y: 0.6,
		},
		width: 0.01,
	}
	ball := ball{
		p: position{
			x: 0.5,
			y: 0.5,
		},
		r: 0.01,
	}
	speed := position{
		x: 0.01,
		y: 0.001,
	}

	s.ball = ball
	s.p1 = p1
	s.p2 = p2
	s.ballSpeed = speed
}

func NewState() GameState {
	s := GameState{}
	s.Reset()
	return s
}
