package main

import (
	"fmt"
	"time"
)

type position struct {
	x float64
	y float64
}

type ball struct {
	p position
	r float64
}

type player struct {
	p1 position
	p2 position
}

const ticks = 60

var done = make(chan struct{})

func main() {
	fmt.Println("Start")
	c := NewCanvas("canvas", 1000, 1000)
	players := []player{{p1: position{
		x: 0,
		y: 0.4,
	}, p2: position{
		x: 0,
		y: 0.6,
	}},
		{p1: position{
			x: 1,
			y: 0.4,
		}, p2: position{
			x: 1,
			y: 0.6,
		}},
	}

	ball := ball{
		p: position{
			x: 0.5,
			y: 0.5,
		},
		r: 0.01,
	}

	fmt.Println("go the canvas")
	go func() {
		fmt.Println("starting loop")
		t := time.Tick(time.Second / ticks)
		for {
			select {
			case <-t:
				fmt.Println("rendering")
				c.render(players, ball)
			}
		}
	}()

	<-done
}
