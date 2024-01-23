package res

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

type Drop struct {
	Data   DropData
	Visual DropVisual
}

func (d *Drop) Sheet() *Sheet {
	return d.Visual.sheet
}

type DropData struct {
	ID          RID
	Name        string
	Description string
}

type DropVisual struct {
	Image         string
	Width         int
	Height        int
	LayerDistance float64 `yaml:"layer_distance"`
	sheet         *Sheet
}

var Drops = make(map[RID]*Drop)

func init() {
	entries, err := fs.ReadDir("drops")
	if err != nil {
		panic(fmt.Errorf("loading drops: %w", err))
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		b, err := fs.ReadFile("drops/" + entry.Name())
		if err != nil {
			panic(fmt.Errorf("could not read drop: %w", err))
		}
		var drop Drop
		err = yaml.Unmarshal(b, &drop)
		if err != nil {
			panic(fmt.Errorf("could not unmarshal drop data: %w", err))
		}
		if drop.Visual.Image != "" {
			w := drop.Visual.Width
			h := drop.Visual.Height
			if w == 0 {
				w = SheetCellWidth
			}
			if h == 0 {
				h = SheetCellHeight
			}
			sheet, err := LoadSheetWithSize(drop.Visual.Image+".png", w, h)
			if err != nil {
				panic(fmt.Errorf("could not load drop sheet: %w", err))
			}
			drop.Visual.sheet = sheet
		}
		Drops[drop.Data.ID] = &drop
	}
}
