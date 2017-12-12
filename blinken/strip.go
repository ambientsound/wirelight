package main

import (
	"fmt"
	"image"
	"time"

	"github.com/ambientsound/wirelight/blinken/lib"
	"github.com/ambientsound/wirelight/blinken/pb"
	"github.com/golang/protobuf/proto"
	"github.com/pebbe/zmq4"
)

// The serial is increased by one every time Blinken sends a LED update.
var serial uint64

// Strip represents a strip of LEDs.
type Strip struct {
	sock        *zmq4.Socket
	refreshRate uint64
	width       int
	height      int
}

// NewStrip returns Strip.
func NewStrip(sock *zmq4.Socket, width, height int, refreshRate uint64) *Strip {
	return &Strip{
		sock:        sock,
		refreshRate: refreshRate, // render all LEDs every 15th update
		width:       width,
		height:      height,
	}
}

// rpcLED transfers one LED value to the remote server.
func (s *Strip) rpcLED(led *pb.LED) error {
	serial++
	led.Serial = serial
	led.Render = (serial%s.refreshRate == 0)

	payload, err := proto.Marshal(led)
	if err != nil {
		return fmt.Errorf("while generating protobuf payload: %s", err)
	}

	_, err = s.sock.SendBytes(payload, 0)
	if err != nil {
		return fmt.Errorf("while sending data using ZeroMQ: %s", err)
	}

	return nil
}

// BitBlit transfers image data from an object implementing the Image interface
// to a remote LED server.
func (s *Strip) BitBlit(img image.Image) error {
	led := &pb.LED{}
	for x := 0; x < s.width; x++ {
		for y := 0; y < s.height; y++ {
			led.Index = uint32(y*s.width + x)
			c := img.At(x, y)
			led.Rgb = lib.RGBA(c)
			err := s.rpcLED(led)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Loop renders the LEDs periodically. This function never returns, so be sure
// to call it in a goroutine.
func (s *Strip) Loop(img image.Image, freq int) {
	c := cycleTime(freq)
	for {
		err := s.BitBlit(img)
		if err != nil {
			fmt.Printf("BitBlit: %s\n", err)
		}
		time.Sleep(c)
	}
}
