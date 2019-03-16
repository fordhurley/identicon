package main

import (
	"flag"
	"image/color"
	"image/png"
	"log"
	"os"

	"github.com/fordhurley/identicon"
)

func main() {
	var themeName string

	flag.StringVar(&themeName, "theme", "rainbow", "select the color theme")
	flag.Parse()

	theme, ok := themes[themeName]
	if !ok {
		log.Fatal("invalid theme name:", themeName)
	}

	gridSize := 5
	scale := 32

	img := identicon.New(os.Stdin, gridSize, scale, theme)
	err := png.Encode(os.Stdout, img)
	if err != nil {
		log.Fatal(err)
	}
}

var themes = map[string][]color.Palette{
	"cool": []color.Palette{
		{
			color.RGBA{50, 178, 255, 255},
			color.RGBA{163, 172, 177, 255},
			color.RGBA{0, 0, 0, 255},
		},
		{
			color.RGBA{255, 255, 255, 255},
		},
	},
	"rainbow": []color.Palette{
		{
			color.RGBA{170, 255, 0, 255},
			color.RGBA{255, 170, 0, 255},
			color.RGBA{255, 0, 170, 255},
			color.RGBA{170, 0, 255, 255},
			color.RGBA{0, 170, 255, 255},
		},
		{
			color.RGBA{255, 255, 255, 255},
		},
	},
}
