package golangci_lint

import (
	"os"
	"os/exec"
	"path/filepath"
)

const (
	binDir          = ".bin"
	golangciBinary  = "golangci-lint"
	golangciInstall = "github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest"
)

// CheckGoInstalled verifies that Go is available in PATH.
func CheckGoInstalled() bool {
	_, err := exec.LookPath("go")
	return err == nil
}

// EnsureBinDir creates the .bin directory if it doesn't exist.
func EnsureBinDir() error {
	return os.MkdirAll(binDir, 0755)
}

// IsBinaryInstalled checks if a binary exists in .bin/.
func IsBinaryInstalled(name string) bool {
	binPath := filepath.Join(binDir, name)
	_, err := os.Stat(binPath)
	return err == nil
}

// IsBinaryAvailable checks if a binary is available either in .bin/ or globally in PATH.
func IsBinaryAvailable(name string) bool {
	// Check .bin/ first
	if IsBinaryInstalled(name) {
		return true
	}
	// Check global PATH
	_, err := exec.LookPath(name)
	return err == nil
}

// InstallBinary installs a Go binary to .bin/ using GOBIN.
func InstallBinary(packagePath string) error {
	absDir, err := filepath.Abs(binDir)
	if err != nil {
		return err
	}

	cmd := exec.Command("go", "install", packagePath)
	cmd.Env = append(os.Environ(), "GOBIN="+absDir)
	return cmd.Run()
}

// RemoveBinary removes a specific binary from .bin/.
func RemoveBinary(name string) error {
	binPath := filepath.Join(binDir, name)
	err := os.Remove(binPath)
	if os.IsNotExist(err) {
		return nil
	}
	return err
}

// AreBinariesInstalled checks if golangci-lint is available (locally or globally).
// Note: task is checked separately as a global requirement on startup.
func AreBinariesInstalled() bool {
	return IsBinaryAvailable(golangciBinary)
}

// InstallBinaries installs golangci-lint if not already available, returns list of installed for rollback.
// Note: task is required globally and not installed locally.
func InstallBinaries() ([]string, error) {
	var installed []string

	// Only install golangci-lint if not available (locally or globally)
	if !IsBinaryAvailable(golangciBinary) {
		if err := InstallBinary(golangciInstall); err != nil {
			return installed, err
		}
		installed = append(installed, golangciBinary)
	}

	return installed, nil
}

// RollbackBinaries removes previously installed binaries.
func RollbackBinaries(installed []string) {
	for _, name := range installed {
		RemoveBinary(name)
	}
	cleanupBinDir()
}

// RemoveAllBinaries removes golangci-lint from .bin/ (for uninstall).
// Note: task is required globally and not managed locally.
func RemoveAllBinaries() error {
	err := RemoveBinary(golangciBinary)
	cleanupBinDir()
	return err
}

// cleanupBinDir removes .bin/ directory if empty.
func cleanupBinDir() {
	entries, err := os.ReadDir(binDir)
	if err != nil {
		return
	}
	if len(entries) == 0 {
		os.Remove(binDir)
	}
}
