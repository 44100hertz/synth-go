package audio

import "math"

// Generate any sort of 16-bit lookup table.
// fn expected to take in range 0.0-1.0, output full-range int16
func makeLUT(fn func(float64) int16, size int) []int16 {
	lut := make([]int16, size)
	for i := 0; i < size; i++ {
		lut[i] = fn(float64(i) / float64(size))
	}
	return lut
}

func sineLUT_maker(off float64) int16 {
	// Set range to 1/4 sine wave
	off = off * math.Pi / 2.0
	// Convert to 16-bit range
	return int16(math.Sin(off) * float64(0x4000*2))
}

// Size of the lookup table to generate
// All waveforms will use by this size * 4
const lutSize int = 0x4000

var lut []int16 = makeLUT(sineLUT_maker, lutSize)

// ∿∿∿∿
func Sine(off int) int16 {
	o := off % lutSize
	switch off / lutSize {
	case 0:
		return lut[o]
	case 1:
		return lut[lutSize-o]
	case 2:
		return -lut[o]
	case 3:
		return -lut[lutSize-o]
	}
	return 0
}

// n_n_
func HalfSine(off int) int16 {
	o := off % lutSize
	switch off / lutSize {
	case 0:
		return lut[o]
	case 1:
		return lut[lutSize-o]
	}
	return 0
}

// nnnn
func CamelSine(off int) int16 {
	o := off % lutSize
	switch off / lutSize {
	case 0, 2:
		return lut[o]
	case 1, 3:
		return lut[lutSize-o]
	}
	return 0
}

// r_r_
func QuarterSine(off int) int16 {
	o := off % lutSize
	switch off / lutSize {
	case 0, 2:
		return lut[o]
	}
	return 0
}

// ΓLΓL
func Pulse(off int) int16 {
	if off > lutSize*2 {
		return 0x7fff
	}
	return -0x8000
}

// /|/|
func Ramp(off int) int16 {
	if off < lutSize*4 {
		return int16(off / 2)
	}
	return 0
}
