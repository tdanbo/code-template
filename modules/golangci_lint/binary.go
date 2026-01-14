package golangci_lint

import (
	"os"
	"os/exec"
	"path/filepath"
)

const (
	binDir          = ".bin"
	golangciBinary  = "golangci-lint"
	taskBinary      = "task"
	golangciInstall = "github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest"
	taskInstall     = "github.com/go-task/task/v3/cmd/task@latest"
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

// AreBinariesInstalled checks if both required binaries are installed.
func AreBinariesInstalled() bool {
	return IsBinaryInstalled(golangciBinary) && IsBinaryInstalled(taskBinary)
}

// InstallBinaries installs all required binaries, returns list of installed for rollback.
func InstallBinaries() ([]string, error) {
	binaries := []struct {
		name string
		pkg  string
	}{
		{golangciBinary, golangciInstall},
		{taskBinary, taskInstall},
	}

	var installed []string
	for _, bin := range binaries {
		if !IsBinaryInstalled(bin.name) {
			if err := InstallBinary(bin.pkg); err != nil {
				return installed, err
			}
			installed = append(installed, bin.name)
		}
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

// RemoveAllBinaries removes both binaries (for uninstall).
func RemoveAllBinaries() error {
	var lastErr error
	if err := RemoveBinary(golangciBinary); err != nil {
		lastErr = err
	}
	if err := RemoveBinary(taskBinary); err != nil {
		lastErr = err
	}
	cleanupBinDir()
	return lastErr
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
