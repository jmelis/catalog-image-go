package catalog

import (
	"io/ioutil"
	"path/filepath"
	"strings"
)

// Catalog represents the bundle repository with Package file, and Bundle dirs.
// The Bundle dirs contain the CSV and may containe sidefiles, generally CRDs
type Catalog struct {
	Operator string
	store    Store
	Bundles  Bundles
}

// LoadCatalog returns a new Catalog with the specified store
func LoadCatalog(store Store) (*Catalog, error) {
	c, err := store.load()
	if err != nil {
		return nil, err
	}
	return c, nil
}

// Save all the bundles to storage
func (c *Catalog) Save() error {
	return c.store.save(c)
}

// AddBundle adds a bundle in the local directory
func (c *Catalog) AddBundle(path string) error {
	var csv CSV
	var sidefiles []SideFile

	latestCSV, err := c.FindLatestCSV()
	if err != nil {
		return err
	}

	files, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}

	for _, file := range files {
		content, err := ioutil.ReadFile(filepath.Join(path, file.Name()))
		if err != nil {
			return err
		}

		if strings.HasSuffix(file.Name(), ".yaml") {
			if strings.HasSuffix(file.Name(), CSVSuffix) {
				csv, err = NewCSV(c.Operator, content)
				if err != nil {
					return err
				}

				err := csv.SetReplaces(latestCSV)
				if err != nil {
					return err
				}
			} else {
				sidefile := SideFile{
					name:    file.Name(),
					content: content,
				}
				sidefiles = append(sidefiles, sidefile)
			}
		}
	}

	bundle := Bundle{
		Operator:  c.Operator,
		CSV:       csv,
		SideFiles: sidefiles,
	}

	c.Bundles = append(c.Bundles, bundle)

	return nil
}

// FindLatestCSV returns the latest CSV
func (c *Catalog) FindLatestCSV() (string, error) {
	return c.Bundles.FindLatestCSV()
}
