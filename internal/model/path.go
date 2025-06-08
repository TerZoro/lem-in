package model

import "sort"

type Path struct {
	Rooms  []*Room
	Length int
}

type Paths struct {
	AllPaths            []*Path
	OptimalDistribution []int // Stores how many ants should use each path
}

func NewPaths() *Paths {
	return &Paths{
		AllPaths: make([]*Path, 0),
	}
}

func (p *Paths) AddPath(path Path) {
	p.AllPaths = append(p.AllPaths, &path)
}

// Calculate how many ants can use this path in parallel
func (p *Path) CalculateCapacity() int {
	// The capacity is the number of rooms that can hold ants
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
	sort.Slice(p.AllPaths, func(i, j int) bool {
		return p.AllPaths[i].Length < p.AllPaths[j].Length
	})
}
