/* This is a testbed main for getting sound out of the mixer.
 * It is designed for ear tests.
 */
package main

import (
	"./audio/"
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
	context := mixer.Context{Seq: getSeq, Wave: waves.Sine}
	context.Init()
	/* Start the mixer running */
	/* Print the mixer's output wave */
}
