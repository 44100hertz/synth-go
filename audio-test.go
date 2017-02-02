package main

import "./audio"

func main() {
	seq := func(m *audio.Mixer) {
		m.Ch[0].Vol = 0
	}
	m := audio.NewMixer(audio.Waves, seq)
	m.Ch[0].TuneSlide = 0xff0
	m.Ch[0].Wave = 3
	m.Ch[0].Vol = 0x10000
	m.Ch[0].DelayNote = 48
	m.Ch[0].DelayLevel = 0x10000
	m.Ch[0].DelayFilter = 0x2
	audio.Start(&m)
}
