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

// Makes ants (from ants list) to move from start to end prints each ant step.
func MakeStep(ants []structs.Ant) {
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

	if allPassed && structs.FARM[structs.ENDROOMID].Ants == structs.ANTCOUNTER {
		return
	} else {
		fmt.Println("")
		MakeStep(ants)
	}
}
