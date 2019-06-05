package catalog

import "fmt"

// Bundle represents a CSV and its sidefiles if any
type Bundle struct {
	Operator  string
	CSV       CSV
	SideFiles []SideFile
}

// Bundles is a collection of bundles
type Bundles []Bundle

// FindLatestCSV returns the latest CSV
func (bundles Bundles) FindLatestCSV() (string, error) {
	var latestCSV string

	setReplaces := make(map[string]bool)
	setCSV := make(map[string]bool)

	for _, b := range bundles {
		version := b.CSV.Version()
		csvName := CSVName(b.Operator, version)
		csvReplaces := b.CSV.Replaces()

		setReplaces[csvReplaces] = true
		setCSV[csvName] = true
	}

	for replaces := range setReplaces {
		delete(setCSV, replaces)
	}

	if len(setCSV) != 1 {
		err := fmt.Errorf("Invalid number of leaves found: %d", len(setCSV))
		return latestCSV, err
	}

	for csv := range setCSV {
		latestCSV = csv
	}

	return latestCSV, nil
}
