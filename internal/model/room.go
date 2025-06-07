package model

import (
	"fmt"
	errhandle "lemin/pkg/errors"
)

type Rooms struct {
	List []*Room
}

func (rms *Rooms) AddRoom(room *Room) {
	rms.List = append(rms.List, room)
}

func (rms *Rooms) searchByID(id int) (*Room, error) {
	if len(rms.List) < 1 {
		return nil, errhandle.FormatError(errhandle.ErrEmptyList)
	}
	for i := range rms.List {
		if rms.List[i].ID == id {
			return rms.List[i], nil
		}
	}
	return nil, errhandle.FormatError(errhandle.ErrRoomNotFound)
}

func (rms *Rooms) SearchStart() (*Room, error) {
	if len(rms.List) < 1 {
		return nil, errhandle.FormatError(errhandle.ErrEmptyList)
	}
	for i := range rms.List {
		if rms.List[i].Flag == "##start" {
			return rms.List[i], nil
		}
	}
	return nil, errhandle.FormatError(errhandle.ErrStartRoom)
}

func (rms *Rooms) SearchEnd() (*Room, error) {
	if len(rms.List) < 1 {
		return nil, errhandle.FormatError(errhandle.ErrEmptyList)
	}
	for i := range rms.List {
		if rms.List[i].Flag == "##end" {
			return rms.List[i], nil
		}
	}
	return nil, errhandle.FormatError(errhandle.ErrEndRoom)
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
	return nil, errhandle.FormatError(errhandle.ErrRoomNotFound)
}

func (r *Rooms) GetStartRoom() (*Room, error) {
	if len(r.List) == 0 {
		return nil, errhandle.FormatError(errhandle.ErrEmptyList)
	}
	for i := range r.List {
		if r.List[i].Flag == "##start" {
			return r.List[i], nil
		}
	}
	return nil, errhandle.FormatError(errhandle.ErrStartRoom)
}

func (r *Rooms) GetEndRoom() (*Room, error) {
	if len(r.List) == 0 {
		return nil, errhandle.FormatError(errhandle.ErrEmptyList)
	}
	for i := range r.List {
		if r.List[i].Flag == "##end" {
			return r.List[i], nil
		}
	}
	return nil, errhandle.FormatError(errhandle.ErrEndRoom)
}

func (r *Rooms) GetRoomByID(id int) (*Room, error) {
	if len(r.List) == 0 {
		return nil, errhandle.FormatError(errhandle.ErrEmptyList)
	}
	room, err := r.searchByID(id)
	if err != nil {
		return nil, err
	}
	if room == nil {
		return nil, errhandle.FormatError(errhandle.ErrRoomNotFound)
	}
	return room, nil
}

type Room struct {
	ID    int
	X     int
	Y     int
	Flag  string // ##start, ##end, ##0
	Links []*Room
}
