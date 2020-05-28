package main

import (
	"errors"
	"fmt"
	"github.com/psykhi/pong/game"
)

type InputManager struct {
	inputs  []game.Inputs
	cursor  int
	maxSize int
	last    int
}

func NewInputManager(size int) *InputManager {
	return &InputManager{
		inputs:  make([]game.Inputs, 0),
		cursor:  0,
		maxSize: size,
		last:    0,
	}
}

func (im *InputManager) Set(in game.Inputs) {
	im.last = in.SequenceID
	if len(im.inputs) < im.maxSize {
		im.inputs = append(im.inputs, in)
		return
	}
	im.inputs[im.cursor] = in
	im.cursor++
	if im.cursor == im.maxSize {
		im.cursor = 0
	}
}

func (im *InputManager) Get(sequenceID int) (game.Inputs, error) {
	if sequenceID > im.last || sequenceID < im.last-len(im.inputs)+1 {
		fmt.Println("not found", sequenceID, im.last)
		return game.Inputs{}, errors.New("sequenceID not found")
	}

	for _, in := range im.inputs {
		if in.SequenceID == sequenceID {
			return in, nil
		}
	}

	// Should not happen
	fmt.Println("no sequence ID", sequenceID, im.inputs)
	panic("Could not find sequence ID")
}
