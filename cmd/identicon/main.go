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
	fg := color.RGBA{49, 90, 125, 255}
	bg := color.RGBA{237, 243, 248, 255}

	img := identicon.New(os.Stdin, gridSize, scale, fg, bg)
	err := png.Encode(os.Stdout, img)
	if err != nil {
		log.Fatal(err)
	}
}
