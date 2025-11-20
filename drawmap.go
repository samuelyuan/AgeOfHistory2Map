package main

import (
	"fmt"
	"math"
	"strconv"

	"github.com/fogleman/gg"
)

type Color struct {
	R, G, B int
}

var (
	// Water color (dark blue)
	waterColor = Color{R: 15, G: 27, B: 41}

	// Water color alternative (lighter blue, used in terrain view)
	waterColorAlt = Color{R: 47, G: 74, B: 93}

	// Default land color (green)
	landColor = Color{R: 105, G: 125, B: 54}

	// Unclaimed land color (dark gray)
	unclaimedLandColor = Color{R: 16, G: 16, B: 16}

	// Outline color (black)
	outlineColor = Color{R: 0, G: 0, B: 0}

	// Label color (white)
	labelColor = Color{R: 255, G: 255, B: 255}
)

func drawScenarioMap(outputFilename string, scenario Scenario) {
	dc := gg.NewContext(int(scenario.GlobalMaxX), int(scenario.GlobalMaxY))

	// water
	dc.SetRGB255(waterColor.R, waterColor.G, waterColor.B)
	dc.Clear()

	fmt.Println("Drawing map...")
	drawScenarioRegionColors(dc, scenario.AllProvinceData, scenario.AllProvinceOwners, scenario.AllCivColors)
	// drawProvinceOutline(dc, scenario.AllProvinceData)

	dc.SavePNG(outputFilename)
	fmt.Println("Saved image to", outputFilename)
}

func drawScenarioRegionColors(dc *gg.Context, allProvinceData [][]ProvinceGameData, allProvinceOwners []int, allCivColors []CivilizationColor) {
	for i := 0; i < len(allProvinceData); i++ {
		for p := 0; p < len(allProvinceData[i]); p++ {
			province := allProvinceData[i][p]

			dc.MoveTo(float64(province.LPointsX[0]), float64(province.LPointsY[0]))
			for j := 1; j < len(province.LPointsX); j++ {
				dc.LineTo(float64(province.LPointsX[j]), float64(province.LPointsY[j]))
			}
			dc.ClosePath()

			provinceOwner := allProvinceOwners[i] - 1
			if province.ProvinceInfo.IContinentID == 0 || province.ProvinceInfo.STerrainTAG == "" {
				// water
				dc.SetRGB255(waterColor.R, waterColor.G, waterColor.B)
			} else if provinceOwner < 0 || provinceOwner >= len(allCivColors) {
				// land doesn't belong to any owner
				fmt.Println("Province owner", provinceOwner, "isn't a valid province")
				dc.SetRGB255(unclaimedLandColor.R, unclaimedLandColor.G, unclaimedLandColor.B)
			} else {
				// belongs to province owner
				provinceColor := allCivColors[provinceOwner]
				fmt.Println("Drawing province", i, "with owner set to", provinceOwner)
				dc.SetRGB255(provinceColor.IRed, provinceColor.IGreen, provinceColor.IBlue)
			}
			dc.Fill()
		}
	}
}

func drawRegionsMap(outputFilename string, regionsMapData RegionsMapData) {
	dc := gg.NewContext(int(regionsMapData.GlobalMaxX), int(regionsMapData.GlobalMaxY))

	// water
	dc.SetRGB255(waterColor.R, waterColor.G, waterColor.B)
	dc.Clear()

	fmt.Println("Drawing map...")
	// drawProvinceTerrain(dc, regionsMapData.AllProvinceData)
	drawProvinceRegionColors(dc, regionsMapData.AllProvinceData, regionsMapData.AllRegionColors)
	drawProvinceOutline(dc, regionsMapData.AllProvinceData)
	// drawProvinceLabel(dc, regionsMapData.AllProvinceData)

	dc.SavePNG(outputFilename)
	fmt.Println("Saved image to", outputFilename)
}

