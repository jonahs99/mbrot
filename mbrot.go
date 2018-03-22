package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math/cmplx"
	"math/rand"
	"os"
)

type point struct {
	pos complex128
	val int
}

func main() {
	minx, miny := -2.0, -1.0
	maxx, maxy := 1.0, 1.0

	scale := 512.0

	w, h := int((maxx-minx)*scale), int((maxy-miny)*scale)

	samples := 20

	points := make([]point, 0, w*h*samples)

	threshhold := 100.0
	iterations := 255

	for i := 0; i < samples; i++ {
		for y := 0; y < h; y++ {
			for x := 0; x < w; x++ {
				c := complex((float64(x)+rand.Float64())/scale+minx, (float64(y)+rand.Float64())/scale+miny)
				m := mbrot(c, threshhold, iterations)
				points = append(points, point{c, m})
			}
		}
		fmt.Printf("%d samples done.\n", i+1)
	}

	pixels := make([]int, w*h)
	for _, p := range points {
		x, y := int((real(p.pos)-minx)*scale), int((imag(p.pos)-miny)*scale)
		pixels[x+y*w] += p.val
	}

	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			var c color.Gray
			c = color.Gray{uint8(pixels[x+y*w] / samples)}
			img.Set(x, y, c)
		}
	}

	file, err := os.Create("out.png")
	defer file.Close()
	if err != nil {
		fmt.Printf("Couldn't create output file!\n")
	}

	png.Encode(file, img)
}

func mbrot(c complex128, threshhold float64, iterations int) int {
	z := complex(0, 0)
	i := 0
	for ; i < iterations; i++ {
		z = cmplx.Pow(z, 2) + c
		if cmplx.Abs(z) > threshhold {
			break
		}
	}
	if i == iterations {
		return 255
	}
	return i
}
