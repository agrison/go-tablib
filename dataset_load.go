package tablib

import (
	"encoding/json"
	"gopkg.in/yaml.v2"
)

// LoadYAML loads a dataset from either a YAML source
func LoadYAML(yamlContent []byte) (*Dataset, error) {
	var input []map[string]interface{}
	if err := yaml.Unmarshal(yamlContent, &input); err != nil {
		return nil, err
	}

	return internalLoadFromDict(input)
}

// LoadJSON loads a dataset from either a JSON source
func LoadJSON(jsonContent []byte) (*Dataset, error) {
	var input []map[string]interface{}
	if err := json.Unmarshal(jsonContent, &input); err != nil {
		return nil, err
	}

	return internalLoadFromDict(input)
}

func internalLoadFromDict(input []map[string]interface{}) (*Dataset, error) {
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
