package out123

/*
#include <stdlib.h>
#include <out123.h>
#cgo LDFLAGS: -lout123
*/
import "C"

import (
	"fmt"
	"unsafe"
)

// All output encoding formats supported by out123
// Contains a handle for and out123 enconder instance
type Encoder struct {
	handle *C.out123_handle
}

// NewEncoder creates a new out123 encoder instance
func NewEncoder() (*Encoder, error) {
	var err C.int
	oh := C.out123_new()
	if oh == nil {
		errstring := C.out123_plain_strerror(err)
		defer C.free(unsafe.Pointer(errstring))
		return nil, fmt.Errorf("error initializing out123 encoder: %s", C.GoString(errstring))
	}
	enc := new(Encoder)
	enc.handle = oh
	return enc, nil
}

// Delete frees an out123 encoder instance
func (e *Encoder) Delete() {
	C.out123_del(e.handle)
}

// strerror returns a string containing the most recent error message corresponding to
// an out123 encoder instance
func (e *Encoder) strerror() string {
	return C.GoString(C.out123_strerror(e.handle))
}

// Start playback with a certain output format
// It might be a good idea to have audio data handy to feed after this
// returns with success.
// Rationale for not taking a pointer to struct mpg123_fmt: This would
// always force you to deal with that type and needlessly enlarge the
// shortest possible program.
// param rate sampling rate
// param encoding sample encoding (values matching libmpg123 API)
// param channels number of channels (1 or 2, usually)
func (e *Encoder) Start(rate int64, channels int, encodings int) error {
	var err C.int
	err = C.out123_start(e.handle, C.long(rate), C.int(channels), C.int(encodings))
	if err < 0 {
		return fmt.Errorf("error getting output format: %s", e.strerror())
	}
	return nil
}

// GetFramesize returns encoder framesize
func (e *Encoder) GetFramesize() (int, error) {
	var err C.int
	var framesize C.int
	err = C.out123_getformat(e.handle, (*C.long)(nil), (*C.int)(nil), (*C.int)(nil), &framesize)
	if err < 0 {
		return 0, fmt.Errorf("error getting output format: %s", e.strerror())
	}
	return int(framesize), nil
}

// Open initializes a encoder to output the mpg123 data read
func (e *Encoder) Open(driver, device string) error {
	var cDriver *C.char
	defer C.free(unsafe.Pointer(cDriver))

	var cDevice *C.char
	defer C.free(unsafe.Pointer(cDevice))

	if device != "" {
		cDevice = C.CString(device)
	}

	if device != "" {
		cDriver = C.CString(driver)
	}

	err := C.out123_open(e.handle, cDriver, cDevice)
	if err < 0 {
		return fmt.Errorf("error opening output encoder: %s", e.strerror())
	}

	return nil
}

// DriverInfo return the driver and device info of the encoder
func (e *Encoder) DriverInfo() (string, string, error) {
	var cDriver *C.char
	defer C.free(unsafe.Pointer(cDriver))

	var cDevice *C.char
	defer C.free(unsafe.Pointer(cDevice))

	err := C.out123_driver_info(e.handle, &cDriver, &cDevice)
	if err < 0 {
		return "", "", fmt.Errorf("error opening output encoder: %s", e.strerror())
	}

	return C.GoString(cDriver), C.GoString(cDevice), nil
}

// EncodeName return the name of enconding
// param encoding code (enum #mpg123_enc_enum)
// This function should probably be in a utils package
// since its dosent need the encoder to exist
func (e *Encoder) EncodeName(encoding int) (string, error) {
	var cEnconding *C.char
	cEnconding = C.out123_enc_name(C.int(encoding))
	if cEnconding == nil {
		return "", fmt.Errorf("error getting encode name")
	}

	return C.GoString(cEnconding), nil
}

// Play hand over data for playback and wait in case audio device is busy.
// So, per default, if you provided a byte count divisible by the PCM frame size,
// it is an error when less bytes than given are played.
// param buffer pointer to raw audio data to be played
// param bytes number of bytes to read from the buffer
// return number of bytes played (might be less than given, even zero)
func (e *Encoder) Play(buffer []byte, readyToRead int) int {
	read := C.out123_play(e.handle, unsafe.Pointer(&buffer[0]), C.size_t(readyToRead))
	return int(read)
}
