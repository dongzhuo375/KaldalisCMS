package service

import (
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
)

func decodeImageConfig(r io.Reader) (image.Config, string, error) {
	return image.DecodeConfig(r)
}
