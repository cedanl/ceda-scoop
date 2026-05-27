# CEDA Store (`ceda-scoop`)

> Geef dit bestand mee aan elke nieuwe Claude Code sessie. Het laadt automatisch.

TUI desktop application for installing and running CEDA-curated R and Python projects without requiring admin rights.

**Status:** In ontwikkeling
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

- Logic lives in the layer that owns it: UI in `internal/ui/`, OS execution in `internal/runner/`, script metadata in `internal/scripts/`
- No screen knows about another screen — only the root routes between them
- No platform detection outside `internal/runner/`
- Side effects (running scripts) go through `tea.ExecProcess`, not raw `exec.Command`
- All `Msg` types in one place: `internal/messages.go`
- All state in `Model` struct — no global state outside Model
- Logging always to `debug.log` via `log.SetOutput(f)` — never `fmt.Println` to stdout
- Script directories use human-readable names: `scripts/mac/`, not `scripts/darwin/` — `runtime.GOOS` returns `"darwin"` but directory names use `mac` for clarity

---

## Folder structure

```
cmd/
  main.go              # entry point — wires runner, registry, root model; Version const here

internal/
  messages.go          # alle Msg types op één plek

  styles/
    styles.go          # alle Lip Gloss styles — Npuls design tokens; nooit inline

  ui/
    root.go            # root model, routes tussen screens, update state machine
    menu/
      model.go         # script list and selection (scaffold)
    runner/
      model.go         # progress and output display (scaffold)
    result/
      model.go         # success / error state (scaffold)

  runner/
    runner.go          # Runner interface + NewRunner() factory + all OS implementations
                       # UI never touches runtime.GOOS — runner resolves it at startup (scaffold)

  scripts/
    scripts.go         # Script struct (name, description, path) + LoadScripts()
                       # walks scripts/<platform>/, returns []Script (scaffold)

  updater/
    updater.go         # GitHub release check + self-update (Unix: rename+exec, Windows: *-update.exe)

scripts/
  windows/             # .ps1 files
  mac/                 # .sh files  ← NOT darwin/ — runtime.GOOS returns "darwin" but dirs use "mac"
  linux/               # .sh files

docs/
  ROADMAP.md           # JTBD-structuur: 18 jobs, 3 personas (eindgebruiker, dev team, maintainer)
  adr/                 # architecture decision records (zie docs/adr/ voor alle ADRs)
  roadmap/             # versie planningen (MoSCoW per versie)
    001-v0.2.0.md
    002-v0.3.0.md
  skills-backlog.md    # developer skill ideeën en backlog
  archive/             # verwijderde CLAUDE.md secties — bewaard als referentie

.claude/
  agents/              # subagents
  skills/              # skills / slash commands

CLAUDE.md
HANDOFF.md             # sessie-overdracht voor volgende Claude sessie
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
- ❌ Geen inline Lip Gloss styles — altijd via `internal/styles/styles.go`
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

**Na elke code change:** altijd `go build ./cmd/main.go && go vet ./...` draaien voor commit — vangt compile errors vroeg.

**Build tool:** Taskfile is canonical (`task <naam>`). Nooit een Makefile aanmaken of suggereren.

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

Version const in `cmd/main.go` — check huidige waarde met `git describe --tags --abbrev=0`.

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

> Alternatieven en toekomstplannen: zie `docs/archive/claude-archived-sections.md`

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

| Agent | Wanneer |
|---|---|
| `code-reviewer` | PR review of voor commit — checkt Go idioom, MVU patronen, Lip Gloss stijl |
| `release-agent` | Versie bumpen, changelog genereren, tag aanmaken |

> Volledige frontmatter: zie `docs/archive/claude-archived-sections.md`

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
| `/context-mode:grill-with-docs` | Tussen plan en bouwen — toetst plan aan ADRs en domeinmodel, update docs inline |
| `/dx:handoff` | Einde van sessie — schrijft handoff doc zodat volgende sessie direct doorgaat |
| `/dx:review-claudemd` | Analyseer sessies en verbeter CLAUDE.md op basis van wat werkte en wat niet |

---

## Geheugen & Context Tips

- Vertel Claude aan het begin van een sessie relevante context:
  `"Remember that we use X for Y"` → slaat op in memory MCP
- Bij lange sessies: gebruik `/compact` om context te comprimeren zonder te verliezen
- Subagents draaien in hun eigen context window — gebruik ze voor exploratiewerk
  zodat je main conversation context schoon blijft

---

## Roadmap & JTBD

`docs/ROADMAP.md` — georganiseerd rond **Jobs to be Done**: 18 jobs verdeeld over 3 persona's (eindgebruiker U1–U11, development team D1–D5, maintainer M1–M4). Features worden gekoppeld aan een job, niet los opgesomd.

Versie-bestanden: `docs/roadmap/00N-vX.Y.Z.md` — MoSCoW per versie (Must / Nice-to-have / Out of scope). Aanmaken met `/plan-versie`.

---

## Architecture Decisions

Zie `docs/adr/` voor alle genomen beslissingen. Aanmaken met `/new-adr`.

---
