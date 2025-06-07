package model

type Ant struct {
	ID        int
	PathIndex int // which path this ant is following
	StepIndex int // current step in the path (0 = start room)
	RoomID    int // current room ID
	Active    bool
	Finished  bool
}
