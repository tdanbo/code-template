package main

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"code-template/autoinit"
	"code-template/helpers"
	"code-template/models"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ViewHeader string

var (
	categories ViewHeader = "categories"
	modules    ViewHeader = "modules"
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
)

type ViewModel struct {
	Categories      map[string]map[string]models.Module
	CategoryKeys    []string
	ModuleKeys      []string
	SelectedIdx     int
	SelectedView    ViewHeader
	CurrentCategory int
	StatusMessage   string
	StatusIsError   bool
	IsLoading       bool
	LoadingMessage  string
	spinner         spinner.Model
}

func (m ViewModel) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m *ViewModel) getCurrentModule() models.Module {
	if m.SelectedView != "items" || len(m.ModuleKeys) == 0 {
		return nil
	}
	categoryName := m.CategoryKeys[m.CurrentCategory]
	moduleName := m.ModuleKeys[m.SelectedIdx]
	return m.Categories[categoryName][moduleName]
}

func (m *ViewModel) updateModuleKeys() {
	categoryName := m.CategoryKeys[m.CurrentCategory]
	modules := m.Categories[categoryName]
	m.ModuleKeys = make([]string, 0, len(modules))
	for key := range modules {
		m.ModuleKeys = append(m.ModuleKeys, key)
	}
	sort.Strings(m.ModuleKeys)
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
			if m.SelectedIdx < len(m.getList())-1 {
				m.SelectedIdx++
			}
		case "enter":
			if m.SelectedView == "categories" {
				m.SelectedView = "items"
				m.CurrentCategory = m.SelectedIdx
				m.updateModuleKeys()
				m.SelectedIdx = 0
			} else if m.SelectedView == "items" {
				module := m.getCurrentModule()
				if module != nil {
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
			}
		case "esc":
			if m.SelectedView == "items" {
				m.SelectedView = "categories"
				m.SelectedIdx = m.CurrentCategory
				m.ModuleKeys = nil
			}
		case "delete", "backspace":
			if m.SelectedView == "items" {
				module := m.getCurrentModule()
				if module != nil && module.IsInstalled() {
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

func (m ViewModel) getList() []string {
	if m.SelectedView == "categories" {
		return m.CategoryKeys
	} else {
		categoryName := m.CategoryKeys[m.CurrentCategory]
		modules := m.Categories[categoryName]
		var names []string
		for _, key := range m.ModuleKeys {
			names = append(names, modules[key].GetName())
		}
		return names
	}
}

func (m ViewModel) View() string {
	var content strings.Builder

	// Title
	var title string
	if m.SelectedView == "categories" {
		title = titleStyle.Render(" Code Template Manager ")
	} else {
		categoryName := m.CategoryKeys[m.CurrentCategory]
		title = titleStyle.Render(fmt.Sprintf(" %s ", categoryName))
	}
	content.WriteString(title + "\n\n")

	// List items
	list := m.getList()
	for i, item := range list {
		var line string

		// Cursor indicator
		cursor := "  "
		if i == m.SelectedIdx {
			cursor = "▸ "
		}

		// For modules view, show checkbox with state
		if m.SelectedView == "items" {
			module := m.Categories[m.CategoryKeys[m.CurrentCategory]][m.ModuleKeys[i]]
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

			line = fmt.Sprintf("%s%s %s%s", cursor, checkbox, item, versionText)
		} else {
			// For categories, show count of modules
			categoryModules := m.Categories[item]
			installedCount := 0
			for _, mod := range categoryModules {
				if mod.IsInstalled() {
					installedCount++
				}
			}
			countText := countStyle.Render(fmt.Sprintf("(%d/%d)", installedCount, len(categoryModules)))
			line = fmt.Sprintf("%s%s %s", cursor, item, countText)
		}

		// Style based on selection
		if i == m.SelectedIdx {
			line = selectedStyle.Render(line)
		}

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
	var helpText string
	if m.SelectedView == "categories" {
		helpText = "↑/↓ navigate • enter select • q quit"
	} else {
		helpText = "↑/↓ navigate • enter install/update • del uninstall • esc back • q quit"
	}
	content.WriteString(helpStyle.Render(helpText))

	// Wrap in container
	return containerStyle.Render(content.String())
}

func main() {
	if err := autoinit.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Initialization failed: %v\n", err)
		os.Exit(1)
	}

	content := helpers.GetContent()

	var categoryKeys []string
	for key := range content {
		categoryKeys = append(categoryKeys, key)
	}
	sort.Strings(categoryKeys)

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(primaryColor)

	m := ViewModel{
		Categories:      content,
		CategoryKeys:    categoryKeys,
		SelectedIdx:     0,
		SelectedView:    "categories",
		CurrentCategory: 0,
		spinner:         s,
	}

	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		panic(err)
	}
}
