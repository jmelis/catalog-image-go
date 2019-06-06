package catalog

// Bundle represents a CSV and its sidefiles if any
type Bundle struct {
	Operator  string
	CSV       CSV
	SideFiles []SideFile
}

// Bundles is a collection of bundles
type Bundles []Bundle

// Name returns the bundle name
func (b Bundle) Name() string {
	return b.CSV.Name()
}

// Hash returns the CSV hash
func (b Bundle) Hash() string {
	return b.CSV.Hash()
}

// Version returns the bundle version
func (b Bundle) Version() string {
	return b.CSV.Version()
}

// Replaces returns .spec.replaces of the CSV in the bundle
func (b Bundle) Replaces() string {
	return b.CSV.Replaces()
}
