# HANDOFF — ceda-store

> Start een nieuwe sessie met: `"Lees HANDOFF.md en ga verder"`

---

## Goal

`ceda-store` (binary: `ceda-scoop`) bouwen: een Go TUI desktop app voor het installeren en uitvoeren van CEDA-gecureerde R en Python projecten zonder adminrechten. Iteratief bouwen via versie-bestanden in `docs/roadmap/`.

---

## Current Progress

### Releases uitgebracht

| Versie | Inhoud | Status |
|---|---|---|
| v0.1.0 | Project scaffolding (directory structuur, Taskfile, GoReleaser) | ✅ |
| v0.2.0 | Foundation: styles, CEDA Store header, resize, graceful exit, GoReleaser CI | ✅ |
| v0.2.1 | Fix: GoReleaser repo naam (`ceda-scoop` → `ceda-store`) | ✅ |
| v0.2.2 | Fix: auto-update API URL, asset filename, status bar scrollbar bug | ✅ |

### Wat er nu draait (v0.2.2)

- Volledig werkende TUI op macOS/Windows/Linux
- `internal/styles/styles.go` — Lip Gloss styles met Npuls design tokens
- `internal/ui/root.go` — root model met header, content zone, status bar, help overlay, resize handling
- `internal/updater/updater.go` — async GitHub API check + zelf-update (Unix: rename+exec, Windows: `*-update.exe`)
- `internal/messages.go` — alle Msg types op één plek
- GitHub Actions CI → GoReleaser → automatische cross-platform binaries op tag push
- Branch: `dev-Sandbox` (nog niet gemerged naar `main`)

### Wat er nog NIET is

- Geen menu screen (script list) — `internal/ui/menu/model.go` is leeg scaffold
- Geen runner screen — `internal/ui/runner/model.go` is leeg scaffold
- Geen scripts in `scripts/` — map bestaat maar geen `.sh`/`.ps1` files
- Geen splash screen (dat is v0.3.0)

---

## What Worked

- **TEA architecture**: Init/Update/View strikt scheiden werkt goed; side effects als `tea.Cmd` teruggeven
- **Async update check**: `checkUpdateCmd` returnt een `UpdateCheckDoneMsg` — blokkeert UI niet
- **Status bar gap**: `strings.Repeat(" ", gap)` in plaats van `styles.Subtle.Render(fmt.Sprintf("%*s", gap, ""))` — ANSI escape codes om spaties breken Lip Gloss breedte-berekening
- **GoReleaser asset names**: `runtime.GOOS` + `runtime.GOARCH` direct gebruiken (lowercase, geen mapping) want GoReleaser produceert bijv. `ceda-scoop_darwin_arm64.tar.gz`
- **SSH remote**: `git remote set-url origin git@github.com:cedanl/ceda-store.git` — HTTPS faalt omdat gh CLI geconfigureerd is voor SSH
- **Skill workflow**: `/plan-versie` → bouwen → `/bouw` checklist → commits → `/github-release`

## What Didn't Work

- **`apiURL` naar `cedanl/ceda-scoop`**: Repo heet `ceda-store` — update check gaf altijd "update onbekend"
- **Asset filename mapping** (`Darwin` → `darwin`, `x86_64` → `amd64`): Niet nodig, GoReleaser gebruikt al `runtime.GOOS`/`GOARCH` lowercase
- **ANSI-wrapped spaces in status bar**: `lipgloss.Render()` om spaties heen → Lip Gloss telt breedte verkeerd → scrollbar zichtbaar bij opstarten
- **`release.github.name: ceda-scoop`** in `.goreleaser.yml`: Moet `ceda-store` zijn (reponaam, niet binary naam)
- **`tea.WithAltScreen()` + min-size enforcement**: Kan terminal niet echt blokkeren op resize — huidige workaround toont foutmelding maar forceert niet

---

## Next Steps

### Direct op te pakken: v0.3.0 Splash Screen

Plan staat in `docs/roadmap/002-v0.3.0.md`. Start met `/bouw`.

**Must-haves:**
1. `docs/adr/003-splash-libraries.md` — ADR voor animatie en layout keuze
2. Splash screen als eerste scherm (`internal/ui/splash/model.go`)
3. Animatie via `charmbracelet/harmonica` — spring physics voor reveal/pulse effect
4. Muisklik support via `tea.WithMouseCellMotion()` — al aanwezig in `main.go` config, alleen activeren
5. Vaste minimale grootte (bouwt voort op v0.2.0 resize guard in `root.go`)
6. Meerdere splash varianten met pijltjestoetsen + muisklik navigatie

**Nice-to-haves:**
- Auto-doorgaan naar menu na X seconden
- Concentric rings border (Npuls vormentaal)
- Chevron navigatie-indicator

### Daarna: v0.4.0+ (nog niet gepland)

- Menu screen met script lijst (`internal/ui/menu/model.go` invullen)
- Runner screen voor script uitvoering
- Echte scripts in `scripts/mac/`, `scripts/windows/`, `scripts/linux/`

---

## Key Files

| Bestand | Doel |
|---|---|
| `CLAUDE.md` | Architectuur, conventies, tooling — lees dit altijd |
| `cmd/main.go` | Entry point, `Version` const |
| `internal/ui/root.go` | Root model — hier routes en state machine |
| `internal/styles/styles.go` | Alle Lip Gloss styles (Npuls tokens) |
| `internal/updater/updater.go` | GitHub release check + download |
| `internal/messages.go` | Alle Msg types |
| `docs/ROADMAP.md` | JTBD-structuur, versie tabel |
| `docs/roadmap/002-v0.3.0.md` | v0.3.0 plan (splash screen) |
| `.goreleaser.yml` | Cross-platform build config |
| `.github/workflows/release.yml` | CI: tag push → GoReleaser |

---

## Branch Status

- **Actieve branch**: `dev-Sandbox`
- **Main**: nog niet gemerged — overweeg PR na v0.3.0 of eerder
- **Laatste tag**: `v0.2.2` op `main`-equivalente commits

---

## Design Tokens (Npuls huisstijl)

```go
ColorOranje  = "#DD784B"  // primaire accent, header
ColorBlauw   = "#3D68EC"  // selected state, links
ColorGroen   = "#00AF81"  // success
ColorGeel    = "#F4D74B"  // warning
ColorDark    = "#111827"  // achtergrond
ColorWit     = "#FFFFFF"  // body tekst
ColorGray500 = "#6B7280"  // subtle / muted
ColorGray700 = "#374151"  // borders
```
