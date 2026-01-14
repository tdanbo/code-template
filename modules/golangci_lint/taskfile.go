package golangci_lint

import (
	yamlhelper "code-template/helpers/yaml"
)

const taskfilePath = "Taskfile.yml"
const taskName = "go-lint"

// HasGoLintTask checks if Taskfile.yml has a go-lint task.
func HasGoLintTask() (bool, error) {
	data, err := yamlhelper.ReadYAML(taskfilePath)
	if err != nil {
		return false, err
	}

	tasks, ok := data["tasks"].(map[string]any)
	if !ok {
		return false, nil
	}

	_, exists := tasks[taskName]
	return exists, nil
}

// AddGoLintTask adds the go-lint task to Taskfile.yml.
// Creates the file if it doesn't exist.
func AddGoLintTask() error {
	data, err := yamlhelper.ReadYAML(taskfilePath)
	if err != nil {
		return err
	}

	// Ensure version exists
	if _, ok := data["version"]; !ok {
		data["version"] = "3"
	}

	// Get or create tasks section
	tasks, ok := data["tasks"].(map[string]any)
	if !ok {
		tasks = make(map[string]any)
	}

	// Add go-lint task using local binary
	tasks[taskName] = map[string]any{
		"desc": "Run golangci-lint",
		"cmds": []string{"./.bin/golangci-lint run ./..."},
	}

	data["tasks"] = tasks
	return yamlhelper.WriteYAML(taskfilePath, data)
}

// RemoveGoLintTask removes the go-lint task from Taskfile.yml.
func RemoveGoLintTask() error {
	data, err := yamlhelper.ReadYAML(taskfilePath)
	if err != nil {
		return err
	}

	tasks, ok := data["tasks"].(map[string]any)
	if !ok {
		return nil // No tasks section, nothing to remove
	}

	delete(tasks, taskName)

	// If tasks is now empty, remove it entirely
	if len(tasks) == 0 {
		delete(data, "tasks")
	} else {
		data["tasks"] = tasks
	}

	return yamlhelper.WriteYAML(taskfilePath, data)
}
