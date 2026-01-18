package tddguard

import (
	yamlhelper "code-template/helpers/yaml"
)

const (
	moduleKey            = "tdd-guard"
	codeTemplateFileName = "code-template.yml"
)

var Module = &TddGuardModule{
	Name:     "tdd-guard",
	Version:  1,
	Category: "claude",
	Path:     "claude/workflow/tdd_guard",
}

type TddGuardModule struct {
	Name     string
	Version  int
	Category string
	Path     string
}

func (m *TddGuardModule) GetName() string {
	return m.Name
}

func (m *TddGuardModule) GetCategory() string {
	return m.Category
}

func (m *TddGuardModule) GetPath() string {
	return m.Path
}

func (m *TddGuardModule) GetVersion() int {
	return m.Version
}

func (m *TddGuardModule) GetKey() string {
	return moduleKey
}

// IsInstalled checks:
// 1. tdd-guard command exists in PATH
// 2. Hooks are configured in .claude/settings.json
// 3. Entry exists in code-template.yml
func (m *TddGuardModule) IsInstalled() bool {
	// Check 1: tdd-guard binary in PATH
	if !IsTddGuardInstalled() {
		return false
	}

	// Check 2: Hooks configured
	if !AreHooksConfigured() {
		return false
	}

	// Check 3: code-template.yml entry
	hasEntry, err := yamlhelper.HasKey(codeTemplateFileName, moduleKey)
	if err != nil || !hasEntry {
		return false
	}

	return true
}

// Install performs installation with rollback on failure
func (m *TddGuardModule) Install() bool {
	// Step 1: Check npm is available
	if !CheckNpmInstalled() {
		return false
	}

	// Step 2: Install tdd-guard via npm (if not already present)
	if !IsTddGuardInstalled() {
		if err := InstallTddGuard(); err != nil {
			return false
		}
	}

	// Step 3: Configure hooks in .claude/settings.json
	if err := AddHooks(); err != nil {
		// Note: We do NOT uninstall npm package per requirements
		return false
	}

	// Step 4: Add entry to code-template.yml
	if err := yamlhelper.SetKey(codeTemplateFileName, moduleKey, m.Version); err != nil {
		RemoveHooks() // Best-effort rollback
		return false
	}

	return true
}

// Uninstall removes configuration (but NOT the npm package)
func (m *TddGuardModule) Uninstall() bool {
	success := true

	// Step 1: Remove hooks from settings.json
	if err := RemoveHooks(); err != nil {
		success = false
	}

	// Step 2: Remove entry from code-template.yml
	if err := yamlhelper.RemoveKey(codeTemplateFileName, moduleKey); err != nil {
		success = false
	}

	// Note: Do NOT uninstall npm package (user might use elsewhere)

	return success
}
