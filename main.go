package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"code-template/autoinit"
	"code-template/helpers"
	"code-template/models"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Messages for async operations
type installResultMsg struct {
	moduleName string
	success    bool
	action     string // "install", "uninstall", or "update"
}

// Style definitions
var (
	primaryColor   = lipgloss.Color("#7D56F4")
	secondaryColor = lipgloss.Color("#04B575")
	warningColor   = lipgloss.Color("#FF6B6B")
	subtleColor    = lipgloss.Color("#626262")

	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(primaryColor).
			Padding(0, 2)

	containerStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(primaryColor).
			Padding(1, 2)

	selectedStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(primaryColor)

	normalStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF"))

	checkboxInstalled = lipgloss.NewStyle().
				Foreground(secondaryColor).
				Bold(true)

	checkboxOutdated = lipgloss.NewStyle().
				Foreground(warningColor).
				Bold(true)

	checkboxNotInstalled = lipgloss.NewStyle().
				Foreground(subtleColor)

	versionStyle = lipgloss.NewStyle().
			Foreground(subtleColor)

	statusSuccessStyle = lipgloss.NewStyle().
				Foreground(secondaryColor).
				Bold(true)

	statusErrorStyle = lipgloss.NewStyle().
				Foreground(warningColor).
				Bold(true)

	helpStyle = lipgloss.NewStyle().
			Foreground(subtleColor)

	countStyle = lipgloss.NewStyle().
			Foreground(subtleColor)

	treeConnectorStyle = lipgloss.NewStyle().
				Foreground(subtleColor)

	expandedStyle = lipgloss.NewStyle().
			Foreground(primaryColor)

	collapsedStyle = lipgloss.NewStyle().
			Foreground(subtleColor)
)

type ViewModel struct {
	Tree           *models.TreeState
	SelectedIdx    int
	StatusMessage  string
	StatusIsError  bool
	IsLoading      bool
	LoadingMessage string
	spinner        spinner.Model
}

func (m ViewModel) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m *ViewModel) getCurrentNode() *models.TreeNode {
	if len(m.Tree.FlatVisible) == 0 || m.SelectedIdx >= len(m.Tree.FlatVisible) {
		return nil
	}
	return m.Tree.FlatVisible[m.SelectedIdx]
}

func (m ViewModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case installResultMsg:
		m.IsLoading = false
		m.LoadingMessage = ""
		if msg.success {
			switch msg.action {
			case "install":
				m.StatusMessage = fmt.Sprintf("✓ Installed %s", msg.moduleName)
			case "uninstall":
				m.StatusMessage = fmt.Sprintf("✓ Uninstalled %s", msg.moduleName)
			case "update":
				m.StatusMessage = fmt.Sprintf("✓ Updated %s", msg.moduleName)
			}
		} else {
			m.StatusIsError = true
			switch msg.action {
			case "install":
				m.StatusMessage = fmt.Sprintf("✗ Failed to install %s", msg.moduleName)
			case "uninstall":
				m.StatusMessage = fmt.Sprintf("✗ Failed to uninstall %s", msg.moduleName)
			case "update":
				m.StatusMessage = fmt.Sprintf("✗ Failed to update %s", msg.moduleName)
			}
		}
		return m, nil

	case tea.KeyMsg:
		// Block input while loading
		if m.IsLoading {
			return m, nil
		}

		// Clear status message on any keypress
		m.StatusMessage = ""
		m.StatusIsError = false

		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit

		case "up", "k":
			if m.SelectedIdx > 0 {
				m.SelectedIdx--
			}

		case "down", "j":
			if m.SelectedIdx < len(m.Tree.FlatVisible)-1 {
				m.SelectedIdx++
			}

		case "right", "l":
			node := m.getCurrentNode()
			if node != nil && node.Type == models.NodeCategory {
				m.Tree.SetExpanded(node.ID, true)
				m.Tree.RebuildFlatVisible()
			}

		case "left", "h":
			node := m.getCurrentNode()
			if node != nil {
				if node.Type == models.NodeCategory && m.Tree.IsExpanded(node.ID) {
					// Collapse current category
					m.Tree.SetExpanded(node.ID, false)
					m.Tree.RebuildFlatVisible()
				} else if node.Parent != nil {
					// Jump to parent
					for i, n := range m.Tree.FlatVisible {
						if n == node.Parent {
							m.SelectedIdx = i
							break
						}
					}
				}
			}

		case "esc":
			// Collapse all categories
			for id := range m.Tree.Expanded {
				m.Tree.Expanded[id] = false
			}
			m.Tree.RebuildFlatVisible()
			m.SelectedIdx = 0

		case "enter":
			node := m.getCurrentNode()
			if node == nil {
				return m, nil
			}

			if node.Type == models.NodeCategory {
				// Toggle expand/collapse
				m.Tree.ToggleExpanded(node.ID)
				m.Tree.RebuildFlatVisible()
			} else if node.Type == models.NodeModule && node.Module != nil {
				// Install/update module
				module := node.Module
				state := helpers.GetModuleState(module)
				moduleName := module.GetName()

				switch state {
				case helpers.StateNotInstalled:
					m.IsLoading = true
					m.LoadingMessage = fmt.Sprintf("Installing %s...", moduleName)
					return m, func() tea.Msg {
						success := module.Install()
						return installResultMsg{
							moduleName: moduleName,
							success:    success,
							action:     "install",
						}
					}
				case helpers.StateOutdated:
					m.IsLoading = true
					m.LoadingMessage = fmt.Sprintf("Updating %s...", moduleName)
					return m, func() tea.Msg {
						success := helpers.UpdateModule(module)
						return installResultMsg{
							moduleName: moduleName,
							success:    success,
							action:     "update",
						}
					}
				case helpers.StateUpToDate:
					m.StatusMessage = fmt.Sprintf("%s is already up to date", moduleName)
				}
			}

		case "delete", "backspace":
			node := m.getCurrentNode()
			if node != nil && node.Type == models.NodeModule && node.Module != nil {
				module := node.Module
				if module.IsInstalled() {
					m.IsLoading = true
					moduleName := module.GetName()
					m.LoadingMessage = fmt.Sprintf("Uninstalling %s...", moduleName)
					return m, func() tea.Msg {
						success := module.Uninstall()
						return installResultMsg{
							moduleName: moduleName,
							success:    success,
							action:     "uninstall",
						}
					}
				}
			}
		}
	}
	return m, nil
}

