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

	allCombinations := paths.FindCombinations(allPaths)
	bestCombination := paths.FindBestComb(allCombinations)

	antsList := ants.SpawnAnts(bestCombination)
	ants.MakeStep(antsList)
}
