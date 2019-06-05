package catalog

// Store implements any backing store for Catalog
type Store interface {
	load() (*Catalog, error)
	save(*Catalog) error
}
