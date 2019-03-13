package identicon

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
)

// New generates a new identicon image. numPixels is the side length of the
// square image in pixels.
func New(hash []byte, numPixels int, palette color.Palette) (image.Image, error) {
	const minBytes = 7 // 3 for fg, 3 for bg, and at least 1 for pixel data
	if len(hash) < minBytes {
		return nil, fmt.Errorf("identicon: need at least %d bytes of hash", minBytes)
	}

	fg := palette.Convert(byteColor(hash[0:3]))
	bg := palette.Convert(byteColor(hash[3:6]))
	pixelBytes := byteRing{bytes: hash[6:len(hash)]}

	rect := image.Rect(0, 0, numPixels, numPixels)
	img := image.NewPaletted(rect, color.Palette{fg, bg})

	for _, cell := range grid(rect, pixelBytes.len()) {
		src := image.NewUniform(byteColor(pixelBytes.nextN(3)))
		draw.Draw(img, cell, src, image.ZP, draw.Src)
	}

	return img, nil
}

func byteColor(bytes []byte) color.Color {
	return color.RGBA{
		R: uint8(bytes[0]),
		G: uint8(bytes[1]),
		B: uint8(bytes[2]),
		A: 255,
	}
}

type byteRing struct {
	bytes []byte
	index int
}

func (r *byteRing) next() byte {
	b := r.bytes[r.index]
	r.index++
	r.index %= len(r.bytes)
	return b
}

func (r *byteRing) nextN(n int) []byte {
	var out []byte
	for i := 0; i < n; i++ {
		out = append(out, r.next())
	}
	return out
}

func (r *byteRing) len() int {
	return len(r.bytes)
}

func grid(rect image.Rectangle, numPixelBytes int) []image.Rectangle {
	return []image.Rectangle{
		image.Rect(0, 0, rect.Size().X/2, rect.Size().Y/2),
		image.Rect(rect.Size().X/2, 0, rect.Size().X, rect.Size().Y/2),
		image.Rect(0, rect.Size().Y/2, rect.Size().X/2, rect.Size().Y),
		image.Rect(rect.Size().X/2, rect.Size().Y/2, rect.Size().X, rect.Size().Y),
	}
}
