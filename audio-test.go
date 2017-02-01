package main

import "./audio"

func main() {
	seq := func(m *audio.Mixer) {
		m.Channel[0].Tune += 0xff0
		m.Channel[2].Tune -= 0xff0
	}
	audio.Start(seq)

}
