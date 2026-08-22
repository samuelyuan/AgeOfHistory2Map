package main

import (
	"testing"

	"github.com/fogleman/gg"
)

func TestIsWaterProvince(t *testing.T) {
	cases := []struct {
		name     string
		province ProvinceGameData
		want     bool
	}{
		{"zero continent id is water", ProvinceGameData{ProvinceInfo: ProvinceInfo{IContinentID: 0, STerrainTAG: "plains"}}, true},
		{"empty terrain tag is water", ProvinceGameData{ProvinceInfo: ProvinceInfo{IContinentID: 1, STerrainTAG: ""}}, true},
		{"land province", ProvinceGameData{ProvinceInfo: ProvinceInfo{IContinentID: 1, STerrainTAG: "plains"}}, false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := isWaterProvince(c.province); got != c.want {
				t.Errorf("isWaterProvince() = %v, want %v", got, c.want)
			}
		})
	}
}

func TestDetermineTerrainColor(t *testing.T) {
	water := ProvinceGameData{ProvinceInfo: ProvinceInfo{IContinentID: 0}}
	land := ProvinceGameData{ProvinceInfo: ProvinceInfo{IContinentID: 1, STerrainTAG: "plains"}}

	if got := determineTerrainColor(water); got != waterColorAlt {
		t.Errorf("water province = %v, want %v", got, waterColorAlt)
	}
	if got := determineTerrainColor(land); got != landColor {
		t.Errorf("land province = %v, want %v", got, landColor)
	}
}

func TestDetermineScenarioProvinceColor(t *testing.T) {
	land := ProvinceGameData{ProvinceInfo: ProvinceInfo{IContinentID: 1, STerrainTAG: "plains"}}
	water := ProvinceGameData{ProvinceInfo: ProvinceInfo{IContinentID: 0}}
	civColors := []CivilizationColor{{IRed: 10, IGreen: 20, IBlue: 30}}

	if color, ok := determineScenarioProvinceColor(water, -1, civColors); color != waterColor || !ok {
		t.Errorf("water province = (%v, %v), want (%v, true)", color, ok, waterColor)
	}
	if color, ok := determineScenarioProvinceColor(land, 0, civColors); color != (Color{R: 10, G: 20, B: 30}) || !ok {
		t.Errorf("valid owner = (%v, %v), want ({10 20 30}, true)", color, ok)
	}
	if color, ok := determineScenarioProvinceColor(land, 5, civColors); color != unclaimedLandColor || ok {
		t.Errorf("out-of-range owner = (%v, %v), want (%v, false)", color, ok, unclaimedLandColor)
	}
	if color, ok := determineScenarioProvinceColor(land, -1, civColors); color != unclaimedLandColor || ok {
		t.Errorf("negative owner = (%v, %v), want (%v, false)", color, ok, unclaimedLandColor)
	}
}

func TestDetermineRegionProvinceColor(t *testing.T) {
	land := ProvinceGameData{ProvinceInfo: ProvinceInfo{IContinentID: 1, STerrainTAG: "plains", IRegionID: 0}}
	water := ProvinceGameData{ProvinceInfo: ProvinceInfo{IContinentID: 0}}
	regionColors := []RegionColor{{FractionRed: 0.5, FractionGreen: 0.25, FractionBlue: 0.1}}

	if got := determineRegionProvinceColor(water, regionColors); got != waterColor.toFraction() {
		t.Errorf("water province = %v, want %v", got, waterColor.toFraction())
	}
	if got := determineRegionProvinceColor(land, regionColors); got != (FractionColor{R: 0.5, G: 0.25, B: 0.1}) {
		t.Errorf("valid region = %v, want {0.5 0.25 0.1}", got)
	}

	noRegion := ProvinceGameData{ProvinceInfo: ProvinceInfo{IContinentID: 1, STerrainTAG: "plains", IRegionID: 99}}
	if got := determineRegionProvinceColor(noRegion, regionColors); got != landColor.toFraction() {
		t.Errorf("out-of-range region = %v, want %v", got, landColor.toFraction())
	}
}

func TestProvinceBounds(t *testing.T) {
	province := ProvinceGameData{
		LPointsX: []int{10, 30, 20},
		LPointsY: []int{5, 15, 0},
	}
	minX, minY, maxX, maxY := provinceBounds(province)
	if minX != 10 || maxX != 30 || minY != 0 || maxY != 15 {
		t.Errorf("bounds = (%v, %v, %v, %v), want (10, 0, 30, 15)", minX, minY, maxX, maxY)
	}
}

func TestProvinceCenter(t *testing.T) {
	province := ProvinceGameData{
		LPointsX: []int{0, 10},
		LPointsY: []int{0, 20},
	}
	x, y := provinceCenter(province)
	if x != 5 || y != 10 {
		t.Errorf("center = (%v, %v), want (5, 10)", x, y)
	}
}

func TestBuildProvincePath(t *testing.T) {
	// Not pure (touches a drawing context), but small enough to smoke-test: draw and fill a
	// square path, then confirm a pixel inside it changed color.
	dc := gg.NewContext(10, 10)
	dc.SetRGB255(255, 255, 255)
	dc.Clear()

	province := ProvinceGameData{
		LPointsX: []int{1, 8, 8, 1},
		LPointsY: []int{1, 1, 8, 8},
	}
	buildProvincePath(dc, province)
	dc.SetRGB255(0, 0, 0)
	dc.Fill()

	r, g, b, _ := dc.Image().At(5, 5).RGBA()
	if r != 0 || g != 0 || b != 0 {
		t.Errorf("expected filled pixel to be black, got (%v, %v, %v)", r, g, b)
	}
}
