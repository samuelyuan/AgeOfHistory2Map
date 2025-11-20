package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

var (
	// Single source of truth for all valid map names
	validMapNamesList = []string{"regions", "1200", "1440", "modernworld"}

	// Derived data structures
	validMapNames     map[string]bool
	supportedMapNames string
)

func init() {
	// Build validMapNames map from the list
	validMapNames = make(map[string]bool, len(validMapNamesList))
	for _, name := range validMapNamesList {
		validMapNames[name] = true
	}

	// Build supportedMapNames string from the list
	supportedMapNames = "[" + strings.Join(validMapNamesList, ", ") + "]"
}

func printUsage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [options]\n\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "Options:\n")
	flag.PrintDefaults()
	fmt.Fprintf(os.Stderr, "\nExamples:\n")
	for _, mapName := range validMapNamesList {
		outputName := mapName + ".png"
		if mapName == "modernworld" {
			outputName = "modern.png"
		}
		fmt.Fprintf(os.Stderr, "  %s -map %s -output %s\n", os.Args[0], mapName, outputName)
	}
}

func main() {
	inputPtr := flag.String("input", "", "Input file (optional)")
	outputPtr := flag.String("output", "output.png", "Output filename")
	mapPtr := flag.String("map", "regions", "Map name "+supportedMapNames)
	flag.Usage = printUsage
	flag.Parse()

	outputFilename := *outputPtr
	mapName := *mapPtr

	// Validate map name
	if !validMapNames[mapName] {
		fmt.Fprintf(os.Stderr, "Error: Invalid map name '%s'\n\n", mapName)
		fmt.Fprintf(os.Stderr, "Supported map names: %s\n\n", supportedMapNames)
		printUsage()
		os.Exit(1)
	}

	fmt.Println("Input filename: ", *inputPtr)
	fmt.Println("Output filename: ", outputFilename)
	fmt.Println("Map name: ", mapName)

	if mapName == "regions" {
		regionsMapData := loadRegionsMap()
		drawRegionsMap(outputFilename, regionsMapData)
	} else {
		// All other map names are scenarios
		scenario := loadScenario(mapName)
		drawScenarioMap(outputFilename, scenario)
	}
}
