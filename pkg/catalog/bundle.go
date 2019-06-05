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

// FindLatestCSV calculates the latest CSV
func (bundles Bundles) FindLatestCSV() (string, error) {
	var latestCSV string

	if len(bundles) == 0 {
		return latestCSV, nil
	}

	setReplaces := make(map[string]bool)
	setCSV := make(map[string]bool)

	for _, b := range bundles {
		csvReplaces := b.CSV.Replaces()

		setReplaces[csvReplaces] = true
		setCSV[b.CSV.Name()] = true
	}

	for replaces := range setReplaces {
		delete(setCSV, replaces)
	}

	if len(setCSV) != 1 {
		err := fmt.Errorf("Invalid number of leaves found: %d", len(setCSV))

		// TODO: REMOVE
		for csv := range setCSV {
			fmt.Println(csv)
		}

		return latestCSV, err
	}

	for csv := range setCSV {
		latestCSV = csv
	}

	return latestCSV, nil
}
