package catalog

// Store implements any backing store for Catalog
type Store interface {
	load() ([]Bundle, error)
	save() error
}
