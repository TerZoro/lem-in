package farm

import (
	"fmt"
	"lemin/internal/model"
	errhandle "lemin/pkg/errors"
	"strconv"
	"strings"
)

func ParseRooms(lines []string) (*model.Rooms, error) {
	if len(lines) == 0 {
		return nil, errhandle.FormatError(errhandle.ErrParseRooms)
	}

	rooms := &model.Rooms{}
	for i := 1; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		if strings.Contains(line, "-") {
			fmt.Println("Debug: Found links section")
			// End of rooms, start of links
			err := parseLinks(rooms, lines[i:])
			if err != nil {
				fmt.Printf("Debug: Error parsing links: %v\n", err)
				return nil, errhandle.FormatError(errhandle.ErrParseRooms)
			}
			return rooms, nil
		}

		parts := strings.Fields(line)
		if len(parts) != 3 {
			fmt.Printf("Debug: Skipping invalid room line: %s\n", line)
			return nil, errhandle.FormatError(errhandle.ErrParseRooms)
		}

		// First part is now a string (room name)
		roomName := parts[0]
		x, err1 := strconv.Atoi(parts[1])
		y, err2 := strconv.Atoi(parts[2])
		if err1 != nil || err2 != nil {
			fmt.Printf("Debug: Invalid coordinates in line: %s\n", line)
			return nil, errhandle.FormatError(errhandle.ErrParseRooms)
		}

		flag := "##0"
		if i > 0 {
			switch strings.TrimSpace(lines[i-1]) {
			case "##start":
				flag = "##start"
				fmt.Printf("Debug: Found start room: %s\n", roomName)
			case "##end":
				flag = "##end"
				fmt.Printf("Debug: Found end room: %s\n", roomName)
			}
		}

		room := &model.Room{ID: roomName, X: x, Y: y, Flag: flag}
		rooms.AddRoom(room)
		fmt.Printf("Debug: Added room: %s at (%d, %d) with flag %s\n", roomName, x, y, flag)
	}
	return rooms, nil
}

func parseLinks(rooms *model.Rooms, lines []string) error {
	fmt.Println("Debug: Starting to parse links")
	for _, line := range lines {
		parts := strings.Split(line, "-")
		if len(parts) != 2 {
			fmt.Printf("Debug: Skipping invalid link line: %s\n", line)
			continue
		}
		from := strings.TrimSpace(parts[0])
		to := strings.TrimSpace(parts[1])
		if err := rooms.AddLink(from, to); err != nil {
			fmt.Printf("Debug: Error adding link %s-%s: %v\n", from, to, err)
			return err
		}
		fmt.Printf("Debug: Added link: %s-%s\n", from, to)
	}
	return nil
}
