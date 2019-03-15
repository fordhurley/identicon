package identicon

import (
	"crypto/sha256"
	"image"
	"image/color"
	"image/draw"
	"io"
)

// New generates a new identicon image. gridSize is the number of grid cells
// horizontally and vertically. scale multiplies the gridSize to size the image
// in pixels. The final image will be [gridSize*scale x gridSize*scale]
func New(r io.Reader, gridSize int, scale int, fgs []color.Color, bg color.Color) image.Image {
	hash := sha256.New()
	_, err := io.Copy(hash, r)
	if err != nil {
		panic(err)
	}

	source := newColorSource(hash.Sum(nil), fgs, bg)

	img := image.NewRGBA(image.Rect(0, 0, gridSize*scale, gridSize*scale))

	// x only needs to go halfway across (+1 if odd) because it will be mirrored
	// horizontally:
	maxX := gridSize/2 + gridSize%2
	for x := 0; x < maxX; x++ {
		for y := 0; y < gridSize; y++ {
			c := image.NewUniform(source.nextColor())

			// Draw this on the left side:
			rect := image.Rect(x*scale, y*scale, (x+1)*scale, (y+1)*scale)
			draw.Draw(img, rect, c, image.ZP, draw.Src)

			// Mirror horizontally to the right side:
			rect = image.Rect((gridSize-x-1)*scale, y*scale, (gridSize-x)*scale, (y+1)*scale)
			draw.Draw(img, rect, c, image.ZP, draw.Src)
		}
	}

	return img
}

type colorSource struct {
	bitSource
	fgs []color.Color
	bg  color.Color
}

func newColorSource(bytes []byte, fgs []color.Color, bg color.Color) *colorSource {
	return &colorSource{
		bitSource: bitSource{bytes: bytes},
		fgs:       fgs,
		bg:        bg,
	}
}

func (cs *colorSource) nextColor() color.Color {
	// 50% chance of picking the background color:
	if cs.nextBool() {
		return cs.bg
	}
	fgIndex := cs.nextUint(uint(len(cs.fgs)))
	return cs.fgs[fgIndex]
}

type bitSource struct {
	bytes []byte

	byteIndex int
	bitIndex  uint
}

func (s *bitSource) nextBool() bool {
	b := s.bytes[s.byteIndex]
	bit := b >> s.bitIndex & 1

	s.bitIndex++
	if s.bitIndex == 8 {
		s.bitIndex = 0
		s.byteIndex++
		if s.byteIndex == len(s.bytes) {
			s.byteIndex = 0
		}
	}

	return bit == 1
}

// readInt reads just enough bits to build an integer between 0 and n
// (exclusive), and reconstructs it as an uint.
func (s *bitSource) nextUint(n uint) uint {
	var x uint
	var i uint

	for m := n - 1; m > 0; m = m >> 1 {
		b := s.nextBool()
		if b {
			x |= 1 << i
		}
		i++
	}

	if x >= n {
		return n - 1
	}
	return x
}
