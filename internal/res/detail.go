package res

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

type Detail struct {
	Data     DetailData
	Visual   DetailVisual
	Behavior DetailBehavior
}

func (d *Detail) Sheet() *Sheet {
	return d.Visual.sheet
}

type DetailData struct {
	ID          RID
	Name        string
	Description string
}

type DetailVisual struct {
	Image          string
	Width          int
	Height         int
	Low            bool
	LayerDistance  float64 `yaml:"layer_distance"`
	RandomRotation bool    `yaml:"random_rotation"`
	sheet          *Sheet
}

type DetailBehavior struct {
	Blocks bool
	Hurts  bool
}

var Details = make(map[RID]*Detail)

func init() {
	entries, err := fs.ReadDir("details")
	if err != nil {
		panic(fmt.Errorf("loading details: %w", err))
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		b, err := fs.ReadFile("details/" + entry.Name())
		if err != nil {
			panic(fmt.Errorf("could not read file: %w", err))
		}
		var detail Detail
		err = yaml.Unmarshal(b, &detail)
		if err != nil {
			panic(fmt.Errorf("could not unmarshal detail data: %w", err))
		}
		if detail.Visual.Image != "" {
			w := detail.Visual.Width
			h := detail.Visual.Height
			if w == 0 {
				w = SheetCellWidth
			}
			if h == 0 {
				h = SheetCellHeight
			}
			sheet, err := LoadSheetWithSize(detail.Visual.Image+".png", w, h)
			if err != nil {
				panic(fmt.Errorf("could not load detail sheet: %w", err))
			}
			detail.Visual.sheet = sheet
		}
		Details[detail.Data.ID] = &detail
	}
}
