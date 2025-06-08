package errhandle

import (
	"fmt"
	"log"
)

// Common error messages
const (
	ErrFileOpen     = "could not open file: %v"
	ErrFileRead     = "error reading file: %v"
	ErrNoStartEnd   = "missing ##start or ##end"
	ErrNoTestFile   = "you should provide test file"
	ErrNoPaths      = "no paths found from start to end"
	ErrInvalidAnts  = "invalid number of ants"
	ErrEmptyList    = "list is empty"
	ErrRoomNotFound = "room not found"
	ErrParseRooms   = "could not parse rooms"
	ErrStartRoom    = "start room not found"
	ErrEndRoom      = "end room not found"
	ErrDrawGraph    = "no rooms to draw"
)

func FormatError(msg string, args ...interface{}) error {
	return fmt.Errorf("ERROR: "+msg, args...)
}

func HandleError(err error) {
	if err != nil {
		fmt.Println("Usage: go run . [filename]")
		log.Fatal(err)
	}
}
