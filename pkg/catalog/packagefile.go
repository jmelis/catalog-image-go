package catalog

import "github.com/ghodss/yaml"

// PackageFile object
type PackageFile struct {
	PackageName string    `json:"packageName"`
	Channels    []Channel `json:"channels"`
}

// Channel object in the Channels parameter of a PackageFile
type Channel struct {
	Name       string `json:"name"`
	CurrentCSV string `json:"currentCSV"`
}

// NewPackageFile creates a new PackageFile
func NewPackageFile(operator, channel, currentCSV string) PackageFile {
	c := Channel{Name: channel, CurrentCSV: currentCSV}
	return PackageFile{PackageName: operator, Channels: []Channel{c}}
}

// YAML returns the PackageFile in YAML format
func (p PackageFile) YAML() ([]byte, error) {
	y, err := yaml.Marshal(p)
	if err != nil {
		return nil, err
	}
	return y, nil
}

// FileName returns the full name of package file
func (p PackageFile) FileName() string {
	return p.PackageName + ".package.yaml"
}
