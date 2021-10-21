package structs

// Ant
type Ant struct {
	Id          int
	Path        []*Room
	CurrentRoom *Room
	RoomsPassed int
}

// Room in farm
type Room struct {
	Name    string
	Ants    int
	X_pos   int
	Y_pos   int
	IsStart bool
	IsEnd   bool
	Links   []*Room
}

// Data to generate future farm
type GenerationData struct {
	Rooms      []string
	Links      []string
	StartIndex int
	EndIndex   int
}

// Structurized paths by links in start room
type PathStuct struct {
	PathName string
	Paths    [][]*Room
}

var ANTCOUNTER int // Amount of ants to spawn
var STARTROOMID int
var ENDROOMID int
var FARM []Room // Farm

var BEST_TURNS_RES int
var BEST_PATH [][]*Room
var BEST_ROOMS_IN_USE_RES int
