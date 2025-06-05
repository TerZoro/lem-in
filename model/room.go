package model

import (
	"errors"
	"fmt"
)

type Rooms struct {
	List []*Room
}

func (rms *Rooms) AddRoom(room *Room) {
	rms.List = append(rms.List, room)
}

func (rms *Rooms) searchByID(id int) (*Room, error) {
	if len(rms.List) < 1 {
		return nil, errors.New("error: list is empty")
	}
	for i := range rms.List {
		if rms.List[i].ID == id {
			return rms.List[i], nil
		}
	}
	return nil, fmt.Errorf("room with ID %d not found", id)
}

func (rms *Rooms) SearchStart() (*Room, error) {
	if len(rms.List) < 1 {
		return nil, errors.New("error: list is empty")
	}
	for i := range rms.List {
		if rms.List[i].Flag == "##start" {
			return rms.List[i], nil
		}
	}
	return nil, fmt.Errorf("start room not found")
}

func (rms *Rooms) SearchEnd() (*Room, error) {
	if len(rms.List) < 1 {
		return nil, errors.New("error: list is empty")
	}
	for i := range rms.List {
		if rms.List[i].Flag == "##end" {
			return rms.List[i], nil
		}
	}
	return nil, fmt.Errorf("end room not found")
}

func (rms *Rooms) AddLink(fromID, toID int) error {
	from, err := rms.searchByID(fromID)
	if err != nil {
		return err
	}
	to, err := rms.searchByID(toID)
	if err != nil {
		return err
	}
	from.Links = append(from.Links, to)
	to.Links = append(to.Links, from)
	return nil
}

func (rms *Rooms) Printer() {
	fmt.Println("Parsed rooms:")
	for _, r := range rms.List {
		fmt.Printf("Room: %d (%d, %d)\n", r.ID, r.X, r.Y)
	}
}

func (rms *Rooms) Search(id int) (*Room, error) {
	for i := range rms.List {
		if rms.List[i].ID == id {
			return rms.List[i], nil
		}
	}
	return nil, errors.New("error: room not found")
}

type Room struct {
	ID    int
	X     int
	Y     int
	Flag  string // ##start, ##end, ##0
	Links []*Room
}

func (r *Room) flagCheck() (string, error) {
	if r.Flag != "" {
		return r.Flag, nil
	}
	return "", errors.New("error: no flag available")
}
