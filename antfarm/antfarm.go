package antfarm

type antFarm struct {
	rooms []room
	ants  ants
}

func (af *antFarm) addRoom(r room) {
	af.rooms = append(af.rooms, r)
}
