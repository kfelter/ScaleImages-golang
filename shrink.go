package main

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"strconv"
)

func main() {
	var w int
	if len(os.Args) < 2 {
		fmt.Println("usage: mode(print,scale) input.png int(width[optional]/scale factor) output.png(optional, defaults to out.png)")
		return
	}
	fmt.Println(os.Args)
	d, err := os.Open(os.Args[2])
	if err != nil {
		log.Fatal(err)
	}
	// This example uses png.Decode which can only decode PNG images.
	// Consider using the general image.Decode as it can sniff and decode any registered image format.
	found := false
	var img image.Image
	img, err = jpeg.Decode(d)
	//fmt.Println(u)
	if err != nil {
		log.Println(img, err)
	} else {
		found = true
	}
	if !found {
		img, err = png.Decode(d)
		if err != nil {
			log.Println(err)
		}
	}

	if os.Args[1] == "print" || os.Args[1] == "p" {
		w = 80
		if len(os.Args) > 3 {
			w, _ = strconv.Atoi(os.Args[3])
		}
		convertToStdOut(img, w)
	} else if os.Args[1] == "scale" || os.Args[1] == "s" {
		i, _ := strconv.Atoi(os.Args[3])
		convertToFile(img, "out.png", i)
	} else if os.Args[1] == "square" {
		CenterSquare(img, "out.png")
	} else if os.Args[1] == "pixel" {
		i, _ := strconv.Atoi(os.Args[3])
		SameResolutionPixelation(img, i)
	}
}

func scaleImg(img image.Image, scale int) image.Image {
	a := image.NewNRGBA(image.Rect(0, 0, img.Bounds().Max.X/scale, img.Bounds().Max.Y/scale))
	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
		if y%scale == 0 {
			for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
				if x%scale == 0 {
					a.Set(x/scale, y/scale, img.At(x, y))
				}
			}
		}
	}
	return a
}
func convertToFile(img image.Image, fname string, scale int) {

	a := scaleImg(img, scale)

	f, err := os.Create(fname)
	if err != nil {
		log.Fatal(err)
	}

	if err := png.Encode(f, a); err != nil {
		f.Close()
		log.Fatal(err)
	}

	if err := f.Close(); err != nil {
		log.Fatal(err)
	}

}

func convertToStdOut(img image.Image, i int) {

	var xmod, ymod float32

	xyratio := float32(25.0) / 35.0
	if i >= img.Bounds().Max.X {
		xmod = 1.0
		ymod = 1.0
	} else {
		xmod = float32(img.Bounds().Max.X) / float32(i)
		ymod = xmod / xyratio
	}

	levels := []string{" ", "░", "▒", "▓", "█"}

	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
		if y%int(ymod) == 0 {
			for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
				if x%int(xmod) == 0 {
					c := color.GrayModel.Convert(img.At(x, y)).(color.Gray)
					level := c.Y / 51 // 51 * 5 = 255
					if level == 5 {
						level--
					}
					fmt.Print(levels[level])

				}
			}
			fmt.Print("\n")
		}

	}
}

// CenterSquare makes an image a square
func CenterSquare(img image.Image, fname string) {
	//get the smaller dimention
	max := img.Bounds().Max.X
	if img.Bounds().Max.Y < max {
		max = img.Bounds().Max.Y
	}
	a := image.NewNRGBA(image.Rect(0, 0, max, max))
	//cut off either side of the bigger dim to get a square
	//could be optimized
	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
			a.Set(x, y, img.At(x, y))
		}
	}

	f, err := os.Create(fname)
	if err != nil {
		log.Fatal(err)
	}

	if err := png.Encode(f, a); err != nil {
		f.Close()
		log.Fatal(err)
	}

	if err := f.Close(); err != nil {
		log.Fatal(err)
	}

}

// IncreaseContrast increases the contrast in an image
func IncreaseContrast() {

}

// SameResolutionPixelation pixelates images but keeps the resolution
func SameResolutionPixelation(img image.Image, scale int) {
	//get smaller side
	width := img.Bounds().Max.X
	height := img.Bounds().Max.Y
	var smallside int
	if width > height {
		smallside = height
	} else {
		smallside = width
	}

	//find the width that will devide into scale well
	for smallside%scale != 0 {
		smallside--
	}

	a := image.NewNRGBA(image.Rect(0, 0, smallside, smallside))
	//cut off either side of the bigger dim to get a square
	//could be optimized
	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
			a.Set(x, y, img.At(x, y))
		}
	}

	//now we have an image that can be devided into a grid of pixels that are size scale*scale
	midpoint := scale / 2
	xoff, yoff := 0, 0
	for y := a.Bounds().Min.Y; y < a.Bounds().Max.Y; y++ {
		for x := a.Bounds().Min.X; x < a.Bounds().Max.X; x++ {
			a.Set(x, y, img.At(midpoint+xoff*scale, midpoint+yoff*scale))
			if x%scale == 0 {
				xoff++
			}
			//a.Set(x, y, img.At(x, y))
		}
		xoff = 0
		if y%scale == 0 {
			yoff++
		}
	}

	f, err := os.Create("out.png")
	if err != nil {
		log.Fatal(err)
	}

	if err := png.Encode(f, a); err != nil {
		f.Close()
		log.Fatal(err)
	}

	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}
