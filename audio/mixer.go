/* Does all of the work of the audio chip, which is really just a
 * mixer with added features.  Everything else is in this package is not
 * considered a part of the chip, but helps to use it.
 *
 * For most uint64's here, top 32 bits are used, bottom 32 are like a counter.
 * This should probably be altered to use all bits of uint32 and just cope
 * with the imprecision and limits.
 */

package audio

import "math"

// The number of channel pairs, or mixer chans
const NumChans int = 4

type Mixer struct {
	// Passed-in values
	srate uint64
	wave  func(int, uint32) int16

	// External values
	Channel   *[NumChans * 2]Channel
	Bpm       uint64
	TickRate  uint64
	TickSpeed uint64

	// Internal values
	chans     *[NumChans](chan int16)
	count     uint64
	tickCount uint64
	nextTick  uint64
}

// Internal channel data
type Channel struct {
	Wave       int
	Note       int
	Period     uint64
	Len, Phase uint64
}

// Create and run a mixer
func NewMixer(wave func(int, uint32) int16) Mixer {
	return Mixer{
		wave:      wave,
		Channel:   new([NumChans * 2]Channel),
		chans:     new([NumChans]chan int16),
		Bpm:       120,
		TickRate:  1,
		TickSpeed: 1,
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
		if m.tickCount == m.TickSpeed {
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
	for i := range m.Channel {
		c := &m.Channel[i]
		c.Period = m.getPointPeriod(c.Len, c.Note)
	}
	m.nextTick = 60*m.srate/m.Bpm/m.TickRate + m.count
	m.tickCount++
}

// Update a pair of Channels
func (m *Mixer) startPair(i int) {
	l := &m.Channel[i*2]
	// basic test code
	l.Len = 0xffff << 32
	l.Note = 60
	for {
		l.Phase = (l.Phase + l.Period) % l.Len
		m.chans[i] <- m.wave(l.Wave, uint32(l.Phase>>32))
	}
}

// Calculate amount to add to phase to produce a given pitch
func (m *Mixer) getPointPeriod(len uint64, note int) uint64 {
	// Find point period for 1hz wave at given length
	rate := len / m.srate
	// Find desired pitch in hertz
	pitch := math.Pow(2, float64(note-60)/12.0) * 440
	return uint64(float64(rate) * pitch)
}
