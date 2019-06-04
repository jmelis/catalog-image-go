package catalog

import "github.com/ghodss/yaml"

// UnstructuredYaml supports any kind of dictionary YAML
type UnstructuredYaml map[string]interface{}

// NewUnstructuredYaml returns a new UnstructuredYaml object
func NewUnstructuredYaml(content []byte) (UnstructuredYaml, error) {
	var uy UnstructuredYaml

	err := yaml.Unmarshal(content, &uy)
	if err != nil {
		return uy, err
	}

	return uy, nil
}

func (uy UnstructuredYaml) String() string {
	y, _ := yaml.Marshal(uy)
	return string(y)
}

// CanonicalizeYaml will unmarshal and re-marshal to obtain a deterministic output
func CanonicalizeYaml(content []byte) ([]byte, error) {
	uy, err := NewUnstructuredYaml(content)
	if err != nil {
		return nil, err
	}

	y, err := yaml.Marshal(uy)
	if err != nil {
		return nil, err
	}

	return y, nil
}
