package wailsdev

import (
	"os/exec"

	"code-template/helpers/taskfile"
	yamlhelper "code-template/helpers/yaml"
)

const (
	moduleKey            = "task-wails-dev"
	codeTemplateFileName = "code-template.yml"
	taskName             = "dev"
	taskDesc             = "Run Wails development server"
)

var taskCommands = []string{"wails dev"}

var Module = &WailsDevTaskModule{
	Name:     "dev",
	Version:  1,
	Category: "tasks",
	Path:     "tasks/wails/wails_dev",
}

type WailsDevTaskModule struct {
	Name     string
	Version  int
	Category string
	Path     string
}

func (m *WailsDevTaskModule) GetName() string {
	return m.Name
}

func (m *WailsDevTaskModule) GetCategory() string {
	return m.Category
}

func (m *WailsDevTaskModule) GetPath() string {
	return m.Path
}

func (m *WailsDevTaskModule) GetVersion() int {
	return m.Version
}

func (m *WailsDevTaskModule) GetKey() string {
	return moduleKey
}

// IsInstalled checks:
// 1. code-template.yml has task-wails-dev entry
// 2. Taskfile.yml has dev task
func (m *WailsDevTaskModule) IsInstalled() bool {
	// Check 1: code-template.yml entry
	hasEntry, err := yamlhelper.HasKey(codeTemplateFileName, moduleKey)
	if err != nil || !hasEntry {
		return false
	}

	// Check 2: Taskfile.yml has dev task
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

// isWailsInstalled checks if wails CLI is available in PATH.
func isWailsInstalled() bool {
	_, err := exec.LookPath("wails")
	return err == nil
}

// Install adds the dev task
func (m *WailsDevTaskModule) Install() bool {
	// Step 1: Check if go-task is installed
	if !isTaskInstalled() {
		return false
	}

	// Step 2: Check if wails is installed
	if !isWailsInstalled() {
		return false
	}

	// Step 3: Add dev task to Taskfile.yml
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

// Uninstall removes the dev task
func (m *WailsDevTaskModule) Uninstall() bool {
	success := true

	// Step 1: Remove dev task from Taskfile.yml
	if err := taskfile.RemoveTask(taskName); err != nil {
		success = false
	}

	// Step 2: Remove entry from code-template.yml
	if err := yamlhelper.RemoveKey(codeTemplateFileName, moduleKey); err != nil {
		success = false
	}

	return success
}
