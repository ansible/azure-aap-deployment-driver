package templates

import (
	"encoding/json"
	"os"
)

// Read in template dependency file
func readDependencyJSON(path string) ([]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var contents []string
	err = json.Unmarshal(data, &contents)
	if err != nil {
		return nil, err
	}
	return contents, nil
}

// Read template content
func readJSON(path string) (map[string]interface{}, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	contents := make(map[string]interface{})
	err = json.Unmarshal(data, &contents)
	if err != nil {
		return nil, err
	}
	return contents, nil
}
