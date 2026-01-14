package autoinit

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const (
	colorReset  = "\033[0m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorRed    = "\033[31m"
	colorCyan   = "\033[36m"
	colorBold   = "\033[1m"
)

func printWelcome() {
	fmt.Println()
	fmt.Printf("%s%s=== Code Template Manager ===%s\n", colorBold, colorCyan, colorReset)
	fmt.Println("First-time setup detected. Running initialization...")
	fmt.Println()
}

func printSuccess(msg string) {
	fmt.Printf("%s[OK]%s %s\n", colorGreen, colorReset, msg)
}

func printWarning(msg string) {
	fmt.Printf("%s[WARN]%s %s\n", colorYellow, colorReset, msg)
}

func printError(msg string) {
	fmt.Printf("%s[ERROR]%s %s\n", colorRed, colorReset, msg)
}

func printInfo(msg string) {
	fmt.Printf("%s[INFO]%s %s\n", colorCyan, colorReset, msg)
}

func promptInstallTask() bool {
	fmt.Println()
	fmt.Println("go-task (task) command not found in PATH.")
	fmt.Println("go-task is required to run module tasks (e.g., task go-lint)")
	fmt.Println()
	fmt.Print("Install go-task globally via 'go install'? [Y/n]: ")

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return false
	}

	input = strings.TrimSpace(strings.ToLower(input))
	return input == "" || input == "y" || input == "yes"
}
