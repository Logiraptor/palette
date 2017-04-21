package palette

import (
	"image"
	"image/color"
	"math"
	"testing"
	"testing/quick"

	"github.com/stretchr/testify/assert"
)

func TestNColorImage(t *testing.T) {
	err := quick.Check(func(numColors uint8) bool {
		var expectedColors = make(color.Palette, numColors)
		for i := range expectedColors {
			expectedColors[i] = color.Gray{Y: uint8(i)}
		}

		// Create an image with N colors in a line
		testImage := image.NewGray(image.Rect(0, 0, int(numColors), 1))
		for i := range expectedColors {
			testImage.Set(i, 0, expectedColors[i])
		}

		// The color palette generated should be the same as the
		c := ColorCut{}
		palette := c.Quantize(make(color.Palette, 0, numColors), testImage)
		return paletteEqual(t, palette, expectedColors)
	}, nil)
	assert.NoError(t, err)
}

func TestPartitioning(t *testing.T) {
	testImage := image.NewGray(image.Rect(0, 0, 6, 1))
	redLower := color.RGBA{A: 255}
	redUpper := color.RGBA{10, 10, 10, 255}
	greenLower := color.RGBA{A: 255}
	greenUpper := color.RGBA{11, 11, 11, 255}
	blueLower := color.RGBA{A: 255}
	blueUpper := color.RGBA{12, 12, 12, 255}
	testImage.Set(0, 0, redLower)
	testImage.Set(1, 0, redUpper)
	testImage.Set(2, 0, greenLower)
	testImage.Set(3, 0, greenUpper)
	testImage.Set(4, 0, blueLower)
	testImage.Set(5, 0, blueUpper)

	c := ColorCut{}
	palette := c.Quantize(make(color.Palette, 0, 2), testImage)
	paletteEqual(t, palette, color.Palette{greenUpper, redLower})
}

func paletteEqual(t *testing.T, a, b color.Palette) bool {
	if !assert.Equal(t, len(a), len(b)) {
		return false
	}

outer:
	for i := range a {
		for j := range b {
			if colorEqual(a[i], b[j]) {
				continue outer
			}
		}
		return assert.True(t, false, "Missing color: %v", a[i])
	}

	return true
}

func colorEqual(a, b color.Color) bool {
	ar, ag, ab, aa := a.RGBA()
	br, bg, bb, ba := b.RGBA()
	return ar == br && ag == bg && ab == bb && aa == ba
}

func TestBoxFit(t *testing.T) {
	b := box{
		colors: []color.RGBA{
			{1, 2, 3, 4},
			{12, 11, 10, 9},
		},
	}
	b.fit()

	assert.Equal(t, b.min, color.RGBA{1, 2, 3, math.MaxUint8})
	assert.Equal(t, b.max, color.RGBA{12, 11, 10, 0})
	assert.Equal(t, b.rng, color.RGBA{11, 9, 7, 0})
}

func TestBoxLongestEdgeAccessor(t *testing.T) {
	distances := []uint8{1, 2, 3}

	for _, rDistance := range distances {
		for _, gDistance := range distances {
			for _, bDistance := range distances {
				b := &box{
					colors: []color.RGBA{
						{
							R: rDistance,
							G: gDistance,
							B: bDistance,
						},
						{R: 0, G: 0, B: 0},
					},
				}
				b.fit()

				var expectedValue uint8
				var dimension string
				if rDistance >= gDistance && rDistance >= bDistance {
					expectedValue = 10
					dimension = "R"
				} else if gDistance >= rDistance && gDistance >= bDistance {
					expectedValue = 20
					dimension = "G"
				} else {
					expectedValue = 30
					dimension = "B"
				}

				accessor := b.longestEdgeAccessor()
				assert.Equal(t, expectedValue,
					accessor(color.RGBA{R: 10, G: 20, B: 30}),
					"Accessor should return the %s component for color volumes: {%v, %v, %v}",
					dimension, rDistance, gDistance, bDistance)
			}
		}
	}
}
