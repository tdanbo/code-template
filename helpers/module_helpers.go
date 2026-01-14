package helpers

import (
	"code-template/models"

	yamlhelper "code-template/helpers/yaml"
)

const codeTemplateFileName = "code-template.yml"

// ModuleState represents the installation state of a module.
type ModuleState int

const (
	StateNotInstalled ModuleState = iota
	StateOutdated
	StateUpToDate
)

// GetInstalledVersion returns the installed version of a module.
// Returns 0 if not installed or on error.
func GetInstalledVersion(m models.Module) int {
	value, exists, err := yamlhelper.GetValue(codeTemplateFileName, m.GetKey())
	if err != nil || !exists {
		return 0
	}

	switch v := value.(type) {
	case int:
		return v
	case int64:
		return int(v)
	case float64:
		return int(v)
	case uint64:
		return int(v)
	default:
		return 0
	}
}

// IsOutdated returns true if the module is installed but at an older version.
func IsOutdated(m models.Module) bool {
	if !m.IsInstalled() {
		return false
	}
	installedVersion := GetInstalledVersion(m)
	return installedVersion < m.GetVersion()
}

// UpdateModule performs an update by uninstalling then reinstalling.
// Returns true if both operations succeed.
func UpdateModule(m models.Module) bool {
	if !m.Uninstall() {
		return false
	}
	return m.Install()
}

// GetModuleState returns the current state of a module.
func GetModuleState(m models.Module) ModuleState {
	if !m.IsInstalled() {
		return StateNotInstalled
	}
	if IsOutdated(m) {
		return StateOutdated
	}
	return StateUpToDate
}
