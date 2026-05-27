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

Gedefinieerd in `.claude/skills/`. Aanroepen met `/skill-naam` of Claude laadt ze automatisch.

### `/new-feature`
`.claude/skills/new-feature/SKILL.md`
```markdown
---
name: new-feature
description: Scaffold een nieuwe TUI feature volgens projectconventies.
invocation: explicit
---
Voer uit in volgorde:
1. Maak `internal/views/<naam>.go` aan met Lip Gloss styles bovenaan als package-level vars
2. Voeg benodigde Msg types toe in `internal/messages.go`
3. Registreer handler in `Model.Update()` switch in `internal/model.go`
4. Voeg view toe in `Model.View()` conditional
5. Conventional commit suggestie: `feat: <beschrijving>`
```

### `/new-adr`
`.claude/skills/new-adr/SKILL.md`
```markdown
---
name: new-adr
description: Maak een nieuwe Architecture Decision Record aan.
invocation: explicit
---
Maak `docs/adr/00N-<titel>.md` aan met:
- Datum
- Status (Proposed / Accepted / Deprecated)
- Context: waarom is deze beslissing nodig?
- Beslissing: wat hebben we gekozen?
- Consequenties: wat zijn de trade-offs?
```

---

## Geheugen & Context Tips

- Vertel Claude aan het begin van een sessie relevante context:
  `"Remember that we use X for Y"` → slaat op in memory MCP
- Bij lange sessies: gebruik `/compact` om context te comprimeren zonder te verliezen
- Subagents draaien in hun eigen context window — gebruik ze voor exploratiewerk
  zodat je main conversation context schoon blijft

---

## Architecture Decisions

Zie `docs/adr/` voor genomen beslissingen. Belangrijkste:

- `001`: Bubble Tea MVU gekozen boven custom event loop
- `002`: Mouse support via `WithMouseCellMotion()`, hit detection in Model

---
