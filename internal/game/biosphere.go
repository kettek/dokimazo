package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kettek/dokimazo/internal/res"
)

type Biosphere struct {
	sneed int64

	time     uint64
	timeSeed uint64

	cloudShader        *ebiten.Shader
	cloudOpts          ebiten.DrawRectShaderOptions
	cloudColor         [3]float32
	cloudWindSpeed     float32
	cloudWindDirection float32
	cloudDensity       float32
	//
	fogShader  *ebiten.Shader
	fogOpts    ebiten.DrawRectShaderOptions
	fogDensity float32
	fogColor   [4]float32
	//
	camera *Camera
}

func NewBiosphere(sneed int64) *Biosphere {
	var err error
	b := &Biosphere{
		sneed:              sneed,
		cloudColor:         [3]float32{0.0, 0.0, 0.03},
		cloudWindSpeed:     3.0,
		cloudWindDirection: 3.0,
		cloudDensity:       1.5,
		//fogDensity:         0.1,
		fogColor: [4]float32{0.6, 0.6, 0.6, 0.6},
	}

	b.cloudShader, err = res.LoadShader("clouds.kage")
	if err != nil {
		panic(err)
	}
	b.cloudOpts = ebiten.DrawRectShaderOptions{
		Uniforms: map[string]interface{}{},
	}

	b.fogShader, err = res.LoadShader("fog.kage")
	if err != nil {
		panic(err)
	}
	b.fogOpts = ebiten.DrawRectShaderOptions{
		Uniforms: map[string]interface{}{},
	}
	return b
}

// RefreshSeed refreshes the current time-based seed.
func (b *Biosphere) RefreshSeed() {
	b.timeSeed = b.time ^ uint64(b.sneed)
}

// Update updates cloud positions, etc.
func (b *Biosphere) Update() {
	b.time++

	b.cloudOpts.Uniforms["Time"] = float32(b.time) / 10.0
	b.cloudOpts.Uniforms["Rotation"] = float32(b.camera.angle)
	b.cloudOpts.Uniforms["Zoom"] = b.camera.Z
	b.cloudOpts.Uniforms["Position"] = []float32{float32(b.camera.X()), float32(b.camera.Y())}
	b.cloudOpts.Uniforms["Color"] = b.cloudColor
	b.cloudOpts.Uniforms["Wind"] = b.cloudWindSpeed
	b.cloudOpts.Uniforms["WindDirection"] = b.cloudWindDirection
	b.cloudOpts.Uniforms["Density"] = b.cloudDensity

	b.fogOpts.Uniforms["Time"] = float32(b.time) / 10.0
	b.fogOpts.Uniforms["Rotation"] = float32(b.camera.angle)
	b.fogOpts.Uniforms["Zoom"] = b.camera.Z
	b.fogOpts.Uniforms["Position"] = []float32{float32(b.camera.X()), float32(b.camera.Y())}
	b.fogOpts.Uniforms["Color"] = b.fogColor
	b.fogOpts.Uniforms["Density"] = b.fogDensity
}
