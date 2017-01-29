/* Does all of the work of the audio chip, which is really just a
 * mixer with added features.  Everything else is in this package is not
 * considered a part of the chip, but helps to use it.
 */

package audio

import "math"

// The number of channel pairs, or mixer chans
const NumChans int = 4

type Mixer struct {
	srate    int
	seq      chan int
	wave     func(int, int) int16
	channel  *[NumChans * 2]Channel
	chans    *[NumChans](chan int16)
	count    int
	nextTick int
	bpm      int
	tickrate int
}

// A single playback channel. Every even channel is "L", odd "R"
type Channel struct {
	instIndex,
	phase, period, len,
	note int
	inst [3]Inst
}

// A public way to modify instrument data
type Inst struct {
	WaveIndex,
	WaveLength,
	Note,
	LoopMode interface{}
}

// Start up the parts of a context that a user needn't touch.
func Init(seq chan int, wave func(int, int) int16, output chan int16) {
	m := Mixer{
		srate:    48000,
		seq:      seq,
		wave:     wave,
		channel:  new([NumChans * 2]Channel),
		chans:    new([NumChans]chan int16),
		bpm:      120,
		tickrate: 1,
	}

	for i := range m.chans {
		m.chans[i] = make(chan int16)
		go m.startPair(i)
	}

	for {
		if m.count == m.nextTick {
			for i := range m.chans {
				m.loadInst(i,1)
			}
			m.nextTick = 60*m.srate/m.bpm/m.tickrate + m.count
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

func instGetInt(i interface{}, modval *int) {
	integer, ok := i.(int)
	if ok {
		*modval = integer
	}
}

func (m *Mixer) loadInst(index int, instpart int) {
	c := &m.channel[index]
	i := &c.inst[instpart]
	c.phase = 0
	c.period = m.getPointPeriod(c.len, c.note)
	instGetInt(i.WaveLength, &c.len)
	c.note = <- m.seq
}
