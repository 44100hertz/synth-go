package main

import (
	"./audio"
	"github.com/veandco/go-sdl2/sdl"
)

func main() {
	sdl.Init(sdl.INIT_AUDIO)

	// Fill the sequence channel with a sequence C4->C5
	instr := []audio.Inst {
		{ Index: 0, Len: 0xffff, },
		{ Index: nil, Len: nil, },
		{ Index: nil, Len: nil, },
	}

	// Initialize mixer data with wave and sequence
	output := make(chan int16)
	go audio.Init(audio.Waves, instr, output)

	/* Start the mixer running */
	for i := 0; i < 480000; i++ {
		<-output
	}
}
