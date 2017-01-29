package audio

import "math"

// The number of channel pairs, or mixer chans
const NumChans int = 4

type Mixer struct {
	srate    int
	seq      func(chan int)
	wave     func(int) int16
	channels *[NumChans * 2]Channel
	chans    *[NumChans](chan int16)
	count    uint64
	nextTick uint64
	bpm      uint16
	tickrate uint16
}

// A single playback channel. Every even channel is "L", odd "R"
type Channel struct {
	phase, period, off, len int
}

// Start up the parts of a context that a user needn't touch.
func Init(seq func(chan int), wave func(int) int16, output chan int16) {
	m := Mixer{
		srate:    48000,
		seq:      seq,
		wave:     wave,
		channels: new([NumChans * 2]Channel),
		chans:    new([NumChans]chan int16),
		bpm:      120,
		tickrate: 4,
	}

	for i := range m.chans {
		m.chans[i] = make(chan int16)
		go m.startPair(i)
	}

	for {
		if m.count == m.nextTick {
			m.nextTick = uint64(60.0*float64(m.srate)/
				float64(m.bpm)/float64(m.tickrate) + float64(m.count))
			for _, c := range m.channels {
				c.period = m.getPointPeriod(c.len, 60)
			}
		}
		var mix int32 = 0
		for i := range m.chans {
			mix += int32(<-m.chans[i])
		}
		output <- int16(mix >> 2)
		m.count++
	}
}

// Update a pair of channels
func (m Mixer) startPair(i int) {
	l := m.channels[i*2]
	l.len = 0x8000
	for {
		l.phase = (l.phase + l.period) % (l.len)
		m.chans[i] <- m.wave(l.phase)
	}
}

// Calculate amount to add to phase to produce a given pitch
func (m Mixer) getPointPeriod(len int, note int) int {
	rate := float64(len) / float64(m.srate)
	pitch := math.Pow(2, float64(note-60)) * 440
	return int(rate * pitch)
}
