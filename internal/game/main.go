package game

import (
	"fmt"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
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

	g.world.biosphere.camera = &g.camera

	return g
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	if g.lastWidth != outsideWidth || g.lastHeight != outsideHeight {
		g.lastWidth, g.lastHeight = outsideWidth, outsideHeight
		g.drawTargets.Ground = ebiten.NewImage(outsideWidth, outsideHeight)
		g.drawTargets.Shadow = ebiten.NewImage(outsideWidth, outsideHeight)
		g.drawTargets.World = ebiten.NewImage(outsideWidth, outsideHeight)
		g.drawTargets.Drops = ebiten.NewImage(outsideWidth, outsideHeight)
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
	distanceX := math.Max(1, float64(screen.Bounds().Dx()/ChunkPixelSize/ChunkTileSize)/g.camera.Z)
	distanceY := math.Max(1, float64(screen.Bounds().Dy()/ChunkPixelSize/ChunkTileSize)/g.camera.Z)
	chunks := g.world.ChunksAround(px, py, int(distanceX), int(distanceY))

	// Collect visuals to render.
	var lowVisuals, medVisuals, dropVisuals, skyVisuals Visuals
	for _, chunk := range chunks {
		lowVisuals = append(lowVisuals, chunk.lowVisuals...)
		medVisuals = append(medVisuals, chunk.medVisuals...)
		dropVisuals = append(dropVisuals, chunk.dropVisuals...)
		skyVisuals = append(skyVisuals, chunk.highVisuals...)
	}

	g.drawTargets.Ground.Clear()
	g.drawTargets.Shadow.Clear()
	g.drawTargets.World.Clear()
	g.drawTargets.Drops.Clear()
	g.drawTargets.Sky.Clear()

	drawOpts := CameraDrawOptions{
		ShadowAngle: g.world.biosphere.daynight * math.Pi * 2,
		// TODO: ShadowAngle... need just a simple daytime float, or extend daynight to be -1 to 1.
	}

	g.camera.Draw(g.drawTargets.Ground, lowVisuals, drawOpts)
	//g.camera.Draw(g.drawTargets.Shadow, medVisuals, CameraDrawOptions{Shadows: true, HideVisuals: true})
	drawOpts.Shadows = true
	g.camera.Draw(g.drawTargets.World, medVisuals, drawOpts)
	g.camera.Draw(g.drawTargets.Drops, dropVisuals, drawOpts)
	drawOpts.Shadows = false
	g.camera.Draw(g.drawTargets.Sky, skyVisuals, drawOpts)

	if g.world.biosphere.fogDensity > 0 {
		w, h := g.drawTargets.World.Bounds().Dx(), g.drawTargets.World.Bounds().Dy()
		g.drawTargets.World.DrawRectShader(w, h, g.world.biosphere.fogShader, &g.world.biosphere.fogOpts)
	}

	if g.world.biosphere.cloudDensity > 0 {
		w, h := g.drawTargets.Sky.Bounds().Dx(), g.drawTargets.Sky.Bounds().Dy()
		g.drawTargets.Sky.DrawRectShader(w, h, g.world.biosphere.cloudShader, &g.world.biosphere.cloudOpts)
	}

	screen.DrawImage(g.drawTargets.Ground, nil)
	//screen.DrawImage(g.drawTargets.Shadow, nil)
	screen.DrawImage(g.drawTargets.Drops, nil)
	screen.DrawImage(g.drawTargets.World, nil)
	screen.DrawImage(g.drawTargets.Sky, nil)

	h, m := math.Modf(g.world.biosphere.clock)
	m *= 100
	k, c, f := g.world.biosphere.Temperatures()
	ebitenutil.DebugPrint(screen, fmt.Sprintf("day %d, %02.0f:%02.0f\nseason %.2f (%s)\ntemp %.2fK | %.2fC | %.2fF\n%.2fm elevation\n%.2f aridity", g.world.biosphere.day, h, m, g.world.biosphere.season, g.world.biosphere.SeasonString(), k, c, f, g.world.biosphere.ElevationAt(g.LocalPlayer().Position())*5000, g.world.biosphere.AridityAt(g.LocalPlayer().Position())))
}

func (g *Game) LocalPlayer() *Player {
	return g.players[0]
}
