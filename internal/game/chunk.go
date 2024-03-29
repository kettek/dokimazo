package game

import (
	"fmt"
	"math"
	"math/rand"
	"os"

	"github.com/kettek/dokimazo/internal/res"
)

type ChunkUpdateRequests struct {
	Chunk    *Chunk
	Requests []ChunkUpdateRequest
}

type ChunkUpdateRequest interface{}

type ChunkUpdateThingRequest struct {
	Thing    Thing
	Requests []Request
}

type ChunkUpdateDropRequest struct {
	Drop     *Drop
	Requests []Request
}

const ChunkPixelSize = 16.0
const ChunkTileSize = 16.0

type VisualLayer uint8

const (
	VisualLayerGround VisualLayer = iota
	VisualLayerDrops
	VisualLayerWorld
	VisualLayerSky
)

// Chunk represents a tile area equal to ChunkTileSize*ChunkTileSize.
type Chunk struct {
	loaded      bool
	loadChan    chan error
	X, Y        int
	Drops       Drops
	Tiles       Tiles
	Things      Things
	lowVisuals  Visuals // Low visuals generally is the ground.
	dropVisuals Visuals // Drop visuals are dropped items.
	medVisuals  Visuals // Med visuals are most game objects.
	highVisuals Visuals // High visuals are a mystery.
}

// NewChunk creates a new chunk.
func NewChunk() *Chunk {
	return &Chunk{
		loadChan: make(chan error),
	}
}

// Width returns the width of the chunk in pixels. This is equal to ChunkPixelSize*ChunkTileSize.
func (c *Chunk) Width() float64 {
	return ChunkPixelSize * ChunkTileSize
}

// Height is like Width.
func (c *Chunk) Height() float64 {
	return ChunkPixelSize * ChunkTileSize
}

func (c *Chunk) AddThing(thing Thing, v VisualLayer) {
	c.Things.Add(thing)
	thing.SetChunk(c)
	if thing, ok := thing.(Visual); ok {
		switch v {
		case VisualLayerGround:
			c.lowVisuals.Add(thing)
		case VisualLayerDrops:
			c.dropVisuals.Add(thing)
		case VisualLayerWorld:
			c.medVisuals.Add(thing)
		case VisualLayerSky:
			c.highVisuals.Add(thing)
		}
	}
}

func (c *Chunk) RemoveThing(thing Thing) {
	c.Things.Remove(thing)
	thing.SetChunk(nil)
	if thing, ok := thing.(Visual); ok {
		c.lowVisuals.Remove(thing)
		c.dropVisuals.Remove(thing)
		c.medVisuals.Remove(thing)
		c.highVisuals.Remove(thing)
	}
}

func (c *Chunk) AddDrop(drop *Drop) {
	c.Drops.Add(drop)
	c.dropVisuals.Add(drop)
}

func (c *Chunk) RemoveDrop(drop *Drop) {
	c.Drops.Remove(drop)
	c.dropVisuals.Remove(drop)
}

// Update updates the chunk, calling Update on all contained things.
func (c *Chunk) Update(w *World) (requests []ChunkUpdateRequest) {
	// Update them thangs.
	for _, thing := range c.Things {
		thingRequests := thing.Update()
		if len(thingRequests) > 0 {
			requests = append(requests, ChunkUpdateThingRequest{
				Thing:    thing,
				Requests: thingRequests,
			})
		}
	}
	// Update them drops.
	for _, drop := range c.Drops {
		dropRequests := drop.Update()
		if len(dropRequests) > 0 {
			requests = append(requests, ChunkUpdateDropRequest{
				Drop:     drop,
				Requests: dropRequests,
			})
		}
	}
	return requests
}

func (c *Chunk) Draw(drawTargets DrawTargets, camera *Camera) {
	x, y := float64(c.X*ChunkTileSize*ChunkPixelSize), float64(c.Y*ChunkTileSize*ChunkPixelSize)
	//x, y = 0, 0
	//w, h := c.Width(), c.Height()

	// Draw ground.
	camera.Draw(drawTargets.Ground, c.lowVisuals, CameraDrawOptions{
		XOffset: x,
		YOffset: y,
	})
	// Draw world.
	camera.Draw(drawTargets.World, c.medVisuals, CameraDrawOptions{
		XOffset: x,
		YOffset: y,
		Shadows: true,
	})
	// Draw drops.
	camera.Draw(drawTargets.Drops, c.dropVisuals, CameraDrawOptions{
		XOffset: x,
		YOffset: y,
	})
	// Draw shadows.
	/*camera.Draw(drawTargets.World, c.medVisuals, CameraDrawOptions{
		XOffset:     x,
		YOffset:     y,
		Shadows:     true,
		HideVisuals: true,
	})*/

	// Draw sky stuff.
	camera.Draw(drawTargets.Sky, c.highVisuals, CameraDrawOptions{
		XOffset: x,
		YOffset: y,
	})
}

