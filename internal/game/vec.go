package game

import "math"

type Vec2 [2]float64

func (v Vec2) X() float64 {
	return v[0]
}

func (v Vec2) Y() float64 {
	return v[1]
}

func (v *Vec2) Assign(v2 Vec2) *Vec2 {
	v[0] = v2[0]
	v[1] = v2[1]
	return v
}

func (v *Vec2) Add(v2 Vec2) *Vec2 {
	v[0] += v2[0]
	v[1] += v2[1]
	return v
}

func (v *Vec2) Sub(v2 Vec2) *Vec2 {
	v[0] -= v2[0]
	v[1] -= v2[1]
	return v
}

func (v *Vec2) Mul(v2 Vec2) *Vec2 {
	v[0] *= v2[0]
	v[1] *= v2[1]
	return v
}

func (v *Vec2) Rotate(angle float64) *Vec2 {
	sin, cos := math.Sincos(angle)
	x := v[0]*cos - v[1]*sin
	y := v[0]*sin + v[1]*cos
	v[0] = x
	v[1] = y
	return v
}

func (v *Vec2) RotateAround(origin Vec2, angle float64) *Vec2 {
	return v.Sub(origin).Rotate(angle).Add(origin)
}

func (v *Vec2) Distance(v2 Vec2) float64 {
	return math.Sqrt(math.Pow(v[0]-v2[0], 2) + math.Pow(v[1]-v2[1], 2))
}

func (v *Vec2) AngleTo(v2 Vec2) float64 {
	return math.Atan2(v2[1]-v[1], v2[0]-v[0])
}

func (v Vec2) Clone() Vec2 {
	return Vec2{v[0], v[1]}
}

type RVec2 struct {
	Vec2
	angle float64
}

func (v RVec2) Angle() float64 {
	return v.angle
}

func (v *RVec2) Rotate(angle float64) *RVec2 {
	v.angle += angle
	return v
}

func (v *RVec2) RotateAround(origin Vec2, angle float64) *RVec2 {
	v.Vec2.RotateAround(origin, v.angle+angle)
	v.angle += angle
	return v
}

func (v *RVec2) Forward() Vec2 {
	return Vec2{math.Cos(v.angle), math.Sin(v.angle)}
}
