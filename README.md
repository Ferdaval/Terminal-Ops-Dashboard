# Custom TUI Dashboard

A high-performance Terminal User Interface (TUI) dashboard built with Go and Bubbletea. Monitor system resources, manage Docker containers, and track GitHub notifications—all without leaving your terminal.

## Project Structure

```
.
├── cmd/
│   └── dashboard/         # Application entry point
│       └── main.go
├── internal/
│   ├── models/           # Shared data types
│   │   └── models.go
│   ├── theme/            # Styling and colors
│   │   └── theme.go
│   └── ui/               # User interface components
│       ├── app.go        # Main Bubbletea model & navigation
│       └── views/        # Individual view modules
│           ├── system.go
│           ├── docker.go
│           └── github.go
├── bin/                  # Build output
│   └── dashboard
├── go.mod               # Dependency manifest
└── README.md
```

## Getting Started

### Prerequisites

- Go 1.21 or later
- A terminal with 256-color support

### Build & Run

```bash
# Build the binary
go build -o bin/dashboard ./cmd/dashboard/main.go

# Run the dashboard
./bin/dashboard
```

Alternatively, run directly without building:
```bash
go run ./cmd/dashboard/main.go
```

## Controls

### Navigation
- **1, 2, 3**: Switch to System, Docker, or GitHub views
- **Tab / Shift+Tab**: Cycle through views
- **j/k**: Navigate lists (Vim-style)
- **h/l**: Scroll horizontally
- **q** or **Ctrl+C**: Quit

### View-Specific Actions
- **System**: Monitor CPU, memory, disk, and network
- **Docker**: Start/stop/restart containers (s, r, l for logs)
- **GitHub**: Mark notifications as read (m) or open in browser (o)

## Step 1: Scaffolding & Navigation (Complete ✓)

This initial implementation includes:

- **Tabbed Navigation**: Switch between 3 main views (System, Docker, GitHub)
- **Global Keybindings**: Vim-style navigation (h, j, k, l) + view shortcuts (1, 2, 3, Tab/Shift+Tab)
- **Theming System**: Consistent colors and styles using Lipgloss
- **Modular View Architecture**: Each module is isolated and easy to extend
- **Error-Safe Design**: Views don't crash the entire app on errors
- **Responsive Layout**: Window resizing is handled automatically

## Step 2: System Monitoring Module (Complete ✓)

Fully implemented real-time system monitoring with:

- **Real-time CPU/Memory/Disk/Network metrics** updated every 1 second
- **Metrics History Tracking**: Last 60 seconds of data captured in ring buffers
- **Sparkline Graphs**: Historical trend visualization using Unicode characters (▁▂▃▄▅▆▇█)
- **Color-Coded Gauges**: Red (≥90%), Orange (≥75%), Blue (≥50%), Green (<50%)
- **Network Throughput**: Upload/Download speeds with sparklines (in MB/s)
- **Disk Space**: Total, used, and available display
- **Non-blocking I/O**: Background worker polls data without blocking the UI thread

## Step 3: Docker Management Module (Complete ✓)

Fully implemented Docker container management with:

- **Docker API Integration**: Uses official Docker SDK with graceful daemon detection
- **Container Listing**: Shows all containers (running and stopped) with status color coding
- **Container Status Display**: 
  - 🟢 Running (green)
  - 🟠 Exited (orange)
  - 🔵 Paused (blue)
- **Container Navigation**: Vim-style j/k keys to select containers
- **Container Actions**:
  - `s` - Start or stop a container (toggles based on current status)
  - `r` - Restart a running container
  - `l` - View last 20 lines of container logs
- **Container Logs Viewer**: Dedicated view showing logs with scrolling
- **Error Handling**: Gracefully handles Docker daemon unavailability
- **Responsive UI**: Real-time updates every 2 seconds, non-blocking I/O

### Features

- Container ID, name, image, status, and port mappings
- Selection highlighting for active container
- Non-blocking Docker API calls with 2-second polling
- Automatic cleanup on exit
- Clear error messages when Docker is unavailable

## Step 4: GitHub Notifications Module (Complete ✓)

Fully implemented GitHub integration with:

- **GitHub API Integration**: Uses official go-github library with OAuth2 authentication
- **PAT Authentication**: Reads `GITHUB_TOKEN` environment variable for secure access
- **Notification Listing**: Displays all unread GitHub notifications with:
  - Repository name
  - Notification title
  - Type (PR, Issue, Release, Mention)
  - Last update time
- **Notification Types with Icons**:
  - 🟣 [PR] - Pull Request notifications
  - 🔴 [Issue] - Issue notifications
  - 🟢 [Release] - Release notifications
  - 🔵 [Mention] - Mention notifications
- **Vim-style Navigation**:
  - `j/k` or `↓/↑` - Navigate notifications
- **Notification Actions**:
  - `m` - Mark as read (marks thread as read)
  - `o` - Open in default web browser
  - `f` - Filter notifications by repository (fuzzy search)
- **Repository Filtering**:
  - Press `f` to enter filter mode
  - Type to search repositories
  - Shows matching repos as you type
  - Press `Esc` to cancel
- **Error Handling**: Gracefully handles:
  - Missing GitHub token
  - API authentication failures
  - Network issues
- **Real-time Updates**: Polling every 3 seconds with non-blocking I/O

### Setup

Set your GitHub Personal Access Token:
```bash
export GITHUB_TOKEN=ghp_xxxxxxxxxxxxxxxxxxxxxxxxxxxx
```

Recommended scopes for the token:
- `notifications` - Read access to notifications
- `read:user` - User data

### Features

- Non-blocking background polling
- Selection highlighting for current notification
- Real-time notification count
- Smart error messaging
- Automatic browser opening (macOS, Linux, Windows)
- Filtered notification viewing

---

## All Steps Complete! 🎉

The Custom TUI Dashboard now has full implementations of:
1. ✅ **System Monitoring** - CPU, Memory, Disk, Network with sparklines
2. ✅ **Docker Management** - List, control, and view logs
3. ✅ **GitHub Notifications** - Fetch, filter, and manage

### Running the Dashboard

```bash
# Set your GitHub token (optional, but required for GitHub module)
export GITHUB_TOKEN=ghp_xxxxxxxxxxxxxxxxxxxxxxxxxxxx

# Run the dashboard
./bin/dashboard

# Or build and run directly
go run ./cmd/dashboard/main.go
```

### Navigation

- **1, 2, 3** - Switch between System, Docker, GitHub views
- **Tab / Shift+Tab** - Cycle through views
- **j/k or ↓/↑** - Navigate lists
- **q or Ctrl+C** - Quit

## Dependencies

- `github.com/charmbracelet/bubbletea` - TUI framework
- `github.com/charmbracelet/lipgloss` - Styling library
- `github.com/shirou/gopsutil/v3` - System metrics (CPU, Memory, Disk, Network)
- `github.com/docker/docker` - Docker API client
- `github.com/google/go-github/v62` - GitHub API client
- `golang.org/x/oauth2` - OAuth2 authentication for GitHub
