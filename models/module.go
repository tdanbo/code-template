package models

type Module interface {
	GetName() string
	GetCategory() string
	GetPath() string // Hierarchical path, e.g., "linting/golangci_lint"
	GetVersion() int
	GetKey() string
	IsInstalled() bool
	Install() bool
	Uninstall() bool
}
