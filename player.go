package main

type player struct {
	bottom position
	top    position
	width  float64
}

func (p player) center() float64 {
	return p.bottom.y + p.height()/2
}

func (p player) offsetRatio(y float64) float64 {
	return (y - p.center()) / (p.height() / 2)

}

func (p player) height() float64 {
	return p.top.y - p.bottom.y
}
