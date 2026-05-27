# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

> Gebruik [Task](https://taskfile.dev) als taakrunner (`brew install go-task` / `scoop install task`).

```bash
task          # build + run
task build    # compileer binary
task run      # start vanuit source
task test     # alle tests
task lint     # golangci-lint
task snapshot # test release lokaal (goreleaser)
task release  # publiceer GitHub release (vereist GITHUB_TOKEN)
task clean    # verwijder build artifacts
```

Directe go-commando's (als Task niet beschikbaar is):
```bash
go run .
go build -o ceda-scoop .
go test ./...
go test ./internal/tui/...
go test ./internal/runner/...
```

## Architecture

**CEDA Store** (`ceda-scoop`) is a TUI desktop application for installing and running CEDA-curated R and Python projects without requiring admin rights. It uses [Bubbletea](https://github.com/charmbracelet/bubbletea) (Elm-inspired MVU pattern), [Bubbles](https://github.com/charmbracelet/bubbles) for UI components, and [Lipgloss](https://github.com/charmbracelet/lipgloss) for styling.

### Layer separation

| Layer | Package | Responsibility |
|-------|---------|----------------|
| TUI (frontend) | `internal/tui/` | Bubbletea model, all views, event handling, styles |
| Runner (backend) | `internal/runner/` | Step definitions, cross-platform exec, PowerShell wrapper |
| Scripts | `scripts/windows/modules/` | PowerShell modules for each installation step |

### Bubbletea model/update/view split

- `internal/tui/model.go` — Central `Model` struct (screen state, selected repo, progress, etc.)
- `internal/tui/update.go` — Handles all `tea.Msg` events (key presses, async completion messages)
- `internal/tui/view.go` — Dispatches rendering to per-screen view files
- `internal/tui/messages.go` — Custom message types (`SplashDoneMsg`, `InstallDoneMsg`, `RunStepDoneMsg`, etc.)
- `internal/tui/catalog.go` — Hardcoded repo catalog (add new projects here)
- `internal/tui/styles.go` — Color palette and all lipgloss styles

Each screen has a dedicated view file: `view_splash.go`, `view_browser.go`, `view_detail.go`, `view_install.go`, `view_run.go`, `view_project_picker.go`, `view_delete_confirm.go`, `view_settings.go`.

### Screen state machine

```
Splash → Store/Library tabs → Detail → Install → Library
                                      → Run → ProjectTypePicker (if ambiguous) → RunSteps
                                      → Delete → ConfirmDialog
```

### Runner / PowerShell integration

- `internal/runner/run.go` — Defines `Step` types and detection logic (R vs uv via lock file presence)
- `internal/runner/exec.go` — Cross-platform command execution (Windows/macOS/Linux)
- `internal/runner/ps.go` — Wraps PowerShell invocation

Each run step maps to a PowerShell module in `scripts/windows/modules/`. Steps for R projects: `scoop-install → r-install → r-config → r-sync → r-run`. Steps for Python/uv projects: `scoop-install → uv-install → uv-config → uv-sync → uv-run`.

### Key data

- Projects are cloned to `~/ceda/{project-id}/` by default (configurable in settings)
- The catalog is hardcoded in `internal/tui/catalog.go`; add new repos there
- Project type is auto-detected by the presence of `renv.lock` (R) or `uv.lock` (Python)
