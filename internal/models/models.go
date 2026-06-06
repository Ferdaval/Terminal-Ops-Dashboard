package models

type ViewType int

const (
	SystemView ViewType = iota
	DockerView
	GitHubView
)

var ViewNames = map[ViewType]string{
	SystemView: "System",
	DockerView: "Docker",
	GitHubView: "GitHub",
}

var ViewShortcuts = map[int]ViewType{
	'1': SystemView,
	'2': DockerView,
	'3': GitHubView,
}

type ContainerStatus string

const (
	StatusRunning    ContainerStatus = "running"
	StatusExited     ContainerStatus = "exited"
	StatusPaused     ContainerStatus = "paused"
	StatusRestarting ContainerStatus = "restarting"
)

type DockerContainer struct {
	ID      string
	Name    string
	Image   string
	Status  ContainerStatus
	Ports   []string
	Error   string
	State   string
}

type DockerMetrics struct {
	Containers []DockerContainer
	Error      string
}

type NotificationType string

const (
	NotificationPR     NotificationType = "PullRequest"
	NotificationIssue  NotificationType = "Issue"
	NotificationMention NotificationType = "Mention"
	NotificationRelease NotificationType = "Release"
)

type Notification struct {
	ID         string
	Title      string
	Repository string
	Type       NotificationType
	URL        string
	Reason     string
	Unread     bool
	UpdatedAt  string
}

type GitHubMetrics struct {
	Notifications []Notification
	Error         string
	Authenticated bool
}
