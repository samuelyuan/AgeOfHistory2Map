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

// Apply sets the drawing context's current color to this color.
func (c Color) Apply(dc *gg.Context) {
	dc.SetRGB255(c.R, c.G, c.B)
}

// toFraction converts a 0-255 Color to the 0-1 fraction form used by FractionColor.
func (c Color) toFraction() FractionColor {
	return FractionColor{R: float64(c.R) / 255, G: float64(c.G) / 255, B: float64(c.B) / 255}
}

// FractionColor is a color expressed as 0-1 fractions, matching the format region data is stored in.
type FractionColor struct {
	R, G, B float64
}

// Apply sets the drawing context's current color to this color.
func (c FractionColor) Apply(dc *gg.Context) {
	dc.SetRGB(c.R, c.G, c.B)
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
	waterColor.Apply(dc)
	dc.Clear()

	fmt.Println("Drawing map...")
	drawScenarioRegionColors(dc, scenario.AllProvinceData, scenario.AllProvinceOwners, scenario.AllCivColors)
	// drawProvinceOutline(dc, scenario.AllProvinceData)

	dc.SavePNG(outputFilename)
	fmt.Println("Saved image to", outputFilename)
}

// buildProvincePath traces a province's boundary as the current path on the drawing context.
func buildProvincePath(dc *gg.Context, province ProvinceGameData) {
	dc.MoveTo(float64(province.LPointsX[0]), float64(province.LPointsY[0]))
	for j := 1; j < len(province.LPointsX); j++ {
		dc.LineTo(float64(province.LPointsX[j]), float64(province.LPointsY[j]))
	}
	dc.ClosePath()
}

// isWaterProvince reports whether a province represents water rather than land. Pure.
func isWaterProvince(province ProvinceGameData) bool {
	return province.ProvinceInfo.IContinentID == 0 || province.ProvinceInfo.STerrainTAG == ""
}

// determineScenarioProvinceColor returns the fill color for a province on the scenario (ownership)
// map, along with whether provinceOwner referred to a valid civilization. Pure: makes no calls to
// dc, so the ownership/coloring decision can be unit tested without a drawing context.
func determineScenarioProvinceColor(province ProvinceGameData, provinceOwner int, allCivColors []CivilizationColor) (color Color, ownerValid bool) {
	if isWaterProvince(province) {
		return waterColor, true
	}
	if provinceOwner < 0 || provinceOwner >= len(allCivColors) {
		return unclaimedLandColor, false
	}
	civColor := allCivColors[provinceOwner]
	return Color{R: civColor.IRed, G: civColor.IGreen, B: civColor.IBlue}, true
}

func drawScenarioRegionColors(dc *gg.Context, allProvinceData [][]ProvinceGameData, allProvinceOwners []int, allCivColors []CivilizationColor) {
	for i := 0; i < len(allProvinceData); i++ {
		for p := 0; p < len(allProvinceData[i]); p++ {
			province := allProvinceData[i][p]
			buildProvincePath(dc, province)

			provinceOwner := allProvinceOwners[i] - 1
			color, ownerValid := determineScenarioProvinceColor(province, provinceOwner, allCivColors)
			if !isWaterProvince(province) {
				if !ownerValid {
					fmt.Println("Province owner", provinceOwner, "isn't a valid province")
				} else {
					fmt.Println("Drawing province", i, "with owner set to", provinceOwner)
				}
			}
			color.Apply(dc)
			dc.Fill()
		}
	}
}

func drawRegionsMap(outputFilename string, regionsMapData RegionsMapData) {
	dc := gg.NewContext(int(regionsMapData.GlobalMaxX), int(regionsMapData.GlobalMaxY))

	// water
	waterColor.Apply(dc)
	dc.Clear()

	fmt.Println("Drawing map...")
	// drawProvinceTerrain(dc, regionsMapData.AllProvinceData)
	drawProvinceRegionColors(dc, regionsMapData.AllProvinceData, regionsMapData.AllRegionColors)
	drawProvinceOutline(dc, regionsMapData.AllProvinceData)
	// drawProvinceLabel(dc, regionsMapData.AllProvinceData)

	dc.SavePNG(outputFilename)
	fmt.Println("Saved image to", outputFilename)
}

// determineTerrainColor returns the fill color for a province on the terrain map. Pure.
func determineTerrainColor(province ProvinceGameData) Color {
	if isWaterProvince(province) {
		return waterColorAlt
	}
	return landColor
}

func drawProvinceTerrain(dc *gg.Context, allProvinceData [][]ProvinceGameData) {
	for i := 0; i < len(allProvinceData); i++ {
		for p := 0; p < len(allProvinceData[i]); p++ {
			province := allProvinceData[i][p]
			buildProvincePath(dc, province)
			determineTerrainColor(province).Apply(dc)
			dc.Fill()
		}
	}
}

// determineRegionProvinceColor returns the fill color for a province on the regions map. Pure.
func determineRegionProvinceColor(province ProvinceGameData, allRegionColors []RegionColor) FractionColor {
	if isWaterProvince(province) {
		return waterColor.toFraction()
	}
	regionId := province.ProvinceInfo.IRegionID
	if regionId < 0 || regionId >= len(allRegionColors) {
		return landColor.toFraction()
	}
	regionColor := allRegionColors[regionId]
	return FractionColor{R: regionColor.FractionRed, G: regionColor.FractionGreen, B: regionColor.FractionBlue}
}

func drawProvinceRegionColors(dc *gg.Context, allProvinceData [][]ProvinceGameData, allRegionColors []RegionColor) {
	for i := 0; i < len(allProvinceData); i++ {
		for p := 0; p < len(allProvinceData[i]); p++ {
			province := allProvinceData[i][p]
			buildProvincePath(dc, province)
			determineRegionProvinceColor(province, allRegionColors).Apply(dc)
			dc.Fill()
		}
	}
}

func drawProvinceOutline(dc *gg.Context, allProvinceData [][]ProvinceGameData) {
	for i := 0; i < len(allProvinceData); i++ {
		for p := 0; p < len(allProvinceData[i]); p++ {
			province := allProvinceData[i][p]
			outlineColor.Apply(dc)

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

// provinceBounds returns the axis-aligned bounding box of a province's points. Pure.
func provinceBounds(province ProvinceGameData) (minX, minY, maxX, maxY float64) {
	minX, minY = math.MaxFloat64, math.MaxFloat64
	for j := 0; j < len(province.LPointsX); j++ {
		x, y := float64(province.LPointsX[j]), float64(province.LPointsY[j])
		minX, maxX = math.Min(minX, x), math.Max(maxX, x)
		minY, maxY = math.Min(minY, y), math.Max(maxY, y)
	}
	return minX, minY, maxX, maxY
}

// provinceCenter returns the midpoint of a province's bounding box. Pure.
func provinceCenter(province ProvinceGameData) (x, y float64) {
	minX, minY, maxX, maxY := provinceBounds(province)
	return (minX + maxX) / 2.0, (minY + maxY) / 2.0
}

func drawProvinceLabel(dc *gg.Context, allProvinceData [][]ProvinceGameData) {
	for i := 0; i < len(allProvinceData); i++ {
		for p := 0; p < len(allProvinceData[i]); p++ {
			province := allProvinceData[i][p]

			labelColor.Apply(dc)
			averageX, averageY := provinceCenter(province)
			fmt.Printf("Province %v center at (%v, %v)\n", i, averageX, averageY)
			dc.DrawString(strconv.Itoa(i), averageX, averageY)
		}
	}
}
