package models

type Module interface {
	GetName() string
	GetCategory() string
	GetVersion() int
	IsInstalled() bool
	Install() bool
	Uninstall() bool
}
