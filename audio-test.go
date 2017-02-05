package main

import "./audio"
import "fmt"

func main() {
	counter := 0
	seq := func(m *audio.Mixer) {
		m.OnPair(0, func(c *audio.Channel) {
			counter++
			if counter == 60 {
				c.NoteOn = false
				fmt.Println("off")
			}
		})
	}
	m := audio.NewMixer(audio.Waves, seq)
	m.OnPair(0, func(c *audio.Channel) {
		//c.Note = 1            // C#-0 as base note
		c.Wave = audio.WaveQSine
		c.NoteOn = true
		c.Attack = 0x1000
		c.Decay = 0x1000
		c.Sustain = 0x4000
		c.Release = 0x500
		// c.DelayVol = 0x8000
		// c.DelayLoop = 0x11000
		// c.DelayNote = 48 << 16
		// c.FilterLen = 0x4
	})

	audio.Start(&m)
	//audio.WaveOut(&m, "out.raw", 48000)
}
