package audio

import (
	"os"
)

func WaveOut(m *Mixer, Filename string, srate uint32) {
	f, err := os.Create(Filename)
	defer f.Close()
	if err != nil {
		panic(err)
	}

	output := make(chan int16)
	go m.Start(output, srate)

	bytes := make([]byte, 4)
	for i := 0; i < 0x100000; i++ {
		left := <-output
		bytes[0] = byte(left & 0xff)
		bytes[1] = byte(left >> 8)
		right := <-output
		bytes[2] = byte(right & 0xff)
		bytes[3] = byte(right >> 8)
		_, _ = f.Write(bytes)
	}
}
