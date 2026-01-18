package gotstwwailsreact

import (
	_ "embed"
	"fmt"
	"os"

	yamlhelper "code-template/helpers/yaml"
)

const (
	wailsJSONFile        = "wails.json"
	frontendPackageJSON  = "frontend/package.json"
	codeTemplateFileName = "code-template.yml"
	moduleKey            = "go-ts-tw-wails-react"
)

var Module = &WailsReactTSModule{
	Name:     "go-ts-tw-wails-react",
	Version:  1,
	Category: "language",
	Path:     "language/go/go_ts_tw_wails_react",
}

type WailsReactTSModule struct {
	Name     string
	Version  int
	Category string
	Path     string
}

func (m *WailsReactTSModule) GetName() string {
	return m.Name
}

func (m *WailsReactTSModule) GetCategory() string {
	return m.Category
}

func (m *WailsReactTSModule) GetPath() string {
	return m.Path
}

func (m *WailsReactTSModule) GetVersion() int {
	return m.Version
}

func (m *WailsReactTSModule) GetKey() string {
	return moduleKey
}

// IsInstalled checks all conditions:
// 1. wails.json exists
// 2. frontend/package.json exists
// 3. code-template.yml has wails-react-ts entry
func (m *WailsReactTSModule) IsInstalled() bool {
	// Check 1: wails.json exists
	if _, err := os.Stat(wailsJSONFile); os.IsNotExist(err) {
		return false
	}

	// Check 2: frontend/package.json exists
	if _, err := os.Stat(frontendPackageJSON); os.IsNotExist(err) {
		return false
	}

	// Check 3: code-template.yml has wails-react-ts entry
	hasEntry, err := yamlhelper.HasKey(codeTemplateFileName, moduleKey)
	if err != nil || !hasEntry {
		return false
	}

	return true
}

// Install performs installation steps with rollback on failure.
func (m *WailsReactTSModule) Install() bool {
	// Step 1: Check Go is installed
	if !CheckGoInstalled() {
		return false
	}

	// Step 2: Check npm is installed
	if !CheckNpmInstalled() {
		return false
	}

	// Step 3: Check Linux system dependencies (gcc, gtk, webkit)
	if ok, missing := CheckLinuxDependencies(); !ok {
		fmt.Printf("Missing Linux dependency: %s\n", missing)
		fmt.Println("Run 'wails doctor' for installation instructions")
		return false
	}

	// Step 4: Install Wails CLI if not present
	if !IsWailsInstalled() {
		if err := InstallWailsCLI(); err != nil {
			return false
		}
	}

	// Step 5: Scaffold Wails project
	projectName := getProjectName()
	if err := ScaffoldProject(projectName); err != nil {
		return false
	}

	// Step 6: Configure webkit2gtk-4.1 if needed (newer Linux distros)
	if NeedsWebkit41BuildTag() {
		if err := ConfigureWebkit41(); err != nil {
			RollbackScaffold()
			return false
		}
	}

	// Step 7: Copy config files to frontend/
	if err := CopyConfigFiles(); err != nil {
		RollbackScaffold()
		return false
	}

	// Step 8: Add Tailwind import to main.tsx
	if err := AddTailwindImport(); err != nil {
		RollbackConfigFiles()
		RollbackScaffold()
		return false
	}

	// Step 9: Upgrade TypeScript to latest (required for tsconfig.json)
	if err := UpgradeTypeScript(); err != nil {
		RollbackConfigFiles()
		RollbackScaffold()
		return false
	}

	// Step 10: Install Tailwind CSS and dependencies
	if err := InstallTailwind(); err != nil {
		RollbackConfigFiles()
		RollbackScaffold()
		return false
	}

	// Step 11: Install frontend dependencies
	if err := InstallFrontendDeps(); err != nil {
		RollbackConfigFiles()
		RollbackScaffold()
		return false
	}

	// Step 12: Add entry to code-template.yml
	if err := yamlhelper.SetKey(codeTemplateFileName, moduleKey, m.Version); err != nil {
		RollbackConfigFiles()
		RollbackScaffold()
		return false
	}

	return true
}

// Uninstall removes only the code-template.yml entry (preserves user code).
func (m *WailsReactTSModule) Uninstall() bool {
	// Only remove entry from code-template.yml
	// Do NOT delete project files as user may have written code
	if err := yamlhelper.RemoveKey(codeTemplateFileName, moduleKey); err != nil {
		return false
	}
	return true
}
