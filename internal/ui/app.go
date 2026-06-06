package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/sadra/tui-dashboard/internal/models"
	"github.com/sadra/tui-dashboard/internal/theme"
	"github.com/sadra/tui-dashboard/internal/ui/views"
)

type Model struct {
	currentView models.ViewType
	systemView  *views.SystemView
	dockerView  *views.DockerView
	githubView  *views.GitHubView
	width       int
	height      int
	quitting    bool
}

func NewApp() *Model {
	return &Model{
		currentView: models.SystemView,
		systemView:  views.NewSystemView(),
		dockerView:  views.NewDockerView(),
		githubView:  views.NewGitHubView(),
	}
}

func (m *Model) Init() tea.Cmd {
	return tea.Batch(
		m.systemView.StartWorker(),
		m.dockerView.StartWorker(),
		m.githubView.StartWorker(),
	)
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			m.systemView.Cleanup()
			m.dockerView.Cleanup()
			m.githubView.Cleanup()
			return m, tea.Quit
		case "1":
			m.currentView = models.SystemView
		case "2":
			m.currentView = models.DockerView
		case "3":
			m.currentView = models.GitHubView
		case "tab":
			m.cycleViewForward()
		case "shift+tab":
			m.cycleViewBackward()
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.updateSizes()
	}

	cmd = m.updateCurrentView(msg)
	return m, cmd
}

func (m *Model) View() string {
	if m.quitting {
		return ""
	}

	tabBar := m.renderTabs()
	viewContent := m.renderCurrentView()

	layout := lipgloss.JoinVertical(
		lipgloss.Left,
		tabBar,
		viewContent,
	)

	return layout
}

func (m *Model) renderTabs() string {
	var tabs []string

	for i := 0; i < 3; i++ {
		view := models.ViewType(i)
		name := models.ViewNames[view]

		if view == m.currentView {
			tabs = append(tabs, theme.TabActiveStyle.Render(name))
		} else {
			tabs = append(tabs, theme.TabInactiveStyle.Render(name))
		}
	}

	return lipgloss.JoinHorizontal(lipgloss.Left, tabs...)
}

func (m *Model) renderCurrentView() string {
	switch m.currentView {
	case models.SystemView:
		return m.systemView.Render()
	case models.DockerView:
		return m.dockerView.Render()
	case models.GitHubView:
		return m.githubView.Render()
	default:
		return theme.PaneStyle.Render("Unknown view")
	}
}

func (m *Model) updateCurrentView(msg tea.Msg) tea.Cmd {
	switch m.currentView {
	case models.SystemView:
		return m.systemView.Update(msg)
	case models.DockerView:
		return m.dockerView.Update(msg)
	case models.GitHubView:
		return m.githubView.Update(msg)
	}
	return nil
}

func (m *Model) updateSizes() {
	viewWidth := m.width
	viewHeight := m.height - 3

	m.systemView.SetSize(viewWidth, viewHeight)
	m.dockerView.SetSize(viewWidth, viewHeight)
	m.githubView.SetSize(viewWidth, viewHeight)
}

func (m *Model) cycleViewForward() {
	m.currentView = models.ViewType((int(m.currentView) + 1) % 3)
}

func (m *Model) cycleViewBackward() {
	m.currentView = models.ViewType((int(m.currentView) - 1 + 3) % 3)
}
