package main

import (
	"./audio"
)

func main() {
	// Fill the sequence channel with a sequence C4->C5
	seq := make(chan int)
	go func() {
		for i := 30; i < 90; i++ {
			seq <- i
		}
		close(seq)
	}()

	// Initialize mixer data with wave and sequence
	output := make(chan int16)
	go audio.Init(seq, audio.Waves, output)

	/* Start the mixer running */
	for i := 0; i < 480000; i++ {
		<-output
	}
}
