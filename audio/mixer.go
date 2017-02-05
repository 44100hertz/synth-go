/* Does all of the work of the audio chip, which is really just a
 * mixer with added features.  Everything else is in this package is not
 * considered a part of the chip, but helps to use it.
 *
 * This works based on channel pairs.
 * Even number channels (including 0) are "left", and odds are "right".
 * For doing things to both channels, use OnPair.
 *
 * Note values are precise, as in 0xNNNNTTTT, N = midi note, T = fine tune.
 * Amplitude values are 100% at 0x10000 or 1<<16.
 * Many are signed, but below-zero values can have undefined behavior.
 * Using more int32's just makes things easier in the long run.
 */

package audio

import "math"

// The number of channel pairs, or mixer chans
const NumChans int = 2

type Mixer struct {
	srate uint32                  // Sample rate
	wave  func(int, uint32) int16 // Function used for sound waves
	seq   func(*Mixer)

	count    uint32 // Point counter
	nextTick uint32 // Location of next tick in points

	Ch        *[NumChans * 2]Channel  // Channels; pairs next to each other
	chans     *[NumChans](chan int32) // Data back from channel pairs
	BPM       uint32                  // Song speed in beats per minute
	TickRate  uint32                  // Ticks per update
	TickSpeed uint32                  // Callback after this many ticks
	tickCount uint32                  // Counts down ticks until callback
}

const ( // Channel pairing modes
	PairStereo = iota // Simple left and right channels
	PairPM            // Phase modulation
	PairAM            // Amplitude modulation
	PairSync          // Phase of left osc overflow = reset phase of right
)

// Internal channel data
type Channel struct {
	Wave     int // Index of wave to use for wave function
	PairMode int // Pair mode. See above.

	Note, Slide int32 // Midi note number

	Vol         int32 // Pre-Volume that affects effects
	MVol        int32 // Mixer volume; after effects
	Fade, MFade int32 // Per-tick volume adjustment

	Len, Phase uint32 // Length of wave and position in wave
	period     uint32 // How much to increment phase for each point

	delay      uint16         // Length of a delay effect in samples
	DelayTicks interface{}    // Length of delay in ticks
	DelayNote  interface{}    // Special delay timing used for guitar pluck
	DelayVol   int32          // Mixing level for delay effect
	DelayLoop  int32          // Feedback mixing level for delay
	FilterLen  uint16         // Rectangular filter added to delay
	hist       [1 << 16]int32 // 64kb of channel history
	histHead   uint16         // Current history location
	delayAvg   int32          // Rolling average tracker for delay
}

func NewMixer(wave func(int, uint32) int16, seq func(*Mixer)) Mixer {
	m := Mixer{
		wave:      wave,
		seq:       seq,
		Ch:        new([NumChans * 2]Channel),
		chans:     new([NumChans]chan int32),
		BPM:       120,
		TickRate:  24,
		TickSpeed: 6,
	}
	// Default params
	for i := range m.Ch {
		c := &m.Ch[i]
		c.MVol = 0x8000
		c.Note = 60 << 16
		c.Len = 0x10000
	}
	return m
}

func (m *Mixer) Start(output chan int16, srate uint32) {
	m.srate = srate

	for i := range m.chans {
		// Go is known to hang for up to 4ms at absolute most.
		// This would put my ideal GC amount at 48*4 = 192 And
		// because of stereo, that's actually 384. This was at
		// 128 before, and was still underrunning. It's
		// important to notice this in addition to the SDL
		// audio buffer.
		m.chans[i] = make(chan int32, 384)
		go m.startPair(i)
	}

	for {
		if m.count == m.nextTick {
			m.tick()
		}
		if m.tickCount >= m.TickSpeed {
			m.seq(m)
			m.tickCount = 0
		}
		var mixL int32 = 0
		var mixR int32 = 0
		for i := range m.chans {
			mixL += <-m.chans[i] * m.Ch[i*2].MVol >> 16
			mixR += <-m.chans[i] * m.Ch[i*2+1].MVol >> 16
		}
		output <- int16(clamp16(mixL))
		output <- int16(clamp16(mixR))
		m.count++
	}

}

