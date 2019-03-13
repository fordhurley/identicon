package identicon

import (
	"bytes"
	"crypto/sha1"
	"image"
	"image/color"
	"image/draw"
	"io"
)

// New generates a new identicon image. gridSize is the number of grid cells
// horizontally and vertically. scale multiplies the gridSize to size the image
// in pixels. The final image will be [gridSize*scale x gridSize*scale]
func New(r io.Reader, gridSize int, scale int, fgs []color.Color, bg color.Color) image.Image {
	hash := sha1.New()
	_, err := io.Copy(hash, r)
	if err != nil {
		panic(err)
	}
	byteReader := bytes.NewReader(hash.Sum(nil))

	bits := bitReader{r: byteReader}

	img := image.NewRGBA(image.Rect(0, 0, gridSize*scale, gridSize*scale))

	// x only needs to go halfway across (+1 if odd) because it will be mirrored
	// horizontally:
	maxX := gridSize/2 + gridSize%2
	for x := 0; x < maxX; x++ {
		for y := 0; y < gridSize; y++ {

			color := bg
			if b, _ := bits.read(); b {
				fgIndex, _ := bits.readInt(uint(len(fgs)))
				color = fgs[fgIndex]
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

type bitReader struct {
	r io.ByteReader

	hasByte  bool
	lastByte byte
	bitIndex uint
}

func (br *bitReader) read() (bool, error) {
	if !br.hasByte {
		b, err := br.r.ReadByte()
		if err != nil {
			return false, err
		}
		br.hasByte = true
		br.lastByte = b
		br.bitIndex = 0
	}

	b := br.lastByte
	bit := b >> br.bitIndex & 1
	br.bitIndex++
	if br.bitIndex == 8 {
		br.hasByte = false
	}

	return bit == 1, nil
}

// readInt reads just enough bits to build an integer between 0 and n
// (exclusive), and reconstructs it as an uint.
func (br *bitReader) readInt(n uint) (uint, error) {
	var x uint
	var i uint

	m := n - 1
	for m > 0 {
		m = m >> 1
		b, err := br.read()
		if err != nil {
			return 0, err
		}
		if b {
			x |= 1 << i
		}
		i++
	}

	if x >= n {
		return n - 1, nil
	}
	return x, nil
}
