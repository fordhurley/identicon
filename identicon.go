package identicon

import (
	"image"
	"image/color"
	"image/draw"
)

// New generates a new identicon image. gridSize is the number of grid cells
// horizontally and vertically. scale multiplies the gridSize to size of the
// image in pixels.
func New(hash []byte, gridSize int, scale int, palette color.Palette) image.Image {
	ring := byteRing{bytes: hash}

	img := image.NewRGBA(image.Rect(0, 0, gridSize*scale, gridSize*scale))

	for x := 0; x < gridSize/2+gridSize%2; x++ {
		for y := 0; y < gridSize; y++ {
			c := int(ring.next()) % len(palette)
			color := image.NewUniform(palette[c])

			rect := image.Rect(x*scale, y*scale, (x+1)*scale, (y+1)*scale)
			draw.Draw(img, rect, color, image.ZP, draw.Src)

			// Mirror horizontally:
			rect = image.Rect((gridSize-x-1)*scale, y*scale, (gridSize-x)*scale, (y+1)*scale)
			draw.Draw(img, rect, color, image.ZP, draw.Src)
		}
	}

	return img
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
