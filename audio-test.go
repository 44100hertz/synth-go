package main

import (
	"./audio"
)

func main() {
	// Fill the sequence channel with a sequence C4->C5
	getSeq := func(seq chan int) {
		for i := 60; i < 72; i++ {
			seq <- i
		}
		close(seq)
	}

	// Initialize mixer data with wave and sequence
	output := make(chan int16)
	go audio.Init(getSeq, audio.Sine, output)

	/* Start the mixer running */
	for i := 0; i < 480000; i++ {
		<-output
	}
}
