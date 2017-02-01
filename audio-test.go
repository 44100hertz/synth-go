package main

import "./audio"

func main() {
	seq := func(m *audio.Mixer) {
		m.Ch[0].Vol = 0
	}
	m := audio.NewMixer(audio.Waves, seq)
	m.Ch[0].TuneSlide = 0xff0
	m.ResetLevels(0.5)
	m.Ch[2].Vol = 0
	m.Ch[0].Delay = 0x1000
	m.Ch[0].DelayLevel = 0xC000
	audio.Start(&m)
}
