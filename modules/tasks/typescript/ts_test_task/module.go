package tstesttask

import (
	"os/exec"

	"code-template/helpers/taskfile"
	yamlhelper "code-template/helpers/yaml"
)

const (
	moduleKey            = "task-ts-test"
	codeTemplateFileName = "code-template.yml"
	taskName             = "ts-test"
	taskDesc             = "Run TypeScript tests"
)

var taskCommands = []string{"npm test"}

var Module = &TSTestTaskModule{
	Name:     "ts-test",
	Version:  1,
	Category: "tasks",
	Path:     "tasks/typescript/ts_test_task",
}

type TSTestTaskModule struct {
	Name     string
	Version  int
	Category string
	Path     string
}

func (m *TSTestTaskModule) GetName() string {
	return m.Name
}

func (m *TSTestTaskModule) GetCategory() string {
	return m.Category
}

func (m *TSTestTaskModule) GetPath() string {
	return m.Path
}

func (m *TSTestTaskModule) GetVersion() int {
	return m.Version
}

func (m *TSTestTaskModule) GetKey() string {
	return moduleKey
}

// IsInstalled checks:
// 1. code-template.yml has task-ts-test entry
// 2. Taskfile.yml has ts-test task
func (m *TSTestTaskModule) IsInstalled() bool {
	// Check 1: code-template.yml entry
	hasEntry, err := yamlhelper.HasKey(codeTemplateFileName, moduleKey)
	if err != nil || !hasEntry {
		return false
	}

	// Check 2: Taskfile.yml has ts-test task
	hasTask, err := taskfile.HasTask(taskName)
	if err != nil || !hasTask {
		return false
	}

	return true
}

// isTaskInstalled checks if go-task (task) is available in PATH.
func isTaskInstalled() bool {
	_, err := exec.LookPath("task")
	return err == nil
}

// isNpmInstalled checks if npm is available in PATH.
func isNpmInstalled() bool {
	_, err := exec.LookPath("npm")
	return err == nil
}

// Install adds the ts-test task
func (m *TSTestTaskModule) Install() bool {
	// Step 1: Check if go-task is installed
	if !isTaskInstalled() {
		return false
	}

	// Step 2: Check if npm is installed
	if !isNpmInstalled() {
		return false
	}

	// Step 3: Add ts-test task to Taskfile.yml
	if err := taskfile.AddTask(taskName, taskDesc, taskCommands); err != nil {
		return false
	}

	// Step 4: Add entry to code-template.yml
	if err := yamlhelper.SetKey(codeTemplateFileName, moduleKey, m.Version); err != nil {
		taskfile.RemoveTask(taskName) // Rollback
		return false
	}

	return true
}

// Uninstall removes the ts-test task
func (m *TSTestTaskModule) Uninstall() bool {
	success := true

	// Step 1: Remove ts-test task from Taskfile.yml
	if err := taskfile.RemoveTask(taskName); err != nil {
		success = false
	}

	// Step 2: Remove entry from code-template.yml
	if err := yamlhelper.RemoveKey(codeTemplateFileName, moduleKey); err != nil {
		success = false
	}

	return success
}
