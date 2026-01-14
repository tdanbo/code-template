package taskfile

import (
	yamlhelper "code-template/helpers/yaml"
)

const taskfilePath = "Taskfile.yml"

// HasTask checks if Taskfile.yml has a specific task.
func HasTask(taskName string) (bool, error) {
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

// AddTask adds a task to Taskfile.yml.
// Creates the file with version "3" if it doesn't exist.
func AddTask(taskName, description string, commands []string) error {
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

	// Add the task
	tasks[taskName] = map[string]any{
		"desc": description,
		"cmds": commands,
	}

	data["tasks"] = tasks
	return yamlhelper.WriteYAML(taskfilePath, data)
}

// RemoveTask removes a task from Taskfile.yml.
func RemoveTask(taskName string) error {
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
