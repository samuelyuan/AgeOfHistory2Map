package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultOutputFilename(t *testing.T) {
	cases := []struct {
		mapName string
		want    string
	}{
		{"regions", "regions.png"},
		{"1200", "1200.png"},
		{"1440", "1440.png"},
		{"modernworld", "modern.png"},
	}
	for _, c := range cases {
		if got := defaultOutputFilename(c.mapName); got != c.want {
			t.Errorf("defaultOutputFilename(%q) = %q, want %q", c.mapName, got, c.want)
		}
	}
}

func TestCheckDataDirectory(t *testing.T) {
	tmp := t.TempDir()
	target := filepath.Join(tmp, "data")

	if err := checkDataDirectory(target); err == nil {
		t.Errorf("checkDataDirectory(%q) = nil, want error when directory is missing", target)
	}

	if err := os.WriteFile(target, []byte("not a directory"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := checkDataDirectory(target); err == nil {
		t.Errorf("checkDataDirectory(%q) = nil, want error when path is a file, not a directory", target)
	}
	if err := os.Remove(target); err != nil {
		t.Fatal(err)
	}

	if err := os.Mkdir(target, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := checkDataDirectory(target); err != nil {
		t.Errorf("checkDataDirectory(%q) = %v, want nil when directory exists", target, err)
	}
}
