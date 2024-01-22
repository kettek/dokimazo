package game

type Thing interface {
	Update() []Request
	HandleRequest(Request, bool)
	Chunk() *Chunk
	SetChunk(*Chunk)
}

type Things []Thing

func (t *Things) Add(thing Thing) {
	*t = append(*t, thing)
}

func (t *Things) Remove(thing Thing) {
	for i, thing2 := range *t {
		if thing2 == thing {
			*t = append((*t)[:i], (*t)[i+1:]...)
			return
		}
	}
}
