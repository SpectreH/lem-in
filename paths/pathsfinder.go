package paths

import (
	"lem-in/structs"
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

// Finds all possible combinations without intersects
func FindCombinations(paths [][]*structs.Room) [][][]*structs.Room {
	var result [][][]*structs.Room

	for i := 0; i < len(paths); i++ {
		var combination [][]*structs.Room
		combination = append(combination, paths[i])
		for k := i + 1; k < len(paths); k++ {
			if !ExcludeIntersect(paths[i][:len(paths[i])-1], paths[k][:len(paths[k])-1]) &&
				!ExcludeIntersectInsideComb(paths[k][:len(paths[k])-1], combination) {
				result = append(result, combination)
				combination = append(combination, paths[k])
			}
		}

		result = append(result, combination)
	}

	return result
}

// Checks intersects between two paths
func ExcludeIntersect(currentPath, pathToCheck []*structs.Room) bool {
	for i := 0; i < len(currentPath); i++ {
		for k := i + 1; k < len(pathToCheck); k++ {
			if currentPath[i].Name == pathToCheck[k].Name {
				return true
			}
		}
	}
	return false
}

// Checks intersects between existing combination and path
func ExcludeIntersectInsideComb(path []*structs.Room, combination [][]*structs.Room) bool {
	for i := 0; i < len(combination); i++ {
		for k := 0; k < len(path); k++ {
			for j := 0; j < len(combination[i]); j++ {
				if path[k] == combination[i][j] {
					return true
				}
			}
		}
	}
	return false
}

// Finds best combination between all combinaitons and calculates best suitable path for each ant
func FindBestComb(c [][][]*structs.Room) [][]*structs.Room {
	var bestScore int
	var bestPath [][]*structs.Room

	for _, paths := range c {
		var pathComb [][]*structs.Room

		antPosTable := make([]int, len(paths))
		var currentIndex = 0
		var nextPathId int
		var updateNextPathId bool = true
		for i := 0; i < structs.ANTCOUNTER; i++ {
			if i == 0 {
				pathComb = append(pathComb, paths[0])
				currentIndex = 0
				antPosTable[currentIndex]++
				continue
			}
			for {
				if updateNextPathId {
					if len(paths) == currentIndex+1 {
						nextPathId = 0
					} else {
						nextPathId = currentIndex + 1
					}
					updateNextPathId = false
				}

				if len(paths) == 1 || paths[currentIndex][0].IsEnd {
					pathComb = append(pathComb, paths[currentIndex])
					antPosTable[currentIndex]++
					break
				}

				if antPosTable[currentIndex]+len(paths[currentIndex]) <= len(paths[nextPathId])+antPosTable[nextPathId] {
					pathComb = append(pathComb, paths[currentIndex])
					antPosTable[currentIndex]++
					break
				} else {
					pathComb = append(pathComb, paths[nextPathId])
					antPosTable[nextPathId]++
					currentIndex = nextPathId
					updateNextPathId = true
					break
				}
			}
		}

		var currentScore int
		for i := 0; i < len(paths); i++ {
			temp := antPosTable[i] + len(paths[i])
			if currentScore == 0 || temp > currentScore {
				currentScore = temp
			}
		}

		if bestScore == 0 || currentScore < bestScore {
			bestScore = currentScore
			bestPath = pathComb
		}
	}

	return bestPath
}
