package workers

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/google/go-github/v62/github"
	"golang.org/x/oauth2"
	"github.com/sadra/tui-dashboard/internal/models"
)

type GitHubWorker struct {
	interval   time.Duration
	client     *github.Client
	token      string
	authenticated bool
}

func NewGitHubWorker(interval time.Duration) *GitHubWorker {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		return &GitHubWorker{
			interval:      interval,
			token:         "",
			authenticated: false,
		}
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	return &GitHubWorker{
		interval:      interval,
		client:        client,
		token:         token,
		authenticated: true,
	}
}

func (gw *GitHubWorker) FetchNotifications(ctx context.Context) (models.GitHubMetrics, error) {
	metrics := models.GitHubMetrics{
		Notifications: []models.Notification{},
		Authenticated: gw.authenticated,
	}

	if !gw.authenticated {
		metrics.Error = "GitHub token not set. Set GITHUB_TOKEN environment variable."
		return metrics, nil
	}

	if gw.client == nil {
		metrics.Error = "GitHub client not initialized"
		return metrics, nil
	}

	opts := &github.NotificationListOptions{
		All: false,
		ListOptions: github.ListOptions{
			PerPage: 50,
		},
	}

	notifications, _, err := gw.client.Activity.ListNotifications(ctx, opts)
	if err != nil {
		metrics.Error = fmt.Sprintf("Failed to fetch notifications: %v", err)
		return metrics, nil
	}

	for _, notif := range notifications {
		notifType := gw.parseNotificationType(notif.Reason)

		id := ""
		if notif.ID != nil {
			id = *notif.ID
		}

		title := ""
		if notif.Subject != nil && notif.Subject.Title != nil {
			title = *notif.Subject.Title
		}

		repo := ""
		if notif.Repository != nil && notif.Repository.Name != nil {
			repo = *notif.Repository.Name
		}

		url := ""
		if notif.Subject != nil && notif.Subject.URL != nil {
			url = *notif.Subject.URL
		}

		reason := ""
		if notif.Reason != nil {
			reason = *notif.Reason
		}

		unread := false
		if notif.Unread != nil {
			unread = *notif.Unread
		}

		notification := models.Notification{
			ID:         id,
			Title:      title,
			Repository: repo,
			Type:       notifType,
			URL:        url,
			Reason:     reason,
			Unread:     unread,
			UpdatedAt:  notif.UpdatedAt.Format("15:04"),
		}

		metrics.Notifications = append(metrics.Notifications, notification)
	}

	return metrics, nil
}

func (gw *GitHubWorker) parseNotificationType(reason *string) models.NotificationType {
	if reason == nil {
		return models.NotificationMention
	}

	switch strings.ToLower(*reason) {
	case "pull_request":
		return models.NotificationPR
	case "issue":
		return models.NotificationIssue
	case "release":
		return models.NotificationRelease
	default:
		return models.NotificationMention
	}
}

func (gw *GitHubWorker) MarkAsRead(ctx context.Context, notificationID string) error {
	if !gw.authenticated || gw.client == nil {
		return fmt.Errorf("not authenticated")
	}

	_, err := gw.client.Activity.MarkThreadRead(ctx, notificationID)
	return err
}

func (gw *GitHubWorker) OpenInBrowser(url string) error {
	// Extract the HTML URL from the API URL if needed
	htmlURL := url
	if strings.Contains(url, "api.github.com") {
		htmlURL = strings.Replace(url, "api.github.com/repos", "github.com", 1)
	}

	var cmd *exec.Cmd
	switch runtime := os.Getenv("GOOS"); runtime {
	case "darwin":
		cmd = exec.Command("open", htmlURL)
	case "linux":
		cmd = exec.Command("xdg-open", htmlURL)
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", htmlURL)
	default:
		return fmt.Errorf("unsupported OS: %s", runtime)
	}

	return cmd.Start()
}

func (gw *GitHubWorker) StartPolling(ctx context.Context, ch chan<- models.GitHubMetrics) {
	ticker := time.NewTicker(gw.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			metrics, _ := gw.FetchNotifications(ctx)
			select {
			case ch <- metrics:
			case <-ctx.Done():
				return
			}
		}
	}
}

func (gw *GitHubWorker) Cleanup() {
}
