package tddguard

import (
	"code-template/services"
)

const (
	tddGuardBinary  = "tdd-guard"
	tddGuardPackage = "tdd-guard"
)

var npmService = services.NPM

// TddGuardPackage is the package definition for tdd-guard.
var TddGuardPackage = services.Package{
	Name:        tddGuardPackage,
	InstallPath: tddGuardPackage,
}

// CheckNpmInstalled verifies npm is available in PATH.
func CheckNpmInstalled() bool {
	return npmService.IsAvailable()
}

// IsTddGuardInstalled checks if tdd-guard is available in PATH.
func IsTddGuardInstalled() bool {
	return npmService.IsInstalled(tddGuardBinary)
}

// InstallTddGuard installs tdd-guard via npm install -g.
func InstallTddGuard() error {
	return npmService.Install(TddGuardPackage)
}

// Note: We intentionally do NOT provide an uninstall function for npm packages
// as they may be used by other projects. The module's Uninstall() only removes
// the configuration, not the npm package itself.
