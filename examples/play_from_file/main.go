package main

import (
	"fmt"
	"os"

	"github.com/vctaragao/go-mpg123/mpg123"
	"github.com/vctaragao/go-mpg123/out123"
)

const SECONDS = 16

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: mp123play <infile.mp3>")
		return
	}

	decoder, err := mpg123.NewDecoder()
	chk(err)
	defer decoder.Delete()

	chk(decoder.Open(os.Args[1]))
	defer decoder.Close()

	rate, chans, enc := decoder.GetFormat()
	fmt.Println("Sample Rate:", rate)
	fmt.Println("Channels:", chans)
	fmt.Println("Enconding:", enc)

	encoder, err := out123.NewEncoder()
	chk(err)
	defer encoder.Delete()

	chk(encoder.Open("", ""))

	driver, device, err := encoder.DriverInfo()
	chk(err)

	if device == "" {
		device = "default"
	}

	fmt.Printf("Effective output: %s\n", device)
	fmt.Printf("Effective output driver: %s\n", driver)

	decoder.FormatNone()
	decoder.Format(rate, chans, enc)

	encoding, err := encoder.EncodeName(enc)
	chk(err)
	fmt.Printf("Playing with %d channels and %d Hz, encoding %s.\n", chans, rate, encoding)

	chk(encoder.Start(rate, chans, enc))

	framesize, err := encoder.GetFramesize()
	chk(err)

	buffer_size_default := decoder.OutBlock()
	buffer_size := buffer_size_default * ((int(rate) * framesize / buffer_size_default) * SECONDS)
	buff := make([]byte, buffer_size)

	for {
		done, err := decoder.Read(buff)
		if err == mpg123.EOF {
			fmt.Println("Song finished")
			break
		}
		chk(err)

		if played := encoder.Play(buff, done); played != done {
			fmt.Printf("Warning: written less than gotten from libmpg123: %d != %d\n", played, done)
		}
	}

}

func chk(err error) {
	if err != nil {
		panic(err)
	}
}
