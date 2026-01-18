package gotstwwailsreact

import (
	_ "embed"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

//go:embed configs/tsconfig.json
var tsconfigJSON []byte

//go:embed configs/tailwind.config.js
var tailwindConfigJS []byte

//go:embed configs/postcss.config.js
var postcssConfigJS []byte

//go:embed configs/index.css
var indexCSS []byte

const frontendDir = "frontend"

// getProjectName returns the current directory name as the project name.
func getProjectName() string {
	wd, err := os.Getwd()
	if err != nil {
		return "myapp"
	}
	return filepath.Base(wd)
}

// ScaffoldProject runs wails init to scaffold a new React+TypeScript project.
func ScaffoldProject(name string) error {
	cmd := exec.Command("wails", "init", "-n", name, "-t", "react-ts", "-d", ".")
	return cmd.Run()
}

// CopyConfigFiles writes the embedded config files to the frontend directory.
func CopyConfigFiles() error {
	// Write tsconfig.json
	tsconfigPath := filepath.Join(frontendDir, "tsconfig.json")
	if err := os.WriteFile(tsconfigPath, tsconfigJSON, 0644); err != nil {
		return err
	}

	// Write tailwind.config.js
	tailwindPath := filepath.Join(frontendDir, "tailwind.config.js")
	if err := os.WriteFile(tailwindPath, tailwindConfigJS, 0644); err != nil {
		os.Remove(tsconfigPath)
		return err
	}

	// Write postcss.config.js
	postcssPath := filepath.Join(frontendDir, "postcss.config.js")
	if err := os.WriteFile(postcssPath, postcssConfigJS, 0644); err != nil {
		os.Remove(tsconfigPath)
		os.Remove(tailwindPath)
		return err
	}

	// Write index.css with Tailwind directives
	cssPath := filepath.Join(frontendDir, "src", "index.css")
	if err := os.WriteFile(cssPath, indexCSS, 0644); err != nil {
		os.Remove(tsconfigPath)
		os.Remove(tailwindPath)
		os.Remove(postcssPath)
		return err
	}

	return nil
}

// UpgradeTypeScript upgrades TypeScript to latest version in the frontend directory.
// Required because tsconfig.json uses TS 5.0+ features (moduleResolution: bundler).
func UpgradeTypeScript() error {
	cmd := exec.Command("npm", "install", "typescript@latest", "--save-dev")
	cmd.Dir = frontendDir
	return cmd.Run()
}

// InstallTailwind installs Tailwind CSS v4 PostCSS plugin in the frontend directory.
func InstallTailwind() error {
	cmd := exec.Command("npm", "install", "@tailwindcss/postcss", "--save-dev")
	cmd.Dir = frontendDir
	return cmd.Run()
}

// InstallFrontendDeps runs npm install in the frontend directory.
func InstallFrontendDeps() error {
	cmd := exec.Command("npm", "install")
	cmd.Dir = frontendDir
	return cmd.Run()
}

// RollbackConfigFiles removes the config files written to frontend/.
func RollbackConfigFiles() {
	os.Remove(filepath.Join(frontendDir, "tsconfig.json"))
	os.Remove(filepath.Join(frontendDir, "tailwind.config.js"))
	os.Remove(filepath.Join(frontendDir, "postcss.config.js"))
	os.Remove(filepath.Join(frontendDir, "src", "index.css"))
}

// RollbackScaffold removes all files created by wails init.
func RollbackScaffold() {
	// Remove wails.json
	os.Remove("wails.json")

	// Remove main.go (wails generated)
	os.Remove("main.go")

	// Remove app.go (wails generated)
	os.Remove("app.go")

	// Remove frontend/ directory
	os.RemoveAll(frontendDir)

	// Remove build/ directory
	os.RemoveAll("build")
}

// ConfigureWebkit41 updates wails.json to use webkit2gtk-4.1 build tags.
// This is required on newer Linux distros that only have webkit2gtk-4.1.
func ConfigureWebkit41() error {
	// Read existing wails.json
	data, err := os.ReadFile("wails.json")
	if err != nil {
		return err
	}

	// Parse as generic map to preserve all fields
	var config map[string]interface{}
	if err := json.Unmarshal(data, &config); err != nil {
		return err
	}

	// Add webkit2_41 build tags (build:tags is the official wails.json field)
	config["build:tags"] = "webkit2_41"

	// Write back with indentation
	output, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile("wails.json", output, 0644)
}

// AddTailwindImport updates main.tsx to import index.css for Tailwind.
func AddTailwindImport() error {
	mainTsxPath := filepath.Join(frontendDir, "src", "main.tsx")

	data, err := os.ReadFile(mainTsxPath)
	if err != nil {
		return err
	}

	content := string(data)

	// Add import for index.css before the existing style.css import
	oldImport := `import './style.css'`
	newImport := `import './index.css'
import './style.css'`

	content = strings.Replace(content, oldImport, newImport, 1)

	return os.WriteFile(mainTsxPath, []byte(content), 0644)
}
