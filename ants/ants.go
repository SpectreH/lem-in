package ants

import (
	"fmt"
	"lem-in/structs"
)

// Makes list of ants with own path, current room and id
func SpawnAnts(paths [][]*structs.Room) []structs.Ant {
	var result []structs.Ant
	for i := 1; i < structs.ANTCOUNTER+1; i++ {
		var antToAppend structs.Ant

		antToAppend.Id = i
		antToAppend.CurrentRoom = &structs.FARM[structs.STARTROOMID]
		antToAppend.Path = paths[i-1]

		result = append(result, antToAppend)
	}

	return result
}

// Makes ants (from ants list) to move from start to end. Mode 1 - to test path, Mode 0 - make final movement with best paths combinations and printing steps
func MakeStep(ants []structs.Ant, mode int, sum *int) {
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

		if mode == 0 {
			fmt.Print("L", ants[i].Id, "-", ants[i].CurrentRoom.Name, " ")
		}
	}

	if allPassed && structs.FARM[structs.ENDROOMID].Ants == structs.ANTCOUNTER {
		return
	} else {
		if mode == 0 {
			fmt.Println("")
		}

		*sum++
		if structs.BEST_TURNS_RES < *sum && structs.BEST_TURNS_RES != 0 {
			return
		}

		MakeStep(ants, mode, sum)
	}
}
