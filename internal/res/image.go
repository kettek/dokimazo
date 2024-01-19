package res

import "github.com/hajimehoshi/ebiten/v2"

type Image struct {
	*ebiten.Image
	width      float64
	height     float64
	halfWidth  float64
	halfHeight float64
}

func (i *Image) HalfWidth() float64 {
	return i.halfWidth
}

func (i *Image) HalfHeight() float64 {
	return i.halfHeight
}

func (i *Image) Width() float64 {
	return i.width
}

func (i *Image) Height() float64 {
	return i.height
}

func NewImage(image *ebiten.Image) *Image {
	w, h := image.Bounds().Size().X, image.Bounds().Size().Y
	return &Image{
		Image:      image,
		width:      float64(w),
		height:     float64(h),
		halfWidth:  float64(w) / 2,
		halfHeight: float64(h) / 2,
	}
}
