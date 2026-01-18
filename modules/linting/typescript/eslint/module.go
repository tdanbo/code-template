package eslint

import (
	_ "embed"
	"os"

	yamlhelper "code-template/helpers/yaml"
	"code-template/services"
)

//go:embed configs/eslint.config.js
var eslintConfig []byte

const (
	eslintConfigFile     = "eslint.config.js"
	codeTemplateFileName = "code-template.yml"
	moduleKey            = "eslint"
)

var Module = &ESLintModule{
	Name:     "eslint",
	Version:  1,
	Category: "linting",
	Path:     "linting/typescript/eslint",
}

type ESLintModule struct {
	Name     string
	Version  int
	Category string
	Path     string
}

func (m *ESLintModule) GetName() string {
	return m.Name
}

func (m *ESLintModule) GetCategory() string {
	return m.Category
}

func (m *ESLintModule) GetPath() string {
	return m.Path
}

func (m *ESLintModule) GetVersion() int {
	return m.Version
}

func (m *ESLintModule) GetKey() string {
	return moduleKey
}

// IsInstalled checks all conditions:
// 1. eslint.config.js exists
// 2. code-template.yml has eslint entry
func (m *ESLintModule) IsInstalled() bool {
	// Check 1: eslint.config.js exists
	if _, err := os.Stat(eslintConfigFile); os.IsNotExist(err) {
		return false
	}

	// Check 2: code-template.yml has eslint entry
	hasEntry, err := yamlhelper.HasKey(codeTemplateFileName, moduleKey)
	if err != nil || !hasEntry {
		return false
	}

	return true
}

// Install performs installation steps with rollback on failure.
func (m *ESLintModule) Install() bool {
	// Step 1: Check npm is available (for npx eslint usage)
	if !services.NPM.IsAvailable() {
		return false
	}

	// Step 2: Write eslint.config.js
	// Note: Dependencies (@eslint/js, typescript-eslint) should be in project's package.json
	// and will be resolved via npx eslint or local node_modules
	if err := os.WriteFile(eslintConfigFile, eslintConfig, 0644); err != nil {
		return false
	}

	// Step 3: Add entry to code-template.yml
	if err := yamlhelper.SetKey(codeTemplateFileName, moduleKey, m.Version); err != nil {
		os.Remove(eslintConfigFile)
		return false
	}

	return true
}

// Uninstall removes the config file and YAML entry.
// Does not uninstall global npm packages as they may be used by other projects.
func (m *ESLintModule) Uninstall() bool {
	success := true

	// Remove eslint.config.js
	if err := os.Remove(eslintConfigFile); err != nil && !os.IsNotExist(err) {
		success = false
	}

	// Remove entry from code-template.yml
	if err := yamlhelper.RemoveKey(codeTemplateFileName, moduleKey); err != nil {
		success = false
	}

	return success
}
