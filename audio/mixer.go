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
const NumChans int = 2

type Mixer struct {
	srate uint64                  // Sample rate
	wave  func(int, uint32) int16 // Function used for sound waves
	seq   func(*Mixer)

	count    uint64 // Point counter
	nextTick uint64 // Location of next tick in points

	Ch        *[NumChans * 2]Channel  // Channels; pairs next to each other
	chans     *[NumChans](chan int32) // Data back from channel pairs
	Bpm       uint64                  // Song speed in beats per minute
	TickRate  uint64                  // Ticks per update
	TickSpeed uint64                  // Callback after this many ticks
	tickCount uint64                  // Counts down ticks until callback
}

// Internal channel data
type Channel struct {
	Wave            int    // Index of wave to use for wave function
	Note            int32  // Midi note number
	Tune, TuneSlide int32  // Fine tuning, one note = 0x8000
	Vol, VolSlide   int32  // Pre-Volume that affects effects
	MVol, MVolSlide int32  // Mixer volume; after effects
	Len, Phase      uint64 // Length of wave and position in wave
	Period          uint64 // How much to increment phase for each point
	DelayTicks      uint64 // Length of delay in ticks
	Delay           uint16 // Length of a delay effect in samples
	DelayLevel      int32  // Level at which to mix in delay effect
	DelayFilter     uint16 // Rectangular filter size 2^n added to delay

	hist     [1 << 16]int16 // 64kb of channel history
	histHead uint16         // Current history location
	delayAvg int32          // Rolling average tracker for delay
}

func NewMixer(wave func(int, uint32) int16, seq func(*Mixer)) Mixer {
	m := Mixer{
		wave:      wave,
		seq:       seq,
		Ch:        new([NumChans * 2]Channel),
		chans:     new([NumChans]chan int32),
		Bpm:       120,
		TickRate:  24,
		TickSpeed: 6,
	}
	m.ForAllCh(func(c *Channel) {
		c.MVol = 0x8000
		c.Note = 60
	})
	return m
}

func (m *Mixer) Start(output chan int16, srate uint64) {
	m.srate = srate

	// Start each channel pair wave output
	for i := range m.chans {
		m.chans[i] = make(chan int32)
		go m.startPair(i)
	}

	// Run the mixer and ticking
	for {
		if m.count == m.nextTick {
			m.tick()
		}
		if m.tickCount >= m.TickSpeed {
			m.seq(m)
			m.tickCount = 0
		}
		var mix int32 = 0
		for i := range m.chans {
			mix += <-m.chans[i]
		}
		switch {
		case mix > 0x7fff:
			mix = 0x7fff
		case mix < -0x8000:
			mix = -0x8000
		}
		output <- int16(mix)
		m.count++
	}

}

// This is ran multiple times per beat in order to update various data.
// It coincides with sequence callbacks.
func (m *Mixer) tick() {
	for i := range m.Ch {
		c := &m.Ch[i]

		// Sliding values
		c.Tune += c.TuneSlide
		c.Vol += c.VolSlide
		c.MVol += c.MVolSlide

		// Limit fine tune range indirectly so that note stays sane
		c.Note = c.Note + (c.Tune / 0x8000)
		c.Tune = c.Tune % 0x8000

		// Cannot delay by amount 0
		if c.DelayFilter == 0 {
			c.DelayFilter = 1
		}

		if c.DelayTicks > 0 {
			c.Delay = uint16(c.DelayTicks * m.srate * 60 /
				m.Bpm / m.TickRate)
		}

		// Set pitch
		c.Period = m.getPointPeriod(c.Len, c.Note, c.Tune)
	}
	m.nextTick = 60*m.srate/m.Bpm/m.TickRate + m.count
	m.tickCount++
}

// Run a pair of Chs
func (m *Mixer) startPair(i int) {
	l := &m.Ch[i*2]
	// basic test code
	l.Len = 0x10000 << 32
	for {
		// Set phase and grab wave
		l.Phase = (l.Phase + l.Period) % l.Len
		wave := int32(m.wave(l.Wave, uint32(l.Phase>>32)))
		wave = int32((wave * l.Vol) >> 16)

		// Apply delay effect
		// Important that this is done first
		var delayStart uint16 = l.histHead - l.Delay
		var delayEnd uint16 = delayStart - l.DelayFilter
		l.delayAvg += int32(l.hist[delayStart]) / int32(l.DelayFilter)
		l.delayAvg -= int32(l.hist[delayEnd]) / int32(l.DelayFilter)
		wave = l.delayAvg*l.DelayLevel>>16 + wave

		// Store history for delay effect
		l.hist[l.histHead] = int16(wave)
		l.histHead++

		m.chans[i] <- int32((wave * l.MVol) >> 16)
	}
}

// Calculate amount to add to phase to produce a given pitch
func (m *Mixer) getPointPeriod(len uint64, note int32, tune int32) uint64 {
	// Find point period for 1hz wave at given length
	rate := float64(len / m.srate)
	// Find desired pitch in hertz
	totalNote := float64(note) + float64(tune)/0x8000
	pitch := math.Pow(2, (totalNote-60)/12.0) * 440
	return uint64(rate * pitch)
}

func (m *Mixer) ForAllCh(op func(*Channel)) {
	for i := range m.Ch {
		op(&m.Ch[i])
	}
}
