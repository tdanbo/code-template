package tslinttask

import (
	"os/exec"

	"code-template/helpers/taskfile"
	yamlhelper "code-template/helpers/yaml"
)

const (
	moduleKey            = "task-ts-lint"
	codeTemplateFileName = "code-template.yml"
	taskName             = "ts-lint"
	taskDesc             = "Run ESLint on TypeScript files"
)

var taskCommands = []string{"npx eslint ."}

var Module = &TSLintTaskModule{
	Name:     "ts-lint",
	Version:  1,
	Category: "tasks",
	Path:     "tasks/typescript/ts_lint_task",
}

type TSLintTaskModule struct {
	Name     string
	Version  int
	Category string
	Path     string
}

func (m *TSLintTaskModule) GetName() string {
	return m.Name
}

func (m *TSLintTaskModule) GetCategory() string {
	return m.Category
}

func (m *TSLintTaskModule) GetPath() string {
	return m.Path
}

func (m *TSLintTaskModule) GetVersion() int {
	return m.Version
}

func (m *TSLintTaskModule) GetKey() string {
	return moduleKey
}

// IsInstalled checks:
// 1. code-template.yml has task-ts-lint entry
// 2. Taskfile.yml has ts-lint task
func (m *TSLintTaskModule) IsInstalled() bool {
	// Check 1: code-template.yml entry
	hasEntry, err := yamlhelper.HasKey(codeTemplateFileName, moduleKey)
	if err != nil || !hasEntry {
		return false
	}

	// Check 2: Taskfile.yml has ts-lint task
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

// Install adds the ts-lint task
func (m *TSLintTaskModule) Install() bool {
	// Step 1: Check if go-task is installed
	if !isTaskInstalled() {
		return false
	}

	// Step 2: Check if npm is installed
	if !isNpmInstalled() {
		return false
	}

	// Step 3: Add ts-lint task to Taskfile.yml
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

// Uninstall removes the ts-lint task
func (m *TSLintTaskModule) Uninstall() bool {
	success := true

	// Step 1: Remove ts-lint task from Taskfile.yml
	if err := taskfile.RemoveTask(taskName); err != nil {
		success = false
	}

	// Step 2: Remove entry from code-template.yml
	if err := yamlhelper.RemoveKey(codeTemplateFileName, moduleKey); err != nil {
		success = false
	}

	return success
}
