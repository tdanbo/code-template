package autoinit

import (
	"os"
	"path/filepath"

	yamlhelper "code-template/helpers/yaml"
)

const (
	binDir               = ".bin"
	golangciBinaryName   = "golangci-lint"
	golangciConfigFile   = ".golangci.yml"
	taskfileFile         = "Taskfile.yml"
	codeTemplateFileName = "code-template.yml"
)

func detectInstalledModules() map[string]int {
	detected := make(map[string]int)

	if isGolangciInstalled() {
		detected["golangci"] = 1
	}

	// Detect task modules
	if hasTask, _ := hasTaskfileTask("go-lint"); hasTask {
		detected["task-go-lint"] = 1
	}
	if hasTask, _ := hasTaskfileTask("go-test"); hasTask {
		detected["task-go-test"] = 1
	}
	if hasTask, _ := hasTaskfileTask("tdd-test"); hasTask {
		detected["task-tdd-test"] = 1
	}

	return detected
}

func isGolangciInstalled() bool {
	if _, err := os.Stat(golangciConfigFile); os.IsNotExist(err) {
		return false
	}

	binPath := filepath.Join(binDir, golangciBinaryName)
	if _, err := os.Stat(binPath); os.IsNotExist(err) {
		return false
	}

	return true
}

func hasTaskfileTask(taskName string) (bool, error) {
	data, err := yamlhelper.ReadYAML(taskfileFile)
	if err != nil {
		return false, err
	}

	tasks, ok := data["tasks"].(map[string]any)
	if !ok {
		return false, nil
	}

	_, exists := tasks[taskName]
	return exists, nil
}

func createConfig(modules map[string]int) error {
	data := make(map[string]any)
	for key, version := range modules {
		data[key] = version
	}
	return yamlhelper.WriteYAML(codeTemplateFileName, data)
}
