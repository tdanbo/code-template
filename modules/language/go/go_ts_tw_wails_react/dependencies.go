package gotstwwailsreact

import (
	"os/exec"
	"runtime"

	"code-template/services"
)

const (
	wailsBinary  = "wails"
	wailsInstall = "github.com/wailsapp/wails/v2/cmd/wails@latest"
)

var goService = services.Go
var npmService = services.NPM

// WailsPackage is the package definition for wails CLI.
var WailsPackage = services.Package{
	Name:        wailsBinary,
	InstallPath: wailsInstall,
}

// CheckGoInstalled verifies that Go is available in PATH.
func CheckGoInstalled() bool {
	return goService.IsAvailable()
}

// CheckNpmInstalled verifies that npm is available in PATH.
func CheckNpmInstalled() bool {
	return npmService.IsAvailable()
}

// IsWailsInstalled checks if wails CLI is available in PATH.
func IsWailsInstalled() bool {
	_, err := exec.LookPath(wailsBinary)
	return err == nil
}

// InstallWailsCLI installs wails CLI globally using go install.
func InstallWailsCLI() error {
	cmd := exec.Command("go", "install", wailsInstall)
	return cmd.Run()
}

// CheckGccInstalled verifies gcc is available in PATH.
func CheckGccInstalled() bool {
	_, err := exec.LookPath("gcc")
	return err == nil
}

// CheckPkgConfig verifies pkg-config is available.
func CheckPkgConfig() bool {
	_, err := exec.LookPath("pkg-config")
	return err == nil
}

// CheckSystemLibrary checks if a library exists via pkg-config.
func CheckSystemLibrary(pkgName string) bool {
	cmd := exec.Command("pkg-config", "--exists", pkgName)
	return cmd.Run() == nil
}

// CheckLinuxDependencies verifies GTK and WebKit are installed.
// Returns (ok, missingDep) where missingDep describes what's missing.
func CheckLinuxDependencies() (bool, string) {
	if runtime.GOOS != "linux" {
		return true, ""
	}

	// Check gcc
	if !CheckGccInstalled() {
		return false, "gcc (install build-essential or base-devel)"
	}

	// Check pkg-config (needed to verify libraries)
	if !CheckPkgConfig() {
		return false, "pkg-config"
	}

	// Check GTK3
	if !CheckSystemLibrary("gtk+-3.0") {
		return false, "libgtk-3-dev"
	}

	// Check WebKit (try 4.1 first for newer distros, then 4.0)
	if !CheckSystemLibrary("webkit2gtk-4.1") && !CheckSystemLibrary("webkit2gtk-4.0") {
		return false, "libwebkit2gtk-4.0-dev or libwebkit2gtk-4.1-dev"
	}

	return true, ""
}

// NeedsWebkit41BuildTag returns true if the system has webkit2gtk-4.1 but not 4.0.
// This means Wails needs the webkit2_41 build tag.
func NeedsWebkit41BuildTag() bool {
	if runtime.GOOS != "linux" {
		return false
	}
	// If 4.0 exists, no special tag needed
	if CheckSystemLibrary("webkit2gtk-4.0") {
		return false
	}
	// If only 4.1 exists, need the build tag
	return CheckSystemLibrary("webkit2gtk-4.1")
}
