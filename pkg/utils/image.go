package utils

import (
	"image"
)

func ImgToRaw(img image.Image, tColor uint16) []uint16 /*rgb565*/ {
	var ret []uint16

	bounds := img.Bounds()
	y := 0
	for y < bounds.Dy() {
		x := 0
		for x < bounds.Dx() {
			r, g, b, a := img.At(x, y).RGBA()
			pixel := uint16((r<<8)&0b1111100000000000 | (g<<3)&0b0000011111100000 | (b>>3)&0b0000000000011111)
			if a == 0 {
				pixel = tColor
			}
			ret = append(ret, pixel)
			x += 1
		}
		y += 1
	}

	return ret
}
