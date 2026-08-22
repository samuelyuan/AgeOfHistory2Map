package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/samuelyuan/AgeOfHistory2Map/jserial"
)

const (
	// Base data directory
	dataDir = "data"

	// Directory paths
	provinceDir = dataDir + "/map/Earth/data/provinces"

	// Region file paths
	regionListFile = dataDir + "/map/data/regions/packges/Earth_AoC2"
	regionFile     = dataDir + "/map/data/regions/packges_data/%v"

	// Scenario file paths
	scenarioDataPath     = dataDir + "/map/Earth/scenarios/%v/%v"
	scenarioProvincePath = dataDir + "/map/Earth/scenarios/%v/%v_PD"

	// Civilization file paths
	civilizationPath = dataDir + "/game/civilizations/%v"

	// Save file paths
	saveDataPath = dataDir + "/saves/%v/%v_4"

	// Image settings
	imageScale = 1
)

type ProvinceBorderGameData struct {
	WithProvinceID int   `json:"withProvinceID"`
	LPointsX       []int `json:"lPointsX"`
	LPointsY       []int `json:"lPointsY"`
}

type ProvinceInfo struct {
	FGrowthRate  float32 `json:"fGrowthRate"`
	STerrainTAG  string  `json:"sTerrainTAG"`
	IContinentID int     `json:"iContinentID"`
	IRegionID    int     `json:"iRegionID"`
	IShiftX      int     `json:"iShiftX"`
	IShiftY      int     `json:"iShiftY"`
}

type ProvinceGameData struct {
	LPointsX               []int                    `json:"lPointsX"`
	LPointsY               []int                    `json:"lPointsY"`
	ProvinceBorderGameData []ProvinceBorderGameData `json:"lProvinceBorder"`
	ProvinceInfo           ProvinceInfo             `json:"provinceInfo"`
}

type RegionList struct {
	LRegionsTags []string `json:"lRegionsTags"`
	SPackageName string   `json:"sPackageName"`
}

type RegionColor struct {
	FractionRed   float64 `json:"fR"`
	FractionGreen float64 `json:"fG"`
	FractionBlue  float64 `json:"fB"`
	SName         string  `json:"sName"`
}

type ScenarioData struct {
	LCivsTags []string `json:"lCivsTags"`
}

type CivilizationColor struct {
	IRed    int    `json:"iR"`
	IGreen  int    `json:"iG"`
	IBlue   int    `json:"iB"`
	SCivTag string `json:"sCivTag"`
}

type ProvinceOwners struct {
	LProvinceOwners []int `json:"lProvinceOwners"`
}

type SaveDataProvinces struct {
	LProvincesData []SaveProvinceInfo `json:"lProvincesData"`
}

type SaveProvinceInfo struct {
	IEconomy int `json:"iEconomy"`
}

type SaveDataOutput struct {
	CivEconomyMap map[int]int
}

type RegionsMapData struct {
	GlobalMaxX      int
	GlobalMaxY      int
	AllProvinceData [][]ProvinceGameData
	AllRegionColors []RegionColor
}

type Scenario struct {
	GlobalMaxX        int
	GlobalMaxY        int
	AllProvinceData   [][]ProvinceGameData
	AllProvinceOwners []int
	AllCivColors      []CivilizationColor
}

// checkDataDirectory verifies that dir exists and is a directory, returning a friendly,
// actionable error otherwise. Called once up front so a missing game data folder produces one
// clear message instead of failing deep inside a specific file loader.
func checkDataDirectory(dir string) error {
	info, err := os.Stat(dir)
	if err != nil || !info.IsDir() {
		return fmt.Errorf(
			"could not find a '%s' folder in the current directory.\n"+
				"Copy your Age of History 2 game data into a folder named '%s' next to this program, then run it again.\n"+
				"See the README for the expected folder layout (data/map/..., data/game/..., data/saves/...)",
			dir, dir)
	}
	return nil
}

