package main

import "./audio"

// Temporary SDL code is in main first as not to clutter things
func main() {
	var instr []audio.Inst = []audio.Inst{
		{Index: 0, Len: 0xffff},
		{Index: nil, Len: nil},
		{Index: nil, Len: nil},
	}

	audio.Start(instr)
}
