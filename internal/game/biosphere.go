package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kettek/dokimazo/internal/res"
)

type Biosphere struct {
	time               float64
	cloudShader        *ebiten.Shader
	cloudOpts          ebiten.DrawRectShaderOptions
	cloudColor         [3]float32
	cloudWindSpeed     float32
	cloudWindDirection float32
	cloudDensity       float32
	//
	fogShader *ebiten.Shader
	fogOpts   ebiten.DrawRectShaderOptions
	camera    *Camera
}

func NewBiosphere() *Biosphere {
	var err error
	b := &Biosphere{
		cloudColor:         [3]float32{0.0, 0.0, 0.03},
		cloudWindSpeed:     3.0,
		cloudWindDirection: 3.0,
		cloudDensity:       0.3,
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

func (b *Biosphere) Update() {
	b.time += 0.1
	b.cloudOpts.Uniforms["Time"] = float32(b.time)
	b.cloudOpts.Uniforms["Rotation"] = float32(b.camera.angle)
	b.cloudOpts.Uniforms["Zoom"] = b.camera.Z
	b.cloudOpts.Uniforms["Position"] = []float32{float32(b.camera.X()), float32(b.camera.Y())}
	b.cloudOpts.Uniforms["Color"] = b.cloudColor
	b.cloudOpts.Uniforms["Wind"] = b.cloudWindSpeed
	b.cloudOpts.Uniforms["WindDirection"] = b.cloudWindDirection
	b.cloudOpts.Uniforms["Density"] = b.cloudDensity
}
