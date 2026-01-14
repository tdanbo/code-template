package autoinit

import (
	"os"
	"os/exec"
)

const (
	taskInstallPkg = "github.com/go-task/task/v3/cmd/task@latest"
)

func isGoInstalled() bool {
	_, err := exec.LookPath("go")
	return err == nil
}

func isTaskGloballyAvailable() bool {
	_, err := exec.LookPath("task")
	return err == nil
}

func installTaskGlobally() error {
	printInfo("Installing go-task globally...")
	cmd := exec.Command("go", "install", taskInstallPkg)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
