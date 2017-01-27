package audio

import "math"

type Context struct {
	srate int
	seq   func(chan int)
	wave  func(int) int16
	chans Channel
}

// A single playback channel. Every even channel is "L", odd "R"
type Channel struct {
	phase, period,
	off, len int
}

// Start up the parts of a context that a user needn't touch.
func Init(seq func(chan int), wave func(int) int16, output chan int16) {
	channel := Channel{0, 1000, 0, 80000}
	c := Context{48000, seq, wave, channel}

	go c.pairUpdate(c.chans, output)
}

// Return increment amount to produce a specific pitch
func (c Context) getPointPeriod(len uint32, note int) uint32 {
	return uint32(float64(len) * math.Pow(float64(note), -2))
}

// Update a pair of channels
func (c Context) pairUpdate(l Channel, output chan int16) {
	for {
		l.phase = (l.phase + l.period) % l.len
		output <- c.wave(l.phase)
	}
}