func (m ViewModel) View() string {
	var content strings.Builder

	// Title
	title := titleStyle.Render(" Code Template Manager ")
	content.WriteString(title + "\n\n")

	// Render tree
	for i, node := range m.Tree.FlatVisible {
		line := m.renderNode(node, i == m.SelectedIdx)
		content.WriteString(line + "\n")
	}

	// Loading indicator or status message
	if m.IsLoading {
		content.WriteString("\n")
		content.WriteString(fmt.Sprintf("%s %s", m.spinner.View(), m.LoadingMessage))
		content.WriteString("\n")
	} else if m.StatusMessage != "" {
		content.WriteString("\n")
		if m.StatusIsError {
			content.WriteString(statusErrorStyle.Render(m.StatusMessage))
		} else {
			content.WriteString(statusSuccessStyle.Render(m.StatusMessage))
		}
		content.WriteString("\n")
	}

	// Help text
	content.WriteString("\n")
	helpText := "↑/↓ navigate • →/l expand • ←/h collapse • enter install • del uninstall • q quit"
	content.WriteString(helpStyle.Render(helpText))

	// Wrap in container
	return containerStyle.Render(content.String())
}

func (m ViewModel) renderNode(node *models.TreeNode, selected bool) string {
	var line strings.Builder

	// Cursor indicator
	if selected {
		line.WriteString("▸ ")
	} else {
		line.WriteString("  ")
	}

	// Tree connectors for depth > 0
	if node.Depth > 0 {
		line.WriteString(m.getTreePrefix(node))
	}

	if node.Type == models.NodeCategory {
		// Category node
		var indicator string
		if m.Tree.IsExpanded(node.ID) {
			indicator = expandedStyle.Render("▼")
		} else {
			indicator = collapsedStyle.Render("▶")
		}

		installed, total := node.GetInstalledCount()
		countText := countStyle.Render(fmt.Sprintf("(%d/%d)", installed, total))

		nodeContent := fmt.Sprintf("%s %s %s", indicator, node.Name, countText)
		if selected {
			line.WriteString(selectedStyle.Render(nodeContent))
		} else {
			line.WriteString(nodeContent)
		}
	} else {
		// Module node
		module := node.Module
		state := helpers.GetModuleState(module)

		var checkbox string
		var versionText string

		switch state {
		case helpers.StateUpToDate:
			checkbox = checkboxInstalled.Render("[✓]")
			versionText = versionStyle.Render(fmt.Sprintf(" (v%d)", module.GetVersion()))
		case helpers.StateOutdated:
			checkbox = checkboxOutdated.Render("[!]")
			installedVer := helpers.GetInstalledVersion(module)
			versionText = versionStyle.Render(fmt.Sprintf(" (v%d → v%d)", installedVer, module.GetVersion()))
		case helpers.StateNotInstalled:
			checkbox = checkboxNotInstalled.Render("[ ]")
			versionText = ""
		}

		nodeContent := fmt.Sprintf("%s %s%s", checkbox, node.Name, versionText)
		if selected {
			line.WriteString(selectedStyle.Render(nodeContent))
		} else {
			line.WriteString(nodeContent)
		}
	}

	return line.String()
}

