package res

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type Sheet struct {
	*ebiten.Image
	cw, ch int
	rows   int
	cols   int
}

func (i *Sheet) At(x, y int) *ebiten.Image {
	x *= i.cw
	y *= i.ch
	if x < 0 || y < 0 || x >= i.cols*i.cw || y >= i.rows*i.ch {
		panic("out of bounds")
	}
	return i.SubImage(image.Rect(x, y, x+i.cw, y+i.ch)).(*ebiten.Image)
}

func (i *Sheet) Width() float64 {
	return float64(i.cw)
}

func (i *Sheet) Height() float64 {
	return float64(i.ch)
}

func (i *Sheet) HalfWidth() float64 {
	return float64(i.cw) / 2
}

func (i *Sheet) HalfHeight() float64 {
	return float64(i.ch) / 2
}

func (i *Sheet) Rows() int {
	return i.rows
}

func (i *Sheet) Cols() int {
	return i.cols
}

func (i *Sheet) Bounds() image.Rectangle {
	return image.Rect(0, 0, i.cw, i.ch)
}

func NewSheet(img *ebiten.Image, cellw, cellh int) *Sheet {
	bounds := img.Bounds()
	w, h := bounds.Size().X, bounds.Size().Y
	cols := w / cellw
	rows := h / cellh
	return &Sheet{
		Image: img,
		cw:    cellw,
		ch:    cellh,
		rows:  rows,
		cols:  cols,
	}
}
