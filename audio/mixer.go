package audio

import "math"

// The number of channel pairs, or mixer chans
const NUM_CHANS int = 4

type Mixer struct {
	srate    int
	seq      func(chan int)
	wave     func(int) int16
	channels *[NUM_CHANS*2]Channel
	chans    *[NUM_CHANS](chan int16)
}

// A single playback channel. Every even channel is "L", odd "R"
type Channel struct {
	phase, period,
	off, len int
}

// Start up the parts of a context that a user needn't touch.
func Init(seq func(chan int), wave func(int) int16, output chan int16) {
	m := Mixer{48000, seq, wave,
		new([NUM_CHANS*2]Channel),
		new([NUM_CHANS]chan int16),
	}

	for i := range m.chans {
		m.chans[i] = make(chan int16)
		go m.startPair(i)
	}

	var mix int16
	for {
		mix = 0
		for i := range m.chans {
			mix += <-m.chans[i]
		}
		output <- mix
	}
}

// Return increment amount to produce a specific pitch
func (m Mixer) getPointPeriod(len uint32, note int) uint32 {
	return uint32(float64(len) * math.Pow(float64(note), -2))
}

// Update a pair of channels
func (m Mixer) startPair(i int) {
	l := m.channels[i*2]
	for {
		l.phase = (l.phase + l.period) % (l.len + 1)
		m.chans[i] <- m.wave(l.phase)
	}
}
