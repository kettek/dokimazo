package game

type Carrier interface {
	AddDrop(drop *Drop)
	RemoveDrop(drop *Drop)
	Drops() Drops
}

type Inventory struct {
	drops Drops
}

func (i *Inventory) Drops() Drops {
	return i.drops
}

func (i *Inventory) AddDrop(drop *Drop) {
	i.drops.Add(drop)
}

func (i *Inventory) RemoveDrop(drop *Drop) {
	i.drops.Remove(drop)
}

func (i *Inventory) CloneInventory() Inventory {
	clone := Inventory{}
	for _, drop := range i.drops {
		clone.drops.Add(drop)
	}
	return clone
}
