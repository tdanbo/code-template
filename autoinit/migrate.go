package autoinit

import (
	yamlhelper "code-template/helpers/yaml"
)

// MigrateTasks detects existing tasks in Taskfile.yml and adds
// corresponding entries to code-template.yml for the new task modules.
func MigrateTasks() error {
	// Check for go-lint task
	if hasTask, _ := hasTaskfileTask("go-lint"); hasTask {
		if hasEntry, _ := yamlhelper.HasKey(codeTemplateFileName, "task-go-lint"); !hasEntry {
			if err := yamlhelper.SetKey(codeTemplateFileName, "task-go-lint", 1); err != nil {
				return err
			}
		}
	}

	// Check for go-test task
	if hasTask, _ := hasTaskfileTask("go-test"); hasTask {
		if hasEntry, _ := yamlhelper.HasKey(codeTemplateFileName, "task-go-test"); !hasEntry {
			if err := yamlhelper.SetKey(codeTemplateFileName, "task-go-test", 1); err != nil {
				return err
			}
		}
	}

	// Check for tdd-test task
	if hasTask, _ := hasTaskfileTask("tdd-test"); hasTask {
		if hasEntry, _ := yamlhelper.HasKey(codeTemplateFileName, "task-tdd-test"); !hasEntry {
			if err := yamlhelper.SetKey(codeTemplateFileName, "task-tdd-test", 1); err != nil {
				return err
			}
		}
	}

	return nil
}
