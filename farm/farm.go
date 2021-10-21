package farm

import (
	"lem-in/structs"
	"lem-in/utils"
	"log"
	"strconv"
	"strings"
)

// Generates farm based on generation data
func GenerateFarm(data structs.GenerationData) {
	var err error
	for i := 0; i < len(data.Rooms); i++ {
		var roomToAppend structs.Room

		if i == data.StartIndex {
			roomToAppend.IsStart = true
			structs.STARTROOMID = i
			roomToAppend.Ants = structs.ANTCOUNTER
		} else if i == data.EndIndex {
			roomToAppend.IsEnd = true
			structs.ENDROOMID = i
		}

		splittedData := strings.Split(data.Rooms[i], " ")

		roomToAppend.Name = splittedData[0]
		roomToAppend.X_pos, err = strconv.Atoi(splittedData[1])
		utils.CheckError(err)
		roomToAppend.Y_pos, err = strconv.Atoi(splittedData[2])
		utils.CheckError(err)

		structs.FARM = append(structs.FARM, roomToAppend)
	}

	ConnectLinks(data.Links)
}

// Connects all rooms based on links from the file
func ConnectLinks(links []string) {
	for i := 0; i < len(links); i++ {
		splittedData := strings.Split(links[i], "-")
		for k := 0; k < len(structs.FARM); k++ {
			if structs.FARM[k].Name == splittedData[0] {
				for m := 0; m < len(structs.FARM); m++ {
					if structs.FARM[m].Name == splittedData[1] {
						if structs.FARM[m].Name == structs.FARM[k].Name {
							log.Fatal("Invalid data format! Self-link is prohibited")
						}
						structs.FARM[k].Links = append(structs.FARM[k].Links, &structs.FARM[m])
						structs.FARM[m].Links = append(structs.FARM[m].Links, &structs.FARM[k])
						break
					}
					if m == len(structs.FARM)-1 {
						log.Fatal("Invalid data format! Room link not found")
					}
				}
				break
			}
			if k == len(structs.FARM)-1 {
				log.Fatal("Invalid data format! Room link not found")
			}
		}
	}
}
