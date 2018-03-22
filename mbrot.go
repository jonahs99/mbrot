package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math/cmplx"
	"math/rand"
	"os"

	"github.com/jonahs99/sobel"
)

type point struct {
	pos complex128
	val int
}

func main() {
	minx := flag.Float64("minx", -2.2, "left-most real coordinate")
	miny := flag.Float64("miny", -1.2, "top-most imaginary coordinate")
	maxx := flag.Float64("maxx", 1.2, "right-most real coordinate")
	maxy := flag.Float64("maxy", 1.2, "bottom-most imaginary coordinate")
	scale := flag.Float64("scale", 32.0, "number of pixels per complex unit")
	samples := flag.Int("samples", 10, "number of samples per pixel")
	threshhold := flag.Float64("thresh", 100.0, "divergence threshhold")
	iterations := flag.Int("iterations", 255, "max number of iterations per sample")
	doSobel := flag.Bool("sobel", false, "perform edge detection on output")
	outPath := flag.String("out", "out.png", "path to output image")

	flag.Parse()

	img := vis(*minx, *miny, *maxx, *maxy, *scale, *threshhold, *iterations, *samples)

	if *doSobel {
		img = sobel.ApplySobel(img)
	}

	file, err := os.Create(*outPath)
	defer file.Close()
	if err != nil {
		fmt.Printf("Couldn't create output file!\n")
	}

	err = png.Encode(file, img)
	if err != nil {
		fmt.Printf("Failed to encode image!\n")
	}
}

func vis(minx, miny, maxx, maxy, scale, threshhold float64, iterations, samples int) image.Image {
	w, h := int((maxx-minx)*scale), int((maxy-miny)*scale)

	points := make([]point, 0, w*h*samples)

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
	return img
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
	return i
}
