package catalog

// Store implements any backing store for Catalog
type Store interface {
	load() (Bundles, error)
	save(Bundles) error
}
