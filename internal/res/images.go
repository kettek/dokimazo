package res

import (
	"bytes"
	"embed"
	"fmt"
	"image"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed images/*
var fs embed.FS

var Images = make(map[string]*Image)
var blankImage *Image

func LoadImage(s string) (*Image, error) {
	if Images[s] != nil {
		return Images[s], nil
	}
	b, err := fs.ReadFile("images/" + s)
	if err != nil {
		return nil, fmt.Errorf("could not read file: %w", err)
	}
	img, _, err := image.Decode(bytes.NewReader(b))
	if err != nil {
		return nil, fmt.Errorf("could not decode image: %w", err)
	}
	Images[s] = NewImage(ebiten.NewImageFromImage(img))
	return Images[s], nil
}

func MustLoadImage(s string) *Image {
	img, err := LoadImage(s)
	if err != nil {
		fmt.Printf("could not load image: %v\n", err)
		return blankImage
	}
	return img
}

func init() {
	blankImage = NewImage(ebiten.NewImage(1, 1))
}
