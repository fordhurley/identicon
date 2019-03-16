package main

import (
	"flag"
	"image/color"
	"image/png"
	"io/ioutil"
	"log"
	"os"

	"github.com/fordhurley/identicon"
)

func main() {
	var input string
	var themeName string
	var gridSize int
	var scale int

	flag.StringVar(&input, "input", "-", "input string (- to read from stdin)")
	flag.StringVar(&themeName, "theme", "rainbow", "select the color theme")
	flag.IntVar(&gridSize, "grid", 5, "number of grid elements horizontally and vertically")
	flag.IntVar(&scale, "scale", 32, "pixels per grid element")
	flag.Parse()

	if input == "-" {
		inputBytes, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			log.Fatal(err)
		}
		input = string(inputBytes)
	}

	theme, ok := themes[themeName]
	if !ok {
		log.Fatal("invalid theme name:", themeName)
	}

	img := identicon.New(input, gridSize, scale, theme)

	err := png.Encode(os.Stdout, img)
	if err != nil {
		log.Fatal(err)
	}
}

var themes = map[string]identicon.Theme{
	"cool": {
		{
			{
				color.RGBA{240, 240, 240, 255},
				color.RGBA{163, 172, 177, 255},
			},
			{
				color.RGBA{50, 178, 255, 255},
				color.RGBA{40, 40, 40, 255},
			},
		},
	},
	"rainbow": {
		{
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
	},
	"solid": {
		{
			{color.RGBA{204, 76, 81, 255}},
			{color.RGBA{240, 240, 240, 255}},
		},
		{
			{color.RGBA{121, 134, 203, 255}},
			{color.RGBA{240, 240, 240, 255}},
		},
		{
			{color.RGBA{100, 181, 246, 255}},
			{color.RGBA{240, 240, 240, 255}},
		},
		{
			{color.RGBA{255, 213, 79, 255}},
			{color.RGBA{240, 240, 240, 255}},
		},
		{
			{color.RGBA{129, 199, 132, 255}},
			{color.RGBA{240, 240, 240, 255}},
		},
		{
			{color.RGBA{240, 98, 146, 255}},
			{color.RGBA{240, 240, 240, 255}},
		},
	},
}
