package main

import "./audio"

func main() {
	seq := func(m *audio.Mixer) {
		m.Ch[0].Vol = 0
	}
	m := audio.NewMixer(audio.Waves, seq)
	m.Ch[0].TuneSlide = 0xff0
	m.Ch[0].Wave = 2
	m.Ch[0].Vol = 0x8000
	m.Ch[0].DelayTicks = 10
	m.Ch[0].DelayLevel = 0xC000
	m.Ch[0].DelayFilter = 0x10
	audio.Start(&m)
}
