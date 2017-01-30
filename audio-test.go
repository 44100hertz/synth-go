package main

import (
	"./audio"
	"time"
)

var instr []audio.Inst = []audio.Inst{
	{Index: 0, Len: 0xffff},
	{Index: nil, Len: nil},
	{Index: nil, Len: nil},
}

// Temporary SDL code is in main first as not to clutter things
func main() {
	m := audio.NewMixer(audio.Waves, instr)
	id := audio.StartMixer(m)
	time.Sleep(time.Second)
	audio.StopMixer(id)
}
