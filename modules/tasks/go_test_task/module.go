package gotesttask

import (
	"code-template/helpers/taskfile"
	yamlhelper "code-template/helpers/yaml"
)

const (
	moduleKey            = "task-go-test"
	codeTemplateFileName = "code-template.yml"
	taskName             = "go-test"
	taskDesc             = "Run Go tests"
)

var taskCommands = []string{"go test ./..."}

var Module = &GoTestTaskModule{
	Name:     "go-test",
	Version:  1,
	Category: "tasks",
}

type GoTestTaskModule struct {
	Name     string
	Version  int
	Category string
}

func (m *GoTestTaskModule) GetName() string {
	return m.Name
}

func (m *GoTestTaskModule) GetCategory() string {
	return m.Category
}

func (m *GoTestTaskModule) GetVersion() int {
	return m.Version
}

func (m *GoTestTaskModule) GetKey() string {
	return moduleKey
}

// IsInstalled checks:
// 1. code-template.yml has task-go-test entry
// 2. Taskfile.yml has go-test task
func (m *GoTestTaskModule) IsInstalled() bool {
	// Check 1: code-template.yml entry
	hasEntry, err := yamlhelper.HasKey(codeTemplateFileName, moduleKey)
	if err != nil || !hasEntry {
		return false
	}

	// Check 2: Taskfile.yml has go-test task
	hasTask, err := taskfile.HasTask(taskName)
	if err != nil || !hasTask {
		return false
	}

	return true
}

// Install adds the go-test task
func (m *GoTestTaskModule) Install() bool {
	// Step 1: Add go-test task to Taskfile.yml
	if err := taskfile.AddTask(taskName, taskDesc, taskCommands); err != nil {
		return false
	}

	// Step 2: Add entry to code-template.yml
	if err := yamlhelper.SetKey(codeTemplateFileName, moduleKey, m.Version); err != nil {
		taskfile.RemoveTask(taskName) // Rollback
		return false
	}

	return true
}

// Uninstall removes the go-test task
func (m *GoTestTaskModule) Uninstall() bool {
	success := true

	// Step 1: Remove go-test task from Taskfile.yml
	if err := taskfile.RemoveTask(taskName); err != nil {
		success = false
	}

	// Step 2: Remove entry from code-template.yml
	if err := yamlhelper.RemoveKey(codeTemplateFileName, moduleKey); err != nil {
		success = false
	}

	return success
}
