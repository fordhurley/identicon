package main

import (
	"crypto/sha1"
	"image/color"
	"image/png"
	"io"
	"log"
	"os"

	"github.com/fordhurley/identicon"
)

func main() {
	h := sha1.New()
	io.Copy(h, os.Stdin)
	hash := h.Sum(nil)

	size := 256

	palette := color.Palette{
		color.RGBA{255, 0, 0, 255},
		color.RGBA{255, 255, 0, 255},
		color.RGBA{0, 0, 255, 255},
		color.RGBA{255, 255, 255, 255},
	}

	img := identicon.New(hash, size, palette)
	err := png.Encode(os.Stdout, img)
	if err != nil {
		log.Fatal(err)
	}
}
