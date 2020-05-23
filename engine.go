package main

type Engine struct {
}

func (e *Engine) Process(s GameState, p1in inputs, p2in inputs) GameState {
	s.p1 = e.process(s.p1, p1in)
	s.p2 = e.process(s.p2, p2in)
	return s
}

func (e *Engine) process(p player, in inputs) player {
	dy := 0.0
	if in.up {
		dy += 0.01
		if p.bottom.y+dy > 1 {
			return p
		}
	}
	if in.down {
		dy -= 0.01
		if p.top.y+dy < 0 {
			return p
		}
	}

	p.bottom.y += dy
	p.top.y += dy

	// Now verify
	//fmt.Println(p.bottom.y)
	return p
}
