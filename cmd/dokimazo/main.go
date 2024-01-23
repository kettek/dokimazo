package main

import (
	"bytes"
	_ "embed"
	"image"
	_ "image/png"

	"github.com/kettek/dokimazo/internal/game"

	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed icon.png
var iconBytes []byte

func main() {
	g := game.New()

	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowTitle("dokimazo")

	icon, _, _ := image.Decode(bytes.NewReader(iconBytes))
	ebiten.SetWindowIcon([]image.Image{icon})

	if err := ebiten.RunGame(g); err != nil {
		panic(err)
	}
}
