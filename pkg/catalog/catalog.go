package catalog

// Catalog represents the bundle repository with Package file, and Bundle dirs.
// The Bundle dirs contain the CSV and may containe sidefiles, generally CRDs
type Catalog struct {
	operator string
	store    Store
	bundles  []Bundle
}

// NewCatalog returns a new Catalog with the specified store
func NewCatalog(operator string, store Store) *Catalog {
	return &Catalog{operator: operator, store: store}
}

// Load all the bundles
func (c *Catalog) Load() error {
	bundles, err := c.store.load()
	if err != nil {
		return err
	}

	c.bundles = bundles
	return nil
}

// Save all the bundles to storage
func (c *Catalog) Save() error {
	return c.store.save()
}
