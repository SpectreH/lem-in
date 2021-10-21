package utils

import (
	"fmt"
	"lem-in/structs"
	"log"
)

func CheckError(err error) {
	if err != nil {
		log.Fatal("Invalid data format!")
	}
}

// Restores farm to defaults, when MakeStep in mode 1 is ended
func RestoreFarm() {
	for i := 0; i < len(structs.FARM); i++ {
		if (structs.FARM)[i].IsStart {
			(structs.FARM)[i].Ants = structs.ANTCOUNTER
		} else {
			(structs.FARM)[i].Ants = 0
		}
	}
}

// Sorts paths by lenght
func SortPaths(paths *[][]*structs.Room) {
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

// Supporting function. Decodes paths list into room names
func PrintPaths(paths [][]*structs.Room) {
	fmt.Println()
	for b := 0; b < len(paths); b++ {
		for c := 0; c < len(paths[b]); c++ {
			fmt.Print(paths[b][c].Name)
			if c != len(paths[b])-1 {
				fmt.Print(" -> ")
			} else {
				fmt.Println()
			}
		}
	}
}
