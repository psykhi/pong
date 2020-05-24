package game

type Player struct {
	Bottom    Position
	Top       Position
	Width     float64
	Score     int
	Connected bool
}

func (p Player) Center() float64 {
	return p.Bottom.Y + p.Height()/2
}

func (p Player) OffsetRatio(y float64) float64 {
	return (y - p.Center()) / (p.Height() / 2)

}

func (p Player) Height() float64 {
	return p.Top.Y - p.Bottom.Y
}
