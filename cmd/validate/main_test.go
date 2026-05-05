package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"giessen-maps-review-removals/internal/mapsreview"
)

func TestParseArgsValidateDefaults(t *testing.T) {
	args, err := parseArgs(nil)
	if err != nil {
		t.Fatal(err)
	}
	if args.Input != mapsreview.ResultsJSON {
		t.Fatalf("Input = %q, want %q", args.Input, mapsreview.ResultsJSON)
	}
	if args.StrictCity {
		t.Fatal("StrictCity should default to false")
	}
}

func TestParseArgsValidateStrict(t *testing.T) {
	args, err := parseArgs([]string{"--strict-city"})
	if err != nil {
		t.Fatal(err)
	}
	if !args.StrictCity {
		t.Fatal("StrictCity should be true")
	}
}

func TestValidateSuccessRows(t *testing.T) {
	rows := []mapsreview.Place{{
		ID:          "test-id-1",
		Name:        "Test Place",
		URL:         "https://www.google.com/maps/place/Test+Place/@50.5840,8.6784,17z",
		Postcode:    mapsreview.StringPtr("35390"),
		Address:     mapsreview.StringPtr("Teststrasse 1, 35390 Giessen"),
		Rating:      mapsreview.FloatPtr(4.5),
		ReviewCount: mapsreview.IntPtr(100),
		Lat:         mapsreview.FloatPtr(50.5840),
		Lng:         mapsreview.FloatPtr(8.6784),
		Status:      "success",
	}}

	tmpFile := writeTempPlacesJSON(t, rows)
	defer os.Remove(tmpFile)

	if err := run(args{Input: tmpFile}); err != nil {
		t.Fatalf("validation failed for valid data: %v", err)
	}
}

func TestValidateDetectsDuplicates(t *testing.T) {
	rows := []mapsreview.Place{
		{
			ID:          "dup-id",
			Name:        "A",
			URL:         "https://maps.example.com/A",
			Rating:      mapsreview.FloatPtr(4.0),
			ReviewCount: mapsreview.IntPtr(10),
			Postcode:    mapsreview.StringPtr("35390"),
			Address:     mapsreview.StringPtr("Str 1, 35390 Giessen"),
			Status:      "success",
		},
		{
			ID:          "dup-id",
			Name:        "B",
			URL:         "https://maps.example.com/B",
			Rating:      mapsreview.FloatPtr(4.0),
			ReviewCount: mapsreview.IntPtr(10),
			Postcode:    mapsreview.StringPtr("35390"),
			Address:     mapsreview.StringPtr("Str 2, 35390 Giessen"),
			Status:      "success",
		},
	}

	tmpFile := writeTempPlacesJSON(t, rows)
	defer os.Remove(tmpFile)

	err := run(args{Input: tmpFile})
	if err == nil {
		t.Fatal("expected validation error for duplicate IDs")
	}
	if !strings.Contains(err.Error(), "validation failed") {
		t.Fatalf("error message does not indicate validation failure: %v", err)
	}
}

func TestValidateOutsidePostcodeWarning(t *testing.T) {
	rows := []mapsreview.Place{{
		ID:          "outside-id",
		Name:        "Outside Place",
		URL:         "https://maps.example.com/Outside",
		Postcode:    mapsreview.StringPtr("80331"),
		Rating:      mapsreview.FloatPtr(4.0),
		ReviewCount: mapsreview.IntPtr(10),
		Status:      "success",
	}}

	tmpFile := writeTempPlacesJSON(t, rows)
	defer os.Remove(tmpFile)

	if err := run(args{Input: tmpFile}); err != nil {
		t.Fatalf("outside postcode should be warning, not error: %v", err)
	}
}

func TestValidateStrictCityErrorsOnOutside(t *testing.T) {
	rows := []mapsreview.Place{{
		ID:          "outside-id",
		Name:        "Outside Place",
		URL:         "https://maps.example.com/Outside",
		Postcode:    mapsreview.StringPtr("80331"),
		Rating:      mapsreview.FloatPtr(4.0),
		ReviewCount: mapsreview.IntPtr(10),
		Status:      "success",
	}}

	tmpFile := writeTempPlacesJSON(t, rows)
	defer os.Remove(tmpFile)

	err := run(args{Input: tmpFile, StrictCity: true})
	if err == nil {
		t.Fatal("expected validation error with --strict-city")
	}
	if !strings.Contains(err.Error(), "validation failed") {
		t.Fatalf("error message does not indicate validation failure: %v", err)
	}
}

func writeTempPlacesJSON(t *testing.T, rows []mapsreview.Place) string {
	t.Helper()
	data, err := json.MarshalIndent(rows, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	file := filepath.Join(t.TempDir(), "places.json")
	if err := os.WriteFile(file, data, 0o644); err != nil {
		t.Fatal(err)
	}
	return file
}
