package farm

import (
	errhandle "lemin/pkg/errors"
	"strconv"
)

func AntsNum(lines []string) (int, error) {
	ants, err := strconv.Atoi(lines[0])
	if err != nil || ants <= 0 {
		return 0, errhandle.FormatError(errhandle.ErrInvalidAnts)
	}
	return ants, nil
}
