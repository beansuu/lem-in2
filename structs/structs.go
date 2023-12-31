package structs

type Ant struct {
	Id          int
	Path        []*Room
	CurrentRoom *Room
	RoomsPassed int
}

type Room struct {
	Name    string
	Ants    int
	X_pos   int
	Y_pos   int
	IsStart bool
	IsEnd   bool
	Links   []*Room
}

type GenerationData struct {
	Rooms      []string
	Links      []string
	StartIndex int
	EndIndex   int
}

type PathStuct struct {
	PathName string
	Paths    [][]*Room
}

var ANTCOUNTER int
var STARTROOMID int
var ENDROOMID int
var FARM []Room

var BEST_TURNS_RES int
var BEST_PATH [][]*Room
var BEST_ROOMS_IN_USE_RES int
