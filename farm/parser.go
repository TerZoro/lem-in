package farm

import (
	"lemin/internal/model"
	"strconv"
	"strings"
)

func ParseRooms(lines []string) (*model.Rooms, error) {
	rooms := &model.Rooms{}
	for i := 1; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		if strings.Contains(line, "-") {
			// End of rooms, start of links
			err := parseLinks(rooms, lines[i:])
			return rooms, err
		}

		parts := strings.Fields(line)
		if len(parts) != 3 {
			continue
		}

		id, err1 := strconv.Atoi(parts[0])
		x, err2 := strconv.Atoi(parts[1])
		y, err3 := strconv.Atoi(parts[2])
		if err1 != nil || err2 != nil || err3 != nil {
			continue
		}

		flag := "##0"
		if i > 0 {
			switch strings.TrimSpace(lines[i-1]) {
			case "##start":
				flag = "##start"
			case "##end":
				flag = "##end"
			}
		}

		room := &model.Room{ID: id, X: x, Y: y, Flag: flag}
		rooms.AddRoom(room)
	}
	return rooms, nil

}

func parseLinks(rooms *model.Rooms, lines []string) error {
	for _, line := range lines {
		parts := strings.Split(line, "-")
		if len(parts) != 2 {
			continue
		}
		from, err1 := strconv.Atoi(parts[0])
		to, err2 := strconv.Atoi(parts[1])
		if err1 != nil || err2 != nil {
			continue
		}
		if err := rooms.AddLink(from, to); err != nil {
			return err
		}
	}
	return nil
}