func parseJsonFile(inputFilename string) ([]byte, error) {
	inputFile, err := os.Open(inputFilename)
	if err != nil {
		log.Fatal("Failed to load map: ", err)
		return []byte{}, err
	}
	defer inputFile.Close()

	fi, err := inputFile.Stat()
	if err != nil {
		log.Fatal(err)
		return []byte{}, err
	}
	fileLength := fi.Size()
	streamReader := io.NewSectionReader(inputFile, int64(0), fileLength)

	sop := jserial.NewSerializedObjectParser(streamReader)

	objects, err := sop.ParseSerializedObjectMinimal()
	if err != nil {
		log.Fatalf("%+v", err)
	}

	jsonBytes, err := json.MarshalIndent(objects, "", "    ")
	if err != nil {
		log.Fatalf("%+v", err)
	}

	// fmt.Println(string(jsonBytes))
	return jsonBytes, nil
}

// loadJsonArray reads and parses a data file that decodes to a JSON array of T.
// It centralizes the parse-then-unmarshal pattern shared by every loader below.
func loadJsonArray[T any](filename string) ([]T, error) {
	jsonBytes, err := parseJsonFile(filename)
	if err != nil {
		return nil, err
	}

	var items []T
	if err := json.Unmarshal(jsonBytes, &items); err != nil {
		return nil, fmt.Errorf("failed to unmarshal %s: %w", filename, err)
	}
	return items, nil
}

// scaleProvincePoints scales every point of every province in place by the given factor.
func scaleProvincePoints(provinces []ProvinceGameData, scale int) {
	for p := range provinces {
		for j := range provinces[p].LPointsX {
			provinces[p].LPointsX[j] *= scale
			provinces[p].LPointsY[j] *= scale
		}
	}
}

// computeGlobalBounds returns the maximum X and Y coordinate across all province points.
func computeGlobalBounds(allProvinceData [][]ProvinceGameData) (globalMaxX int, globalMaxY int) {
	for _, provinces := range allProvinceData {
		for _, province := range provinces {
			for j := range province.LPointsX {
				if province.LPointsX[j] > globalMaxX {
					globalMaxX = province.LPointsX[j]
				}
				if province.LPointsY[j] > globalMaxY {
					globalMaxY = province.LPointsY[j]
				}
			}
		}
	}
	return globalMaxX, globalMaxY
}

func loadAllProvinces() ([][]ProvinceGameData, int, int) {
	files, err := os.ReadDir(provinceDir)
	if err != nil {
		log.Fatalf("Failed to read provinces directory '%s': %v\nMake sure the data directory structure is correct.", provinceDir, err)
	}
	maxProvinces := len(files)
	fmt.Println("Number of provinces:", maxProvinces)
	fmt.Println("Loading provinces...")

	allProvinceData := make([][]ProvinceGameData, maxProvinces)
	for i := 0; i < len(allProvinceData); i++ {
		provinceFileName := fmt.Sprintf("%v/%v", provinceDir, i)
		provinces, err := loadJsonArray[ProvinceGameData](provinceFileName)
		if err != nil {
			log.Fatal("Failed to read input file: ", err)
		}
		fmt.Println("Province", i, ":", provinces[0].ProvinceInfo)

		scaleProvincePoints(provinces, imageScale)
		allProvinceData[i] = provinces
	}

	globalMaxX, globalMaxY := computeGlobalBounds(allProvinceData)
	return allProvinceData, globalMaxX, globalMaxY
}

func loadRegionsMap() RegionsMapData {
	regionFilenames, err := loadJsonArray[RegionList](regionListFile)
	if err != nil {
		log.Fatal("Failed to read input file: ", err)
	}
	fmt.Println("Region filenames:", regionFilenames)

	numberRegions := len(regionFilenames[0].LRegionsTags)
	fmt.Println("Number regions in regions file:", numberRegions)

	allRegionColors := make([]RegionColor, numberRegions)
	for i := 0; i < numberRegions; i++ {
		regionData, err := loadJsonArray[RegionColor](fmt.Sprintf(regionFile, regionFilenames[0].LRegionsTags[i]))
		if err != nil {
			log.Fatal("Failed to read input file: ", err)
		}
		fmt.Println("Region", i, "data:", regionData)
		allRegionColors[i] = regionData[0]
	}

	allProvinceData, globalMaxX, globalMaxY := loadAllProvinces()
	fmt.Println("Max x:", globalMaxX, ", max y:", globalMaxY)

	return RegionsMapData{
		GlobalMaxX:      globalMaxX,
		GlobalMaxY:      globalMaxY,
		AllProvinceData: allProvinceData,
		AllRegionColors: allRegionColors,
	}
}

