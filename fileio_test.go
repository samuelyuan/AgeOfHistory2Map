package main

import (
	"reflect"
	"testing"
)

func TestScaleProvincePoints(t *testing.T) {
	provinces := []ProvinceGameData{
		{LPointsX: []int{1, 2, 3}, LPointsY: []int{4, 5, 6}},
		{LPointsX: []int{10}, LPointsY: []int{20}},
	}

	scaleProvincePoints(provinces, 2)

	want := []ProvinceGameData{
		{LPointsX: []int{2, 4, 6}, LPointsY: []int{8, 10, 12}},
		{LPointsX: []int{20}, LPointsY: []int{40}},
	}
	if !reflect.DeepEqual(provinces, want) {
		t.Errorf("scaleProvincePoints() = %+v, want %+v", provinces, want)
	}
}

func TestComputeGlobalBounds(t *testing.T) {
	allProvinceData := [][]ProvinceGameData{
		{
			{LPointsX: []int{0, 100}, LPointsY: []int{0, 50}},
		},
		{
			{LPointsX: []int{200, 10}, LPointsY: []int{75, 5}},
		},
	}

	maxX, maxY := computeGlobalBounds(allProvinceData)
	if maxX != 200 || maxY != 75 {
		t.Errorf("computeGlobalBounds() = (%v, %v), want (200, 75)", maxX, maxY)
	}
}

func TestComputeGlobalBoundsEmpty(t *testing.T) {
	maxX, maxY := computeGlobalBounds(nil)
	if maxX != 0 || maxY != 0 {
		t.Errorf("computeGlobalBounds(nil) = (%v, %v), want (0, 0)", maxX, maxY)
	}
}

func TestStripKnownCivTagSuffix(t *testing.T) {
	cases := []struct {
		name       string
		civTag     string
		wantTag    string
		wantSuffix string
		wantOK     bool
	}{
		{"strips romance suffix", "french_r", "french", "_r", true},
		{"strips first matching suffix only", "french_r_c", "french", "_r", true},
		{"no known suffix", "french", "french", "", false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			gotTag, gotSuffix, gotOK := stripKnownCivTagSuffix(c.civTag)
			if gotTag != c.wantTag || gotSuffix != c.wantSuffix || gotOK != c.wantOK {
				t.Errorf("stripKnownCivTagSuffix(%q) = (%q, %q, %v), want (%q, %q, %v)",
					c.civTag, gotTag, gotSuffix, gotOK, c.wantTag, c.wantSuffix, c.wantOK)
			}
		})
	}
}

func TestComputeCivEconomyMap(t *testing.T) {
	allProvinceOwners := []int{1, 1, 2, 0}
	provincesData := []SaveProvinceInfo{
		{IEconomy: 10},
		{IEconomy: 5},
		{IEconomy: 20},
		{IEconomy: 1},
	}

	got := computeCivEconomyMap(allProvinceOwners, provincesData)
	want := map[int]int{
		0:  15, // owner 1 -> key 0: provinces 0 and 1
		1:  20, // owner 2 -> key 1: province 2
		-1: 1,  // owner 0 -> key -1: province 3
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("computeCivEconomyMap() = %v, want %v", got, want)
	}
}

func TestComputeCivEconomyMapShorterProvincesData(t *testing.T) {
	// Provinces beyond the bounds of provincesData should be skipped, not panic.
	allProvinceOwners := []int{1, 1, 2}
	provincesData := []SaveProvinceInfo{
		{IEconomy: 10},
	}

	got := computeCivEconomyMap(allProvinceOwners, provincesData)
	want := map[int]int{0: 10}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("computeCivEconomyMap() = %v, want %v", got, want)
	}
}