// process loads details from RID and loads accordingly.
func (c *Chunk) process() {
	// NOTE: We store x,y as absolute world coordinates since it's easier to work with during visuals collection + rendering.
	x, y := float64(c.X*ChunkTileSize*ChunkPixelSize), float64(c.Y*ChunkTileSize*ChunkPixelSize)
	for i := range c.Tiles {
		for j := range c.Tiles[i] {
			tile := &c.Tiles[i][j]
			for k := range tile.Details {
				td := &tile.Details[k]
				td.detail = res.Details[td.RID]
				sprite := NewSpriteStackFromSheet(td.detail.Sheet())
				sprite.Assign(Vec2{x + float64(i*ChunkPixelSize), y + float64(j*ChunkPixelSize)})
				if td.detail.Visual.LayerDistance != 0 {
					sprite.LayerDistance = td.detail.Visual.LayerDistance
				}
				if td.detail.Visual.RandomRotation {
					sprite.angle = rand.Float64() * math.Pi * 2
				}
				td.visual = sprite
				if td.detail.Visual.Low {
					c.lowVisuals.Add(td.visual)
				} else {
					c.medVisuals.Add(td.visual)
				}
				if td.detail.Behavior.Blocks {
					tile.State |= TileStateSolid
				}
			}
		}
	}
	for _, cd := range c.Drops {
		cd.drop = res.Drops[cd.RID]
		sprite := NewSpriteStackFromSheet(cd.drop.Sheet())
		sprite.Assign(Vec2{cd.x, cd.y})
		if cd.drop.Visual.LayerDistance != 0 {
			sprite.LayerDistance = cd.drop.Visual.LayerDistance
		}
		cd.SpriteStack = sprite
		c.dropVisuals.Add(cd)
	}
}

// Load either loads the chunk from disk or generates a new chunk. It is intended to be run as a goroutine within World.
func (c *Chunk) Load(b *Biosphere) {
	var chunkHash uint64
	chunkHash = uint64(c.X)
	chunkHash ^= uint64(c.Y) << 32
	chunkHash ^= uint64(c.Y) >> 32
	chunkSeed := chunkHash ^ uint64(b.sneed)

	randy := rand.New(rand.NewSource(int64(chunkSeed)))
	for done := false; !done; {
		var err error
		chunkPath := fmt.Sprintf("chunks/%d/%d", c.X, c.Y)
		// Check if a file at chunkPath exists.
		if _, err = os.Stat(chunkPath); err != nil {
			if os.IsNotExist(err) {
				// TODO: Create the file?
				c.Tiles = make(Tiles, 16)
				for i := range c.Tiles {
					c.Tiles[i] = make([]Tile, 16)
					for j := range c.Tiles[i] {
						t := &c.Tiles[i][j]
						px := (float64(c.X*ChunkTileSize) + float64(i))
						py := (float64(c.Y*ChunkTileSize) + float64(j))

						r := b.noiseGenerator.Eval64(px/60, py/60, 0)
						e := b.ElevationAt(Vec2{px * ChunkPixelSize, py * ChunkPixelSize})
						a := b.AridityAt(Vec2{px * ChunkPixelSize, py * ChunkPixelSize})

						// Randomly throw some pierogies in there.
						if randy.Intn(100) < 5 {
							d := &Drop{}
							rid, _ := res.RIDFromString("edibles:pierogi")
							d.x = px*ChunkPixelSize + ChunkPixelSize/2
							d.y = py*ChunkPixelSize + ChunkPixelSize/2
							d.RID = rid
							c.Drops = append(c.Drops, d)
						}

						if randy.Intn(100) < 5 {
							//if i == 0 || j == 0 || i == 15 || j == 15 {
							rid, _ := res.RIDFromString("tree:taint")
							t.Details = append(t.Details, TileDetail{
								RID: rid,
							})
						}
						var rid res.RID
						r += e
						//fmt.Println(r, e, a)
						if (r < -0.2 || e < -0.01) && a < 0.3 {
							rid, _ = res.RIDFromString("liquid:water")
						} else {
							//r -= a
							if r < -0.1 && a > 0.5 {
								rid, _ = res.RIDFromString("ground:taintsand")
							} else if r < 0.2 {
								rid, _ = res.RIDFromString("ground:taintsand")
							} else if r < 0.8 {
								rid, _ = res.RIDFromString("ground:taintground")
							} else if e > 0.5 {
								rid, _ = res.RIDFromString("ground:taintsand")
							} else {
								rid, _ = res.RIDFromString("ground:taintground")
							}
						}
						t.Details = append(t.Details, TileDetail{
							RID: rid,
						})
					}
				}
				err = nil
			}
		}
		if err == nil {
			c.process()
		}

		done = true
		c.loadChan <- err
	}
}
