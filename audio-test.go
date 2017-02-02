package main

import "./audio"

func main() {
	seq := func(m *audio.Mixer) {
		m.Ch[1].DelayNote++
		// m.OnPair(0, func(c *audio.Channel) {
		// 	c.Vol = 0
		// })
	}
	m := audio.NewMixer(audio.Waves, seq)
	m.OnPair(0, func(c *audio.Channel) {
		c.Vol = 0x1000 // Full volume (single channel)
		// c.Note = 1            // C#-0 as base note
		// c.DelayLevel = 0xFFF0 // A bit of delay attenuation
		// c.Filter = 0x3        // Use a delay averaged by 3 samples
	})
	m.Ch[0].Wave = 0
	m.Ch[0].PairMode = audio.PAIR_PM
	m.Ch[0].VolSlide = 0x200

	m.Ch[1].Wave = 1
	m.Ch[1].DelayLevel = 0x10000
	m.Ch[1].DelayNote = 30
	m.Ch[1].Filter = 3
	m.Ch[1].VolSlide = -0x1

	//	audio.Start(&m)
	audio.WaveOut(&m, "out.raw", 48000)
}
