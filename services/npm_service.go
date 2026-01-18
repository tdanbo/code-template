package services

import (
	"os/exec"
)

// NPMService manages NPM packages installed globally via `npm install -g`.
type NPMService struct{}

// NewNPMService creates a new NPMService.
func NewNPMService() *NPMService {
	return &NPMService{}
}

// Name returns "npm".
func (s *NPMService) Name() string {
	return "npm"
}

// IsAvailable checks if npm is installed and available in PATH.
func (s *NPMService) IsAvailable() bool {
	_, err := exec.LookPath("npm")
	return err == nil
}

// IsInstalled checks if a package is installed by checking if its binary is in PATH.
func (s *NPMService) IsInstalled(binaryName string) bool {
	_, err := exec.LookPath(binaryName)
	return err == nil
}

// Install installs a package globally using `npm install -g`.
// The pkg.Name is the npm package name, pkg.InstallPath can be used for specific versions.
func (s *NPMService) Install(pkg Package) error {
	installName := pkg.Name
	if pkg.InstallPath != "" {
		installName = pkg.InstallPath
	}
	cmd := exec.Command("npm", "install", "-g", installName)
	return cmd.Run()
}

// Uninstall removes a package globally using `npm uninstall -g`.
// Note: This is typically not called as we don't want to remove global packages
// that might be used by other projects.
func (s *NPMService) Uninstall(binaryName string) error {
	cmd := exec.Command("npm", "uninstall", "-g", binaryName)
	return cmd.Run()
}

// Global instance for convenience
var NPM = NewNPMService()
