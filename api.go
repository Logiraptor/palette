package palette

import (
	"image"
	"image/color"
)

func GenerateColors(img image.Image, numColors int) color.Palette {
	c := ColorCut{}
	p := make(color.Palette, 0, numColors)
	return c.Quantize(p, img)
}
