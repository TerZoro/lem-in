package pathfinding

import (
	"lemin/internal/model"
	"sort"
)

// DistributionResult represents the result of ant distribution optimization
type DistributionResult struct {
	Paths        []*model.Path
	Distribution []int
	TotalTurns   int
}

// optimizeAntDistribution finds the optimal distribution of ants across paths
func optimizeAntDistribution(paths []*model.Path, numAnts int) *DistributionResult {
	// Sort paths by length for optimal distribution
	sort.Slice(paths, func(i, j int) bool {
		return paths[i].Length < paths[j].Length
	})

	// Initialize DP table
	// dp[i][j] represents the minimum turns needed to move j ants using paths[0...i]
	dp := make([][]int, len(paths))
	for i := range dp {
		dp[i] = make([]int, numAnts+1)
		for j := range dp[i] {
			dp[i][j] = -1
		}
	}

	// Initialize distribution array
	distribution := make([]int, len(paths))

	// Base case: using only the first path
	for j := 0; j <= numAnts; j++ {
		dp[0][j] = paths[0].Length + j - 1
	}

	// Fill DP table
	for i := 1; i < len(paths); i++ {
		for j := 0; j <= numAnts; j++ {
			minTurns := dp[i-1][j] // Don't use current path
			for k := 1; k <= j; k++ {
				// Try using k ants on current path
				turns := max(dp[i-1][j-k], paths[i].Length+k-1)
				if minTurns == -1 || turns < minTurns {
					minTurns = turns
				}
			}
			dp[i][j] = minTurns
		}
	}

	// Reconstruct optimal distribution
	remainingAnts := numAnts
	for i := len(paths) - 1; i >= 0; i-- {
		if i == 0 {
			distribution[i] = remainingAnts
			break
		}

		bestK := 0
		bestTurns := dp[i][remainingAnts]
		for k := 1; k <= remainingAnts; k++ {
			turns := max(dp[i-1][remainingAnts-k], paths[i].Length+k-1)
			if turns == bestTurns {
				bestK = k
				break
			}
		}

		distribution[i] = bestK
		remainingAnts -= bestK
	}

	return &DistributionResult{
		Paths:        paths,
		Distribution: distribution,
		TotalTurns:   dp[len(paths)-1][numAnts],
	}
}

// max returns the maximum of two integers
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
