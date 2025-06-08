package main

import (
	"lemin/farm"
	"lemin/internal/drawing"
	"lemin/internal/pathfinding"
	"lemin/internal/simulation"
	errhandle "lemin/pkg/errors"
	"log"
)

func main() {
	filename, err := simulation.ReadArgs()
	errhandle.HandleError(err)

	lines, err := simulation.FileReader(filename)
	errhandle.HandleError(err)

	if err := simulation.CheckStartEnd(lines); err != nil {
		log.Fatal(err)
	}

	// Parse ants count
	ants, err := farm.AntsNum(lines)
	if err != nil {
		log.Fatal(err)
	}

	rooms, err := farm.ParseRooms(lines)
	if err != nil {
		log.Printf("error: Could not parse rooms")
	}

	paths, err := pathfinding.PathStarter(rooms, ants)
	if err != nil {
		log.Fatal(err)
	}

	// Simulate ant movement
	simulation.SimulateAnts(paths, ants, rooms, lines)

	// Draw the graph visualization
	if err := drawing.DrawGraph(rooms); err != nil {
		log.Printf("Warning: Could not draw graph: %v", err)
	}
}
