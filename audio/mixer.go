/* Does all of the work of the audio chip, which is really just a
 * mixer with added features.  Everything else is in this package is not
 * considered a part of the chip, but helps to use it.
 */

package audio

import "math"

// The number of channel pairs, or mixer chans
const NumChans int = 4

type Mixer struct {
	// Passed-in, user-controlled params
	wave func(int, uint32) int16

	// Calculated values
	srate                               uint64
	channel                             *[NumChans * 2]channel
	chans                               *[NumChans](chan int16)
	tickCount, bpm, tickRate, tickSpeed uint64
	count, nextTick                     uint64
}

// Internal channel data
type channel struct {
	wave, note         int
	len, phase, period uint64 // Top 32 bits are used, bottom 32 are like a counter
}

// Create and run a mixer
func NewMixer(wave func(int, uint32) int16) Mixer {
	return Mixer{
		wave:      wave,
		channel:   new([NumChans * 2]channel),
		chans:     new([NumChans]chan int16),
		bpm:       120,
		tickRate:  1,
		tickSpeed: 1,
	}
}

func (m *Mixer) Start(output chan int16, srate uint64) {
	m.srate = srate
	for i := range m.chans {
		m.chans[i] = make(chan int16)
		go m.startPair(i)
	}

	for {
		if m.count == m.nextTick {
			m.tick()
		}
		if m.tickCount == m.tickSpeed {
			// Load sequence data here
		}
		var mix int32 = 0
		for i := range m.chans {
			mix += int32(<-m.chans[i])
		}
		output <- int16(mix >> 2)
		m.count++
	}

}

func (m *Mixer) tick() {
	for i := range m.channel {
		c := &m.channel[i]
		c.period = m.getPointPeriod(c.len, c.note)
	}
	m.nextTick = 60*m.srate/m.bpm/m.tickRate + m.count
	m.tickCount++
}

// Update a pair of channels
func (m *Mixer) startPair(i int) {
	l := &m.channel[i*2]
	// basic test code
	l.len = 0xffff << 32
	l.note = 60
	for {
		l.phase = (l.phase + l.period) % l.len
		m.chans[i] <- m.wave(l.wave, uint32(l.phase>>32))
	}
}

// Calculate amount to add to phase to produce a given pitch
func (m *Mixer) getPointPeriod(len uint64, note int) uint64 {
	rate := float64(len>>11) / float64(m.srate) // point period for 1hz wave
	pitch := math.Pow(2, float64(note-60)/12.0) * 440
	return uint64(rate*pitch) << 11
}
