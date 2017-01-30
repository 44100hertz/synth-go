package main

/*
typedef unsigned char Uint8;
void callback(void *userdata, Uint8 *stream, int len);
#include <SDL2/SDL.h>
#ifdef _WIN32

#else
  #cgo CFLAGS : -I/usr/include/SDL2 -D_REENTRANT
  #cgo LDFLAGS : -lSDL2
#endif
*/
import "C"
import (
	"reflect"
	"time"
	"unsafe"

	"./audio"
)

var instr []audio.Inst = []audio.Inst{
	{Index: 0, Len: 0xffff},
	{Index: nil, Len: nil},
	{Index: nil, Len: nil},
}
var output chan int16 = make(chan int16)

//export callback
func callback(userdata unsafe.Pointer, stream *C.Uint8, length C.int) {
	// Lifted partially from SDL2 audio.go
	n := int(length)
	hdr := reflect.SliceHeader{Data: uintptr(unsafe.Pointer(stream)), Len: n, Cap: n}
	buf := *(*[]C.Uint8)(unsafe.Pointer(&hdr))

	for i := 0; i < n; i += 2 {
		nextSamp := <-output
		buf[i] = C.Uint8(nextSamp & 0xff)
		buf[i+1] = C.Uint8(nextSamp >> 8)
	}
}

// Temporary SDL code is in main first as not to clutter things
func main() {
	C.SDL_Init(C.SDL_INIT_AUDIO)
	defer C.SDL_Quit()

	const bufSize uint16 = 1024

	want := C.SDL_AudioSpec{
		freq:     48000,
		format:   C.AUDIO_S16,
		samples:  C.Uint16(bufSize),
		channels: 1,
		callback: C.SDL_AudioCallback(C.callback),
	}
	var have C.SDL_AudioSpec
	go audio.Init(audio.Waves, instr, output)

	dev := C.SDL_OpenAudioDevice(nil, 0, &want, &have, 0)
	C.SDL_PauseAudioDevice(dev, 0)
	time.Sleep(time.Second)
	C.SDL_CloseAudioDevice(dev)
}
