package main

import (
	"lem-in/ant"
	"lem-in/antfarm"
	"lem-in/parser"
	"lem-in/pathfinding"
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

	data := parser.LoadData(os.Args[1])
	generationData := parser.ReadData(data)
	antfarm.GenerateFarm(generationData)

	var allPaths [][]*structs.Room
	pathfinding.FindAllPossiblePaths(make([]*structs.Room, 0), structs.FARM[structs.STARTROOMID], 0, &allPaths, &structs.FARM[structs.STARTROOMID])
	utils.SortPaths(&allPaths)

	allCombinations := pathfinding.FindCombinations(allPaths)
	bestCombination := pathfinding.FindBestComb(allCombinations)

	antsList := ant.SpawnAnts(bestCombination)
	ant.MakeStep(antsList)
}
