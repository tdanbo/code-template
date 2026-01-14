package tddtesttask

import (
	"code-template/helpers/taskfile"
	yamlhelper "code-template/helpers/yaml"
)

const (
	moduleKey            = "task-tdd-test"
	codeTemplateFileName = "code-template.yml"
	taskName             = "tdd-test"
	taskDesc             = "Run Go tests with TDD Guard"
)

var taskCommands = []string{"go test -json ./... 2>&1 | tdd-guard-go -project-root $(pwd)"}

var Module = &TddTestTaskModule{
	Name:     "tdd-test",
	Version:  1,
	Category: "tasks",
}

type TddTestTaskModule struct {
	Name     string
	Version  int
	Category string
}

func (m *TddTestTaskModule) GetName() string {
	return m.Name
}

func (m *TddTestTaskModule) GetCategory() string {
	return m.Category
}

func (m *TddTestTaskModule) GetVersion() int {
	return m.Version
}

func (m *TddTestTaskModule) GetKey() string {
	return moduleKey
}

// IsInstalled checks:
// 1. code-template.yml has task-tdd-test entry
// 2. Taskfile.yml has tdd-test task
func (m *TddTestTaskModule) IsInstalled() bool {
	// Check 1: code-template.yml entry
	hasEntry, err := yamlhelper.HasKey(codeTemplateFileName, moduleKey)
	if err != nil || !hasEntry {
		return false
	}

	// Check 2: Taskfile.yml has tdd-test task
	hasTask, err := taskfile.HasTask(taskName)
	if err != nil || !hasTask {
		return false
	}

	return true
}

// Install adds the tdd-test task
func (m *TddTestTaskModule) Install() bool {
	// Step 1: Add tdd-test task to Taskfile.yml
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

// Uninstall removes the tdd-test task
func (m *TddTestTaskModule) Uninstall() bool {
	success := true

	// Step 1: Remove tdd-test task from Taskfile.yml
	if err := taskfile.RemoveTask(taskName); err != nil {
		success = false
	}

	// Step 2: Remove entry from code-template.yml
	if err := yamlhelper.RemoveKey(codeTemplateFileName, moduleKey); err != nil {
		success = false
	}

	return success
}
