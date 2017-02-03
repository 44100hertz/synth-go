package main

import "./audio"

func main() {
	counter := int32(0)
	seq := func(m *audio.Mixer) {
		m.OnPair(0, func(c *audio.Channel) {
			if counter%12 < 11 {
				c.Vol = 0
			} else {
				c.DelayNote += 1 << 16
				c.Vol = 0x10000
			}
		})
		counter++
	}
	m := audio.NewMixer(audio.Waves, seq)
	m.OnPair(0, func(c *audio.Channel) {
		c.Vol = 0x8000 // Full volume (single channel)
		// c.Note = 1            // C#-0 as base note
		// c.DelayLevel = 0xFFF0 // A bit of delay attenuation
		// c.Filter = 0x3        // Use a delay averaged by 3 samples
		c.Wave = 2
		c.Note = 1
		c.DelayLevel = 0xFFF0
		c.DelayNote = 48 << 16
		c.Filter = 3
	})

	audio.Start(&m)
	//audio.WaveOut(&m, "out.raw", 48000)
}
