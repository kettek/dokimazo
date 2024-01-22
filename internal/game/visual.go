package game

type Visual interface {
	Draw(ops DrawOpts)
	Position() Vec2
	Size() Vec2
	Angle() float64
	Z() float64
	SetZ(float64)
}

type VisualShadow interface {
	DrawShadow(ops DrawOpts)
}

type Visuals []Visual

func (v *Visuals) Add(visual Visual) {
	*v = append(*v, visual)
}

func (v *Visuals) Remove(visual Visual) {
	for i, visual2 := range *v {
		if visual2 == visual {
			*v = append((*v)[:i], (*v)[i+1:]...)
			return
		}
	}
}
