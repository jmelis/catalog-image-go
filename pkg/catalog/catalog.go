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

// NewCatalog returns a new Catalog with the specified store
func NewCatalog(operator string, store Store) *Catalog {
	return &Catalog{Operator: operator, store: store}
}

// Load all the bundles
func (c *Catalog) Load() error {
	bundles, err := c.store.load()
	if err != nil {
		return err
	}

	c.Bundles = bundles
	return nil
}

// Save all the bundles to storage
func (c *Catalog) Save() error {
	return c.store.save(c.Bundles)
}

// AddBundle adds a bundle in the local directory
func (c *Catalog) AddBundle(path string) error {
	var csv CSV
	var sidefiles []SideFile

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
				csv, err = NewCSV(content)
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
		bundle := Bundle{
			Operator:  c.Operator,
			CSV:       csv,
			SideFiles: sidefiles,
		}

		c.Bundles = append(c.Bundles, bundle)

		return nil
	}

	return nil
}
