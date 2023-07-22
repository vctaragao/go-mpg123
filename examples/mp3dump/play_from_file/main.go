package main

import (
	"fmt"
	"github.com/vctaragao/go-mpg123/mpg123"
	"os"
)

func main() {
	// check command-line arguments
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: mp123play <infile.mp3> <outfile.raw>")
		return
	}

	// create mpg123 decoder instance
	decoder, err := mpg123.NewDecoder("")
	if err != nil {
		panic("could not initialize mpg123")
	}

	// open a file with decoder
	err = decoder.Open(os.Args[1])
	if err != nil {
		panic("error opening mp3 file")
	}
	defer decoder.Close()

	// get audio format information
    rate, chans, enc :=  decoder.GetFormat()
	fmt.Fprintln(os.Stderr, "Encoding: Signed 16bit")
	fmt.Fprintln(os.Stderr, "Sample Rate:", rate)
	fmt.Fprintln(os.Stderr, "Channels:", chans)
    fmt.Fprintln(os.Stderr, "Enconding:", enc)

	// make sure output format does not change
	decoder.FormatNone()
	decoder.Format(rate, chans, mpg123.ENC_SIGNED_16)

	// open output file
	o, err := os.Create(os.Args[2])
	if err != nil {
		panic("error opening output file")
	}
	defer o.Close()

	// decode mp3 file and dump output
	buf := make([]byte, 2048*16)
	for {
		len, err := decoder.Read(buf)
		o.Write(buf[0:len])
		if err != nil {
			break
		}
	}
	o.Close()
	decoder.Delete()
}
