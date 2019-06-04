package catalog

// SideFile represents any YAML file included in the bundle
type SideFile struct {
	name    string
	content []byte
}
