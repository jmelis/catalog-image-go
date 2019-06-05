package catalog

import (
	"fmt"
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
	c, err := store.Load()
	if err != nil {
		return nil, err
	}
	return c, nil
}

// Save all the bundles to storage
func (c *Catalog) Save() error {
	return c.store.Save(c)
}

// AddBundle adds a bundle in the local directory
func (c *Catalog) AddBundle(path string) error {
	var csv CSV
	var sidefiles []SideFile

	latestBundle, err := c.FindLatestBundle()
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

				err := (&csv).SetReplaces(latestBundle.Name())
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

// FindLatestBundle returns the latest bundle
func (c *Catalog) FindLatestBundle() (Bundle, error) {
	var latestBundleName string
	var latestBundle Bundle

	if len(c.Bundles) == 0 {
		return latestBundle, fmt.Errorf("no bundles exist")
	}

	setReplaces := make(map[string]bool)
	setCSV := make(map[string]bool)

	for _, b := range c.Bundles {
		if csvReplaces := b.Replaces(); csvReplaces != "" {
			setReplaces[csvReplaces] = true
		}
		setCSV[b.CSV.Name()] = true
	}

	for csvReplaces := range setReplaces {
		delete(setCSV, csvReplaces)
	}

	if len(setCSV) != 1 {
		err := fmt.Errorf("invalid number of leaves found: %d", len(setCSV))

		return latestBundle, err
	}

	for csv := range setCSV {
		latestBundleName = csv
	}

	for _, b := range c.Bundles {
		if b.Name() == latestBundleName {
			return b, nil
		}
	}

	return latestBundle, fmt.Errorf("latest bundle not found")
}

// RemoveBundle removes a bundle from Bundles without updating the replaces field
func (c *Catalog) RemoveBundle(csvName string) error {
	csvIndex := -1
	for i, b := range c.Bundles {
		if b.CSV.Name() == csvName {
			csvIndex = i
			break
		}
	}

	if csvIndex == -1 {
		return fmt.Errorf("CSV %s not found", csvName)
	}

	c.Bundles = append(c.Bundles[:csvIndex], c.Bundles[csvIndex+1:]...)

	return nil
}

// PruneAfterCSV TODO
func (c *Catalog) PruneAfterCSV(csvName string) error {
	var bundle Bundle

	// ensure csvName exists
	if _, err := c.GetBundle(csvName); err != nil {
		return err
	}

	// get latest bundle
	bundle, err := c.FindLatestBundle()
	if err != nil {
		return err
	}

	// start with latest bundle and remove each one until we found the bundle
	for bundle.Name() != csvName {
		parent, err := c.GetBundle(bundle.Replaces())
		if err != nil {
			return err
		}

		if err := c.RemoveBundle(bundle.Name()); err != nil {
			return err
		}

		bundle = parent
	}

	return nil
}

// GetBundle returns the Bundle object that matches a given name
func (c *Catalog) GetBundle(csvName string) (Bundle, error) {
	var bundle Bundle

	for _, b := range c.Bundles {
		if b.Name() == csvName {
			return b, nil
		}
	}

	return bundle, fmt.Errorf("Bundle %s not found", csvName)
}
