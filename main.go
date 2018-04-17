package main

import (
	"fmt"
	"github.com/blackjack/webcam"
	"io/ioutil"
	"os"
)

const timeout = 96

func check(e error) {
	if e != nil {
		panic(e)
	}
}
func main() {

	d1 := []byte("hello\ngo\n")
	err := ioutil.WriteFile("/tmp/dat1", d1, 0644)
	check(err)

	f, err := os.Create("/tmp/dat2")
	check(err)
	defer f.Close()

	// Open webcam
	cam, err := webcam.Open("/dev/video0")
	if err != nil {
		panic(err.Error())
	}
	defer cam.Close()

	err = cam.StartStreaming()
	if err != nil {
		panic(err.Error())
	}
	for {
		err = cam.WaitForFrame(timeout)

		switch err.(type) {
		case nil:
		case *webcam.Timeout:
			fmt.Fprint(os.Stderr, err.Error())
			continue
		default:
			panic(err.Error())
		}

		frame, err := cam.ReadFrame()
		if len(frame) != 0 {
			_, err = f.Write(frame)
			check(err)
			f.Sync()
			break
		} else if err != nil {
			panic(err.Error())
		}
	}
}
