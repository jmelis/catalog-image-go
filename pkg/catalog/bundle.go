package catalog

// Bundle represents a CSV and its sidefiles if any
type Bundle struct {
	csv       CSV
	sidefiles []SideFile
}

// Bundles is a collection of bundles
type Bundles []Bundle
