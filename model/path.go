package model

type Path struct {
	Rooms  []*Room
	Length int
}

type Paths struct {
	AllPaths []Path
}

func NewPaths() *Paths {
	return &Paths{
		AllPaths: make([]Path, 0),
	}
}

func (p *Paths) AddPath(path Path) {
	p.AllPaths = append(p.AllPaths, path)
}

// Calculate how many ants can use this path in parallel
func (p *Path) CalculateCapacity() int {
	// The capacity is the number of rooms that can hold ants
	// Start and end rooms can hold multiple ants
	capacity := 0
	for _, room := range p.Rooms {
		if room.Flag == "##start" || room.Flag == "##end" {
			continue
		}
		capacity++
	}
	return capacity
}

// Sort paths by length for optimal distribution
func (p *Paths) SortByLength() {
	for i := 0; i < len(p.AllPaths)-1; i++ {
		for j := i + 1; j < len(p.AllPaths); j++ {
			if p.AllPaths[i].Length > p.AllPaths[j].Length {
				p.AllPaths[i], p.AllPaths[j] = p.AllPaths[j], p.AllPaths[i]
			}
		}
	}
}
