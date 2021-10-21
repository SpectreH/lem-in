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

type PathStuct struct {
	PathName string
	Paths    [][]*Room
}

var ANTCOUNTER int
var STARTROOMID int
var ENDROOMID int

var RESULT int
var BESTPATH [][]*Room
var PATHSINUSE int

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Missing file name")
		os.Exit(1)
	}

	data := LoadData(os.Args[1])
	generationData := ReadData(data)
	farm := GenerateFarm(generationData)

	var path []*Room
	var paths [][]*Room
	CalculatePath(path, farm[STARTROOMID], 0, &paths, &farm[STARTROOMID])
	SortPaths(&paths)
	test := PotentialPaths(paths, farm[STARTROOMID])
	tempVar := make([][]*Room, len(farm[STARTROOMID].Links))
	FindBestPaths(&farm[STARTROOMID], &farm[ENDROOMID], test, tempVar, 0, len(farm[STARTROOMID].Links))

	// for b := 0; b < len(BESTPATH); b++ {
	// 	for c := 0; c < len(BESTPATH[b]); c++ {
	// 		fmt.Print(BESTPATH[b][c].Name)

	// 		if c != len(BESTPATH[b])-1 {
	// 			fmt.Print(" -> ")
	// 		} else {
	// 			fmt.Println()
	// 		}
	// 	}
	// }

	// pathsForAnts := make([][]*Room, ANTCOUNTER)
	// var tempVar [][]*Room

	// _, _ = Test(&farm, test, pathsForAnts, tempVar, &farm[STARTROOMID], &farm[ENDROOMID], 0, 0)
	// fmt.Println(BESTPATH)
	// cvb := [][]*Room{paths[3], paths[2], paths[4], paths[3], paths[2], paths[4], paths[3], paths[2], paths[4], paths[3]}
	// cc := [][]*Room{paths[6], paths[6]}
	RestoreFarm(&farm)

	ants := SpawnAnts(&farm[STARTROOMID], BESTPATH)

	MakeStep(ants, &farm[ENDROOMID])
}

func Foo(paths [][]*Room, startRoom *Room, endRoom *Room) ([][]*Room, int, int) {
	var resultPath [][]*Room

	sortedByLenghtPaths := make([][]*Room, 0)
	sortedByLenghtPaths = append(sortedByLenghtPaths, paths...)
	SortPaths(&sortedByLenghtPaths)

	var roomsInUse []*Room
	checkedPath := make([][]*Room, 0)
	for i := 0; i < len(sortedByLenghtPaths); i++ {
		var skipPath bool = false
		for k := 0; k < len(sortedByLenghtPaths[i]); k++ {
			startLenght := len(roomsInUse)
			for m := 0; m < startLenght; m++ {
				if roomsInUse[m].IsEnd {
					continue
				}
				if roomsInUse[m].Name == sortedByLenghtPaths[i][k].Name {
					skipPath = true
					break
				}
				roomsInUse = append(roomsInUse, sortedByLenghtPaths[i][k])
			}
			if skipPath {
				break
			}
			if len(roomsInUse) == 0 {
				roomsInUse = append(roomsInUse, sortedByLenghtPaths[i]...)
				break
			}
		}
		if skipPath {
			sortedByLenghtPaths = append(sortedByLenghtPaths[:i], sortedByLenghtPaths[i+1:]...)
		} else {
			checkedPath = append(checkedPath, sortedByLenghtPaths[i])
		}
	}

	// fmt.Println()
	// for b := 0; b < len(checkedPath); b++ {
	// 	for c := 0; c < len(checkedPath[b]); c++ {
	// 		fmt.Print(checkedPath[b][c].Name)

	// 		if c != len(checkedPath[b])-1 {
	// 			fmt.Print(" -> ")
	// 		} else {
	// 			fmt.Println()
	// 		}
	// 	}
	// }

	antPosTable := make([]int, len(checkedPath))
	for i := 0; i < ANTCOUNTER; i++ {
		var minimumPath []*Room
		var indexForPosTable int = -1

		for k := 0; k < len(checkedPath); k++ {
			if len(checkedPath) == 1 {
				resultPath = append(resultPath, checkedPath[k])
				continue
			}

			if checkedPath[k][0].IsEnd {
				resultPath = append(resultPath, checkedPath[k])
				break
			}

			if indexForPosTable == -1 {
				minimumPath = checkedPath[k]
				indexForPosTable = k
				continue
			}

			if antPosTable[indexForPosTable]+len(minimumPath) > len(checkedPath[k])+antPosTable[k] {
				minimumPath = checkedPath[k]
				indexForPosTable = k
			}

			if k == len(checkedPath)-1 {
				antPosTable[indexForPosTable]++
				resultPath = append(resultPath, minimumPath)
			}
		}

	}

	var stepsCounter int
	for i := 0; i < len(checkedPath); i++ {
		if i == 0 {
			stepsCounter = len(checkedPath[i]) + antPosTable[i]
			continue
		}

		if len(checkedPath[i])+antPosTable[i] > stepsCounter {
			stepsCounter = len(checkedPath[i]) + antPosTable[i]
		}
	}

	// for b := 0; b < len(resultPath); b++ {
	// 	for c := 0; c < len(resultPath[b]); c++ {
	// 		fmt.Print(resultPath[b][c].Name)

	// 		if c != len(resultPath[b])-1 {
	// 			fmt.Print(" -> ")
	// 		} else {
	// 			fmt.Println()
	// 		}
	// 	}
	// }
	res := Food(resultPath, startRoom, endRoom)
	return resultPath, res, len(checkedPath)
}