// This is ran multiple times per beat in order to update various data.
// It coincides with sequence callbacks.
func (m *Mixer) tick() {
	for i := range m.Ch {
		c := &m.Ch[i]

		c.Note += c.Slide
		c.Vol = max(c.Vol+c.Fade, 0)
		c.MVol = max(c.MVol+c.MFade, 0)

		delayNote, ok := c.DelayNote.(int32)
		if ok {
			c.delay = uint16(float64(m.srate) / Note(delayNote))
			c.DelayNote = nil
		}

		delayTicks, ok := c.DelayTicks.(uint32)
		if ok {
			c.delay = uint16(delayTicks * m.srate * 60 /
				m.BPM / m.TickRate)
			c.DelayTicks = nil
		}

		// This avoids division by 0
		if c.FilterLen == 0 {
			c.FilterLen = 1
		}

		// Set pitch
		c.period = uint32(float64(c.Len/m.srate) * Note(c.Note))
	}
	m.nextTick = 60*m.srate/m.BPM/m.TickRate + m.count
	m.tickCount++
}

// Run a pair of Chs
func (m *Mixer) startPair(i int) {
	// Update the phase of a given channel in a canonical way
	phase := func(c *Channel) {
		c.Phase = (c.Phase + c.period) % c.Len
	}
	// Update the internal channel and return a mixer-usable wave
	wave := func(c *Channel, phase uint32) int32 {
		// Calculate delay
		var delayStart uint16 = c.histHead - c.delay
		var delayEnd uint16 = delayStart - c.FilterLen
		c.delayAvg += int32(c.hist[delayStart]) / int32(c.FilterLen)
		c.delayAvg -= int32(c.hist[delayEnd]) / int32(c.FilterLen)
		c.delayAvg = clamp16(c.delayAvg)

		// Get a wave output
		dry := int32(m.wave(c.Wave, phase))

		// Store history for delay effect
		c.hist[c.histHead] = dry + c.delayAvg*c.DelayLoop>>16
		c.histHead++
		return dry*c.Vol>>16 + c.delayAvg*c.DelayVol>>16
	}

	l := &m.Ch[i*2]
	r := &m.Ch[i*2+1]
	for {
		switch l.PairMode {
		case PairSync:
			// On new left osc cycle, new right osc cycle
			phase(l)
			phase(r)
			if l.Phase+l.period >= l.Len {
				r.Phase = 0
			}
			rwave := wave(r, r.Phase)
			m.chans[i] <- rwave
			m.chans[i] <- rwave
		case PairPM:
			// Use the wave of the left osc as the
			// phase of the right one.
			phase(l)
			lwave := uint32(wave(l, l.Phase)) + 0x8000
			rwave := wave(r, lwave)
			m.chans[i] <- rwave
			m.chans[i] <- rwave
		case PairAM:
			// Modulate amplitude of both waves
			phase(l)
			phase(r)
			lwave := wave(l, l.Phase)
			rwave := wave(r, r.Phase)
			total := lwave * rwave >> 16
			m.chans[i] <- total
			m.chans[i] <- total
		default:
			// Straight stereo left/right
			phase(l)
			phase(r)
			m.chans[i] <- wave(l, l.Phase)
			m.chans[i] <- wave(r, r.Phase)
		}
	}
}

func Note(note int32) float64 {
	fnote := float64(note) / (1 << 16)
	return math.Pow(2, (fnote-69)/12.0) * 440
}

func (m *Mixer) OnPair(i int, op func(*Channel)) {
	op(&m.Ch[i*2])
	op(&m.Ch[i*2+1])
}

func clamp16(a int32) int32 {
	if a < -0x8000 {
		return -0x8000
	} else if a > 0x7fff {
		return 0x7fff
	}
	return a
}

func max(a int32, max int32) int32 {
	if a < max {
		return max
	}
	return a
}
