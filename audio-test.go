/* This is a testbed main for getting sound out of the mixer.
 * It is designed for ear tests.
 */
package main

func main() {
	// Return a pulse for producing square waves.
	getWave := func(offset int32) int8 {
		if offset < 127 {
			return -128
		}
		return 127
	}

	// Fill the sequence channel with a sequence C4->C5
	getSeq := func(seq chan int) {
		for i := 60; i < 72; i++ {
			seq <- i
		}
		close(seq)
	}

	/* Initialize mixer data with wave and sequence */
	/* Start the mixer running */
	/* Print the mixer's output wave */
}
