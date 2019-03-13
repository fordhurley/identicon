package main

import (
	"image/color"
	"image/png"
	"log"
	"os"

	"github.com/fordhurley/identicon"
)

func main() {
	gridSize := 5
	scale := 32
	fgs := []color.Color{
		color.RGBA{50, 178, 255, 255},
		color.RGBA{163, 172, 177, 255},
		color.RGBA{0, 0, 0, 255},
	}
	bg := color.RGBA{255, 255, 255, 255}

	img := identicon.New(os.Stdin, gridSize, scale, fgs, bg)
	err := png.Encode(os.Stdout, img)
	if err != nil {
		log.Fatal(err)
	}
}
