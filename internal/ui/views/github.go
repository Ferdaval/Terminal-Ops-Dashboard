package views

import (
	"context"
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/sadra/tui-dashboard/internal/models"
	"github.com/sadra/tui-dashboard/internal/theme"
	"github.com/sadra/tui-dashboard/internal/workers"
)

type GitHubMetricsMsg models.GitHubMetrics

type GitHubView struct {
	width            int
	height           int
	metrics          models.GitHubMetrics
	worker           *workers.GitHubWorker
	ctx              context.Context
	cancel           context.CancelFunc
	metricsChad      chan models.GitHubMetrics
	started          bool
	selectedIndex    int
	filterMode       bool
	filterText       string
	filteredNotifs   []models.Notification
	markingAsRead    bool
}

func NewGitHubView() *GitHubView {
	ctx, cancel := context.WithCancel(context.Background())
	return &GitHubView{
		worker:         workers.NewGitHubWorker(3 * time.Second),
		ctx:            ctx,
		cancel:         cancel,
		metricsChad:    make(chan models.GitHubMetrics, 1),
		metrics:        models.GitHubMetrics{Notifications: []models.Notification{}},
		filteredNotifs: []models.Notification{},
	}
}

func (gv *GitHubView) StartWorker() tea.Cmd {
	if gv.started {
		return nil
	}
	gv.started = true

	go gv.worker.StartPolling(gv.ctx, gv.metricsChad)
	return gv.waitForMetricsCmd()
}

func (gv *GitHubView) Update(msg interface{}) tea.Cmd {
	switch msg := msg.(type) {
	case GitHubMetricsMsg:
		gv.metrics = models.GitHubMetrics(msg)
		gv.updateFilteredNotifications()
		if gv.selectedIndex >= len(gv.filteredNotifs) && len(gv.filteredNotifs) > 0 {
			gv.selectedIndex = len(gv.filteredNotifs) - 1
		}
		return gv.waitForMetricsCmd()

	case tea.KeyMsg:
		if gv.filterMode {
			switch msg.String() {
			case "esc":
				gv.filterMode = false
				gv.filterText = ""
				gv.selectedIndex = 0
				gv.updateFilteredNotifications()
			case "backspace":
				if len(gv.filterText) > 0 {
					gv.filterText = gv.filterText[:len(gv.filterText)-1]
					gv.selectedIndex = 0
					gv.updateFilteredNotifications()
				}
			default:
				if len(msg.String()) == 1 && msg.String() >= "a" && msg.String() <= "z" {
					gv.filterText += msg.String()
					gv.selectedIndex = 0
					gv.updateFilteredNotifications()
				} else if len(msg.String()) == 1 && msg.String() >= "A" && msg.String() <= "Z" {
					gv.filterText += strings.ToLower(msg.String())
					gv.selectedIndex = 0
					gv.updateFilteredNotifications()
				} else if msg.String() == " " {
					gv.filterText += " "
					gv.selectedIndex = 0
					gv.updateFilteredNotifications()
				}
			}
			return nil
		}

		switch msg.String() {
		case "j", "down":
			if gv.selectedIndex < len(gv.filteredNotifs)-1 {
				gv.selectedIndex++
			}
		case "k", "up":
			if gv.selectedIndex > 0 {
				gv.selectedIndex--
			}
		case "f":
			gv.filterMode = true
			gv.filterText = ""
		case "m":
			if gv.selectedIndex < len(gv.filteredNotifs) {
				notif := gv.filteredNotifs[gv.selectedIndex]
				gv.markingAsRead = true
				go func() {
					gv.worker.MarkAsRead(gv.ctx, notif.ID)
					time.Sleep(500 * time.Millisecond)
					gv.markingAsRead = false
				}()
			}
		case "o":
			if gv.selectedIndex < len(gv.filteredNotifs) {
				notif := gv.filteredNotifs[gv.selectedIndex]
				go gv.worker.OpenInBrowser(notif.URL)
			}
		}
	}
	return nil
}

func (gv *GitHubView) updateFilteredNotifications() {
	if gv.filterText == "" {
		gv.filteredNotifs = gv.metrics.Notifications
		return
	}

	gv.filteredNotifs = []models.Notification{}
	filterLower := strings.ToLower(gv.filterText)
	for _, notif := range gv.metrics.Notifications {
		if strings.Contains(strings.ToLower(notif.Repository), filterLower) {
			gv.filteredNotifs = append(gv.filteredNotifs, notif)
		}
	}
}

func (gv *GitHubView) waitForMetricsCmd() tea.Cmd {
	return func() tea.Msg {
		select {
		case m := <-gv.metricsChad:
			return GitHubMetricsMsg(m)
		case <-gv.ctx.Done():
			return nil
		}
	}
}

