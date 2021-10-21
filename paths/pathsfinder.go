package paths

import (
	"lem-in/ants"
	"lem-in/structs"
	"lem-in/utils"
)

// Finds all paths from start to end
func FindAllPossiblePaths(path []*structs.Room, currentRoom structs.Room, step int, paths *[][]*structs.Room, previousRoom *structs.Room) {
	if currentRoom.IsEnd {
		var skipPath bool
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
			FindAllPossiblePaths(pathToPass, *currentRoom.Links[i], step+1, paths, &currentRoom)
			pathToPass = path
		}
	}

	for i := 0; i < len(*paths); i++ {
		if (*paths)[i] == nil {
			*paths = append((*paths)[:i], (*paths)[i+1:]...)
		}
	}
}

// Sturcturize all paths by links in start room and saves only 3 shortest paths for each link in start room
func StructurizePaths(paths [][]*structs.Room) []structs.PathStuct {
	structurizedPaths := make([]structs.PathStuct, len(structs.FARM[structs.STARTROOMID].Links))
	for k := 0; k < len(structs.FARM[structs.STARTROOMID].Links); k++ {
		var amount int
		var pathToAppend [][]*structs.Room
		for i := 0; i < len(paths); i++ {
			if amount == 3 {
				break
			}

			if paths[i][0].Name == structs.FARM[structs.STARTROOMID].Links[k].Name {
				pathToAppend = append(pathToAppend, paths[i])
				amount++
			}
		}

		structurizedPaths[k].PathName = structs.FARM[structs.STARTROOMID].Links[k].Name
		structurizedPaths[k].Paths = append(structurizedPaths[k].Paths, pathToAppend...)
	}

	for i := 0; i < len(structurizedPaths); i++ {
		utils.SortPaths(&structurizedPaths[i].Paths)
	}

	return structurizedPaths
}

// Parses best paths combinations and saves best one
func FindBestPathsComb(paths []structs.PathStuct, pathToCheck [][]*structs.Room, index int, maxLinks int) {
	for i := 0; i < len(paths[index].Paths); i++ {
		pathToCheck[index] = paths[index].Paths[i]

		if maxLinks-1 == index {
			possibleBestPath, possibleBestSteps, pathsInUse := TryPathsComb(pathToCheck)
			utils.RestoreFarm()

			if structs.BEST_PATH == nil || possibleBestSteps < structs.BEST_TURNS_RES {
				structs.BEST_PATH = possibleBestPath
				structs.BEST_TURNS_RES = possibleBestSteps
				structs.BEST_ROOMS_IN_USE_RES = pathsInUse
			}

			if possibleBestSteps == structs.BEST_TURNS_RES {
				if pathsInUse < structs.BEST_ROOMS_IN_USE_RES {
					structs.BEST_PATH = possibleBestPath
					structs.BEST_TURNS_RES = possibleBestSteps
				}
			}
		} else {
			FindBestPathsComb(paths, pathToCheck, index+1, maxLinks)
		}
	}
}

// Tries best possible path combination
func TryPathsComb(paths [][]*structs.Room) ([][]*structs.Room, int, int) {
	var resultPath [][]*structs.Room
	checkedPath := ExcludeIntersections(paths)

	antPosTable := make([]int, len(checkedPath))
	for i := 0; i < structs.ANTCOUNTER; i++ {
		var minimumPath []*structs.Room
		var indexForPosTable int = -1

		for k := 0; k < len(checkedPath); k++ {
			if len(checkedPath) == 1 || checkedPath[k][0].IsEnd {
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

	antsList := ants.SpawnAnts(resultPath)
	var stepsCounter int = 0
	ants.MakeStep(antsList, 1, &stepsCounter)

	return resultPath, stepsCounter, len(checkedPath)
}

// Excludes that path from path combination, what can interrupt shortest path from this combination
func ExcludeIntersections(paths [][]*structs.Room) [][]*structs.Room {
	sortedByLenghtPaths := make([][]*structs.Room, 0)
	sortedByLenghtPaths = append(sortedByLenghtPaths, paths...)
	utils.SortPaths(&sortedByLenghtPaths)

	var roomsInUse []*structs.Room
	checkedPath := make([][]*structs.Room, 0)
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
	return checkedPath
}
