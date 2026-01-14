package models

type Module interface {
	GetName() string
	GetCategory() string
	GetVersion() int
	GetKey() string
	IsInstalled() bool
	Install() bool
	Uninstall() bool
}
