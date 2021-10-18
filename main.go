package main

import (
	"bytes"
	"log"
	"os"
	"strconv"
	"strings"
)

type Room struct {
	Name    string
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

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Missing file name")
		os.Exit(1)
	}

	data := LoadData(os.Args[1])
	generationData := ReadData(data)
	farm := GenerateFarm(generationData)

	_ = farm
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

	var commentsCounter int = 1
	for i := 1; i < len(data); i++ {
		if strings.Contains(string(data[i]), "##") {
			if string(data[i]) == "##start" {
				result.StartIndex = i - commentsCounter
			} else if string(data[i]) == "##end" {
				result.EndIndex = i - commentsCounter
			} else {
				log.Fatal("Invalid data format!")
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
			log.Fatal("Invalid data format!")
		}
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
						(farm)[k].Links = append((farm)[k].Links, &(farm)[m])
						break
					}

					if m == len(farm)-1 {
						log.Fatal("Invalid data format!")
					}
				}
				break
			}

			if k == len(farm)-1 {
				log.Fatal("Invalid data format!")
			}
		}
	}

	return farm
}

func CheckError(err error) {
	if err != nil {
		log.Fatal("Invalid data format!")
	}
}
