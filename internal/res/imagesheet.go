package res

import (
	"fmt"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type ImageSheet struct {
	*ebiten.Image
	cw, ch int
	rows   int
	cols   int
}

func (i *ImageSheet) At(x, y int) *ebiten.Image {
	x *= i.cw
	y *= i.ch
	if x < 0 || y < 0 || x >= i.cols*i.cw || y >= i.rows*i.ch {
		panic("out of bounds")
	}
	return i.SubImage(image.Rect(x, y, x+i.cw, y+i.ch)).(*ebiten.Image)
}

func (i *ImageSheet) HalfWidth() float64 {
	return float64(i.cw) / 2
}

func (i *ImageSheet) HalfHeight() float64 {
	return float64(i.ch) / 2
}

func (i *ImageSheet) Rows() int {
	return i.rows
}

func (i *ImageSheet) Cols() int {
	return i.cols
}

func (i *ImageSheet) Bounds() image.Rectangle {
	return image.Rect(0, 0, i.cw, i.ch)
}

func NewImageSheet(image *Image, cellw, cellh int) *ImageSheet {
	bounds := image.Bounds()
	w, h := bounds.Size().X, bounds.Size().Y
	cols := w / cellw
	rows := h / cellh
	fmt.Println("wah", cellw, cellh, w, h, rows, cols)
	return &ImageSheet{
		Image: image.Image,
		cw:    cellw,
		ch:    cellh,
		rows:  rows,
		cols:  cols,
	}
}
