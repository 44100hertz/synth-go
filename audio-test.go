package main

import "./audio"

func main() {
	seq := func(m *audio.Mixer) {
		// m.OnPair(0, func(c *audio.Channel) {
		// 	c.Vol = 0
		// })
	}
	m := audio.NewMixer(audio.Waves, seq)
	m.OnPair(0, func(c *audio.Channel) {
		c.Wave = 0      // Quarter sine
		c.Vol = 0x10000 // Full volume (single channel)
		// c.Note = 1            // C#-0 as base note
		// c.DelayLevel = 0xFFF0 // A bit of delay attenuation
		// c.Filter = 0x3        // Use a delay averaged by 3 samples
	})
	m.Ch[0].Wave = 1
	m.Ch[0].PairMode = audio.PAIR_PM
	m.Ch[0].VolSlide = 0x400
	audio.Start(&m)
}
