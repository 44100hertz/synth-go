package main

import "./audio"

func main() {
	seq := func(m *audio.Mixer) {
		m.Ch[0].Tune += 0xff0
		m.Ch[2].Tune -= 0xff0
	}
	m := audio.NewMixer(audio.Waves, seq)
	m.Ch[0].Vol = 0x8000
	m.Ch[2].Vol = 0x8000
	audio.Start(&m)
}
