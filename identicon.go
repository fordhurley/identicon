package identicon

import (
	"image"
	"image/color"
	"image/draw"
)

// New generates a new identicon image. gridSize is the number of grid cells
// horizontally and vertically. scale multiplies the gridSize to size the image
// in pixels. The final image will be [gridSize*scale x gridSize*scale]
func New(hash []byte, gridSize int, scale int, palette color.Palette) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, gridSize*scale, gridSize*scale))

	i := 0

	// x only needs to go halfway across (+1 if odd) because it will be mirrored
	// horizontally:
	maxX := gridSize/2 + gridSize%2
	for x := 0; x < maxX; x++ {
		for y := 0; y < gridSize; y++ {
			// Grab the next byte:
			b := hash[i]
			i++
			i %= len(hash)

			// Use the byte to pick a color from the palette:
			color := palette[int(b)%len(palette)]
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
