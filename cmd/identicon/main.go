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

	palette := color.Palette{
		color.RGBA{237, 243, 248, 255},
		color.RGBA{70, 130, 180, 255},
		color.RGBA{49, 90, 125, 255},
		color.RGBA{21, 38, 53, 255},
	}

	img := identicon.New(hash, 5, 36, palette)
	err := png.Encode(os.Stdout, img)
	if err != nil {
		log.Fatal(err)
	}
}
