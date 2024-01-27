package game

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/kettek/dokimazo/internal/res"
)

// TileState is a bitfield of potential tile states.
type TileState uint8

// TileState constants, wow.
const (
	TileStateSolid TileState = 1 << iota
	TileStateLiquid
	TileStateHurts
	TileStatePoisons
)

// Tile represents a single tile within a chunk.
type Tile struct {
	State TileState
	// Some sorta image or ID reference.
	Details []TileDetail // loaded details.
}

type Tiles [][]Tile

func (t Tiles) At(x, y int) *Tile {
	if x < 0 || y < 0 || x >= len(t) || y >= len((t)[x]) {
		return nil
	}
	return &(t)[x][y]
}

type TileDetail struct {
	RID    res.RID
	State  interface{} // I guess for now
	detail *res.Detail
	visual Visual
}

type World struct {
	sneed         int64
	loadingChunks []*Chunk
	loadedChunks  []*Chunk
	biosphere     *Biosphere
}

func NewWorld() *World {
	sneed := rand.Int63()
	return &World{
		sneed:     sneed,
		biosphere: NewBiosphere(sneed),
	}
}

func (w *World) Update() error {
	// Check on our loading chunks.
	chunks := w.loadingChunks[:0]
	for _, chunk := range w.loadingChunks {
		select {
		case err := <-chunk.loadChan:
			if err != nil {
				panic(fmt.Errorf("failed to load chunk: %w", err))
			}
			//fmt.Println("loaded chunk", chunk.X, chunk.Y)
			chunk.loaded = true
		default:
		}
		if chunk.loaded {
			w.loadedChunks = append(w.loadedChunks, chunk)
		} else {
			chunks = append(chunks, chunk)
		}
	}
	w.loadingChunks = chunks

	// Update our loaded chunks.
	var chunkUpdateRequests []ChunkUpdateRequests
	chunks = w.loadedChunks[:0]
	for _, chunk := range w.loadedChunks {
		chunkUpdateRequests = append(chunkUpdateRequests, ChunkUpdateRequests{
			Chunk:    chunk,
			Requests: chunk.Update(w),
		})
		// TODO: Check if chunk should be unloaded (probably distance from player).
		chunks = append(chunks, chunk)
	}
	w.loadedChunks = chunks

	// Process chunk updates.
	for _, chunkUpdate := range chunkUpdateRequests {
		for _, chunkRequest := range chunkUpdate.Requests {
			switch chunkRequest := chunkRequest.(type) {
			case ChunkUpdateThingRequest:
				for _, thingRequest := range chunkRequest.Requests {
					switch thingRequest := thingRequest.(type) {
					case RequestMove:
						targetChunk := chunkUpdate.Chunk
						cx, cy := int(math.Floor(thingRequest.To.X()/ChunkPixelSize/ChunkTileSize)), int(math.Floor(thingRequest.To.Y()/ChunkPixelSize/ChunkTileSize))
						if cx != chunkRequest.Thing.Chunk().X || cy != chunkRequest.Thing.Chunk().Y {
							targetChunk = w.LoadChunk(cx, cy)
						}
						// Check the tile the thing wants to move to.
						// Get the tile x,y from the request position.
						tx, ty := int(math.Floor(thingRequest.To.X()/ChunkPixelSize))-cx*ChunkPixelSize, int(math.Floor(thingRequest.To.Y()/ChunkPixelSize))-cy*ChunkPixelSize
						if tile := targetChunk.Tiles.At(tx, ty); tile != nil {
							if tile.State&TileStateSolid == 0 {
								if targetChunk != chunkUpdate.Chunk {
									chunkUpdate.Chunk.RemoveThing(chunkRequest.Thing)
									targetChunk.AddThing(chunkRequest.Thing, VisualLayerWorld)
								}
							} else {
								tvec := Vec2{
									float64(tx+cx*ChunkPixelSize) * ChunkPixelSize,
									float64(ty+cy*ChunkPixelSize) * ChunkPixelSize,
								}
								dvec := tvec.Clone()
								// Get the distance from the tile to the thing.
								dvec.Sub(thingRequest.To)
								// Get the distance to move to the tile.
								var moveX, moveY float64
								if math.Abs(dvec.X()) > math.Abs(dvec.Y()) {
									moveX = dvec.X() / math.Abs(dvec.X())
								} else {
									moveY = dvec.Y() / math.Abs(dvec.Y())
								}
								// Set the new position.
								thingRequest.To.Assign(Vec2{thingRequest.To.X() + moveX, thingRequest.To.Y() + moveY})
							}
						}
						chunkRequest.Thing.HandleRequest(thingRequest, true)

						// If the thing is a carrier, check for drops in range and have them move towards the thing and/or be picked up.
						if carrier, ok := chunkRequest.Thing.(Carrier); ok {
							chunks := w.ChunksInRange(thingRequest.To.X(), thingRequest.To.Y(), 20.0)
							// Check if any drops are in range.
							for _, chunk := range chunks {
								for _, drop := range chunk.Drops {
									// Gradually move the drop towards the player.
									dx := thingRequest.To.X() - drop.Position().X()
									dy := thingRequest.To.Y() - drop.Position().Y()
									if math.Abs(dx) < 20.0 && math.Abs(dy) < 20.0 {
										if math.Abs(dx) < 5.0 && math.Abs(dy) < 5.0 {
											chunk.RemoveDrop(drop)
											carrier.AddDrop(drop)
										} else {
											if dx < 0 {
												dx = -3
											} else if dx > 0 {
												dx = 3
											}
											if dy < 0 {
												dy = -3
											} else if dy > 0 {
												dy = 3
											}
											drop.Add(Vec2{dx, dy})
										}
									}
								}
							}
						}
					case RequestRotate:
						chunkRequest.Thing.HandleRequest(thingRequest, true)
					}
					/*case ChunkUpdateMoveThing:
					switch thing := result.Thing.(type) {
					case *Mover:
						if result.To.X() < 0 || result.To.Y() < 0 || result.To.X() >= chunkUpdate.Chunk.Width() || result.To.Y() >= chunkUpdate.Chunk.Height() {
							x := 0
							y := 0
							if result.To.X() < 0 {
								x = -1
							} else if result.To.X() >= chunkUpdate.Chunk.Width() {
								x = 1
							}
							if result.To.Y() < 0 {
								y = -1
							} else if result.To.Y() >= chunkUpdate.Chunk.Height() {
								y = 1
							}
							// Get/load target chunk.
							chunk := w.LoadChunk(thing.Chunk().X+x, thing.Chunk().Y+y)

							// Calculate position in new chunk and assign.
							var rx, ry float64
							if x < 0 {
								rx = chunk.Width() - 1
							} else if x > 0 {
								rx = 0
							}
							if y < 0 {
								ry = chunk.Height() - 1
							} else if y > 0 {
								ry = 0
							}
							rx += result.To.X()
							ry += result.To.Y()

							// TODO: See if rx, ry is open in target chunk.
							if true {
								// Remove thing from current chunk and add it to target chunk.
								chunkUpdate.Chunk.Things.Remove(thing)
								chunk.Things.Add(thing)
								thing.SetChunk(chunk)
								thing.Vec2.Assign(Vec2{rx, ry})
							}
						}*/
				}
			}
		}
	}

	// Update the biosphere.
	w.biosphere.Update()

	return nil
}

