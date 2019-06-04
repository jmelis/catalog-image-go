package catalog

import "fmt"

// Bundle represents a CSV and its sidefiles if any
type Bundle struct {
	operator  string
	csv       CSV
	sidefiles []SideFile
}

// Bundles is a collection of bundles
type Bundles []Bundle

// FindLatestCSV returns the latest CSV
func (bundles Bundles) FindLatestCSV() (string, error) {
	setReplaces := make(map[string]bool)
	setCSV := make(map[string]bool)

	for _, b := range bundles {
		version := b.csv.version
		csvName := CSVName(b.operator, version)
		csvReplaces := b.csv.GetReplaces()

		setReplaces[csvReplaces] = true
		setCSV[csvName] = true
	}

	for replaces := range setReplaces {
		delete(setCSV, replaces)
	}

	if len(setCSV) != 1 {
		return "", fmt.Errorf("Invalid number of leaves found: %d", len(setCSV))
	}

	var latestCSV string
	for csv := range setCSV {
		latestCSV = csv
	}

	return latestCSV, nil
}
