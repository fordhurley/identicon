package identicon

import (
	"image"
	"image/color"
	"image/draw"
)

// New generates a new identicon image. gridSize is the number of grid cells
// horizontally and vertically. scale multiplies the gridSize to size the image
// in pixels. The final image will be [gridSize*scale x gridSize*scale]
func New(hash []byte, gridSize int, scale int, fg color.Color, bg color.Color) image.Image {
	bits := bitSource{bytes: hash}
	img := image.NewRGBA(image.Rect(0, 0, gridSize*scale, gridSize*scale))

	// x only needs to go halfway across (+1 if odd) because it will be mirrored
	// horizontally:
	maxX := gridSize/2 + gridSize%2
	for x := 0; x < maxX; x++ {
		for y := 0; y < gridSize; y++ {
			color := fg
			if bits.next() {
				color = bg
			}

			src := image.NewUniform(color)

			// Draw this on the left side:
			rect := image.Rect(x*scale, y*scale, (x+1)*scale, (y+1)*scale)
			draw.Draw(img, rect, src, image.ZP, draw.Src)

			// Mirror horizontally to the right side:
			rect = image.Rect((gridSize-x-1)*scale, y*scale, (gridSize-x)*scale, (y+1)*scale)
			draw.Draw(img, rect, src, image.ZP, draw.Src)
		}
	}

	return img
}

type bitSource struct {
	bytes     []byte
	byteIndex uint
	bitIndex  uint
}

func (s *bitSource) next() bool {
	b := s.bytes[s.byteIndex]
	bit := b >> s.bitIndex & 1
	s.bitIndex++
	if s.bitIndex == 8 {
		s.bitIndex = 0
		s.byteIndex++
		s.byteIndex %= uint(len(s.bytes))
	}
	if bit == 1 {
		return true
	}
	return false
}
