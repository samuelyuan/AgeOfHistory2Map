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

func parseJsonFile(inputFilename string) ([]byte, error) {
	inputFile, err := os.Open(inputFilename)
	defer inputFile.Close()
	if err != nil {
		log.Fatal("Failed to load map: ", err)
		return []byte{}, err
	}
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

func loadAllProvinces() ([][]ProvinceGameData, int, int) {
	files, err := os.ReadDir(provinceDir)
	if err != nil {
		log.Fatalf("Failed to read provinces directory '%s': %v\nMake sure the data directory structure is correct.", provinceDir, err)
	}
	maxProvinces := len(files)
	fmt.Println("Number of provinces:", maxProvinces)

	globalMaxX := 0
	globalMaxY := 0
	fmt.Println("Loading provinces...")

	allProvinceData := make([][]ProvinceGameData, maxProvinces)
	for i := 0; i < len(allProvinceData); i++ {
		provinceFileName := fmt.Sprintf("%v/%v", provinceDir, i)
		jsonBytes, err := parseJsonFile(provinceFileName)
		if err != nil {
			log.Fatal("Failed to read input file: ", err)
		}

		var provinces []ProvinceGameData
		err = json.Unmarshal(jsonBytes, &provinces)
		if err != nil {
			fmt.Println("error:", err)
		}
		fmt.Println("Province", i, ":", provinces[0].ProvinceInfo)
		allProvinceData[i] = provinces

		// Scale all points
		for p := 0; p < len(allProvinceData[i]); p++ {
			province := allProvinceData[i][p]
			for j := 0; j < len(province.LPointsX); j++ {
				allProvinceData[i][p].LPointsX[j] *= imageScale
				allProvinceData[i][p].LPointsY[j] *= imageScale
			}
		}

		for p := 0; p < len(allProvinceData[i]); p++ {
			province := allProvinceData[i][p]
			for j := 0; j < len(province.LPointsX); j++ {
				currentIndex := j
				if int(province.LPointsX[currentIndex]) > globalMaxX {
					globalMaxX = int(province.LPointsX[currentIndex])
				}
				if int(province.LPointsY[currentIndex]) > globalMaxY {
					globalMaxY = int(province.LPointsY[currentIndex])
				}
			}
		}
	}
	return allProvinceData, globalMaxX, globalMaxY
}

func loadRegionsMap() RegionsMapData {
	regionListBytes, err := parseJsonFile(regionListFile)
	if err != nil {
		log.Fatal("Failed to read input file: ", err)
	}
	var regionFilenames []RegionList
	err = json.Unmarshal(regionListBytes, &regionFilenames)
	if err != nil {
		fmt.Println("Error:", err)
	}
	fmt.Println("Region filenames:", regionFilenames)

	numberRegions := len(regionFilenames[0].LRegionsTags)
	fmt.Println("Number regions in regions file:", numberRegions)

	allRegionColors := make([]RegionColor, numberRegions)

	for i := 0; i < numberRegions; i++ {
		regionDataBytes, err := parseJsonFile(fmt.Sprintf(regionFile, regionFilenames[0].LRegionsTags[i]))

		if err != nil {
			log.Fatal("Failed to read input file: ", err)
		}
		var regionData []RegionColor
		err = json.Unmarshal(regionDataBytes, &regionData)
		if err != nil {
			fmt.Println("Error:", err)
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

func loadSavedProvincesData(saveFolder string, allProvinceOwners []int) SaveDataOutput {
	saveDataBytes, err := parseJsonFile(fmt.Sprintf(saveDataPath, saveFolder, saveFolder))
	if err != nil {
		log.Fatal("Failed to read input file: ", err)
	}

	var saveDataProvinces []SaveDataProvinces
	err = json.Unmarshal(saveDataBytes, &saveDataProvinces)
	if err != nil {
		fmt.Println("Error:", err)
	}
	fmt.Println("saveDataProvinces:", saveDataProvinces)

	civEconomyMap := make(map[int]int)
	for i := 0; i < len(allProvinceOwners); i++ {
		provinceOwner := allProvinceOwners[i]
		mapKey := provinceOwner - 1
		_, ok := civEconomyMap[mapKey]
		if !ok {
			civEconomyMap[mapKey] = 0
		}
		civEconomyMap[mapKey] += saveDataProvinces[0].LProvincesData[i].IEconomy
	}

	fmt.Println("Civ economy map:", civEconomyMap)
	return SaveDataOutput{
		CivEconomyMap: civEconomyMap,
	}
}

func loadScenario(scenario string) Scenario {
	scenarioDataBytes, err := parseJsonFile(fmt.Sprintf(scenarioDataPath, scenario, scenario))
	if err != nil {
		log.Fatal("Failed to read input file: ", err)
	}
	fmt.Println("Sceario data raw:", string(scenarioDataBytes))

	var scenarioData []ScenarioData
	err = json.Unmarshal(scenarioDataBytes, &scenarioData)
	if err != nil {
		fmt.Println("Error:", err)
	}
	fmt.Println("Scenario data:", scenarioData)

	provinceOwnersBytes, err := parseJsonFile(fmt.Sprintf(scenarioProvincePath, scenario, scenario))
	if err != nil {
		log.Fatal("Failed to read input file: ", err)
	}

	var provinceOwners []ProvinceOwners
	err = json.Unmarshal(provinceOwnersBytes, &provinceOwners)
	if err != nil {
		fmt.Println("Error:", err)
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
			tagEndings := []string{"_r", "_c", "_m", "_s", "_h", "_t"}
			for j := 0; j < len(tagEndings); j++ {
				if strings.Contains(civTag, tagEndings[j]) {
					civTag = civTag[0:strings.Index(civTag, tagEndings[j])]
					fmt.Println(fmt.Sprintf("Removing '%v' from civ %v tag", tagEndings[j], i))
				}
			}
		}

		civilizationDataBytes, err := parseJsonFile(fmt.Sprintf(civilizationPath, civTag))
		if err != nil {
			log.Fatal("Failed to read input file: ", err)
		}
		var civilizationColor []CivilizationColor
		err = json.Unmarshal(civilizationDataBytes, &civilizationColor)
		if err != nil {
			fmt.Println("Error:", err)
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
