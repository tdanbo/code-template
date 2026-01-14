package tddguard

import (
	"os/exec"
)

const tddGuardPackage = "tdd-guard"

// CheckNpmInstalled verifies npm is available in PATH
func CheckNpmInstalled() bool {
	_, err := exec.LookPath("npm")
	return err == nil
}

// IsTddGuardInstalled checks if tdd-guard is available in PATH
func IsTddGuardInstalled() bool {
	_, err := exec.LookPath("tdd-guard")
	return err == nil
}

// InstallTddGuard runs: npm install -g tdd-guard
func InstallTddGuard() error {
	cmd := exec.Command("npm", "install", "-g", tddGuardPackage)
	return cmd.Run()
}
