package services

// Package represents a dependency that can be installed.
type Package struct {
	Name        string // Binary/command name (e.g., "golangci-lint")
	InstallPath string // Install path (e.g., "github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest")
}

// PackageService defines the interface for managing packages.
type PackageService interface {
	// Name returns the service name (e.g., "go", "npm").
	Name() string

	// IsAvailable checks if the package manager is installed and available.
	IsAvailable() bool

	// IsInstalled checks if a package is installed (by binary name).
	IsInstalled(binaryName string) bool

	// Install installs a package. Returns error if installation fails.
	Install(pkg Package) error

	// Uninstall removes a package. Returns error if uninstallation fails.
	Uninstall(binaryName string) error
}
