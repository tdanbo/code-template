package golangci_lint

import (
	_ "embed"
	"os"

	yamlhelper "code-template/helpers/yaml"
)

//go:embed golangci.yml
var golangciConfig []byte

const (
	golangciFileName     = ".golangci.yml"
	codeTemplateFileName = "code-template.yml"
	moduleKey            = "golangci"
)

var Module = &GolangciLintModule{
	Name:     "golangci",
	Version:  1,
	Category: "go",
}

type GolangciLintModule struct {
	Name     string
	Version  int
	Category string
}

func (m *GolangciLintModule) GetName() string {
	return m.Name
}

func (m *GolangciLintModule) GetCategory() string {
	return m.Category
}

func (m *GolangciLintModule) GetVersion() int {
	return m.Version
}

func (m *GolangciLintModule) GetKey() string {
	return moduleKey
}

// IsInstalled checks all conditions:
// 1. .golangci.yml exists
// 2. code-template.yml has golangci entry
// 3. Required binaries are installed in .bin/
func (m *GolangciLintModule) IsInstalled() bool {
	// Check 1: .golangci.yml exists
	if _, err := os.Stat(golangciFileName); os.IsNotExist(err) {
		return false
	}

	// Check 2: code-template.yml has golangci entry
	hasEntry, err := yamlhelper.HasKey(codeTemplateFileName, moduleKey)
	if err != nil || !hasEntry {
		return false
	}

	// Check 3: Binaries are installed
	if !AreBinariesInstalled() {
		return false
	}

	return true
}

// Install performs installation steps with rollback on failure.
func (m *GolangciLintModule) Install() bool {
	// Step 0: Check Go is installed
	if !CheckGoInstalled() {
		return false
	}

	// Step 1: Create .bin directory
	if err := EnsureBinDir(); err != nil {
		return false
	}

	// Step 2: Install binaries
	installed, err := InstallBinaries()
	if err != nil {
		RollbackBinaries(installed)
		return false
	}

	// Step 3: Add .bin/ to .gitignore
	if err := AddToGitignore(); err != nil {
		RollbackBinaries(installed)
		return false
	}

	// Step 4: Copy golangci.yml to .golangci.yml
	if err := os.WriteFile(golangciFileName, golangciConfig, 0644); err != nil {
		RollbackBinaries(installed)
		RemoveFromGitignore()
		return false
	}

	// Step 5: Add entry to code-template.yml
	if err := yamlhelper.SetKey(codeTemplateFileName, moduleKey, m.Version); err != nil {
		os.Remove(golangciFileName)
		RollbackBinaries(installed)
		RemoveFromGitignore()
		return false
	}

	return true
}

// Uninstall removes all installed components (best-effort cleanup).
func (m *GolangciLintModule) Uninstall() bool {
	success := true

	// Step 1: Remove .golangci.yml
	if err := os.Remove(golangciFileName); err != nil && !os.IsNotExist(err) {
		success = false
	}

	// Step 2: Remove entry from code-template.yml
	if err := yamlhelper.RemoveKey(codeTemplateFileName, moduleKey); err != nil {
		success = false
	}

	// Step 3: Remove binaries from .bin/
	if err := RemoveAllBinaries(); err != nil {
		success = false
	}

	// Step 4: Remove .bin/ from .gitignore
	if err := RemoveFromGitignore(); err != nil {
		success = false
	}

	return success
}
