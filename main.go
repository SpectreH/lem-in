package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

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

var ANTCOUNTER int
var STARTROOMID int

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Missing file name")
		os.Exit(1)
	}

	data := LoadData(os.Args[1])
	generationData := ReadData(data)
	farm := GenerateFarm(generationData)

	//_ = farm
	var path []*Room
	var paths [][]*Room
	CalculatePath(path, farm[STARTROOMID], 0, &paths, &farm[STARTROOMID])
	fmt.Println(paths)
	ants := SpawnAnts(&farm[STARTROOMID], paths)

	MakeStep(ants)
}

func SpawnAnts(startRoom *Room, paths [][]*Room) []Ant {
	var result []Ant

	pathIndex := 0
	for i := 1; i < ANTCOUNTER+1; i++ {
		var antToAppend Ant

		antToAppend.Id = i
		antToAppend.CurrentRoom = startRoom

		if pathIndex == len(paths) {
			pathIndex = 0
		}

		antToAppend.Path = paths[pathIndex]
		pathIndex++

		result = append(result, antToAppend)
	}

	return result
}

func MakeStep(ants []Ant) {
	var allPassed bool = true

	for i := 0; i < len(ants); i++ {
		if ants[i].CurrentRoom.IsEnd {
			continue
		}

		nextRoomId := ants[i].RoomsPassed

		if ants[i].Path[nextRoomId].Ants != 0 {
			if !ants[i].Path[nextRoomId].IsEnd {
				continue
			}
		}

		ants[i].CurrentRoom.Ants--
		ants[i].CurrentRoom = ants[i].Path[nextRoomId]
		ants[i].CurrentRoom.Ants++
		ants[i].RoomsPassed++
		allPassed = false

		fmt.Print("L", ants[i].Id, "-", ants[i].CurrentRoom.Name, " ")
	}

	if allPassed {
		return
	} else {
		fmt.Println("")
		MakeStep(ants)
	}
}

func LoadData(fileName string) [][]byte {
	data, err := os.ReadFile(os.Args[1])

	if err != nil {
		log.Fatalf("failed to open: %s", fileName)
	}

	sep := []byte{13, 10}
	transformedData := bytes.Split(data, sep)

	return transformedData
}

func ReadData(data [][]byte) GenerationData {
	var result GenerationData

	var err error
	ANTCOUNTER, err = strconv.Atoi(string(data[0]))
	CheckError(err)

	if ANTCOUNTER <= 0 {
		log.Fatal("Invalid number of Ants!")
	}

	var startFound, endFound bool

	var commentsCounter int = 1
	for i := 1; i < len(data); i++ {
		if strings.Contains(string(data[i]), "##") {
			if string(data[i]) == "##start" {
				startFound = true
				result.StartIndex = i - commentsCounter
			} else if string(data[i]) == "##end" {
				endFound = true
				result.EndIndex = i - commentsCounter
			} else {
				log.Fatal("Invalid start or end data format!")
			}
			commentsCounter++
			continue
		} else if strings.Contains(string(data[i]), "#") {
			commentsCounter++
			continue
		}

		if strings.Count(string(data[i]), " ") == 2 {
			result.Rooms = append(result.Rooms, string(data[i]))
		} else if strings.Count(string(data[i]), "-") == 1 {
			result.Links = append(result.Links, string(data[i]))
		} else {
			log.Fatal("Invalid link data format!")
		}
	}

	if !startFound || !endFound {
		log.Fatal("Invalid data format, no start or end room found")
	}

	return result
}

func GenerateFarm(data GenerationData) []Room {
	var farm []Room
	var err error

	for i := 0; i < len(data.Rooms); i++ {
		var roomToAppend Room

		if i == data.StartIndex {
			roomToAppend.IsStart = true
			STARTROOMID = i
			roomToAppend.Ants = ANTCOUNTER
		} else if i == data.EndIndex {
			roomToAppend.IsEnd = true
		}

		splittedData := strings.Split(data.Rooms[i], " ")

		roomToAppend.Name = splittedData[0]
		roomToAppend.X_pos, err = strconv.Atoi(splittedData[1])
		CheckError(err)
		roomToAppend.Y_pos, err = strconv.Atoi(splittedData[2])
		CheckError(err)

		farm = append(farm, roomToAppend)
	}

	farm = ConnectLinks(farm, data.Links)

	return farm
}

func ConnectLinks(farm []Room, links []string) []Room {
	for i := 0; i < len(links); i++ {
		splittedData := strings.Split(links[i], "-")
		for k := 0; k < len(farm); k++ {
			if (farm)[k].Name == splittedData[0] {
				for m := 0; m < len(farm); m++ {
					if (farm)[m].Name == splittedData[1] {
						if (farm)[m].Name == (farm)[k].Name {
							log.Fatal("Invalid data format! Self-link is prohibited")
						}
						(farm)[k].Links = append((farm)[k].Links, &(farm)[m])
						break
					}
					if m == len(farm)-1 {
						log.Fatal("Invalid data format! Room link not found")
					}
				}
				break
			}
			if k == len(farm)-1 {
				log.Fatal("Invalid data format! Room link not found")
			}
		}
	}

	return farm
}

func CalculatePath(path []*Room, currentRoom Room, step int, paths *[][]*Room, previousRoom *Room) {
	if currentRoom.IsEnd {
		var skipPath bool

		for k := 0; k < len(*paths); k++ {
			if (*paths)[k] == nil {
				continue
			}

			if path[0].Name == (*paths)[k][0].Name {
				if len(path) < len((*paths)[k]) {
					(*paths)[k] = nil
					break
				} else {
					skipPath = true
					break
				}
			}
		}

		if len(*paths) == 0 {
			*paths = append((*paths), nil)
		} else if (*paths)[len(*paths)-1] != nil {
			*paths = append((*paths), nil)
		}

		for i := 0; i < len(path); i++ {
			if !skipPath {
				(*paths)[len(*paths)-1] = append((*paths)[len(*paths)-1], path[i])
			} else {
				break
			}
		}
	}

	for i := 0; i < len(currentRoom.Links); i++ {
		var toContinue bool

		for k := 0; k < len(path); k++ {
			if path[k].Name == currentRoom.Links[i].Name {
				toContinue = true
				break
			}
		}

		if !toContinue {
			pathToPass := path
			pathToPass = append(pathToPass, currentRoom.Links[i])
			CalculatePath(pathToPass, *currentRoom.Links[i], step+1, paths, &currentRoom)
			pathToPass = path
		}
	}

	for i := 0; i < len(*paths); i++ {
		if (*paths)[i] == nil {
			*paths = append((*paths)[:i], (*paths)[i+1:]...)
		}
	}
}

func CheckError(err error) {
	if err != nil {
		log.Fatal("Invalid data format!")
	}
}