func (w *World) ChunksAround(x, y, dx, dy int) (chunks []*Chunk) {
	for x1 := x - dx; x1 <= x+dx; x1++ {
		for y1 := y - dy; y1 <= y+dy; y1++ {
			chunks = append(chunks, w.LoadChunk(x1, y1))
		}
	}
	return chunks
}

func (w *World) ChunksInRange(x, y float64, distance float64) (chunks []*Chunk) {
	cx := (x / ChunkPixelSize / ChunkTileSize)
	cy := (y / ChunkPixelSize / ChunkTileSize)
	//cx := x
	//cy := y
	dx := cx - math.Floor(cx)
	dy := cy - math.Floor(cy)
	distance = distance / ChunkPixelSize / ChunkTileSize

	chunks = append(chunks, w.LoadChunk(int(math.Floor(cx)), int(math.Floor(cy))))

	if dx-distance < 0 {
		chunks = append(chunks, w.LoadChunk(int(math.Floor(cx-1)), int(math.Floor(cy))))
	}
	if dx+distance > 1 {
		chunks = append(chunks, w.LoadChunk(int(math.Floor(cx+1)), int(math.Floor(cy))))
	}
	if dy-distance < 0 {
		chunks = append(chunks, w.LoadChunk(int(math.Floor(cx)), int(math.Floor(cy-1))))
	}
	if dy+distance > 1 {
		chunks = append(chunks, w.LoadChunk(int(math.Floor(cx)), int(math.Floor(cy+1))))
	}
	if dx-distance < 0 && dy-distance < 0 {
		chunks = append(chunks, w.LoadChunk(int(math.Floor(cx-1)), int(math.Floor(cy-1))))
	}
	if dx+distance > 1 && dy-distance < 0 {
		chunks = append(chunks, w.LoadChunk(int(math.Floor(cx+1)), int(math.Floor(cy-1))))
	}
	if dx-distance < 0 && dy+distance > 1 {
		chunks = append(chunks, w.LoadChunk(int(math.Floor(cx-1)), int(math.Floor(cy+1))))
	}
	if dx+distance > 1 && dy+distance > 1 {
		chunks = append(chunks, w.LoadChunk(int(math.Floor(cx+1)), int(math.Floor(cy+1))))
	}

	return
}

func (w *World) LoadChunk(x, y int) *Chunk {
	for _, chunk := range w.loadingChunks {
		if chunk.X == x && chunk.Y == y {
			return chunk
		}
	}
	for _, chunk := range w.loadedChunks {
		if chunk.X == x && chunk.Y == y {
			return chunk
		}
	}
	chunk := NewChunk()
	chunk.X = x
	chunk.Y = y
	w.loadingChunks = append(w.loadingChunks, chunk)
	go chunk.Load(w.biosphere)
	return chunk
}
