package palette

import (
	"image/color"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateTextColors(t *testing.T) {
	type args struct {
		rgb color.RGBA
	}
	tests := []struct {
		name string
		args args
		want color.RGBA
	}{
		{
			name: "Black background gets white text",
			args: args{rgb: white},
			want: black,
		},
		{
			name: "White background gets black text",
			args: args{rgb: black},
			want: white,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TextColor(tt.args.rgb)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GenerateTextColors() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_compositeComponent(t *testing.T) {
	type args struct {
		fgC uint8
		fgA uint8
		bgC uint8
		bgA uint8
		a   uint8
	}
	tests := []struct {
		name string
		args args
		want uint8
	}{
		{
			name: "0 alpha yields 0 color",
			args: args{1, 1, 1, 1, 0},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := compositeComponent(tt.args.fgC, tt.args.fgA, tt.args.bgC, tt.args.bgA, tt.args.a); got != tt.want {
				t.Errorf("compositeComponent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_calculateContrast(t *testing.T) {
	type args struct {
		foreground color.RGBA
		background color.RGBA
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			args: args{foreground: color.RGBA{255, 255, 255, 128}, background: black},
			want: 5.317210002277984,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := calculateContrast(tt.args.foreground, tt.args.background); got != tt.want {
				t.Errorf("calculateContrast() = %v, want %v", got, tt.want)
			}
		})
	}

	assert.Panics(t, func() {
		calculateContrast(white, color.RGBA{
			A: 0,
		})
	})
}

func Test_gammaCorrect(t *testing.T) {
	type args struct {
		x float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			args: args{x: 100},
			want: 0.12743768043564743,
		},
		{
			args: args{x: 2},
			want: 0.000607053967097675,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := gammaCorrect(tt.args.x); got != tt.want {
				t.Errorf("f() = %v, want %v", got, tt.want)
			}
		})
	}
}
