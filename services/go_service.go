package services

import (
	"os"
	"os/exec"
	"path/filepath"
)

const (
	defaultBinDir = ".bin"
)

// GoService manages Go packages installed via `go install`.
// Packages are installed to a local .bin/ directory by default.
type GoService struct {
	// BinDir is the directory where binaries are installed.
	// Defaults to ".bin" if empty.
	BinDir string
}

// NewGoService creates a new GoService with default settings.
func NewGoService() *GoService {
	return &GoService{
		BinDir: defaultBinDir,
	}
}

// Name returns "go".
func (s *GoService) Name() string {
	return "go"
}

// IsAvailable checks if Go is installed and available in PATH.
func (s *GoService) IsAvailable() bool {
	_, err := exec.LookPath("go")
	return err == nil
}

// IsInstalled checks if a binary is installed.
// Checks both the local bin directory and global PATH.
func (s *GoService) IsInstalled(binaryName string) bool {
	// Check local .bin/ first
	if s.isInstalledLocally(binaryName) {
		return true
	}
	// Check global PATH
	_, err := exec.LookPath(binaryName)
	return err == nil
}

// IsInstalledLocally checks if a binary exists in the local bin directory.
func (s *GoService) IsInstalledLocally(binaryName string) bool {
	return s.isInstalledLocally(binaryName)
}

func (s *GoService) isInstalledLocally(binaryName string) bool {
	binPath := filepath.Join(s.getBinDir(), binaryName)
	_, err := os.Stat(binPath)
	return err == nil
}

// Install installs a Go package to the local bin directory.
// Creates the bin directory if it doesn't exist.
func (s *GoService) Install(pkg Package) error {
	binDir := s.getBinDir()

	// Ensure bin directory exists
	if err := os.MkdirAll(binDir, 0755); err != nil {
		return err
	}

	absDir, err := filepath.Abs(binDir)
	if err != nil {
		return err
	}

	cmd := exec.Command("go", "install", pkg.InstallPath)
	cmd.Env = append(os.Environ(), "GOBIN="+absDir)
	return cmd.Run()
}

// Uninstall removes a binary from the local bin directory.
// Also cleans up the bin directory if it becomes empty.
func (s *GoService) Uninstall(binaryName string) error {
	binDir := s.getBinDir()
	binPath := filepath.Join(binDir, binaryName)

	err := os.Remove(binPath)
	if os.IsNotExist(err) {
		err = nil // Not an error if file doesn't exist
	}

	// Clean up empty bin directory
	s.cleanupBinDir()

	return err
}

// EnsureBinDir creates the bin directory if it doesn't exist.
func (s *GoService) EnsureBinDir() error {
	return os.MkdirAll(s.getBinDir(), 0755)
}

// GetBinPath returns the full path to a binary in the bin directory.
func (s *GoService) GetBinPath(binaryName string) string {
	return filepath.Join(s.getBinDir(), binaryName)
}

func (s *GoService) getBinDir() string {
	if s.BinDir == "" {
		return defaultBinDir
	}
	return s.BinDir
}

func (s *GoService) cleanupBinDir() {
	binDir := s.getBinDir()
	entries, err := os.ReadDir(binDir)
	if err != nil {
		return
	}
	if len(entries) == 0 {
		os.Remove(binDir)
	}
}

// Global instance for convenience
var Go = NewGoService()
