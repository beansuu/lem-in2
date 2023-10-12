package antfarm

type room struct {
	name string
	x, y int
	ants ants
}

func newRoom(name string, x, y int) room {
	return room{name: name, x: x, y: y}
}
