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
	wave func(int, int) int16
	inst []Inst

	// Calculated values
	srate   int
	channel *[NumChans * 2]Channel
	chans   *[NumChans](chan int16)
	count, tickCount,
	bpm, nextTick,
	tickRate, tickSpeed int
}

// Internal channel data
type Channel struct {
	index, len,
	phase, period,
	note int
	inst *[3]int
}

// A public way to modify instrument data
type Inst struct {
	Index, Len interface{}
}

// Create and run a mixer
func Init(wave func(int, int) int16, inst []Inst, output chan int16) {
	m := Mixer{
		wave:      wave,
		channel:   new([NumChans * 2]Channel),
		chans:     new([NumChans]chan int16),
		inst:      inst,
		srate:     48000,
		bpm:       120,
		tickRate:  1,
		tickSpeed: 1,
	}

	for i := range m.chans {
		m.chans[i] = make(chan int16)
		go m.startPair(i)
	}

	for {
		if m.count == m.nextTick {
			for i := range m.chans {
				m.loadInst(i, 0)
			}
			m.nextTick = 60*m.srate/m.bpm/m.tickRate + m.count
			m.tickCount++
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

// Update a pair of channels
func (m *Mixer) startPair(i int) {
	l := m.channel[i*2]
	l.len = 0x8000
	l.note = 60
	l.inst = new([3]int)
	for {
		l.phase = (l.phase + l.period) % (l.len)
		m.chans[i] <- m.wave(0, l.phase) // Sine wave
	}
}

// Calculate amount to add to phase to produce a given pitch
func (m *Mixer) getPointPeriod(len int, note int) int {
	rate := float64(len) / float64(m.srate)
	pitch := math.Pow(2, float64(note-60)) * 440
	return int(rate * pitch)
}

// Load instrument data once
func (m *Mixer) loadInst(index int, instpart int) {
	c := &m.channel[index]
	i := m.inst[0]
//	i := &m.inst[c.inst[instpart]]
	c.phase = 0
	c.period = m.getPointPeriod(c.len, c.note)
	c.len = i.Len.(int) // Todo: make optional
}