func (m ViewModel) getTreePrefix(node *models.TreeNode) string {
	var prefix strings.Builder

	// Build prefix from ancestors
	ancestors := make([]*models.TreeNode, 0)
	current := node.Parent
	for current != nil {
		ancestors = append([]*models.TreeNode{current}, ancestors...)
		current = current.Parent
	}

	// Add vertical lines for each ancestor level
	for _, ancestor := range ancestors {
		if ancestor.IsLastChild() {
			prefix.WriteString("   ")
		} else {
			prefix.WriteString(treeConnectorStyle.Render("│  "))
		}
	}

	// Add connector for this node
	if node.IsLastChild() {
		prefix.WriteString(treeConnectorStyle.Render("└─ "))
	} else {
		prefix.WriteString(treeConnectorStyle.Render("├─ "))
	}

	return prefix.String()
}

// CLI flags
var (
	installFlag   string
	uninstallFlag string
	versionFlag   string
	listFlag      bool
	debugTreeFlag bool
)

func init() {
	flag.StringVar(&installFlag, "install", "", "Install a module by name")
	flag.StringVar(&installFlag, "i", "", "Install a module by name (shorthand)")
	flag.StringVar(&uninstallFlag, "uninstall", "", "Uninstall a module by name")
	flag.StringVar(&uninstallFlag, "u", "", "Uninstall a module by name (shorthand)")
	flag.StringVar(&versionFlag, "version", "", "Show version info for a module")
	flag.StringVar(&versionFlag, "v", "", "Show version info for a module (shorthand)")
	flag.BoolVar(&listFlag, "list", false, "List all available modules")
	flag.BoolVar(&listFlag, "l", false, "List all available modules (shorthand)")
	flag.BoolVar(&debugTreeFlag, "debug-tree", false, "Debug: show tree structure")
}

// findModule finds a module by name or key.
func findModule(modules []models.Module, name string) models.Module {
	name = strings.ToLower(name)
	for _, m := range modules {
		if strings.ToLower(m.GetName()) == name || strings.ToLower(m.GetKey()) == name {
			return m
		}
	}
	return nil
}

// runInstall installs a module by name.
func runInstall(modules []models.Module, name string) int {
	module := findModule(modules, name)
	if module == nil {
		fmt.Fprintf(os.Stderr, "Error: module '%s' not found\n", name)
		fmt.Fprintln(os.Stderr, "Use --list to see available modules")
		return 1
	}

	state := helpers.GetModuleState(module)
	switch state {
	case helpers.StateUpToDate:
		fmt.Printf("Module '%s' is already installed (v%d)\n", module.GetName(), module.GetVersion())
		return 0
	case helpers.StateOutdated:
		fmt.Printf("Updating '%s' from v%d to v%d...\n",
			module.GetName(), helpers.GetInstalledVersion(module), module.GetVersion())
		if helpers.UpdateModule(module) {
			fmt.Printf("✓ Updated '%s' to v%d\n", module.GetName(), module.GetVersion())
			return 0
		}
		fmt.Fprintf(os.Stderr, "✗ Failed to update '%s'\n", module.GetName())
		return 1
	case helpers.StateNotInstalled:
		fmt.Printf("Installing '%s' v%d...\n", module.GetName(), module.GetVersion())
		if module.Install() {
			fmt.Printf("✓ Installed '%s' v%d\n", module.GetName(), module.GetVersion())
			return 0
		}
		fmt.Fprintf(os.Stderr, "✗ Failed to install '%s'\n", module.GetName())
		return 1
	}
	return 1
}

