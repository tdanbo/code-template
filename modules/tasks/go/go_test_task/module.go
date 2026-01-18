package gotesttask

import (
	"os/exec"

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
	Path:     "tasks/go/go_test_task",
}

type GoTestTaskModule struct {
	Name     string
	Version  int
	Category string
	Path     string
}

func (m *GoTestTaskModule) GetName() string {
	return m.Name
}

func (m *GoTestTaskModule) GetCategory() string {
	return m.Category
}

func (m *GoTestTaskModule) GetPath() string {
	return m.Path
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

// isTaskInstalled checks if go-task (task) is available in PATH.
func isTaskInstalled() bool {
	_, err := exec.LookPath("task")
	return err == nil
}

// Install adds the go-test task
func (m *GoTestTaskModule) Install() bool {
	// Step 1: Check if go-task is installed
	if !isTaskInstalled() {
		return false
	}

	// Step 2: Add go-test task to Taskfile.yml
	if err := taskfile.AddTask(taskName, taskDesc, taskCommands); err != nil {
		return false
	}

	// Step 3: Add entry to code-template.yml
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
