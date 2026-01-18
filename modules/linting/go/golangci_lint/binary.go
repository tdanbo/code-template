package golangci_lint

import (
	"code-template/services"
)

const (
	golangciBinary  = "golangci-lint"
	golangciInstall = "github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest"
)

var goService = services.Go

// GolangciPackage is the package definition for golangci-lint.
var GolangciPackage = services.Package{
	Name:        golangciBinary,
	InstallPath: golangciInstall,
}

// CheckGoInstalled verifies that Go is available in PATH.
func CheckGoInstalled() bool {
	return goService.IsAvailable()
}

// EnsureBinDir creates the .bin directory if it doesn't exist.
func EnsureBinDir() error {
	return goService.EnsureBinDir()
}

// AreBinariesInstalled checks if golangci-lint is available (locally or globally).
func AreBinariesInstalled() bool {
	return goService.IsInstalled(golangciBinary)
}

// InstallBinaries installs golangci-lint if not already available, returns list of installed for rollback.
func InstallBinaries() ([]string, error) {
	var installed []string

	// Only install golangci-lint if not available (locally or globally)
	if !goService.IsInstalled(golangciBinary) {
		if err := goService.Install(GolangciPackage); err != nil {
			return installed, err
		}
		installed = append(installed, golangciBinary)
	}

	return installed, nil
}

// RollbackBinaries removes previously installed binaries.
func RollbackBinaries(installed []string) {
	for _, name := range installed {
		goService.Uninstall(name)
	}
}

// RemoveAllBinaries removes golangci-lint from .bin/ (for uninstall).
func RemoveAllBinaries() error {
	return goService.Uninstall(golangciBinary)
}
