package main

import (
	"flag"
	"fmt"
	"log"
)

const (
	supportedMapNames = "[regions, 1440, modernworld]"
)

func main() {
	inputPtr := flag.String("input", "", "Input file")
	outputPtr := flag.String("output", "output.png", "Output filename")
	mapPtr := flag.String("map", "regions", "Map name "+supportedMapNames)
	flag.Parse()

	fmt.Println("Input filename: ", *inputPtr)
	fmt.Println("Output filename: ", *outputPtr)
	fmt.Println("Map name: ", *mapPtr)

	outputFilename := *outputPtr
	mapName := *mapPtr
	if mapName == "regions" {
		regionsMapData := loadRegionsMap()
		drawRegionsMap(outputFilename, regionsMapData)
	} else if mapName == "1440" {
		scenario := loadScenario("1440")
		drawScenarioMap(outputFilename, scenario)
	} else if mapName == "modernworld" {
		scenario := loadScenario("modernworld")
		drawScenarioMap(outputFilename, scenario)
	} else {
		log.Fatal("Map name " + mapName + " is unsupported. Supported map names: " + supportedMapNames)
	}
}
