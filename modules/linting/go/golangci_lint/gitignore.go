package golangci_lint

import (
	"bufio"
	"os"
	"strings"
)

const gitignoreFile = ".gitignore"
const binDirEntry = ".bin/"

// AddToGitignore adds .bin/ to .gitignore if not already present.
func AddToGitignore() error {
	if HasGitignoreEntry() {
		return nil
	}

	content, err := os.ReadFile(gitignoreFile)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	var prefix string
	if len(content) > 0 && content[len(content)-1] != '\n' {
		prefix = "\n"
	}

	f, err := os.OpenFile(gitignoreFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(prefix + binDirEntry + "\n")
	return err
}

// HasGitignoreEntry checks if .bin/ is already in .gitignore.
func HasGitignoreEntry() bool {
	f, err := os.Open(gitignoreFile)
	if err != nil {
		return false
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == binDirEntry || line == ".bin" {
			return true
		}
	}
	return false
}

// RemoveFromGitignore removes .bin/ from .gitignore.
func RemoveFromGitignore() error {
	content, err := os.ReadFile(gitignoreFile)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}

	lines := strings.Split(string(content), "\n")
	var newLines []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != binDirEntry && trimmed != ".bin" {
			newLines = append(newLines, line)
		}
	}

	for len(newLines) > 0 && newLines[len(newLines)-1] == "" {
		newLines = newLines[:len(newLines)-1]
	}

	newContent := strings.Join(newLines, "\n")
	if len(newLines) > 0 {
		newContent += "\n"
	}

	return os.WriteFile(gitignoreFile, []byte(newContent), 0644)
}
