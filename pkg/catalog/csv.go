package catalog

import "fmt"

// CSV represents ClusterServiceVersion
type CSV struct {
	version string
	content []byte
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

// GetReplaces returns .spec.replaces. Empty string if not present.
func (c CSV) GetReplaces() string {
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

	if replaces, ok = spec["replaces"].(string); !ok {
		return ""
	}

	return replaces
}
