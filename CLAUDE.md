# CEDA Store (`ceda-scoop`)

> Geef dit bestand mee aan elke nieuwe Claude Code sessie. Het laadt automatisch.

TUI desktop application for installing and running CEDA-curated R and Python projects without requiring admin rights.

**Status:** In ontwikkeling
**Versie:** 0.1.0
**Go versie:** 1.23+

---

## Tech Stack

| Library | Repo | Doel |
|---|---|---|
| Bubble Tea | `charmbracelet/bubbletea` | MVU core + mouse support |
| Bubbles | `charmbracelet/bubbles` | Pre-built componenten (tabs, viewport, spinner, etc.) |
| Lip Gloss | `charmbracelet/lipgloss` | Styling, borders, layout |
| Huh | `charmbracelet/huh` | Forms / user input (embeddable in bubbletea) |
| Glamour | `charmbracelet/glamour` | Markdown rendering |
| Harmonica | `charmbracelet/harmonica` | Physics-based spring animaties |
| Log | `charmbracelet/log` | Leveled logger — altijd naar file, nooit stdout |

---

## Architecture

**The Elm Architecture (TEA)** — enforced by BubbleTea. Every screen implements `Init / Update / View`. Side effects never happen directly inside `Update`; they return as `tea.Cmd` and come back as messages.

**Nested models** for multi-screen navigation. The root model owns the current screen and delegates `Update`/`View` to it. Screen transitions happen when the root catches a typed message (e.g. `ScriptSelectedMsg`). Each screen is self-contained and knows nothing about other screens.

**Loose coupling via interfaces.** The UI only knows about `runner.Runner` — never about the OS. `runner.NewRunner()` detects the platform at startup and injects the right implementation. Adding a new platform means adding a file, not touching the UI.

---

## Philosophy

- Logic lives in the layer that owns it: UI in `ui/`, OS execution in `runner/`, script metadata in `scripts/`
- No screen knows about another screen — only the root routes between them
- No platform detection outside `runner/`
- Side effects (running scripts) go through `tea.ExecProcess`, not raw `exec.Command`
- All `Msg` types in one place: `internal/messages.go`
- All state in `Model` struct — no global state outside Model
- Logging always to `debug.log` via `log.SetOutput(f)` — never `fmt.Println` to stdout

---

## Folder structure

```
cmd/
  main.go              # entry point — wires runner, registry, root model

internal/
  ui/
    root.go            # root model, routes between screens
    menu/
      model.go         # script list and selection
    runner/
      model.go         # progress and output display
    result/
      model.go         # success / error state

  runner/
    runner.go          # Runner interface + NewRunner() factory + all OS implementations
                       # UI never touches runtime.GOOS — runner resolves it at startup

  scripts/
    scripts.go         # Script struct (name, description, path) + LoadScripts()
                       # walks scripts/<platform>/, returns []Script for the menu to display

scripts/
  windows/             # .ps1 files
  mac/                 # .sh files
  linux/               # .sh files

docs/
  adr/                 # architecture decision records
  roadmap/             # versie planningen (rolling window)
    001-bubbletea-mvu.md

.claude/
  agents/              # subagents
  skills/              # skills / slash commands
  commands/            # legacy slash commands

CLAUDE.md
CHANGELOG.md
.mcp.json              # project-scoped MCP servers
```

---

## How it wires together

`main.go` is the only place that knows about everything. It runs once at startup:

```go
scripts, _ := scripts.LoadScripts()  // discovers scripts/<platform>/*.sh|.ps1
runner     := runner.NewRunner()      // picks Windows/Mac/Linux impl based on OS

tea.NewProgram(ui.New(scripts, runner)).Run()
```

From that point the UI only works with the results:
- `menu/model.go` renders `[]Script` as a list — no filesystem knowledge
- `runner/model.go` calls `m.runner.Command(script.Path)` — no OS knowledge
- neither package imports `runner` or `scripts` directly; they receive them via the root model

---

## Key types

```go
// Runner interface — only thing UI knows about
type Runner interface {
    Command(script string) *exec.Cmd
}

// Screen transition messages — all Msg types live in internal/messages.go
type ScriptSelectedMsg struct { Script scripts.Script }
type ScriptDoneMsg     struct { Err error }

// Root model
type Model struct {
    currentScreen screen
    menu          menu.Model
    runner        runner.Model
    result        result.Model
    r             runner.Runner
}
```

---

## Mouse Support

```go
// main.go
p := tea.NewProgram(model, tea.WithMouseCellMotion())
```

- `tea.MouseMsg` afvangen in `Update()`
- Hit detection via `msg.X` / `msg.Y` — posities bijhouden in Model struct
- Klikbare gebieden als `lipgloss.Rect` of simpele int-bounds opslaan in Model

---

## Don'ts

- ❌ Geen `fmt.Println` / `fmt.Fprintf` naar stdout of stderr — breekt de TUI
- ❌ Geen global state buiten `Model` struct
- ❌ Geen `panic()` — altijd error teruggeven
- ❌ Geen inline Lip Gloss styles — altijd via `styles/styles.go`
- ❌ Geen `exec.Command` direct in `Update` — gebruik `tea.ExecProcess`
- ❌ Geen platform detection buiten `runner/`
- ❌ Geen HTTP calls zonder `context` + timeout
- ❌ Geen tests op TUI rendering — alleen pure functies testen

