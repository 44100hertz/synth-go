package main

import "./audio"

func main() {
	seq := func(m *audio.Mixer) {
	}
	m := audio.NewMixer(audio.Waves, seq)
	m.Ch[0].TuneSlide = 0xff0
	m.Ch[2].TuneSlide = -0xff0
	audio.Start(&m)
}