func drawProvinceTerrain(dc *gg.Context, allProvinceData [][]ProvinceGameData) {
	for i := 0; i < len(allProvinceData); i++ {
		for p := 0; p < len(allProvinceData[i]); p++ {
			province := allProvinceData[i][p]
			dc.MoveTo(float64(province.LPointsX[0]), float64(province.LPointsY[0]))
			for j := 1; j < len(province.LPointsX); j++ {
				dc.LineTo(float64(province.LPointsX[j]), float64(province.LPointsY[j]))
			}
			dc.ClosePath()

			if province.ProvinceInfo.IContinentID == 0 || province.ProvinceInfo.STerrainTAG == "" {
				// water
				dc.SetRGB255(waterColorAlt.R, waterColorAlt.G, waterColorAlt.B)
			} else {
				// land
				dc.SetRGB255(landColor.R, landColor.G, landColor.B)
			}
			dc.Fill()
		}
	}
}

func drawProvinceRegionColors(dc *gg.Context, allProvinceData [][]ProvinceGameData, allRegionColors []RegionColor) {
	for i := 0; i < len(allProvinceData); i++ {
		for p := 0; p < len(allProvinceData[i]); p++ {
			province := allProvinceData[i][p]

			dc.MoveTo(float64(province.LPointsX[0]), float64(province.LPointsY[0]))
			for j := 1; j < len(province.LPointsX); j++ {
				dc.LineTo(float64(province.LPointsX[j]), float64(province.LPointsY[j]))
			}
			dc.ClosePath()

			if province.ProvinceInfo.IContinentID == 0 || province.ProvinceInfo.STerrainTAG == "" {
				// water
				dc.SetRGB255(waterColor.R, waterColor.G, waterColor.B)
			} else if province.ProvinceInfo.IRegionID < 0 || province.ProvinceInfo.IRegionID >= len(allRegionColors) {
				// land doesn't belong to valid region
				dc.SetRGB255(landColor.R, landColor.G, landColor.B)
			} else {
				// region color
				regionId := province.ProvinceInfo.IRegionID
				regionColor := allRegionColors[regionId]
				dc.SetRGB(regionColor.FractionRed, regionColor.FractionGreen, regionColor.FractionBlue)
			}
			dc.Fill()
		}
	}
}

func drawProvinceOutline(dc *gg.Context, allProvinceData [][]ProvinceGameData) {
	for i := 0; i < len(allProvinceData); i++ {
		for p := 0; p < len(allProvinceData[i]); p++ {
			province := allProvinceData[i][p]
			dc.SetRGB255(outlineColor.R, outlineColor.G, outlineColor.B)

			for j := 0; j < len(province.LPointsX); j++ {
				currentIndex := j % len(province.LPointsX)
				nextIndex := (j + 1) % len(province.LPointsX)

				dc.DrawLine(float64(province.LPointsX[currentIndex]), float64(province.LPointsY[currentIndex]),
					float64(province.LPointsX[nextIndex]), float64(province.LPointsY[nextIndex]))
				dc.Stroke()
			}
		}
	}
}

func drawProvinceLabel(dc *gg.Context, allProvinceData [][]ProvinceGameData) {
	for i := 0; i < len(allProvinceData); i++ {
		for p := 0; p < len(allProvinceData[i]); p++ {
			province := allProvinceData[i][p]

			provinceMinX := math.MaxFloat64
			provinceMinY := math.MaxFloat64
			provinceMaxX := 0.0
			provinceMaxY := 0.0

			for j := 0; j < len(province.LPointsX); j++ {
				currentIndex := j % len(province.LPointsX)

				if float64(province.LPointsX[currentIndex]) > provinceMaxX {
					provinceMaxX = float64(province.LPointsX[currentIndex])
				}
				if float64(province.LPointsY[currentIndex]) > provinceMaxY {
					provinceMaxY = float64(province.LPointsY[currentIndex])
				}
				if float64(province.LPointsX[currentIndex]) < provinceMinX {
					provinceMinX = float64(province.LPointsX[currentIndex])
				}
				if float64(province.LPointsY[currentIndex]) < provinceMinY {
					provinceMinY = float64(province.LPointsY[currentIndex])
				}
			}

			dc.SetRGB255(labelColor.R, labelColor.G, labelColor.B)
			averageX := (provinceMinX + provinceMaxX) / 2.0
			averageY := (provinceMinY + provinceMaxY) / 2.0
			fmt.Println(fmt.Sprintf("Province %v bounds min(%v, %v), max (%v, %v), average(%v, %v)",
				i, provinceMinX, provinceMinY, provinceMaxX, provinceMaxY, averageX, averageY))
			dc.DrawString(strconv.Itoa(i), averageX, averageY)
		}
	}
}