---

## Tooling & Commands

```bash
go run ./cmd/main.go                    # starten
go build -o ceda-scoop ./cmd/main.go   # builden
go test ./...                           # testen
go vet ./...                            # statische analyse
go mod tidy                             # dependencies opruimen
go mod download                         # dependencies downloaden (CI)
golangci-lint run                       # linting
goreleaser release --snapshot --clean  # lokale release testen
goreleaser release                      # release bouwen
git-cliff -o CHANGELOG.md             # changelog genereren
tail -f debug.log                       # logs volgen in tweede terminal
```

---

## Tests

- Unit tests alleen op pure functies zonder I/O
- `Model.Update()` is een pure functie → goed testbaar
- Testbestand naast source: `model_test.go`
- Geen integration tests voor TUI rendering

---

## Commits

Follow Conventional Commits. Format: `type(scope): description`

```
feat: add linux runner
fix: script path resolution on windows
refactor: collapse runner into single file
chore: update dependencies
docs: update CLAUDE.md
test: add runner factory tests
ci: add release workflow
```

Allowed types: `feat`, `fix`, `refactor`, `chore`, `docs`, `test`, `ci`. Scope is optional but use it for clarity (`feat(runner)`, `fix(menu)`).

Breaking changes: append `!` to the type and add a `BREAKING CHANGE:` footer. Use sparingly — this triggers a major version bump:
```
feat!: change Runner interface to accept Script struct

BREAKING CHANGE: Command(string) replaced by Command(Script)
```

---

## Versioning

Follows Semantic Versioning: `vMAJOR.MINOR.PATCH`

Version const in `cmd/main.go`: `const Version = "0.1.0"`

| Commit type | Version bump |
|---|---|
| `fix:` | patch — `v1.0.0 → v1.0.1` |
| `feat:` | minor — `v1.0.0 → v1.1.0` |
| `feat!:` / `BREAKING CHANGE` | major — `v1.0.0 → v2.0.0` |
| `chore:`, `docs:`, `refactor:` | none |

Release flow: bump versie → `git-cliff` → commit → tag → push → GoReleaser bouwt cross-platform binaries + GitHub Release assets.

---

## Development Workflow

### Werkmethode

Een combinatie van **iterative development**, **MoSCoW prioritering**, en **timeboxing per versie**:

1. **Brainstorm** — vrij nadenken over wat de volgende versie moet bevatten
2. **Versie-bestand** — kort planning document in `docs/roadmap/` met MoSCoW structuur:
   - **Must:** wat er absoluut in moet voor deze release
   - **Nice-to-have:** wat erbij kan als er tijd over is
   - **Out of scope:** bewust uitgesteld naar een volgende versie
3. **Bouwen** — één ding tegelijk afmaken (timeboxing); geen nieuwe scope toevoegen
4. **Commits tussendoor** — conventional commits tijdens het bouwen, niet aan het einde
5. **Release** — zodra alles uit de Must-lijst klaar is: `git-cliff` → tag → GoReleaser

> Geen officiële naam, maar de combinatie: iterative development (rolling focus per versie) + MoSCoW (scope afbakening) + timeboxing (één ding tegelijk).

### Versie-bestand formaat

Zie `docs/roadmap/` voor voorbeelden. Kort, scanbaar, geen essays:

```markdown
# v0.2.0 — <thema>

## Must
- [ ] feature A
- [ ] fix B

## Nice-to-have
- [ ] feature C

## Out of scope
- feature D → v0.3.0
```

Current flow (samengevat): `brainstorm → versie-bestand → bouwen → commits tussendoor → release`

### Modern alternatives (reference)

| Pattern | Flow | Best for |
|---|---|---|
| **GitHub Flow** | issue → branch → PR → merge main → auto-release | Small teams, continuous delivery |
| **Conventional Commits + semantic-release** | commit with type prefix → CI bumps version + changelog + release | Fully automated release pipelines |
| **ADR-first** | brainstorm → ADR → plan → build → commits → release | Projects with significant architecture decisions |
| **Ship/It loop** | feature flag off → build → merge → flag on → observe | Large teams, continuous deployment |

### Automation gap (next step)

This project already has `git-cliff` + GoReleaser + Conventional Commits. One GitHub Actions workflow closes the loop:

```
feat/fix commit → push → CI:
  git-cliff → CHANGELOG.md
  goreleaser → binaries + GitHub Release
  version tag bumped automatically
```

Tools to consider: **`release-please`** (Google) or **`semantic-release`** — both integrate with GitHub Actions and automate the current manual release flow.

---

## Future Recommendations (Out of Scope)

Things worth considering when the project matures — not needed now.

### Branching strategie
- **Trunk-based development** — direct op `main` committen, geen langlevende feature branches; past bij solo/kleine teams met snelle iteraties
- **GitHub Flow** — korte feature branches + PR naar `main`; voegt peer review toe zodra het team groeit
- **GitFlow** — aparte `develop`, `feature/*`, `release/*` branches; overkill voor nu maar relevant als releases strikt gescheiden moeten zijn van development

