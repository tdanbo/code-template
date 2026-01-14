package golinttask

import (
	"code-template/helpers/taskfile"
	yamlhelper "code-template/helpers/yaml"
)

const (
	moduleKey            = "task-go-lint"
	codeTemplateFileName = "code-template.yml"
	taskName             = "go-lint"
	taskDesc             = "Run golangci-lint"
)

var taskCommands = []string{"./.bin/golangci-lint run ./..."}

var Module = &GoLintTaskModule{
	Name:     "go-lint",
	Version:  1,
	Category: "tasks",
}

type GoLintTaskModule struct {
	Name     string
	Version  int
	Category string
}

func (m *GoLintTaskModule) GetName() string {
	return m.Name
}

func (m *GoLintTaskModule) GetCategory() string {
	return m.Category
}

func (m *GoLintTaskModule) GetVersion() int {
	return m.Version
}

func (m *GoLintTaskModule) GetKey() string {
	return moduleKey
}

// IsInstalled checks:
// 1. code-template.yml has task-go-lint entry
// 2. Taskfile.yml has go-lint task
func (m *GoLintTaskModule) IsInstalled() bool {
	// Check 1: code-template.yml entry
	hasEntry, err := yamlhelper.HasKey(codeTemplateFileName, moduleKey)
	if err != nil || !hasEntry {
		return false
	}

	// Check 2: Taskfile.yml has go-lint task
	hasTask, err := taskfile.HasTask(taskName)
	if err != nil || !hasTask {
		return false
	}

	return true
}

// Install adds the go-lint task
func (m *GoLintTaskModule) Install() bool {
	// Step 1: Add go-lint task to Taskfile.yml
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

// Uninstall removes the go-lint task
func (m *GoLintTaskModule) Uninstall() bool {
	success := true

	// Step 1: Remove go-lint task from Taskfile.yml
	if err := taskfile.RemoveTask(taskName); err != nil {
		success = false
	}

	// Step 2: Remove entry from code-template.yml
	if err := yamlhelper.RemoveKey(codeTemplateFileName, moduleKey); err != nil {
		success = false
	}

	return success
}
