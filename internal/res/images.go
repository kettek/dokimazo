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
//go:embed all:images
//go:embed details/*
var fs embed.FS

const SheetCellWidth = 16
const SheetCellHeight = 16

type sizeDef struct {
	s    string
	w, h int
}

var Sheets = make(map[sizeDef]*Sheet)
var blankSheet *Sheet

func LoadSheet(s string) (*Sheet, error) {
	return LoadSheetWithSize(s, SheetCellWidth, SheetCellHeight)
}

func LoadSheetWithSize(s string, w, h int) (*Sheet, error) {
	sd := sizeDef{s, w, h}
	if Sheets[sd] != nil {
		return Sheets[sd], nil
	}
	b, err := fs.ReadFile("images/" + s)
	if err != nil {
		return nil, fmt.Errorf("could not read file: %w", err)
	}
	img, _, err := image.Decode(bytes.NewReader(b))
	if err != nil {
		return nil, fmt.Errorf("could not decode image: %w", err)
	}
	Sheets[sd] = NewSheet(ebiten.NewImageFromImage(img), w, h)
	return Sheets[sd], nil
}

func MustLoadSheet(s string) *Sheet {
	return MustLoadSheetWithSize(s, SheetCellWidth, SheetCellHeight)
}

func MustLoadSheetWithSize(s string, w, h int) *Sheet {
	img, err := LoadSheetWithSize(s, w, h)
	if err != nil {
		fmt.Printf("could not load image: %v\n", err)
		return blankSheet
	}
	return img
}

func init() {
	blankSheet = NewSheet(ebiten.NewImage(1, 1), 1, 1)
}
