package frontenddesign

import (
	_ "embed"
	"os"
	"path/filepath"

	yamlhelper "code-template/helpers/yaml"
)

//go:embed SKILL.md
var skillContent []byte

const (
	moduleKey            = "skill-frontend-design"
	codeTemplateFileName = "code-template.yml"
	skillFileName        = "SKILL.md"
	skillDir             = ".claude/skills/frontend-design"
)

var Module = &FrontendDesignModule{
	Name:     "frontend-design",
	Version:  1,
	Category: "claude",
	Path:     "claude/skills/frontend-design",
}

type FrontendDesignModule struct {
	Name     string
	Version  int
	Category string
	Path     string
}

func (m *FrontendDesignModule) GetName() string {
	return m.Name
}

func (m *FrontendDesignModule) GetCategory() string {
	return m.Category
}

func (m *FrontendDesignModule) GetPath() string {
	return m.Path
}

func (m *FrontendDesignModule) GetVersion() int {
	return m.Version
}

func (m *FrontendDesignModule) GetKey() string {
	return moduleKey
}

// getSkillPath returns the full path to the skill file
func getSkillPath() string {
	return filepath.Join(skillDir, skillFileName)
}

// IsInstalled checks:
// 1. .claude/skills/frontend-design/SKILL.md exists
// 2. Entry exists in code-template.yml
func (m *FrontendDesignModule) IsInstalled() bool {
	// Check 1: Skill file exists
	if _, err := os.Stat(getSkillPath()); os.IsNotExist(err) {
		return false
	}

	// Check 2: code-template.yml entry
	hasEntry, err := yamlhelper.HasKey(codeTemplateFileName, moduleKey)
	if err != nil || !hasEntry {
		return false
	}

	return true
}

// Install copies the skill to .claude/skills/frontend-design/ and registers in code-template.yml
func (m *FrontendDesignModule) Install() bool {
	// Step 1: Create .claude/skills/frontend-design directory if it doesn't exist
	if err := os.MkdirAll(skillDir, 0755); err != nil {
		return false
	}

	// Step 2: Write skill file
	if err := os.WriteFile(getSkillPath(), skillContent, 0644); err != nil {
		return false
	}

	// Step 3: Add entry to code-template.yml
	if err := yamlhelper.SetKey(codeTemplateFileName, moduleKey, m.Version); err != nil {
		os.Remove(getSkillPath()) // Rollback
		return false
	}

	return true
}

// Uninstall removes the skill file and code-template.yml entry
func (m *FrontendDesignModule) Uninstall() bool {
	success := true

	// Step 1: Remove skill file
	if err := os.Remove(getSkillPath()); err != nil && !os.IsNotExist(err) {
		success = false
	}

	// Step 2: Remove entry from code-template.yml
	if err := yamlhelper.RemoveKey(codeTemplateFileName, moduleKey); err != nil {
		success = false
	}

	return success
}
