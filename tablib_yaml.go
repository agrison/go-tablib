package tablib

import "gopkg.in/yaml.v2"

// LoadYAML loads a dataset from a YAML source.
func LoadYAML(yamlContent []byte) (*Dataset, error) {
	var input []map[string]interface{}
	if err := yaml.Unmarshal(yamlContent, &input); err != nil {
		return nil, err
	}

	return internalLoadFromDict(input)
}

// YAML returns a YAML representation of the Dataset as string.
func (d *Dataset) YAML() (string, error) {
	back := d.Dict()

	b, err := yaml.Marshal(back)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// YAML returns a YAML representation of the Databook as string.
func (d *Databook) YAML() (string, error) {
	y := make([]map[string]interface{}, len(d.sheets))
	i := 0
	for _, s := range d.sheets {
		y[i] = make(map[string]interface{})
		y[i]["title"] = s.title
		y[i]["dataset"] = s.dataset.Dict()
	}
	b, err := yaml.Marshal(y)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
