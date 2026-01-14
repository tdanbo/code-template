package helpers

import (
	"code-template/models"
	"code-template/modules/golangci_lint"
	golinttask "code-template/modules/tasks/go_lint_task"
	gotesttask "code-template/modules/tasks/go_test_task"
	tddtesttask "code-template/modules/tasks/tdd_test_task"
	tddguard "code-template/modules/tdd_guard"
)

// GetContent sets up the view based on the added modules.
func GetContent() map[string]map[string]models.Module {
	view := map[string]map[string]models.Module{}

	moduleList := []models.Module{
		tddguard.Module,
		golangci_lint.Module,
		golinttask.Module,
		gotesttask.Module,
		tddtesttask.Module,
	}

	for _, module := range moduleList {
		if view[module.GetCategory()] == nil {
			view[module.GetCategory()] = map[string]models.Module{}
		}
		view[module.GetCategory()][module.GetName()] = module
	}

	return view
}