// runUninstall uninstalls a module by name.
func runUninstall(modules []models.Module, name string) int {
	module := findModule(modules, name)
	if module == nil {
		fmt.Fprintf(os.Stderr, "Error: module '%s' not found\n", name)
		fmt.Fprintln(os.Stderr, "Use --list to see available modules")
		return 1
	}

	if !module.IsInstalled() {
		fmt.Printf("Module '%s' is not installed\n", module.GetName())
		return 0
	}

	fmt.Printf("Uninstalling '%s'...\n", module.GetName())
	if module.Uninstall() {
		fmt.Printf("✓ Uninstalled '%s'\n", module.GetName())
		return 0
	}
	fmt.Fprintf(os.Stderr, "✗ Failed to uninstall '%s'\n", module.GetName())
	return 1
}

// runVersion shows version info for a module.
func runVersion(modules []models.Module, name string) int {
	module := findModule(modules, name)
	if module == nil {
		fmt.Fprintf(os.Stderr, "Error: module '%s' not found\n", name)
		fmt.Fprintln(os.Stderr, "Use --list to see available modules")
		return 1
	}

	fmt.Printf("Module: %s\n", module.GetName())
	fmt.Printf("  Path:     %s\n", module.GetPath())
	fmt.Printf("  Category: %s\n", module.GetCategory())
	fmt.Printf("  Version:  v%d\n", module.GetVersion())

	state := helpers.GetModuleState(module)
	switch state {
	case helpers.StateUpToDate:
		fmt.Printf("  Status:   installed (up to date)\n")
	case helpers.StateOutdated:
		fmt.Printf("  Status:   installed (outdated, v%d → v%d)\n",
			helpers.GetInstalledVersion(module), module.GetVersion())
	case helpers.StateNotInstalled:
		fmt.Printf("  Status:   not installed\n")
	}
	return 0
}

// runDebugTree prints the tree structure for debugging.
func runDebugTree(modules []models.Module) int {
	tree := helpers.BuildTree(modules)

	var printNode func(node *models.TreeNode, indent string)
	printNode = func(node *models.TreeNode, indent string) {
		if node.Type == models.NodeCategory {
			fmt.Printf("%s[Category] %s (id: %s)\n", indent, node.Name, node.ID)
		} else {
			path := ""
			if node.Module != nil {
				path = node.Module.GetPath()
			}
			fmt.Printf("%s[Module] %s (id: %s, path: %s)\n", indent, node.Name, node.ID, path)
		}
		for _, child := range node.Children {
			printNode(child, indent+"  ")
		}
	}

	fmt.Println("Tree structure:")
	for _, root := range tree.Roots {
		printNode(root, "")
	}
	return 0
}

// runList lists all available modules.
func runList(modules []models.Module) int {
	fmt.Println("Available modules:")
	fmt.Println()

	// Group by category
	categories := make(map[string][]models.Module)
	for _, m := range modules {
		cat := m.GetCategory()
		categories[cat] = append(categories[cat], m)
	}

	for cat, mods := range categories {
		fmt.Printf("  %s:\n", cat)
		for _, m := range mods {
			state := helpers.GetModuleState(m)
			var status string
			switch state {
			case helpers.StateUpToDate:
				status = fmt.Sprintf("[✓] v%d", m.GetVersion())
			case helpers.StateOutdated:
				status = fmt.Sprintf("[!] v%d → v%d", helpers.GetInstalledVersion(m), m.GetVersion())
			case helpers.StateNotInstalled:
				status = "[ ]"
			}
			fmt.Printf("    %-20s %s\n", m.GetName(), status)
		}
		fmt.Println()
	}
	return 0
}

// runTUI runs the interactive terminal UI.
func runTUI(modules []models.Module) int {
	tree := helpers.BuildTree(modules)

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(primaryColor)

	m := ViewModel{
		Tree:        tree,
		SelectedIdx: 0,
		spinner:     s,
	}

	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return 1
	}
	return 0
}

func main() {
	flag.Parse()

	if err := autoinit.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Initialization failed: %v\n", err)
		os.Exit(1)
	}

	modules := helpers.GetModules()

	// Handle CLI commands
	if installFlag != "" {
		os.Exit(runInstall(modules, installFlag))
	}
	if uninstallFlag != "" {
		os.Exit(runUninstall(modules, uninstallFlag))
	}
	if versionFlag != "" {
		os.Exit(runVersion(modules, versionFlag))
	}
	if listFlag {
		os.Exit(runList(modules))
	}
	if debugTreeFlag {
		os.Exit(runDebugTree(modules))
	}

	// No CLI flags, run TUI
	os.Exit(runTUI(modules))
}
