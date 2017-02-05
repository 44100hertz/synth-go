package main

import "./audio"

func main() {
	counter := int32(120)
	seq := func(m *audio.Mixer) {
		m.OnPair(0, func(c *audio.Channel) {
			if counter%6 == 0 {
				c.Note = (counter << 16) * 16 / 19 / 3
				c.DelayTicks = 1
				c.Vol = 0x8000
			}
		})
		counter++
	}
	m := audio.NewMixer(audio.Waves, seq)
	m.OnPair(0, func(c *audio.Channel) {
		c.DelayVol = 0x8000 // A bit of delay attenuation
		c.Filter = 0x5      // Use a delay averaged by 3 samples
		c.Wave = 3
	})

	m.Ch[0].Fade = -0x200
	m.Ch[1].Fade = -0x200
	m.Ch[0].PairMode = audio.PAIR_AM
	audio.Start(&m)
	//audio.WaveOut(&m, "out.raw", 48000)
}
