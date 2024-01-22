package game

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	camera Camera

	world World

	players []*Player

	//
	drawTargets DrawTargets
	//
	lastWidth, lastHeight int
}

func New() *Game {
	g := &Game{
		camera: *NewCamera(),
		world:  *NewWorld(),
	}

	chunk := g.world.LoadChunk(0, 0)

	p := NewPlayer()

	px, py := float64(chunk.X)*ChunkPixelSize*ChunkTileSize, float64(chunk.Y)*ChunkPixelSize*ChunkTileSize
	p.Assign(Vec2{px, py})
	chunk.AddThing(p, VisualLayerWorld)

	g.players = append(g.players, p)

	g.camera.Target = p

	return g
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	if g.lastWidth != outsideWidth || g.lastHeight != outsideHeight {
		g.lastWidth, g.lastHeight = outsideWidth, outsideHeight
		g.drawTargets.Ground = ebiten.NewImage(outsideWidth, outsideHeight)
		g.drawTargets.Shadow = ebiten.NewImage(outsideWidth, outsideHeight)
		g.drawTargets.World = ebiten.NewImage(outsideWidth, outsideHeight)
		g.drawTargets.Sky = ebiten.NewImage(outsideWidth, outsideHeight)
		g.camera.Layout(outsideWidth, outsideHeight)
	}
	return outsideWidth, outsideHeight
}

func (g *Game) Update() error {
	if err := g.world.Update(); err != nil {
		return err
	}
	g.camera.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	px, py := g.LocalPlayer().Chunk().X, g.LocalPlayer().Chunk().Y
	chunks := g.world.ChunksAround(px, py)

	// Collect visuals to render.
	var lowVisuals, medVisuals, skyVisuals Visuals
	for _, chunk := range chunks {
		lowVisuals = append(lowVisuals, chunk.lowVisuals...)
		medVisuals = append(medVisuals, chunk.medVisuals...)
		skyVisuals = append(skyVisuals, chunk.highVisuals...)
	}

	g.drawTargets.Ground.Clear()
	g.drawTargets.Shadow.Clear()
	g.drawTargets.World.Clear()
	g.drawTargets.Sky.Clear()

	g.camera.Draw(g.drawTargets.Ground, lowVisuals, CameraDrawOptions{})
	//g.camera.Draw(g.drawTargets.Shadow, medVisuals, CameraDrawOptions{Shadows: true, HideVisuals: true})
	g.camera.Draw(g.drawTargets.World, medVisuals, CameraDrawOptions{Shadows: true})
	g.camera.Draw(g.drawTargets.Sky, skyVisuals, CameraDrawOptions{})

	screen.DrawImage(g.drawTargets.Ground, nil)
	//screen.DrawImage(g.drawTargets.Shadow, nil)
	screen.DrawImage(g.drawTargets.World, nil)
	screen.DrawImage(g.drawTargets.Sky, nil)
}

func (g *Game) LocalPlayer() *Player {
	return g.players[0]
}
