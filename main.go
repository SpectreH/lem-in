package main

import (
	"lem-in/ants"
	dataparser "lem-in/data-parser"
	"lem-in/farm"
	"lem-in/paths"
	"lem-in/structs"
	"lem-in/utils"
	"log"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Missing file name")
		os.Exit(1)
	}

	data := dataparser.LoadData(os.Args[1])
	generationData := dataparser.ReadData(data)
	farm.GenerateFarm(generationData)

	var allPaths [][]*structs.Room
	paths.FindAllPossiblePaths(make([]*structs.Room, 0), structs.FARM[structs.STARTROOMID], 0, &allPaths, &structs.FARM[structs.STARTROOMID])
	utils.SortPaths(&allPaths)

	potentialPaths := paths.StructurizePaths(allPaths)
	paths.FindBestPathsComb(potentialPaths, make([][]*structs.Room, len(structs.FARM[structs.STARTROOMID].Links)), 0, len(structs.FARM[structs.STARTROOMID].Links))

	antsList := ants.SpawnAnts(structs.BEST_PATH)
	var stepsCounter int
	ants.MakeStep(antsList, 0, &stepsCounter)
}
