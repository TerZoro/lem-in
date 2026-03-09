# Lem-in

A Go implementation of the lem-in algorithm that finds the optimal path for ants to travel from start to end through a graph without intersections.

## Description

Lem-in is an algorithm that simulates ant colony movement through a graph. The goal is to find the optimal set of non-intersecting paths from start to end, and then simulate the movement of ants through these paths in the minimum number of steps.

## Features

- **Optimal Path Finding**: Uses BFS to find shortest paths from each neighbor of start to end
- **Non-intersecting Path Selection**: Finds the best combination of paths that don't share nodes (except start/end)
- **Ant Distribution**: Optimally distributes ants across selected paths
- **Movement Simulation**: Simulates step-by-step ant movement through the graph
- **Clean Output**: Matches the required output format exactly

## Usage

```bash
go run main.go <input_file>
```

### Example

```bash
go run main.go tests/example01
```

## Input Format

The input file should contain:

1. **Number of ants** (first line)
2. **Rooms** (nodes) with coordinates:
   - `##start` - marks the start room
   - `##end` - marks the end room
   - `room_name x_coord y_coord` - regular rooms
3. **Connections** (edges):
   - `room1-room2` - bidirectional connection between rooms

### Example Input
```
10
##start
start 1 6
0 4 8
o 6 8
n 6 6
e 8 4
t 1 9
E 5 9
a 8 9
m 8 6
h 4 6
A 5 2
c 8 1
k 11 2
##end
end 11 6
start-t
n-e
a-m
A-c
0-o
E-a
k-end
start-h
o-n
m-end
t-E
start-0
h-A
e-end
c-k
n-m
h-n
```

## Output Format

The program outputs:
1. **Input file contents** (exactly as provided)
2. **Ant movements** (one step per line):
   - `L1-room1 L2-room2` - ants moving to rooms
   - Ants are numbered L1, L2, L3, etc.
   - Movements are sorted by ant ID

### Example Output
```
10
##start
start 1 6
...
start-t
n-e
...
h-n

L1-t L2-h L3-0
L1-E L2-A L3-o L4-t L5-h L6-0
L1-a L2-c L3-n L4-E L5-A L6-o L7-t L8-h L9-0
...
```

## Algorithm

### 1. Path Finding
- Uses BFS to find shortest paths from each neighbor of start to end
- Ensures paths are optimal (shortest possible)

### 2. Flow Selection
- Generates all possible combinations of non-intersecting paths
- Selects the combination that minimizes total steps for all ants
- Uses backtracking to explore all valid combinations

### 3. Ant Distribution
- Distributes ants optimally across selected paths
- Accounts for path length differences
- Ensures fair distribution (extra ants go to shorter paths)

### 4. Movement Simulation
- Simulates ant movement step by step
- Ants move one room per turn
- New ants are sent only when there's space
- Ants reaching the end are removed from simulation

## Project Structure

```
lem-in/
├── main.go          # Main program entry point
├── graph/
│   ├── graph.go     # Graph implementation and algorithms
│   └── node.go      # Node structure and parsing
└── tests/           # Test input files
    ├── example01
    ├── example02
    └── ...
```

## Key Functions

- `findShortestPaths()`: Finds shortest paths from each neighbor to end
- `bfsShortestPath()`: BFS implementation for shortest path finding
- `flows()`: Generates all non-intersecting path combinations
- `execute()`: Simulates ant movement through selected paths
- `calculateSteps()`: Calculates total steps for a given flow

## Performance

- **Time Complexity**: O(V + E) for BFS path finding
- **Space Complexity**: O(V) for visited nodes and queue
- **Optimal**: Always finds the best non-intersecting path combination

## Testing

Test with various input files:

```bash
go run main.go tests/example01
go run main.go tests/example02
go run main.go tests/example03
```

## Requirements

- Go 1.18+ (for generics support)
- Standard library only (no external dependencies)

## License

This project is part of the 42 school curriculum. 
