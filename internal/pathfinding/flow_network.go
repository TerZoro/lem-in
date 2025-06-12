package pathfinding

import (
	"fmt"
	"lemin/internal/model"
	"math"
)

// FlowNetwork represents the flow network for the ant farm
type FlowNetwork struct {
	capacity  map[string]map[string]int // Capacity matrix
	flow      map[string]map[string]int // Flow matrix
	residual  map[string]map[string]int // Residual capacity matrix
	parent    map[string]string         // Parent array for BFS
	visited   map[string]bool           // Visited array for BFS
	rooms     *model.Rooms
	roomCount int
}

// NewFlowNetwork creates a new flow network from the given rooms
func NewFlowNetwork(rooms *model.Rooms) *FlowNetwork {
	network := &FlowNetwork{
		capacity:  make(map[string]map[string]int),
		flow:      make(map[string]map[string]int),
		residual:  make(map[string]map[string]int),
		parent:    make(map[string]string),
		visited:   make(map[string]bool),
		rooms:     rooms,
		roomCount: len(rooms.List),
	}

	// Initialize maps for each room
	for _, room := range rooms.List {
		network.capacity[room.ID] = make(map[string]int)
		network.flow[room.ID] = make(map[string]int)
		network.residual[room.ID] = make(map[string]int)
	}

	// Set up capacities and residual capacities
	for _, room := range rooms.List {
		for _, link := range room.Links {
			// Set capacity to 1 for each edge
			network.capacity[room.ID][link.ID] = 1
			network.capacity[link.ID][room.ID] = 1
			network.residual[room.ID][link.ID] = 1
			network.residual[link.ID][room.ID] = 1
		}
	}

	startRoom, _ := rooms.GetStartRoom()
	endRoom, _ := rooms.GetEndRoom()
	fmt.Printf("Debug: Created flow network with %d rooms\n", network.roomCount)
	if startRoom != nil && endRoom != nil {
		fmt.Printf("Debug: Start room: %s, End room: %s\n", startRoom.ID, endRoom.ID)
	} else {
		fmt.Println("Debug: Could not find start or end room")
	}
	return network
}

// findAugmentingPath uses BFS to find an augmenting path in the residual network
func (fn *FlowNetwork) findAugmentingPath(source, sink string) []string {
	// Reset visited and parent maps
	for id := range fn.visited {
		fn.visited[id] = false
		fn.parent[id] = ""
	}

	// Initialize queue for BFS
	queue := []string{source}
	fn.visited[source] = true

	// BFS to find augmenting path
	for len(queue) > 0 {
		u := queue[0]
		queue = queue[1:]

		for v := range fn.residual[u] {
			if !fn.visited[v] && fn.residual[u][v] > 0 {
				fn.parent[v] = u
				fn.visited[v] = true
				queue = append(queue, v)

				if v == sink {
					// Reconstruct path
					path := []string{v}
					for v != source {
						v = fn.parent[v]
						path = append([]string{v}, path...)
					}
					fmt.Printf("Debug: Found augmenting path: %v\n", path)
					return path
				}
			}
		}
	}

	fmt.Printf("Debug: No augmenting path found from %s to %s\n", source, sink)
	return nil
}

// findMinCapacity finds the minimum capacity in the given path
func (fn *FlowNetwork) findMinCapacity(path []string) int {
	minCapacity := math.MaxInt32
	for i := 0; i < len(path)-1; i++ {
		u, v := path[i], path[i+1]
		if fn.residual[u][v] < minCapacity {
			minCapacity = fn.residual[u][v]
		}
	}
	fmt.Printf("Debug: Min capacity on path: %d\n", minCapacity)
	return minCapacity
}

// updateFlow updates the flow and residual capacities along the given path
func (fn *FlowNetwork) updateFlow(path []string, flow int) {
	for i := 0; i < len(path)-1; i++ {
		u, v := path[i], path[i+1]
		fn.flow[u][v] += flow
		fn.flow[v][u] -= flow
		fn.residual[u][v] -= flow
		fn.residual[v][u] += flow
	}
	fmt.Printf("Debug: Updated flow by %d along path %v\n", flow, path)
}

// EdmondsKarp implements the Edmonds-Karp algorithm to find the maximum flow
func (fn *FlowNetwork) EdmondsKarp(source, sink string) int {
	maxFlow := 0
	for {
		path := fn.findAugmentingPath(source, sink)
		if path == nil {
			break
		}
		flow := fn.findMinCapacity(path)
		fn.updateFlow(path, flow)
		maxFlow += flow
	}
	fmt.Printf("Debug: Edmonds-Karp completed with max flow: %d\n", maxFlow)
	return maxFlow
}

// extractPaths extracts all possible paths from the residual network
func (fn *FlowNetwork) extractPaths(source, sink string) [][]string {
	var paths [][]string
	used := make(map[string]bool)

	// Create a copy of the flow network for path extraction
	flowCopy := make(map[string]map[string]int)
	for u := range fn.flow {
		flowCopy[u] = make(map[string]int)
		for v, flow := range fn.flow[u] {
			flowCopy[u][v] = flow
		}
	}

	for {
		// Try to find a path from source to sink
		path := []string{source}
		current := source
		visited := make(map[string]bool)
		visited[source] = true

		for current != sink {
			found := false
			for next := range flowCopy[current] {
				if !visited[next] && !used[next] && flowCopy[current][next] > 0 {
					path = append(path, next)
					current = next
					visited[next] = true
					found = true
					break
				}
			}
			if !found {
				break
			}
		}

		if current == sink {
			// Mark all nodes in the path as used (except source and sink)
			for i := 1; i < len(path)-1; i++ {
				used[path[i]] = true
			}
			paths = append(paths, path)
			fmt.Printf("Debug: Found path: %v\n", path)

			// Update flow for this path
			for i := 0; i < len(path)-1; i++ {
				u, v := path[i], path[i+1]
				flowCopy[u][v] = 0
				flowCopy[v][u] = 0
			}
		} else {
			break
		}
	}

	fmt.Printf("Debug: Extracted %d paths\n", len(paths))
	return paths
}

// FindPaths finds all possible paths from source to sink
func (fn *FlowNetwork) FindPaths(source, sink string) [][]string {
	// First run Edmonds-Karp to find the maximum flow
	maxFlow := fn.EdmondsKarp(source, sink)
	if maxFlow == 0 {
		fmt.Println("Debug: No flow found from source to sink")
		return nil
	}

	// Then extract the paths
	paths := fn.extractPaths(source, sink)
	if len(paths) == 0 {
		fmt.Println("Debug: No paths extracted after finding flow")
		return nil
	}

	return paths
}