### CI/CD
- **GitHub Actions** for automated test + lint + release on push to `main`
- `golangci-lint` in CI before merge — catches issues earlier than local runs
- Cross-platform build matrix in CI (`windows`, `darwin`, `linux`) to catch platform regressions

### Release automation
- `release-please` or `semantic-release` to automate version bumping + changelog + GitHub Release from Conventional Commits
- Signed binaries via GoReleaser + `cosign` for supply chain integrity

### Distribution
- Homebrew tap (`homebrew-ceda`) for `brew install ceda-scoop` on Mac/Linux
- Winget / Scoop manifest for Windows distribution without admin rights
- Auto-update mechanism in the TUI (check GitHub Releases API on startup)

### Quality
- `govulncheck` in CI for dependency vulnerability scanning
- Integration tests that run the actual TUI against a mock script registry
- E2E smoke test: install the binary, run it, verify it exits cleanly

### Developer experience
- `devcontainer.json` for reproducible dev environment (Go + tools pre-installed)
- `pre-commit` hooks for `go vet`, `golangci-lint`, and commit message linting
- VS Code / Zed workspace settings committed in `.vscode/` / `.zed/`

---

## MCP Servers

Geconfigureerd in `.mcp.json` (project-scoped, gecommit in repo):

```json
{
  "mcpServers": {
    "memory": {
      "command": "npx",
      "args": ["-y", "@modelcontextprotocol/server-memory"],
      "env": {
        "MEMORY_FILE_PATH": ".claude/memory.json"
      }
    }
  }
}
```

> Memory MCP geeft Claude Code persistente kennis over architectuurbeslissingen,
> bekende bugs, en projectcontext tussen sessies. Gebruik het als:
> "Remember that we use X pattern for Y"

---

## Subagents

Gedefinieerd in `.claude/agents/`. Claude delegeert automatisch op basis van de `description`.

### code-reviewer
`.claude/agents/code-reviewer.md`
```markdown
---
name: code-reviewer
description: Reviewt Go code op correctheid, idioom, en Charm conventies. Gebruik bij PR review of voor je commit.
model: claude-haiku-4-5
tools: Read, Grep, Glob
memory: project
---
Review de aangewezen bestanden op:
- Go idioom en best practices
- Correcte Bubble Tea MVU patronen (geen side effects in View)
- Lip Gloss styles buiten inline
- Logging naar file, niet stdout
- Conventional commit suggestie voor de wijzigingen
```

### release-agent
`.claude/agents/release-agent.md`
```markdown
---
name: release-agent
description: Begeleidt het releaseproces: versie bumpen, changelog genereren, tag aanmaken.
model: claude-sonnet-4-6
tools: Read, Write, Bash
---
Voer het volgende uit in volgorde:
1. Bepaal nieuwe versie op basis van commits sinds laatste tag (semver)
2. Bump versie in cmd/main.go
3. Draai git-cliff om CHANGELOG.md te updaten
4. Commit: `chore: release vX.Y.Z`
5. Tag: `git tag vX.Y.Z`
6. Geef instructies voor `git push origin vX.Y.Z`
```

---

## Skills (Slash Commands)

Gedefinieerd in `.claude/skills/`. Aanroepen met `/skill-naam`.

| Skill | Wanneer gebruiken |
|---|---|
| `/plan-versie` | Begin van een nieuwe iteratie — brainstorm → versie-bestand schrijven |
| `/bouw` | Start van het bouwen — leest versie-bestand, begeleidt met commits |
| `/new-feature` | Scaffold een nieuwe TUI feature (model, Msg types, root koppeling) |
| `/new-adr` | Nieuwe Architecture Decision Record aanmaken |
| `/conventional-commit` | Commit message schrijven tussendoor |
| `/github-release` | Release publiceren na afgeronde versie |
| `/git-workflow` | PR voorbereiden, branches opruimen |

---

## Geheugen & Context Tips

- Vertel Claude aan het begin van een sessie relevante context:
  `"Remember that we use X for Y"` → slaat op in memory MCP
- Bij lange sessies: gebruik `/compact` om context te comprimeren zonder te verliezen
- Subagents draaien in hun eigen context window — gebruik ze voor exploratiewerk
  zodat je main conversation context schoon blijft

---

## Roadmap & JTBD

`docs/ROADMAP.md` — georganiseerd rond **Jobs to be Done**: welke taak probeert de gebruiker bereiken, in welke situatie, en waarom? Features worden gekoppeld aan een job, niet los opgesomd.

Versie-bestanden: `docs/roadmap/00N-vX.Y.Z.md` — MoSCoW per versie (Must / Nice-to-have / Out of scope).

---

## Architecture Decisions

Zie `docs/adr/` voor genomen beslissingen. Belangrijkste:

- `001`: Bubble Tea MVU gekozen boven custom event loop
- `002`: Mouse support via `WithMouseCellMotion()`, hit detection in Model

---
