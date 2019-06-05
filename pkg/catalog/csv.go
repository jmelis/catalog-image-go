package catalog

import "fmt"

// CSV represents ClusterServiceVersion
type CSV struct {
	content []byte
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
func NewCSV(content []byte) (CSV, error) {
	content, err := CanonicalizeYaml(content)
	if err != nil {
		return CSV{}, err
	}

	return CSV{content: content}, nil
}

// SetReplaces returns a new CSV with a modified .spec.replaces
func (c CSV) SetReplaces(replaces string) (CSV, error) {
	var spec map[string]interface{}
	var ok bool

	uy, err := NewUnstructuredYaml(c.content)
	if err != nil {
		return c, err
	}

	if spec, ok = uy["spec"].(map[string]interface{}); !ok {
		return c, fmt.Errorf(".spec not readable")
	}

	spec["replaces"] = replaces

	c.content = []byte(uy.String())

	return c, nil
}

// Replaces returns .spec.replaces. Empty string if not present.
func (c CSV) Replaces() string {
	return c.GetSpecStringParameter("replaces")
}

// Version returns .spec.version. Empty string if not present.
func (c CSV) Version() string {
	return c.GetSpecStringParameter("version")
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
