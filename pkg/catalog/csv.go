package catalog

import "fmt"

// CSV represents ClusterServiceVersion
type CSV struct {
	operator string
	content  []byte
}

// CSVSuffix all CSVs must end with this suffix
const CSVSuffix = ".clusterserviceversion.yaml"

// CSVName generates the name of a CSV file
func CSVName(operator, version string) string {
	return fmt.Sprintf("%s-operator.v%s", operator, version)
}

// CSVFileName generates the name of a CSV file
func CSVFileName(operator, version string) string {
	return CSVName(operator, version) + CSVSuffix
}

// NewCSV returns a new CSV
func NewCSV(operator string, content []byte) (CSV, error) {
	return CSV{operator: operator, content: content}, nil
}

// SetReplaces changes .spec.replaces
func (c *CSV) SetReplaces(replaces string) error {
	if replaces == "" {
		return nil
	}

	uy, err := NewUnstructuredYaml(c.content)
	if err != nil {
		return err
	}

	spec, ok := uy["spec"].(map[string]interface{})
	if !ok {
		return fmt.Errorf(".spec not readable")
	}

	spec["replaces"] = replaces

	c.content = []byte(uy.String())
	return nil
}

// SetCatalogHash sets the value for .metadata.annotations."catalog-image/hash"
func (c *CSV) SetCatalogHash(hash string) error {
	if hash == "" {
		return nil
	}

	uy, err := NewUnstructuredYaml(c.content)
	if err != nil {
		return err
	}

	metadata, ok := uy["metadata"].(map[string]interface{})
	if !ok {
		return fmt.Errorf(".metadata not readable")
	}

	annotations, ok := metadata["annotations"].(map[string]interface{})
	if !ok {
		annotations = make(map[string]interface{})
	}

	annotations["catalog-image/hash"] = hash

	metadata["annotations"] = annotations

	c.content = []byte(uy.String())
	return nil
}

// Name returns the full name
func (c CSV) Name() string {
	return CSVName(c.operator, c.Version())
}

// FileName returns the file name
func (c CSV) FileName() string {
	return CSVFileName(c.operator, c.Version())
}

// Replaces returns .spec.replaces. Empty string if not present.
func (c CSV) Replaces() string {
	return c.GetSpecStringParameter("replaces")
}

// Version returns .spec.version. Empty string if not present.
func (c CSV) Version() string {
	return c.GetSpecStringParameter("version")
}

// Hash returns .metadata.annotations."catalog-image/hash"
func (c CSV) Hash() string {
	uy, err := NewUnstructuredYaml(c.content)
	if err != nil {
		return ""
	}

	metadata, ok := uy["metadata"].(map[string]interface{})
	if !ok {
		return ""
	}

	annotations, ok := metadata["annotations"].(map[string]interface{})
	if !ok {
		return ""
	}

	if hash, ok := annotations["catalog-image/hash"].(string); ok {
		return hash
	}

	return ""
}

// GetSpecStringParameter returns .spec.<param>. Empty string if not present.
func (c CSV) GetSpecStringParameter(param string) string {
	var spec map[string]interface{}
	var replaces string
	var ok bool

	uy, err := NewUnstructuredYaml(c.content)
	if err != nil {
		return ""
	}

	if spec, ok = uy["spec"].(map[string]interface{}); !ok {
		return ""
	}

	if replaces, ok = spec[param].(string); !ok {
		return ""
	}

	return replaces
}
