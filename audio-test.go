package main

import "./audio"

// Temporary SDL code is in main first as not to clutter things
func main() {
	var instr []audio.Inst = []audio.Inst{
		{Index: 0, Len: 0xffff},
	}

	audio.Start(instr)
}
