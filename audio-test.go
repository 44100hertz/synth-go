package main

import "./audio"

func main() {
	seq := func(m *audio.Mixer) {
		m.OnPair(0, func(c *audio.Channel) {
		})
	}
	m := audio.NewMixer(audio.Waves, seq)
	m.OnPair(0, func(c *audio.Channel) {
		//c.Note = 1            // C#-0 as base note
		c.Peak = 0x8000
		c.Wave = audio.WaveQSine
		c.NoteOn = true
		c.Attack = 0x100
		c.Decay = 0x100
		c.Sustain = 0x4000
		c.Release = 0x80
		c.Tremolo = (1 << 16) / 4
		c.TremoloRate = (1 << 32) / 8
		// c.DelayVol = 0x8000
		// c.DelayLoop = 0x11000
		// c.DelayNote = 48 << 16
		// c.FilterLen = 0x4
	})

	audio.Start(&m)
	//audio.WaveOut(&m, "out.raw", 48000)
}
