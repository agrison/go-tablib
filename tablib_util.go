package tablib

// internalLoadFromDict creates a Dataset from an array of map representing columns.
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
