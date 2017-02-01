package main

import "./audio"

func main() {
	seq := func(m *audio.Mixer) {
		m.Channel[0].Note -= 1
	}
	audio.Start(seq)

}
