package catalog

// Store implements any backing store for Catalog
type Store interface {
	Load() (*Catalog, error)
	Save(*Catalog) error
}
