package helpers

import (
	"code-template/models"
	tddguard "code-template/modules/claude/workflow/tdd_guard"
	gotstwwailsreact "code-template/modules/language/go/go_ts_tw_wails_react"
	golangcilint "code-template/modules/linting/go/golangci_lint"
	eslint "code-template/modules/linting/typescript/eslint"
	golinttask "code-template/modules/tasks/go/go_lint_task"
	gotesttask "code-template/modules/tasks/go/go_test_task"
	tslinttask "code-template/modules/tasks/typescript/ts_lint_task"
	tstesttask "code-template/modules/tasks/typescript/ts_test_task"
	wailsdev "code-template/modules/tasks/wails/wails_dev"
)

// GetModules returns the list of all available modules.
func GetModules() []models.Module {
	return []models.Module{
		tddguard.Module,
		gotstwwailsreact.Module,
		golangcilint.Module,
		eslint.Module,
		golinttask.Module,
		gotesttask.Module,
		tslinttask.Module,
		tstesttask.Module,
		wailsdev.Module,
	}
}

// GetContent sets up the view based on the added modules.
// Deprecated: Use GetModules() with BuildTree() instead.
func GetContent() map[string]map[string]models.Module {
	view := map[string]map[string]models.Module{}

	moduleList := GetModules()

	for _, module := range moduleList {
		if view[module.GetCategory()] == nil {
			view[module.GetCategory()] = map[string]models.Module{}
		}
		view[module.GetCategory()][module.GetName()] = module
	}

	return view
}
