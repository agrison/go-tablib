package tablib

import (
	"encoding/json"
)

// LoadJSON loads a dataset from a YAML source.
func LoadJSON(jsonContent []byte) (*Dataset, error) {
	var input []map[string]interface{}
	if err := json.Unmarshal(jsonContent, &input); err != nil {
		return nil, err
	}

	return internalLoadFromDict(input)
}

// LoadDatabookJSON loads a Databook from a JSON source.
func LoadDatabookJSON(jsonContent []byte) (*Databook, error) {
	var input []map[string]interface{}
	var internalInput []map[string]interface{}
	if err := json.Unmarshal(jsonContent, &input); err != nil {
		return nil, err
	}

	db := NewDatabook()
	for _, d := range input {
		b, err := json.Marshal(d["data"])
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal(b, &internalInput); err != nil {
			return nil, err
		}

		if ds, err := internalLoadFromDict(internalInput); err == nil {
			db.AddSheet(d["title"].(string), ds)
		} else {
			return nil, err
		}
	}

	return db, nil
}

// JSON returns a JSON representation of the Dataset as string.
func (d *Dataset) JSON() (string, error) {
	back := d.Dict()

	b, err := json.Marshal(back)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// JSON returns a JSON representation of the Databook as string.
func (d *Databook) JSON() (string, error) {
	str := "["
	for _, s := range d.sheets {
		str += "{\"title\": \"" + s.title + "\", \"data\": "
		js, err := s.dataset.JSON()
		if err != nil {
			return "", err
		}
		str += js + "},"
	}
	str = str[:len(str)-1] + "]"
	return str, nil
}
