package main

// typedef unsigned char Uint8;
// void callback(void *userdata, Uint8 *stream, int len);
import "C"
import (
	"reflect"
	"time"
	"unsafe"

	"./audio"
	"github.com/veandco/go-sdl2/sdl"
)

//export callback
func callback(userdata unsafe.Pointer, stream *C.Uint8, length C.int) {
	// Lifted partially from SDL2 audio.go
	n := int(length)
	hdr := reflect.SliceHeader{Data: uintptr(unsafe.Pointer(stream)), Len: n, Cap: n}
	buf := *(*[]C.Uint8)(unsafe.Pointer(&hdr))

	channel := *(*chan int)(userdata)
	for i := 0; i < n; i++ {
		buf[i] = C.Uint8(<-channel >> 8)
		buf[i+1] = C.Uint8(<-channel & 0x00ff)
	}
}

// Temporary SDL code is in main first as not to clutter things
func main() {
	// Fill the sequence channel with a sequence C4->C5
	instr := []audio.Inst{
		{Index: 0, Len: 0xffff},
		{Index: nil, Len: nil},
		{Index: nil, Len: nil},
	}

	output := make(chan int16)

	// Start the mixer running
	go audio.Init(audio.Waves, instr, output)
	sdl.Init(sdl.INIT_AUDIO)
	defer sdl.Quit()

	const bufSize uint16 = 1024

	want := sdl.AudioSpec{
		Freq:     48000,
		Format:   sdl.AUDIO_S16,
		Samples:  bufSize,
		Channels: 1,
		Callback: sdl.AudioCallback(C.callback),
		UserData: unsafe.Pointer(&shitdicks),
	}
	var have sdl.AudioSpec

	dev, err := sdl.OpenAudioDevice("", false, &want, &have, 0)
	if err != nil {
		panic(err)
	}
	sdl.PauseAudioDevice(dev, false)
	time.Sleep(1 * time.Second)
	sdl.CloseAudioDevice(dev)

	// ***** Deprecated method attempt *****
	// callback := func(userdata, new([1024]uint8), bufSize) {
	// }

	// // Initialize mixer data with wave and sequence
	// desired := sdl.AudioSpec{
	// 	freq: 48000,
	// 	format: sdl.AUDIO_S16,
	// 	samples: bufSize,
	// 	callback: callback
	// 	userdata: userdata
	// }
	// var obtained *sdl.AudioSpec
	// sdl.OpenAudio(&desired, obtained)
}
