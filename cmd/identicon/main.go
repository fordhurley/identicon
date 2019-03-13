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

	fg := color.RGBA{49, 90, 125, 255}
	bg := color.RGBA{237, 243, 248, 255}

	img := identicon.New(hash, 5, 32, fg, bg)
	err := png.Encode(os.Stdout, img)
	if err != nil {
		log.Fatal(err)
	}
}
