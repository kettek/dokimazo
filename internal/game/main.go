package game

import (
	"fmt"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	input "github.com/quasilyte/ebitengine-input"
)

type Game struct {
	camera Camera

	world World

	players []*Player

	input input.System

	//
	nightShader     *ebiten.Shader
	nightShaderOpts ebiten.DrawRectShaderOptions

	//
	drawTargets DrawTargets
	//
	lastWidth, lastHeight int
}

const (
	InputTurnLeft input.Action = iota
	InputTurnRight
	InputMoveForward
	InputMoveBackward
	InputInteract
	//
	InputRotateCameraLeft
	InputRotateCameraRight
	InputZoomCameraIn
	InputZoomCameraOut
)

func New() *Game {
	var err error
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

	g.input.Init(input.SystemConfig{
		DevicesEnabled: input.MouseDevice | input.KeyboardDevice | input.GamepadDevice,
	})
	keymap := input.Keymap{
		InputTurnLeft:     {input.KeyA, input.KeyGamepadLeft},
		InputTurnRight:    {input.KeyD, input.KeyGamepadRight},
		InputMoveForward:  {input.KeyW, input.KeyGamepadUp},
		InputMoveBackward: {input.KeyS, input.KeyGamepadDown},
		InputInteract:     {input.KeyF, input.KeyMouseLeft, input.KeyGamepadA},
	}
	p.input = g.input.NewHandler(0, keymap)
	g.camera.input = g.input.NewHandler(0, input.Keymap{
		InputRotateCameraLeft:  {input.KeyQ, input.KeyGamepadL1},
		InputRotateCameraRight: {input.KeyE, input.KeyGamepadR1},
		InputZoomCameraIn:      {input.KeyZ, input.KeyGamepadL2},
		InputZoomCameraOut:     {input.KeyX, input.KeyGamepadR2},
	})

	/*g.nightShader, err = res.LoadShader("night.kage")
	if err != nil {
		panic(err)
	}
	g.nightShaderOpts = ebiten.DrawRectShaderOptions{
		Uniforms: map[string]interface{}{
			"DayNight": 0.5,
		},
	}*/

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
	g.input.Update()
	if err := g.world.Update(); err != nil {
		return err
	}
	g.camera.Update()

	/*g.nightShaderOpts.Uniforms["Rotation"] = float32(g.camera.angle)
	g.nightShaderOpts.Uniforms["Zoom"] = g.camera.Z
	g.nightShaderOpts.Uniforms["Position"] = []float32{float32(g.camera.X()), float32(g.camera.Y())}*/

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

	// Darken the ground, drops, and world for nighttime.
	// Calculate darkness with 0.4 to 0.6 to being broad daylight. 0.0 to 0.4 and 0.6 to 1.0 is night.
	/*darkness := math.Max(0.3, 1.0-math.Abs(g.world.biosphere.daynight-0.5)*2)

	// Get visible lights.
	{
		lights := make([]float32, 16*4)
		if p := g.LocalPlayer(); p != nil {
			x := p.Position().X()
			y := p.Position().Y()
			// x, y, radius, intensity
			lights[0] = float32(x)
			lights[1] = float32(y)
			lights[2] = 100.0
			lights[3] = 1.0
			// rgba
			lights[4] = 1.0
			lights[5] = 1.0
			lights[6] = 1.0
			lights[7] = 1.0
		}
		g.nightShaderOpts.Uniforms["Lights"] = lights
	}*/

	/*g.nightShaderOpts.Images[0] = g.drawTargets.Ground
	g.nightShaderOpts.Uniforms["DayNight"] = darkness
	screen.DrawRectShader(g.drawTargets.Ground.Bounds().Dx(), g.drawTargets.Ground.Bounds().Dy(), g.nightShader, &g.nightShaderOpts)
	g.nightShaderOpts.Images[0] = g.drawTargets.Drops
	screen.DrawRectShader(g.drawTargets.Drops.Bounds().Dx(), g.drawTargets.Drops.Bounds().Dy(), g.nightShader, &g.nightShaderOpts)
	g.nightShaderOpts.Images[0] = g.drawTargets.World
	screen.DrawRectShader(g.drawTargets.World.Bounds().Dx(), g.drawTargets.World.Bounds().Dy(), g.nightShader, &g.nightShaderOpts)

	g.nightShaderOpts.Images[0] = g.drawTargets.Sky
	screen.DrawRectShader(g.drawTargets.Sky.Bounds().Dx(), g.drawTargets.Sky.Bounds().Dy(), g.nightShader, &g.nightShaderOpts)*/
	// TODO: FOR NOW we're just drawing without lighting postprocessing.
	screen.DrawImage(g.drawTargets.Ground, nil)
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
