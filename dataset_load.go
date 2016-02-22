package tablib

import (
	"gopkg.in/yaml.v2"
)

// LoadDataset loads a dataset from either a YAML, JSON, CSV, TSV or XML file
func LoadYAML(yamlContent []byte) (*Dataset, error) {
	var input []map[string]interface{}
	if err := yaml.Unmarshal(yamlContent, &input); err != nil {
		return nil, err
	}

	// retrieve columns
	headers := make([]string, 0, 10)
	for h := range input[0] {
		headers = append(headers, h)
	}

	ds := NewDataset(headers)
	for _, e := range input {
		row := make([]interface{}, 0, len(headers))
		for _, h := range headers {
			row = append(row, e[h])
		}
		ds.AppendValues(row...)
	}

	return ds, nil
}
