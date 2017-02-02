package main

import "./audio"

func main() {
	seq := func(m *audio.Mixer) {
		m.Ch[0].Vol = 0 // On first tick, mute the wave
	}
	m := audio.NewMixer(audio.Waves, seq)
	m.Ch[0].Wave = 3            // Quarter sine
	m.Ch[0].Vol = 0x10000       // Full volume (single channel)
	m.Ch[0].Note = 1            // C#-0 as base note
	m.Ch[0].DelayNote = 48      // C-5 as note to pluck
	m.Ch[0].DelayLevel = 0xFFF0 // A bit of delay attenuation
	m.Ch[0].Filter = 0x3        // Use a delay averaged by 3 samples
	audio.Start(&m)
}
