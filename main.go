package main

import (
	"code-template/helpers"
	"code-template/models"
	"fmt"
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ViewHeader string

var (
	categories ViewHeader = "categories"
	modules    ViewHeader = "modules"
)

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

	checkboxNotInstalled = lipgloss.NewStyle().
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
}

func (m ViewModel) Init() tea.Cmd {
	return nil
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
	case tea.KeyMsg:
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
			}
		case "esc", "backspace":
			if m.SelectedView == "items" {
				m.SelectedView = "categories"
				m.SelectedIdx = m.CurrentCategory
				m.ModuleKeys = nil
			}
		case " ":
			if m.SelectedView == "items" {
				module := m.getCurrentModule()
				if module != nil {
					if module.IsInstalled() {
						if module.Uninstall() {
							m.StatusMessage = fmt.Sprintf("✓ Uninstalled %s", module.GetName())
						} else {
							m.StatusMessage = fmt.Sprintf("✗ Failed to uninstall %s", module.GetName())
							m.StatusIsError = true
						}
					} else {
						if module.Install() {
							m.StatusMessage = fmt.Sprintf("✓ Installed %s", module.GetName())
						} else {
							m.StatusMessage = fmt.Sprintf("✗ Failed to install %s", module.GetName())
							m.StatusIsError = true
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

		// For modules view, show checkbox
		if m.SelectedView == "items" {
			module := m.Categories[m.CategoryKeys[m.CurrentCategory]][m.ModuleKeys[i]]
			var checkbox string
			if module.IsInstalled() {
				checkbox = checkboxInstalled.Render("[✓]")
			} else {
				checkbox = checkboxNotInstalled.Render("[ ]")
			}
			line = fmt.Sprintf("%s%s %s", cursor, checkbox, item)
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

	// Status message
	if m.StatusMessage != "" {
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
		helpText = "↑/↓ navigate • space toggle • esc back • q quit"
	}
	content.WriteString(helpStyle.Render(helpText))

	// Wrap in container
	return containerStyle.Render(content.String())
}

func main() {
	content := helpers.GetContent()

	var categoryKeys []string
	for key := range content {
		categoryKeys = append(categoryKeys, key)
	}
	sort.Strings(categoryKeys)

	m := ViewModel{
		Categories:      content,
		CategoryKeys:    categoryKeys,
		SelectedIdx:     0,
		SelectedView:    "categories",
		CurrentCategory: 0,
	}

	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		panic(err)
	}
}
