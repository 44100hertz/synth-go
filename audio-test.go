package main

import "./audio"

func main() {
	seq := func(m *audio.Mixer) {
		m.Ch[0].Tune += 0xff0
		m.Ch[2].Tune -= 0xff0
	}
	m := audio.NewMixer(audio.Waves, seq)
	audio.Start(&m)
}
