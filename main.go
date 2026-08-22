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

// defaultOutputFilename returns the output filename to use when the user doesn't pass -output,
// so e.g. `-map modernworld` alone produces "modern.png" instead of a generic "output.png". Pure.
func defaultOutputFilename(mapName string) string {
	if mapName == "modernworld" {
		return "modern.png"
	}
	return mapName + ".png"
}

func printUsage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [options]\n\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "Options:\n")
	flag.PrintDefaults()
	fmt.Fprintf(os.Stderr, "\nThis program renders maps from Age of History 2 game data. It expects a '%s' folder\n", dataDir)
	fmt.Fprintf(os.Stderr, "with a copy of the game's data files to exist in the current directory (see the README\n")
	fmt.Fprintf(os.Stderr, "for the expected folder layout).\n")
	fmt.Fprintf(os.Stderr, "\nExamples:\n")
	for _, mapName := range validMapNamesList {
		fmt.Fprintf(os.Stderr, "  %s -map %s -output %s\n", os.Args[0], mapName, defaultOutputFilename(mapName))
	}
}

func main() {
	outputPtr := flag.String("output", "", "Output filename (default: \"<map name>.png\")")
	mapPtr := flag.String("map", "regions", "Map name "+supportedMapNames)
	flag.Usage = printUsage
	flag.Parse()

	mapName := *mapPtr

	// Validate map name
	if !validMapNames[mapName] {
		fmt.Fprintf(os.Stderr, "Error: Invalid map name '%s'\n\n", mapName)
		fmt.Fprintf(os.Stderr, "Supported map names: %s\n\n", supportedMapNames)
		printUsage()
		os.Exit(1)
	}

	outputFilename := *outputPtr
	if outputFilename == "" {
		outputFilename = defaultOutputFilename(mapName)
	}

	// Check for the game data folder up front, so a missing/misplaced data/ directory produces
	// one clear, actionable message instead of a confusing failure deep inside file loading.
	if err := checkDataDirectory(dataDir); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		fmt.Fprintln(os.Stderr)
		printUsage()
		os.Exit(1)
	}

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
