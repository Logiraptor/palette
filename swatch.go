package palette

import (
	"image/color"
	"math"
)

// This file translated from https://android.googlesource.com/platform/frameworks/support/+/master/v7/palette/src/main/java/android/support/v7/graphics/Palette.java

var white = color.RGBAModel.Convert(color.White).(color.RGBA)
var black = color.RGBAModel.Convert(color.Black).(color.RGBA)

// TextColor determines whether white or black should
// be used for text against a given background
func TextColor(background color.RGBA) color.RGBA {
	whiteContrast := calculateContrast(white, background)
	blackContrast := calculateContrast(black, background)
	if whiteContrast > blackContrast {
		return white
	}
	return black
}

func compositeColors(foreground color.RGBA, background color.RGBA) color.RGBA {
	a := compositeAlpha(foreground.A, background.A)
	r := compositeComponent(foreground.R, foreground.A, background.R, background.A, a)
	g := compositeComponent(foreground.G, foreground.A, background.G, background.A, a)
	b := compositeComponent(foreground.B, foreground.A, background.B, background.A, a)
	return color.RGBA{
		R: r, G: g, B: b, A: a,
	}
}
func compositeAlpha(foregroundAlpha, backgroundAlpha uint8) uint8 {
	return uint8(0xFF - (((0xFF - uint64(backgroundAlpha)) * (0xFF - uint64(foregroundAlpha))) / 0xFF))
}
func compositeComponent(fgC, fgA, bgC, bgA, a uint8) uint8 {
	var (
		fgC64 = uint64(fgC)
		fgA64 = uint64(fgA)
		bgC64 = uint64(bgC)
		bgA64 = uint64(bgA)
		a64   = uint64(a)
	)
	if a == 0 {
		return 0
	}
	return uint8(((0xFF * fgC64 * fgA64) + (bgC64 * bgA64 * (0xFF - fgA64))) / (a64 * 0xFF))
}

func calculateContrast(foreground, background color.RGBA) float64 {
	if background.A != 255 {
		panic("background cannot be translucent")
	}
	if foreground.A < 255 {
		// If the foreground is translucent, composite the foreground over the background
		foreground = compositeColors(foreground, background)
	}
	luminance1 := calculateLuminance(foreground) + 0.05
	luminance2 := calculateLuminance(background) + 0.05
	// Now return the lighter luminance divided by the darker luminance
	return math.Max(luminance1, luminance2) / math.Min(luminance1, luminance2)
}

func calculateLuminance(c color.RGBA) float64 {
	red := gammaCorrect(float64(c.R))
	green := gammaCorrect(float64(c.G))
	blue := gammaCorrect(float64(c.B))
	return (0.2126 * red) + (0.7152 * green) + (0.0722 * blue)
}

func gammaCorrect(x float64) float64 {
	x /= 255
	if x < 0.03928 {
		return x / 12.92
	}
	return math.Pow((x+0.055)/1.055, 2.4)
}
