package tddguard

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const (
	claudeDir        = ".claude"
	settingsFileName = "settings.json"
)

// HookEntry represents a single hook command
type HookEntry struct {
	Type    string `json:"type"`
	Command string `json:"command"`
}

// HookConfig represents a hook configuration with optional matcher
type HookConfig struct {
	Matcher string      `json:"matcher,omitempty"`
	Hooks   []HookEntry `json:"hooks"`
}

// getSettingsPath returns the full path to .claude/settings.json
func getSettingsPath() string {
	return filepath.Join(claudeDir, settingsFileName)
}

// getTddGuardHookConfig returns the hook configuration for tdd-guard
func getTddGuardHookConfig() map[string][]HookConfig {
	return map[string][]HookConfig{
		"PreToolUse": {
			{
				Matcher: "Write|Edit|MultiEdit|TodoWrite",
				Hooks:   []HookEntry{{Type: "command", Command: "tdd-guard"}},
			},
		},
		"UserPromptSubmit": {
			{
				Hooks: []HookEntry{{Type: "command", Command: "tdd-guard"}},
			},
		},
		"SessionStart": {
			{
				Matcher: "startup|resume|clear",
				Hooks:   []HookEntry{{Type: "command", Command: "tdd-guard"}},
			},
		},
	}
}

// ReadSettings reads .claude/settings.json and returns hooks and other fields separately
// Returns empty structures if file doesn't exist
func ReadSettings() (map[string][]HookConfig, map[string]any, error) {
	path := getSettingsPath()
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return make(map[string][]HookConfig), make(map[string]any), nil
	}
	if err != nil {
		return nil, nil, err
	}

	// Handle empty file
	if len(data) == 0 {
		return make(map[string][]HookConfig), make(map[string]any), nil
	}

	// Parse into raw map to preserve unknown fields
	var rawMap map[string]any
	if err := json.Unmarshal(data, &rawMap); err != nil {
		return nil, nil, err
	}

	// Extract and parse hooks
	hooks := make(map[string][]HookConfig)
	if hooksRaw, ok := rawMap["hooks"]; ok {
		hooksBytes, err := json.Marshal(hooksRaw)
		if err == nil {
			json.Unmarshal(hooksBytes, &hooks)
		}
		delete(rawMap, "hooks")
	}

	return hooks, rawMap, nil
}

// WriteSettings writes hooks and other fields back to .claude/settings.json
func WriteSettings(hooks map[string][]HookConfig, otherFields map[string]any) error {
	// Ensure .claude directory exists
	if err := os.MkdirAll(claudeDir, 0755); err != nil {
		return err
	}

	// Build output map
	outputMap := make(map[string]any)
	for k, v := range otherFields {
		outputMap[k] = v
	}
	if len(hooks) > 0 {
		outputMap["hooks"] = hooks
	}

	data, err := json.MarshalIndent(outputMap, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(getSettingsPath(), data, 0644)
}

// isTddGuardHook checks if a hook entry list contains a tdd-guard command
func isTddGuardHook(hooks []HookEntry) bool {
	for _, h := range hooks {
		if h.Type == "command" && h.Command == "tdd-guard" {
			return true
		}
	}
	return false
}

// AreHooksConfigured checks if all tdd-guard hooks are present in settings.json
func AreHooksConfigured() bool {
	hooks, _, err := ReadSettings()
	if err != nil {
		return false
	}

	requiredHooks := getTddGuardHookConfig()

	for hookType, required := range requiredHooks {
		existing, ok := hooks[hookType]
		if !ok {
			return false
		}

		// Check each required hook config exists
		for _, reqConfig := range required {
			found := false
			for _, existConfig := range existing {
				// Match by matcher and check for tdd-guard command
				if existConfig.Matcher == reqConfig.Matcher && isTddGuardHook(existConfig.Hooks) {
					found = true
					break
				}
			}
			if !found {
				return false
			}
		}
	}

	return true
}

// AddHooks adds tdd-guard hooks to settings.json, merging with existing hooks
func AddHooks() error {
	hooks, otherFields, err := ReadSettings()
	if err != nil {
		return err
	}

	tddGuardHooks := getTddGuardHookConfig()

	// Merge hooks - add tdd-guard hooks without removing existing ones
	for hookType, newConfigs := range tddGuardHooks {
		existing := hooks[hookType]

		for _, newConfig := range newConfigs {
			// Check if this exact config already exists
			alreadyExists := false
			for _, ex := range existing {
				if ex.Matcher == newConfig.Matcher && isTddGuardHook(ex.Hooks) {
					alreadyExists = true
					break
				}
			}
			if !alreadyExists {
				existing = append(existing, newConfig)
			}
		}
		hooks[hookType] = existing
	}

	return WriteSettings(hooks, otherFields)
}

// RemoveHooks removes only tdd-guard hooks from settings.json, preserving other hooks
func RemoveHooks() error {
	hooks, otherFields, err := ReadSettings()
	if err != nil {
		return err
	}

	// For each hook type, filter out tdd-guard hooks
	for hookType, configs := range hooks {
		var filtered []HookConfig
		for _, config := range configs {
			if !isTddGuardHook(config.Hooks) {
				filtered = append(filtered, config)
			}
		}
		if len(filtered) > 0 {
			hooks[hookType] = filtered
		} else {
			delete(hooks, hookType)
		}
	}

	return WriteSettings(hooks, otherFields)
}