// computeCivEconomyMap sums each province's economy value into its owning civilization's total
// (keyed by owner id - 1, matching the save format). Provinces beyond the bounds of provincesData
// are skipped rather than panicking. Pure: safe to unit test without any save files on disk.
func computeCivEconomyMap(allProvinceOwners []int, provincesData []SaveProvinceInfo) map[int]int {
	civEconomyMap := make(map[int]int)
	for i, provinceOwner := range allProvinceOwners {
		if i >= len(provincesData) {
			continue
		}
		mapKey := provinceOwner - 1
		civEconomyMap[mapKey] += provincesData[i].IEconomy
	}
	return civEconomyMap
}

func loadSavedProvincesData(saveFolder string, allProvinceOwners []int) SaveDataOutput {
	saveDataProvinces, err := loadJsonArray[SaveDataProvinces](fmt.Sprintf(saveDataPath, saveFolder, saveFolder))
	if err != nil {
		log.Fatal("Failed to read input file: ", err)
	}
	fmt.Println("saveDataProvinces:", saveDataProvinces)

	civEconomyMap := computeCivEconomyMap(allProvinceOwners, saveDataProvinces[0].LProvincesData)
	fmt.Println("Civ economy map:", civEconomyMap)
	return SaveDataOutput{
		CivEconomyMap: civEconomyMap,
	}
}

// knownCivTagSuffixes lists the suffixes stripped from a civ tag when its data file is missing,
// e.g. "french_r" falls back to the base civ "french".
var knownCivTagSuffixes = []string{"_r", "_c", "_m", "_s", "_h", "_t"}

// stripKnownCivTagSuffix removes the first matching known suffix from a civ tag, if any.
// Pure: separated from the filesystem check in loadScenario so the fallback logic can be
// unit tested directly.
func stripKnownCivTagSuffix(civTag string) (stripped string, suffixRemoved string, ok bool) {
	for _, suffix := range knownCivTagSuffixes {
		if idx := strings.Index(civTag, suffix); idx != -1 {
			return civTag[:idx], suffix, true
		}
	}
	return civTag, "", false
}

func loadScenario(scenario string) Scenario {
	scenarioData, err := loadJsonArray[ScenarioData](fmt.Sprintf(scenarioDataPath, scenario, scenario))
	if err != nil {
		log.Fatal("Failed to read input file: ", err)
	}
	fmt.Println("Scenario data:", scenarioData)

	provinceOwners, err := loadJsonArray[ProvinceOwners](fmt.Sprintf(scenarioProvincePath, scenario, scenario))
	if err != nil {
		log.Fatal("Failed to read input file: ", err)
	}
	allProvinceOwners := provinceOwners[0].LProvinceOwners
	fmt.Println("Province owners:", allProvinceOwners)

	numCivs := len(scenarioData[0].LCivsTags)
	fmt.Println("Number of civilizations:", numCivs)

	allCivColors := make([]CivilizationColor, numCivs)
	for i := 0; i < numCivs; i++ {
		civTag := scenarioData[0].LCivsTags[i]

		if _, err := os.Stat(fmt.Sprintf(civilizationPath, civTag)); errors.Is(err, os.ErrNotExist) {
			fmt.Println("File for civ tag", civTag, "doesn't exist")
			if stripped, suffix, ok := stripKnownCivTagSuffix(civTag); ok {
				fmt.Printf("Removing '%v' from civ %v tag\n", suffix, i)
				civTag = stripped
			}
		}

		civilizationColor, err := loadJsonArray[CivilizationColor](fmt.Sprintf(civilizationPath, civTag))
		if err != nil {
			log.Fatal("Failed to read input file: ", err)
		}
		fmt.Println("Civilization color:", civilizationColor)
		allCivColors[i] = civilizationColor[0]
	}

	allProvinceData, globalMaxX, globalMaxY := loadAllProvinces()
	fmt.Println("Max x:", globalMaxX, ", max y:", globalMaxY)

	return Scenario{
		GlobalMaxX:        globalMaxX,
		GlobalMaxY:        globalMaxY,
		AllProvinceData:   allProvinceData,
		AllProvinceOwners: allProvinceOwners,
		AllCivColors:      allCivColors,
	}
}
