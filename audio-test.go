package main

import "./audio"

func main() {
	count := 0
	seq := func(m *audio.Mixer) {
		m.OnPair(0, func(c *audio.Channel) {
			count++
		})
	}
	m := audio.NewMixer(audio.Waves, seq)
	m.OnPair(0, func(c *audio.Channel) {
		//c.Note = 1            // C#-0 as base note
		c.Peak = 0x10000
		c.Wave = audio.WaveQSine
		c.NoteOn = true
		c.Attack = 0x100
		c.Decay = 0x100
		c.Sustain = 0x0
		c.Release = 0x80
		c.Vibrato = (1 << 16) / 4
		c.VibratoRate = (1 << 32) / 8
		c.DryLevel = 0x1000
		c.WetLevel = 0x8000
		c.FilterLen = 0x8
	})

	audio.Start(&m)
	//audio.WaveOut(&m, "out.raw", 48000)
}
