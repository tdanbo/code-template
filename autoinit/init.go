package autoinit

import (
	"fmt"
	"os"
)

type InitError struct {
	Step    string
	Message string
	Err     error
}

func (e *InitError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func configExists() bool {
	_, err := os.Stat(codeTemplateFileName)
	return err == nil
}

func Run() error {
	// Always check for task - it's a hard requirement
	if !isTaskGloballyAvailable() {
		printError("go-task is not installed")
		return &InitError{
			Step:    "task_check",
			Message: "go-task is required but not found in PATH. Install with: go install github.com/go-task/task/v3/cmd/task@latest",
		}
	}

	if configExists() {
		// Run migration for existing installations
		if err := MigrateTasks(); err != nil {
			printWarning("Failed to migrate tasks: " + err.Error())
		}
		return nil
	}

	printWelcome()

	if !isGoInstalled() {
		printError("Go compiler not found in PATH")
		return &InitError{
			Step:    "go_check",
			Message: "Go compiler not found in PATH. Please install Go first: https://go.dev/dl/",
		}
	}
	printSuccess("Go compiler found")

	if !isTaskGloballyAvailable() {
		if promptInstallTask() {
			if err := installTaskGlobally(); err != nil {
				printError("Failed to install go-task")
				return &InitError{
					Step:    "task_install",
					Message: "Failed to install go-task globally",
					Err:     err,
				}
			}
			printSuccess("go-task installed globally")
		} else {
			printWarning("Skipping go-task installation. You'll need to use ./.bin/task instead of task")
		}
	} else {
		printSuccess("go-task found in PATH")
	}

	detected := detectInstalledModules()
	if len(detected) > 0 {
		printInfo(fmt.Sprintf("Detected %d installed module(s)", len(detected)))
	}

	if err := createConfig(detected); err != nil {
		printError("Failed to create code-template.yml")
		return &InitError{
			Step:    "config_create",
			Message: "Failed to create code-template.yml",
			Err:     err,
		}
	}
	printSuccess("Created code-template.yml")

	fmt.Println()
	printSuccess("Initialization complete!")
	fmt.Println()

	return nil
}
