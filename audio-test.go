package main

import "./audio"

func main() {
	seq := func(m *audio.Mixer) {
		m.OnPair(0, func(c *audio.Channel) {
			c.Vol = 0
		})
	}
	m := audio.NewMixer(audio.Waves, seq)
	m.OnPair(0, func(c *audio.Channel) {
		c.Wave = 3            // Quarter sine
		c.Vol = 0x10000       // Full volume (single channel)
		c.Note = 1            // C#-0 as base note
		c.DelayLevel = 0xFFF0 // A bit of delay attenuation
		c.Filter = 0x3        // Use a delay averaged by 3 samples
	})
	m.Ch[0].DelayNote = 48
	m.Ch[1].DelayNote = 49
	audio.Start(&m)
}
