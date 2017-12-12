package main

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"math/rand"
	"time"

	"github.com/ambientsound/wirelight/blinken/lib"
	colorful "github.com/lucasb-eyer/go-colorful"
)

func fill(canvas *image.RGBA, col color.Color) {
	b := canvas.Bounds()
	for x := b.Min.X; x < b.Max.X; x++ {
		for y := b.Min.Y; y < b.Max.Y; y++ {
			canvas.Set(x, y, col)
		}
	}
}

func northernLights(canvas *image.RGBA) {
	b := canvas.Bounds()
	old := make([]colorful.Color, b.Max.X*b.Max.Y)
	for {
		for angle := 0.0; angle < 360.0; angle++ {
			for x := b.Min.X; x < b.Max.X; x++ {
				for y := b.Min.Y; y < b.Max.Y; y++ {
					i := (y * b.Max.X) + x
					col := colorful.Hsl(angle+rand.Float64()*50.0, 1, rand.Float64()*0.1)
					step := col.BlendHcl(old[i], 0.92).Clamped()
					canvas.Set(x, y, step)
					old[i] = step
				}
			}
			time.Sleep(time.Millisecond * 100)
		}
	}
}

func black(canvas *image.RGBA) {
	for {
		fill(canvas, colorful.Hsv(0, 0, 0))
	}
}

func white(canvas *image.RGBA) {
	for {
		hue := rand.Float64() * 360.0
		for deg := 0.0; deg <= 180.0; deg += 1 {
			l := math.Sin(lib.Rad(deg))
			col := colorful.Hsv(hue, 1.0, l*0.5).Clamped()
			fill(canvas, col)
			time.Sleep(time.Microsecond * 1500)
		}
		time.Sleep(time.Millisecond * 185)
	}
}

// directionTest draws up a gradient on each strip.
func directionTest(canvas *image.RGBA) {
	c := 1.0
	l := 0.05

	src := colorful.Hcl(0.0, c, l)
	dst := colorful.Hcl(160.0, c, l)
	b := canvas.Bounds()
	count := b.Max.X - b.Min.X
	step := float64(1.0) / float64(count)

	for {
		for y := b.Min.Y; y < b.Max.Y; y++ {
			n := 0.0
			for x := b.Min.X; x < b.Max.X; x++ {
				n += step
				col := src.BlendHcl(dst, n).Clamped()
				canvas.Set(x, y, col)
			}
		}
		time.Sleep(time.Millisecond * 1000)
	}
}

func gradients(canvas *image.RGBA) {
	var h, c, l float64
	h = 0.0
	c = 0.8
	l = 0.5
	_, _ = c, l
	src := colorful.Hsv(h, 1, 1)
	dst := colorful.Hsv(h, 1, 1)

	for {
		src = dst
		h += 30
		if h >= 360 {
			h = 0
		}
		dst = colorful.Hsv(h, 1, 1)
		fmt.Printf("hue=%.2f, blend %#v %#v\n", h, src, dst)

		// interpolate between the two colors.
		for n := 0.0; n < 1.0; n += 0.01 {
			col := src.BlendHcl(dst, n).Clamped()
			fill(canvas, col)
			time.Sleep(time.Millisecond * 20)
		}
	}
}

func staccatoWheel(canvas *image.RGBA) {
	var h float64
	for {
		h += 31
		if h > 360 {
			h -= 360
		}
		col := colorful.Hsv(h, 1, 0.25).Clamped()
		fill(canvas, col)
		time.Sleep(time.Millisecond * 250)
	}
}

func wheelHCL(canvas *image.RGBA) {
	var h float64
	for {
		h += 1
		if h > 360 {
			h = 0
		}
		col := colorful.Hcl(h, 0.2, 0).Clamped()
		fill(canvas, col)
		time.Sleep(time.Millisecond * 10)
	}
}

func wheelHSV(canvas *image.RGBA) {
	var h float64
	for {
		h += 1
		if h > 360 {
			h = 0
		}
		col := colorful.Hsv(h, 1, 1)
		fill(canvas, col)
		time.Sleep(time.Millisecond * 10)
	}
}
