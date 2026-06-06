# 🖥️ Custom TUI Dashboard

A high-performance Terminal User Interface (TUI) dashboard built with Go and Bubbletea. It centralizes your daily workflow by combining real-time system resource monitoring, interactive Docker container management, and live GitHub notification tracking into a single command-line tool, ensuring your hands never have to leave the keyboard.

## ✨ Features

### 📊 System Monitoring
* **Real-time Metrics:** Monitor CPU, Memory, Disk, and Network throughput updated every second.
* **Visualizations:** Sparkline graphs for historical trend visualization and color-coded gauges (Green to Red based on usage).
* **Non-blocking:** Background workers poll data without blocking the UI thread.

### 🐳 Docker Management
* **Container Control:** View all running and stopped containers with color-coded statuses.
* **Quick Actions:** Start, stop, or restart containers instantly using keyboard shortcuts.
* **Log Streaming:** Dedicated log viewer to stream the last 20 lines of container logs.
* **Graceful Handling:** Smart detection and clear error messaging if the Docker daemon is unavailable.

### 🐙 GitHub Notifications
* **Live Tracking:** View unread notifications categorized by type (PR, Issue, Release, Mention) with clear icons.
* **Interactive Management:** Mark notifications as read, open them directly in your OS web browser, or filter by repository using fuzzy search.
* **Secure:** Authenticates via your GitHub Personal Access Token (PAT).

### ⌨️ Global UI & Navigation
* **Vim-Style Bindings:** Use `h`, `j`, `k`, `l` to navigate lists and scroll.
* **Tabbed Views:** Quickly switch between System, Docker, and GitHub views using numbers `1`, `2`, `3` or `Tab` / `Shift+Tab`.
* **Responsive:** Adapts automatically to terminal window resizing.

## 🚀 Getting Started

### Prerequisites
* Go 1.21 or later
* Docker daemon running (for the Docker module)
* A terminal with 256-color support

### Setup & Installation

1. **Clone the repository and build the binary:**
 ```bash
         go build -o bin/dashboard ./cmd/dashboard/main.go
 ```  
2. Set your GitHub Token (Required for GitHub module):
    ```Bash
           export GITHUB_TOKEN=ghp_xxxxxxxxxxxxxxxxxxxxxxxxxxxx
    ```
     Recommended token scopes: notifications, read:user
5. Run the dashboard:
     ```Bash
     ./bin/dashboard
     ```
(Alternatively, run directly without building: go run ./cmd/dashboard/main.go)🎮 
##Controls
Key           Action
- 1, 2, 3   ->   Switch to System, Docker, or GitHub views
- Tab / Shift+Tab  ->  Cycle through views
- j / k or ↓ / ↑   ->   Navigate lists (Vim-style)
- q or Ctrl+C   ->   Quit application

#### Docker Specific:
- s: Start/Stop selected container
- r: Restart selected container
- l: View container logs
#### GitHub Specific:
- m: Mark notification as read
- o: Open notification in default web browser
- f: Filter notifications by repository (fuzzy search)
- Esc: Cancel filter mode

### 📁 Project Structure

├── cmd/
│   └── dashboard/         # Application entry point
│       └── main.go
├── internal/
│   ├── models/           # Shared data types
│   ├── theme/            # Styling and colors (Lipgloss)
│   └── ui/               # User interface components
│       ├── app.go        # Main Bubbletea model & navigation
│       └── views/        # Individual view modules (System, Docker, GitHub)
├── bin/                  # Build output
├── go.mod                # Dependency manifest
└── README.md

### 📦 Dependencies
- Bubbletea - TUI framework
- Lipgloss - Terminal styling
- gopsutil - System metrics
- Docker SDK for Go - Docker API integration
- go-github - GitHub API integration
- OAuth2 - Authentication
### 📄 License
This project is open-source and available under the MIT License.