func FindBestPaths(startRoom *Room, endRoom *Room, paths []PathStuct, pathToCheck [][]*Room, index int, maxLinks int) {
	for i := 0; i < len(paths[index].Paths); i++ {
		pathToCheck[index] = paths[index].Paths[i]

		if maxLinks-1 == index {
			// fmt.Println()
			// for b := 0; b < len(pathToCheck); b++ {
			// 	for c := 0; c < len(pathToCheck[b]); c++ {
			// 		fmt.Print(pathToCheck[b][c].Name)

			// 		if c != len(pathToCheck[b])-1 {
			// 			fmt.Print(" -> ")
			// 		} else {
			// 			fmt.Println()
			// 		}
			// 	}
			// }

			possibleBestPath, possibleBestSteps, pathsInUse := Foo(pathToCheck, startRoom, endRoom)

			if BESTPATH == nil || possibleBestSteps < RESULT {
				BESTPATH = possibleBestPath
				RESULT = possibleBestSteps
				PATHSINUSE = pathsInUse
			}

			if possibleBestSteps == RESULT {
				if pathsInUse < PATHSINUSE {
					BESTPATH = possibleBestPath
					RESULT = possibleBestSteps
				}
			}

		} else {
			FindBestPaths(startRoom, endRoom, paths, pathToCheck, index+1, maxLinks)
		}
	}
}

func PotentialPaths(paths [][]*Room, startRoom Room) []PathStuct {
	sortedPaths := make([]PathStuct, len(startRoom.Links))
	for k := 0; k < len(startRoom.Links); k++ {
		var amount int
		var pathToAppend [][]*Room
		for i := 0; i < len(paths); i++ {
			if amount == 3 {
				break
			}

			if paths[i][0].Name == startRoom.Links[k].Name {
				pathToAppend = append(pathToAppend, paths[i])
				amount++
			}
		}

		sortedPaths[k].PathName = startRoom.Links[k].Name
		sortedPaths[k].Paths = append(sortedPaths[k].Paths, pathToAppend...)
	}

	for i := 0; i < len(sortedPaths); i++ {
		SortPaths(&sortedPaths[i].Paths)
	}

	// for a := 0; a < len(sortedPaths); a++ {
	// 	for i := 0; i < len(sortedPaths[a].Paths); i++ {
	// 		for k := 0; k < len(sortedPaths[a].Paths[i]); k++ {
	// 			fmt.Print(sortedPaths[a].Paths[i][k].Name)

	// 			if k != len(sortedPaths[a].Paths[i])-1 {
	// 				fmt.Print(" -> ")
	// 			} else {
	// 				fmt.Println()
	// 			}
	// 		}
	// 	}
	// }

	return sortedPaths
}

func Food(pathToTest [][]*Room, startRoom *Room, endRoom *Room) int {
	ants := SpawnAnts(startRoom, pathToTest)

	var sum int = 0

	endRoom.Ants = 0
	startRoom.Ants = ANTCOUNTER

	Asdf(ants, startRoom, endRoom, &sum)
	return sum
}

func Asdf(ants []Ant, startRoom *Room, endRoom *Room, sum *int) {
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
	}

	if allPassed && endRoom.Ants == ANTCOUNTER {
		return
	} else {
		*sum++
		Asdf(ants, startRoom, endRoom, sum)
	}
}

func RestoreFarm(farm *[]Room) {
	for i := 0; i < len(*farm); i++ {
		if (*farm)[i].IsStart {
			(*farm)[i].Ants = ANTCOUNTER
		} else {
			(*farm)[i].Ants = 0
		}
	}
}

func SortPaths(paths *[][]*Room) {
	var n = len(*paths)
	for i := 0; i < n; i++ {
		var minIdx = i
		for j := i; j < n; j++ {
			if len((*paths)[j]) < len((*paths)[minIdx]) {
				minIdx = j
			}
		}
		(*paths)[i], (*paths)[minIdx] = (*paths)[minIdx], (*paths)[i]
	}
}

func SpawnAnts(startRoom *Room, paths [][]*Room) []Ant {
	var result []Ant
	for i := 1; i < ANTCOUNTER+1; i++ {
		var antToAppend Ant

		antToAppend.Id = i
		antToAppend.CurrentRoom = startRoom
		antToAppend.Path = paths[i-1]

		result = append(result, antToAppend)
	}

	return result
}

func MakeStep(ants []Ant, endRoom *Room) {
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

	if allPassed && endRoom.Ants == ANTCOUNTER {
		return
	} else {
		fmt.Println("")
		MakeStep(ants, endRoom)
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
			ENDROOMID = i
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
						(farm)[m].Links = append((farm)[m].Links, &(farm)[k])
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
				// if len(path) < len((*paths)[k]) {
				// 	(*paths)[k] = nil
				// 	break
				// } else {
				// 	skipPath = true
				// 	break
				// }
			}
		}

		for i := 0; i < len(path); i++ {
			if path[i].IsStart {
				skipPath = true
				break
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
