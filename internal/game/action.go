package game

// Action represents actions that can be requested by things.
type Action interface{}

// ActionMove represents a move action.
type ActionMove struct {
	Position Vec2
}

type ActionRotate struct {
	Rotation float64
}

type Request interface{}

type RequestMove struct {
	From Vec2
	To   Vec2
}

type RequestRotate struct {
	Rotation float64
}