func (gv *GitHubView) SetSize(width, height int) {
	gv.width = width
	gv.height = height
}

func (gv *GitHubView) notificationIcon(notifType models.NotificationType) string {
	switch notifType {
	case models.NotificationPR:
		return "[PR]"
	case models.NotificationIssue:
		return "[Issue]"
	case models.NotificationRelease:
		return "[Release]"
	default:
		return "[Mention]"
	}
}

func (gv *GitHubView) notificationColor(notifType models.NotificationType) lipgloss.Color {
	switch notifType {
	case models.NotificationPR:
		return theme.Color.Primary
	case models.NotificationIssue:
		return theme.Color.Error
	case models.NotificationRelease:
		return theme.Color.Success
	default:
		return theme.Color.Accent
	}
}

func (gv *GitHubView) Render() string {
	if gv.filterMode {
		return gv.renderFilterMode()
	}

	if !gv.metrics.Authenticated {
		return theme.PaneStyle.Width(gv.width - 4).Render(
			fmt.Sprintf("%s\n\n%s\n\n%s",
				theme.HeaderStyle.Render("🔔 GitHub Notifications"),
				theme.HelpStyle.Render("⚠️  GitHub token not set"),
				theme.HelpStyle.Render("Set GITHUB_TOKEN environment variable to use this feature"),
			),
		)
	}

	if gv.metrics.Error != "" {
		return theme.PaneStyle.Width(gv.width - 4).Render(
			fmt.Sprintf("%s\n\n%s",
				theme.HeaderStyle.Render("🔔 GitHub Notifications"),
				theme.HelpStyle.Render("⚠️  "+gv.metrics.Error),
			),
		)
	}

	if len(gv.filteredNotifs) == 0 {
		msg := "No unread notifications"
		if gv.filterText != "" {
			msg = fmt.Sprintf("No notifications match: %s", gv.filterText)
		}
		return theme.PaneStyle.Width(gv.width - 4).Render(
			fmt.Sprintf("%s\n\n%s",
				theme.HeaderStyle.Render("🔔 GitHub Notifications"),
				theme.HelpStyle.Render(msg),
			),
		)
	}

	notifLines := []string{
		theme.HeaderStyle.Render("🔔 GitHub Notifications"),
		fmt.Sprintf("Showing %d unread notification(s)", len(gv.filteredNotifs)),
		"",
	}

	for i, n := range gv.filteredNotifs {
		icon := gv.notificationIcon(n.Type)
		color := gv.notificationColor(n.Type)
		styledIcon := lipgloss.NewStyle().Foreground(color).Render(icon)

		title := n.Title
		if len(title) > 40 {
			title = title[:37] + "..."
		}

		repo := n.Repository
		if len(repo) > 20 {
			repo = repo[:17] + "..."
		}

		line := fmt.Sprintf("%s %-22s %s %s", styledIcon, repo, title, n.UpdatedAt)

		if i == gv.selectedIndex {
			line = lipgloss.NewStyle().
				Background(theme.Color.Primary).
				Foreground(theme.Color.Background).
				Render(line)
		}

		notifLines = append(notifLines, line)
	}

	notifLines = append(notifLines, "")
	notifLines = append(notifLines, theme.HeaderStyle.Render("Actions"))
	if gv.markingAsRead && gv.selectedIndex < len(gv.filteredNotifs) {
		notifLines = append(notifLines, "  (j/k) Navigate  (m) Marking as read...  (o) Open  (f) Filter")
	} else {
		notifLines = append(notifLines, "  (j/k) Navigate  (m) Mark as read  (o) Open in browser  (f) Filter by repo")
	}

	content := strings.Join(notifLines, "\n")
	return theme.PaneStyle.Width(gv.width - 4).Render(content)
}

func (gv *GitHubView) renderFilterMode() string {
	content := fmt.Sprintf(`%s

Filter repositories:
  %s_

(Type to filter, Esc to cancel)

Matching repositories:
`,
		theme.HeaderStyle.Render("🔍 Filter Notifications"),
		gv.filterText,
	)

	count := 0
	repos := make(map[string]bool)
	for _, notif := range gv.metrics.Notifications {
		if strings.Contains(strings.ToLower(notif.Repository), strings.ToLower(gv.filterText)) {
			if !repos[notif.Repository] {
				repos[notif.Repository] = true
				content += fmt.Sprintf("\n  • %s", notif.Repository)
				count++
				if count >= 10 {
					content += "\n  ..."
					break
				}
			}
		}
	}

	if count == 0 {
		content += "\n  (no matches)"
	}

	return theme.PaneStyle.Width(gv.width - 4).Render(content)
}

func (gv *GitHubView) Cleanup() {
	gv.cancel()
	gv.worker.Cleanup()
}
