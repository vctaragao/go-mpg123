package out123

/*
#include <stdlib.h>
#include <out123.h>
#cgo LDFLAGS: -lout123
*/
import "C"

import (
	"errors"
	"fmt"
	"unsafe"
)

var EOF = errors.New("EOF")

// All output encoding formats supported by out123
// Contains a handle for and out123 enconder instance
type Encoder struct {
	handle *C.out123_handle
}

///////////////////////////
// DECODER INSTANCE CODE //
///////////////////////////

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

// returns a string containing the most recent error message corresponding to
// an out123 encoder instance
func (e *Encoder) strerror() string {
	return C.GoString(C.out123_strerror(e.handle))
}

func (e *Encoder) Start(rate int64, channels int, encodings int)  error {
    var err C.int
	err = C.out123_start(e.handle, C.long(rate), C.int(channels), C.int(encodings))
    if err < 0{
        return  fmt.Errorf("error getting output format: %s", e.strerror())
    }
	return nil
}


// GetFormat returns current output format
func (e *Encoder) GetFramesize() (int, error) {
    var err C.int
    var framesize *C.int
	err = C.out123_getformat(e.handle, (*C.long)(nil), (*C.int)(nil), (*C.int)(nil), framesize)
    if err < 0{
        return 0, fmt.Errorf("error getting output format: %s", e.strerror())
    }
	return int(*framesize), nil
}


/////////////////////////////
// INPUT AND DECODING CODE //
/////////////////////////////

// Open initializes a encoder to output the mpg123 data read
func (e *Encoder) Open(driver, device string) error {
    var cDriver *C.char
	defer C.free(unsafe.Pointer(cDriver))

    var cDevice *C.char
	defer C.free(unsafe.Pointer(cDevice))

	cDriver = C.CString(driver)
    cDevice = C.CString(device)
	err := C.out123_open(e.handle, cDriver, cDevice)
	if err < 0 {
        return fmt.Errorf("error opening output encoder: %s", e.strerror())
	}

	return nil
}

// Open initializes a encoder to output the mpg123 data read
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

// this function should probably be in a utils package
// since its dosent need the encoder to exist
func (e *Encoder) EncodeName(encoding int) (string, error) {
    var cEnconding *C.char
	cEnconding = C.out123_enc_name(C.int(encoding))
	if cEnconding == nil{
        return "", fmt.Errorf("error getting encode name")
	}

	return C.GoString(cEnconding), nil
}

func (e *Encoder) Play(buffer []byte, readyToRead int)  (int) {
    read := C.out123_play(e.handle, unsafe.Pointer(&buffer), C.size_t(readyToRead))
    return int(read)
}
