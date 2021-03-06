package identicon

import (
	"crypto/sha256"
	"image"
	"image/color"
	"image/draw"
)

// New generates an identicon image. gridSize is the number of grid cells
// horizontally and vertically. scale multiplies the gridSize to size the image
// in pixels. The final image will be [gridSize*scale x gridSize*scale]
func New(input string, theme Theme, gridSize int, scale int) image.Image {
	h := sha256.New()
	h.Write([]byte(input))
	hash := h.Sum(nil)

	bitSource := BitSource{bytes: hash}

	colorSource := ColorSource{
		BitSource: bitSource,
		Palettes:  theme[bitSource.NextUint(uint(len(theme)))],
	}

	img := image.NewRGBA(image.Rect(0, 0, gridSize*scale, gridSize*scale))

	// x only needs to go halfway across (+1 if odd) because it will be mirrored
	// horizontally:
	maxX := gridSize/2 + gridSize%2
	for x := 0; x < maxX; x++ {
		for y := 0; y < gridSize; y++ {
			c := image.NewUniform(colorSource.NextColor())

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

// A Theme is a set of sets of color palettes.
type Theme [][]color.Palette

// ColorSource uses a BitSource to determine colors for the identicon. Colors
// are chosen from the available palettes.
type ColorSource struct {
	BitSource
	Palettes []color.Palette
}

// NextColor chooses a color based on the BitSource, first selecting a palette
// uniformly from the available palettes, then selecting uniformly from the
// colors in that palette.
func (cs *ColorSource) NextColor() color.Color {
	index := cs.NextUint(uint(len(cs.Palettes)))
	palette := cs.Palettes[index]
	index = cs.NextUint(uint(len(palette)))
	return palette[index]
}

// BitSource uses a slice of bytes to produce bits for determining features of
// the identicon. The methods provided consume bytes as conservatively as
// possible. If the bytes are exhausted, it will reset and begin consuming bytes
// from the beginning of the slice again.
type BitSource struct {
	bytes     []byte
	byteIndex int
	bitIndex  uint
}

// NewBitSource constructs a BitSource drawing from bytes. Panics if
// len(bytes) < 1.
func NewBitSource(bytes []byte) *BitSource {
	if len(bytes) < 1 {
		panic("identicon: not enough bytes")
	}
	return &BitSource{bytes: bytes}
}

// NextBit consumes a single bit.
func (s *BitSource) NextBit() uint {
	b := s.bytes[s.byteIndex]
	shift := 7 - s.bitIndex // so that we start from MSB
	bit := b >> shift & 1

	s.bitIndex++
	if s.bitIndex == 8 {
		s.bitIndex = 0
		s.byteIndex++
		if s.byteIndex == len(s.bytes) {
			s.byteIndex = 0
		}
	}

	return uint(bit)
}

// NextBool consumes the next bit and returns it as a bool.
func (s *BitSource) NextBool() bool {
	return s.NextBit() == 1
}

// NextUint consumes just enough bits to build an integer between 0 and n
// (exclusive), and reconstructs it as a uint.
func (s *BitSource) NextUint(n uint) uint {
	// Count the number of bits needed:
	nBits := 0
	for m := n - 1; m > 0; m = m >> 1 {
		nBits++
	}

	var x uint
	for shift := nBits - 1; shift >= 0; shift-- {
		x |= s.NextBit() << uint(shift)
	}

	if x > n-1 {
		x = n - 1
	}
	return x
}
