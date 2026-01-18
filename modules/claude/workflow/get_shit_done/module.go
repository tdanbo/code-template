package getshitdone

import (
	"os"
	"os/exec"

	yamlhelper "code-template/helpers/yaml"
)

const (
	moduleKey            = "get-shit-done"
	codeTemplateFileName = "code-template.yml"
	gsdPackage           = "get-shit-done-cc"
	gsdVersionFile       = ".claude/get-shit-done/VERSION"
)

var Module = &GetShitDoneModule{
	Name:     "get-shit-done",
	Version:  1,
	Category: "claude",
	Path:     "claude/workflow/get_shit_done",
}

type GetShitDoneModule struct {
	Name     string
	Version  int
	Category string
	Path     string
}

func (m *GetShitDoneModule) GetName() string {
	return m.Name
}

func (m *GetShitDoneModule) GetCategory() string {
	return m.Category
}

func (m *GetShitDoneModule) GetPath() string {
	return m.Path
}

func (m *GetShitDoneModule) GetVersion() int {
	return m.Version
}

func (m *GetShitDoneModule) GetKey() string {
	return moduleKey
}

// isNpxAvailable checks if npx is available in PATH
func isNpxAvailable() bool {
	_, err := exec.LookPath("npx")
	return err == nil
}

// isGsdInstalled checks if gsd is installed locally by checking for VERSION file
func isGsdInstalled() bool {
	_, err := os.Stat(gsdVersionFile)
	return err == nil
}

// installGsd runs npx get-shit-done-cc --local to install gsd locally
func installGsd() error {
	cmd := exec.Command("npx", gsdPackage, "--local")
	return cmd.Run()
}

// IsInstalled checks:
// 1. gsd command exists in PATH
// 2. Entry exists in code-template.yml
func (m *GetShitDoneModule) IsInstalled() bool {
	// Check 1: gsd binary in PATH
	if !isGsdInstalled() {
		return false
	}

	// Check 2: code-template.yml entry
	hasEntry, err := yamlhelper.HasKey(codeTemplateFileName, moduleKey)
	if err != nil || !hasEntry {
		return false
	}

	return true
}

// Install runs npx get-shit-done-cc to install gsd
func (m *GetShitDoneModule) Install() bool {
	// Step 1: Check npx is available
	if !isNpxAvailable() {
		return false
	}

	// Step 2: Install gsd via npx (if not already present)
	if !isGsdInstalled() {
		if err := installGsd(); err != nil {
			return false
		}
	}

	// Step 3: Add entry to code-template.yml
	if err := yamlhelper.SetKey(codeTemplateFileName, moduleKey, m.Version); err != nil {
		return false
	}

	return true
}

// Uninstall removes the code-template.yml entry (does NOT uninstall gsd)
func (m *GetShitDoneModule) Uninstall() bool {
	success := true

	// Remove entry from code-template.yml
	// Note: We do NOT uninstall gsd as it may be used elsewhere
	if err := yamlhelper.RemoveKey(codeTemplateFileName, moduleKey); err != nil {
		success = false
	}

	return success
}
