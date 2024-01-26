package game

import (
	"math"

	"github.com/KEINOS/go-noise"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kettek/dokimazo/internal/res"
)

const DayLength = 60 * 60 * 24 // 24 minutes to the day.
const DaysPerYear = 12.0       // 12 days to the year.

type Biosphere struct {
	sneed int64

	time     uint64
	timeSeed uint64

	//
	noiseGenerator noise.Generator

	daytime uint64  // 0-DayLength, representing a day/night cycle in ticks.
	day     uint64  // represents the current day, as a result of time / DayLength.
	clock   float64 // 0-24.00, representing the current time of day.
	season  float64 // 0-1.0, representing the current season.

	temperature float64 // 0-1.0, representing the current temperature.

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

	// Always start time after "winter"
	b.time = DaysPerYear / 4.0 * DayLength

	b.noiseGenerator, err = noise.New(noise.OpenSimplex, int64(b.sneed))
	if err != nil {
		panic(err)
	}

	b.RefreshSeason()

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

func (b *Biosphere) UpdateTime() {
	b.time++
	//b.time += 100
	b.RefreshSeed()
	// Convert b.time into day/night cycle, divided by
	b.daytime = b.time % DayLength
	b.day = b.time / DayLength
	daynight := float64(b.daytime) / float64(DayLength)

	// Convert daynight to 24-hour clock time.
	h, m := math.Modf(daynight * 24.0)
	b.clock = h + m*0.59

	// Refresh temp/season every 10 seconds.
	if b.time%(60*10) == 0 {
		b.RefreshSeason()
	}
}

func (b *Biosphere) RefreshSeason() {
	// Season is based on the cosine of the day of the year, with the center of the parabola being summer and wider than winter. (roughly winter lasts for 35% of the year, with temperatures decreasing/increasing faster the closer to the outside of the parabola)
	b.season = 1.0 - (math.Pow((1+math.Cos(float64(b.time)/DayLength/DaysPerYear*math.Pi*2.0)), 2) / 4.0)

	// Temperature is a bit cheaty, but simply uses the season and adds some noise to it.
	b.temperature = b.season + b.noiseGenerator.Eval64(float64(b.day))/6
	// convert b.temperature to 260-310 Kelvin
	b.temperature = b.temperature*50 + 260

	// TODO: Modify temperature by time of day.
}

func (b *Biosphere) TemperatureKelvin() float64 {
	return b.temperature
}

func (b *Biosphere) TemperatureCelsius() float64 {
	return b.temperature - 273.15
}

func (b *Biosphere) TemperatureFahrenheit() float64 {
	return (b.temperature-273.15)*1.8 + 32
}

func (b *Biosphere) Temperatures() (kelvin, celsius, fahrenheit float64) {
	return b.TemperatureKelvin(), b.TemperatureCelsius(), b.TemperatureFahrenheit()
}

func (b *Biosphere) SeasonString() string {
	switch {
	case b.season > 0.50:
		return "summer"
	case b.season > 0.30 && b.day%DaysPerYear >= DaysPerYear/2:
		return "fall"
	case b.season > 0.30:
		return "spring"
	default:
		return "winter"
	}
}

// Update updates cloud positions, etc.
func (b *Biosphere) Update() {
	b.UpdateTime()

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
