package main

import (
	"fmt"
	"image"
	"image/color"
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
	} else {
		w = 80
	}
	fmt.Println(os.Args)
	d, _ := os.Open(os.Args[2])
	// This example uses png.Decode which can only decode PNG images.
	// Consider using the general image.Decode as it can sniff and decode any registered image format.
	img, err := png.Decode(d)
	if err != nil {
		log.Fatal(err)
	}
	if os.Args[1] == "print" || os.Args[1] == "p" {
		if len(os.Args) > 3 {
			w, _ = strconv.Atoi(os.Args[3])
		}
		convertToStdOut(img, w)
	} else if os.Args[1] == "scale" || os.Args[1] == "s" {
		i, _ := strconv.Atoi(os.Args[3])
		convertToFile(img, "out.png", i)
	}
}

func convertToFile(img image.Image, fname string, scale int) {

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

	xyratio := float32(25.0) / 25.0
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
