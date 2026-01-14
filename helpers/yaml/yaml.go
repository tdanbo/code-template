package yamlhelper

import (
	"os"

	"github.com/goccy/go-yaml"
)

// ReadYAML reads a YAML file into a map, returning empty map if file doesn't exist.
func ReadYAML(path string) (map[string]any, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return make(map[string]any), nil
	}
	if err != nil {
		return nil, err
	}

	var result map[string]any
	if len(data) == 0 {
		return make(map[string]any), nil
	}
	if err := yaml.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	if result == nil {
		result = make(map[string]any)
	}
	return result, nil
}

// WriteYAML writes a map to a YAML file.
func WriteYAML(path string, data map[string]any) error {
	out, err := yaml.Marshal(data)
	if err != nil {
		return err
	}
	return os.WriteFile(path, out, 0644)
}

// HasKey checks if a YAML file contains a specific top-level key.
func HasKey(path string, key string) (bool, error) {
	data, err := ReadYAML(path)
	if err != nil {
		return false, err
	}
	_, exists := data[key]
	return exists, nil
}

// SetKey sets a top-level key in a YAML file, preserving other keys.
func SetKey(path string, key string, value any) error {
	data, err := ReadYAML(path)
	if err != nil {
		return err
	}
	data[key] = value
	return WriteYAML(path, data)
}

// RemoveKey removes a top-level key from a YAML file.
func RemoveKey(path string, key string) error {
	data, err := ReadYAML(path)
	if err != nil {
		return err
	}
	delete(data, key)
	return WriteYAML(path, data)
}
