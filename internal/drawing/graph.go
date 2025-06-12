package drawing

import (
	"fmt"
	"lemin/internal/model"
	errhandle "lemin/pkg/errors"
	"math"
	"strings"
)

const (
	scale   = 4 // Scales coordinates to make the graph bigger.
	padding = 3 // Adds padding around the graph.
	offset  = 2 // The offset used to draw tunnels around obstacles.
)

func DrawGraph(rooms *model.Rooms) error {
	if rooms == nil || len(rooms.List) == 0 {
		return errhandle.FormatError(errhandle.ErrDrawGraph)
	}

	fmt.Println("\nGraph Visualization:")
	minX, maxX, minY, maxY := findBounds(rooms)

	// Create a grid large enough to hold the scaled graph and offsets.
	width := (maxX-minX)*scale + 2*padding + 2*offset
	height := (maxY-minY)*scale + 2*padding + 2*offset
	grid := make([][]rune, height)
	for i := range grid {
		grid[i] = make([]rune, width)
		for j := range grid[i] {
			grid[i][j] = ' '
		}
	}

	drawTunnels(grid, rooms, minX, minY)

	drawRooms(grid, rooms, minX, minY)

	printGrid(grid)

	return nil
}

func drawTunnels(grid [][]rune, rooms *model.Rooms, minX, minY int) {
	drawn := make(map[string]bool)
	connectionAttempts := make(map[string]int)

	for _, startRoom := range rooms.List {
		for _, endRoom := range startRoom.Links {
			key := fmt.Sprintf("%s-%s", minString(startRoom.ID, endRoom.ID), maxString(startRoom.ID, endRoom.ID))
			if drawn[key] {
				continue
			}
			drawn[key] = true

			x1 := (startRoom.X-minX)*scale + padding
			y1 := (startRoom.Y-minY)*scale + padding
			x2 := (endRoom.X-minX)*scale + padding
			y2 := (endRoom.Y-minY)*scale + padding

			// Find a clear path and draw it.
			findAndDrawPath(grid, x1, y1, x2, y2, connectionAttempts, key, rooms, startRoom, endRoom, minX, minY)
		}
	}
}

// findAndDrawPath finds the best route for a tunnel, avoiding obstacles.
// If the direct path is blocked, it finds an offset path with smooth corners.
func findAndDrawPath(grid [][]rune, x1, y1, x2, y2 int, attempts map[string]int, key string, allRooms *model.Rooms, startRoom, endRoom *model.Room, minX, minY int) {
	// 1. Check if the direct path is clear.
	if isPathClear(grid, x1, y1, x2, y2, allRooms, startRoom, endRoom, minX, minY) {
		drawLine(grid, x1, y1, x2, y2)
		return
	}

	// 2. If blocked, find an alternative path using an offset.
	attemptNum := attempts[key]
	attempts[key]++

	currentOffset := offset
	if attemptNum%2 != 0 {
		currentOffset = -offset // Alternate to try both directions (e.g., up/down).
	}

	// Try two offset directions (+offset and -offset).
	for i := 0; i < 2; i++ {
		var ox1, oy1, ox2, oy2 int

		// For horizontal tunnels, try vertical offsets. For others, try horizontal offsets.
		isHorizTunnel := y1 == y2
		if isHorizTunnel {
			oy1, oy2 = y1+currentOffset, y2+currentOffset
			ox1, ox2 = x1, x2
		} else {
			ox1, ox2 = x1+currentOffset, x2+currentOffset
			oy1, oy2 = y1, y2
		}

		// Check if all three segments of the potential offset path are clear.
		pathIsClear := isPathClear(grid, x1, y1, ox1, oy1, allRooms, startRoom, endRoom, minX, minY) &&
			isPathClear(grid, ox1, oy1, ox2, oy2, allRooms, startRoom, endRoom, minX, minY) &&
			isPathClear(grid, ox2, oy2, x2, y2, allRooms, startRoom, endRoom, minX, minY)

		if pathIsClear {
			// Draw the three segments to form a path with smooth corners.
			drawLine(grid, x1, y1, ox1, oy1)   // Corner from start room
			drawLine(grid, ox1, oy1, ox2, oy2) // Main offset path
			drawLine(grid, ox2, oy2, x2, y2)   // Corner to end room
			return
		}

		currentOffset = -currentOffset // Try the opposite offset direction.
	}

	// 3. If no clear path is found, draw the direct path as a last resort.
	drawLine(grid, x1, y1, x2, y2)
}

// Collision Detection and Primitives

