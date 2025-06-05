package farm

import (
	"log"
	"strconv"
)

func AntsNum(lines []string) int {
	ants, err := strconv.Atoi(lines[0])
	if err != nil || ants <= 0 {
		log.Fatalf("ERROR: invalid number of ants")
	}
	return ants
}
