package main

import "./audio"

func main() {
	counter := int32(0)
	seq := func(m *audio.Mixer) {
		m.OnPair(0, func(c *audio.Channel) {
			if counter%6 == 0 {
				note := (counter / 4) << 16
				//				c.Note = note / 2
				c.DelayNote = note
			}
		})
		counter++
	}
	m := audio.NewMixer(audio.Waves, seq)
	m.OnPair(0, func(c *audio.Channel) {
		//c.Vol = 0x8000 // Full volume (single channel)
		// c.Note = 1            // C#-0 as base note
		// c.DelayLevel = 0xFFF0 // A bit of delay attenuation
		// c.Filter = 0x3        // Use a delay averaged by 3 samples
		c.Wave = audio.WaveQSine
		c.Fade = -0x1000
		c.Note = 1
		c.DelayVol = 0x8000
		c.DelayLoop = 0x11000
		c.DelayNote = 48 << 16
		c.FilterLen = 0x4
	})

	m.BPM = 300
	audio.Start(&m)
	//audio.WaveOut(&m, "out.raw", 48000)
}
