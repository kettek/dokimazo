package res

import (
	"github.com/hajimehoshi/ebiten/v2"
)

var Shaders = make(map[string]*ebiten.Shader)

func LoadShader(s string) (*ebiten.Shader, error) {
	if Shaders[s] != nil {
		return Shaders[s], nil
	}
	b, err := fs.ReadFile("shaders/" + s)
	if err != nil {
		return nil, err
	}
	Shaders[s], err = ebiten.NewShader(b)
	return Shaders[s], err
}
