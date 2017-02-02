package audio

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
	"fmt"
	"reflect"
	"time"
	"unsafe"
)

var output chan int16 = make(chan int16)

//export callback
func callback(userdata unsafe.Pointer, stream *C.Uint8, length C.int) {
	// Lifted partially from SDL2 audio.go
	n := int(length)
	hdr := reflect.SliceHeader{Data: uintptr(unsafe.Pointer(stream)), Len: n, Cap: n}
	buf := *(*[]C.Uint8)(unsafe.Pointer(&hdr))

	for i := 0; i < n; i += 4 {
		left := <-output
		buf[i] = C.Uint8(left & 0xff)
		buf[i+1] = C.Uint8(left >> 8)
		right := <-output
		buf[i+2] = C.Uint8(right & 0xff)
		buf[i+3] = C.Uint8(right >> 8)
	}
}

// Temporary SDL code is in main first as not to clutter things
func Start(m *Mixer) {
	// Create and set up SDL context
	C.SDL_Init(C.SDL_INIT_AUDIO)
	defer C.SDL_Quit()

	want := C.SDL_AudioSpec{
		freq:     48000,
		format:   C.AUDIO_S16,
		samples:  512,
		channels: 2,
		callback: C.SDL_AudioCallback(C.callback),
	}
	var have C.SDL_AudioSpec
	dev := C.SDL_OpenAudioDevice(nil, 0, &want, &have, 0)
	fmt.Println("Buffer length: ", have.samples)

	// Initialize a mixer
	go m.Start(output, uint32(have.freq))

	// Play 1 second of audio
	C.SDL_PauseAudioDevice(dev, 0)
	time.Sleep(time.Second * 30)
	C.SDL_CloseAudioDevice(dev)
}
