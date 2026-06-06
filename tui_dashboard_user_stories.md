# Custom TUI Dashboard - User Stories

Building a custom TUI dashboard using a library like Bubbletea (Go) or Ratatui (Rust) is a fantastic project. It hits the sweet spot of systems programming while providing immediate, highly satisfying visual feedback. 

Here is a comprehensive set of user stories broken down by the core modules and overall experience.

## 1. System Monitoring Module

* **As a** system admin, **I want to** view real-time CPU and memory usage graphs **so that** I can quickly spot performance spikes or resource exhaustion.
    * *Acceptance Criteria:* Data updates at a configurable interval (e.g., every 1 second). Visual indicators (like color changes to red) trigger when usage exceeds 90%.
* **As a** developer, **I want to** see a breakdown of my disk space across major partitions **so that** I know when I need to clean up old files or caches.
    * *Acceptance Criteria:* Displays total, used, and available space in GB/TB. 
* **As a** power user, **I want to** monitor active network throughput (upload/download speeds) **so that** I can verify my connection stability.
    * *Acceptance Criteria:* Speeds are displayed in Kbps/Mbps with a rolling historical sparkline graph.

## 2. Docker Management Module

* **As a** backend engineer, **I want to** view a list of all Docker containers (running and stopped) **so that** I have a high-level overview of my local environments.
    * *Acceptance Criteria:* List displays container name, ID, status (color-coded), and mapped ports.
* **As a** developer, **I want to** use keyboard shortcuts to start, stop, or restart a highlighted container **so that** I can manage my stack without typing long CLI commands.
    * *Acceptance Criteria:* Pressing 's' stops a running container; 'r' restarts it. A confirmation or loading state is shown during the action.
* **As a** debugger, **I want to** select a container and stream its most recent logs in a dedicated pane **so that** I can troubleshoot application errors on the fly.
    * *Acceptance Criteria:* Logs auto-scroll, with an option to pause scrolling or filter by keyword.

## 3. GitHub Notifications Module

* **As a** maintainer, **I want to** see a list of my unread GitHub notifications **so that** I can stay on top of issues and pull requests needing my attention.
    * *Acceptance Criteria:* Authenticates via GitHub Personal Access Token (PAT). Displays repository name, notification title, and type (Issue, PR, Mention).
* **As a** collaborator, **I want to** filter my notifications by repository **so that** I can focus on the project I am currently working on.
    * *Acceptance Criteria:* Pressing a hotkey opens a filter prompt to type and fuzzy-search repository names.
* **As a** user, **I want to** press a key to mark a notification as read or open it in my default web browser **so that** I can seamlessly transition to taking action.
    * *Acceptance Criteria:* Pressing 'm' marks as read and removes it from the list; pressing 'o' opens the default OS browser to the specific GitHub URL.

## 4. Overall TUI Experience & Navigation

* **As a** terminal enthusiast, **I want to** navigate entirely using Vim-style keybindings (h, j, k, l) **so that** my hands never have to leave the keyboard.
    * *Acceptance Criteria:* Standard Vim bindings move focus between panes and scroll through lists.
* **As a** user, **I want to** switch between different dashboard views (Stats, Docker, GitHub) using tabs **so that** the screen doesn't become cluttered with too much information at once.
    * *Acceptance Criteria:* Numbers (1, 2, 3) or Tab/Shift-Tab cycle through the different main views.
* **As a** developer, **I want to** configure the dashboard's layout and default startup view via a YAML/JSON configuration file **so that** I can tailor the tool to my specific workflow.
    * *Acceptance Criteria:* App reads a config file on startup to determine which modules to load and what colors/themes to use.