// isPathClear checks if a path is free of obstacles (other tunnels or rooms).
func isPathClear(grid [][]rune, x1, y1, x2, y2 int, allRooms *model.Rooms, startRoom, endRoom *model.Room, minX, minY int) bool {
	dx, dy := float64(x2-x1), float64(y2-y1)
	steps := int(math.Max(math.Abs(dx), math.Abs(dy)))
	if steps == 0 {
		return true
	}
	xInc, yInc := dx/float64(steps), dy/float64(steps)

	fx, fy := float64(x1), float64(y1)
	for i := 0; i <= steps; i++ {
		x, y := int(math.Round(fx)), int(math.Round(fy))

		// Check for collision, ignoring the exact start and end points of the path.
		isEndpoint := (i == 0 || i == steps)
		if !isEndpoint && y >= 0 && y < len(grid) && x >= 0 && x < len(grid[0]) {
			// 1. Check for collision with already drawn tunnels.
			if grid[y][x] != ' ' {
				return false
			}
			// 2. Check for collision with other rooms.
			for _, otherRoom := range allRooms.List {
				if otherRoom.ID == startRoom.ID || otherRoom.ID == endRoom.ID {
					continue
				}
				otherX := (otherRoom.X-minX)*scale + padding
				otherY := (otherRoom.Y-minY)*scale + padding
				labelWidth := len(fmt.Sprintf("[%s]", otherRoom.ID))
				if x >= otherX && x < otherX+labelWidth && y == otherY {
					return false // Path is blocked by another room.
				}
			}
		}
		fx += xInc
		fy += yInc
	}
	return true
}

// drawLine draws a line using the appropriate characters (-, |, /, \).
func drawLine(grid [][]rune, x1, y1, x2, y2 int) {
	dx, dy := float64(x2-x1), float64(y2-y1)
	steps := int(math.Max(math.Abs(dx), math.Abs(dy)))
	if steps == 0 {
		return
	}
	xInc, yInc := dx/float64(steps), dy/float64(steps)

	fx, fy := float64(x1), float64(y1)
	for i := 0; i <= steps; i++ {
		x, y := int(math.Round(fx)), int(math.Round(fy))
		if y >= 0 && y < len(grid) && x >= 0 && x < len(grid[0]) && grid[y][x] == ' ' {
			isDiag := math.Abs(xInc) > 0.4 && math.Abs(yInc) > 0.4
			if isDiag {
				grid[y][x] = '\\'
				if (xInc > 0 && yInc < 0) || (xInc < 0 && yInc > 0) {
					grid[y][x] = '/'
				}
			} else if math.Abs(dx) > math.Abs(dy) {
				grid[y][x] = '-'
			} else {
				grid[y][x] = '|'
			}
		}
		fx += xInc
		fy += yInc
	}
}

// --- Display and Utility Functions ---

// drawRooms places room labels like "[room_name]" on the grid.
func drawRooms(grid [][]rune, rooms *model.Rooms, minX, minY int) {
	for _, room := range rooms.List {
		x := (room.X-minX)*scale + padding
		y := (room.Y-minY)*scale + padding
		label := fmt.Sprintf("[%s]", room.ID)
		// Clear the area for the room label before drawing it.
		for i := 0; i < len(label); i++ {
			if y < len(grid) && x+i < len(grid[0]) {
				grid[y][x+i] = ' '
			}
		}
		// Draw the label.
		for i, char := range label {
			if y < len(grid) && x+i < len(grid[0]) {
				grid[y][x+i] = char
			}
		}
	}
}

// printGrid outputs the final ASCII art grid to the console.
func printGrid(grid [][]rune) {
	fmt.Println(strings.Repeat("#", len(grid[0])))
	for _, row := range grid {
		fmt.Println(strings.TrimRight(string(row), " "))
	}
	fmt.Println(strings.Repeat("#", len(grid[0])))
}

// findBounds determines the min/max X and Y coordinates to size the grid.
func findBounds(rooms *model.Rooms) (minX, maxX, minY, maxY int) {
	if len(rooms.List) == 0 {
		return 0, 0, 0, 0
	}
	minX, maxX, minY, maxY = rooms.List[0].X, rooms.List[0].X, rooms.List[0].Y, rooms.List[0].Y
	for _, r := range rooms.List {
		if r.X < minX {
			minX = r.X
		}
		if r.X > maxX {
			maxX = r.X
		}
		if r.Y < minY {
			minY = r.Y
		}
		if r.Y > maxY {
			maxY = r.Y
		}
	}
	return minX, maxX, minY, maxY
}

func minString(a, b string) string {
	if a < b {
		return a
	}
	return b
}

func maxString(a, b string) string {
	if a > b {
		return a
	}
	return b
}
