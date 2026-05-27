# CEDA Store — Roadmap

Georganiseerd rond **Jobs to be Done (JTBD)**: wat probeert de gebruiker bereiken, in welke situatie, en waarom?

Format: *When [situatie], I want to [actie], so I can [uitkomst].*

Status: ✅ Voldaan | 🚧 Gedeeltelijk | 🔲 Open

Versies dragen bij aan jobs — elke versie beschrijft wat er gebouwd wordt en welke job(s) daardoor dichterbij komen.

---

## Eindgebruiker

### U1 — Installeren zonder IT `🔲 Open`
> When I need to run an R or Python project from CEDA on my work laptop,
> I want to install it without needing admin rights or IT support,
> so I can start working immediately without bureaucratic delay.

- [ ] Script discovery per platform (Windows / Mac / Linux)
- [ ] Scripts uitvoeren zonder admin rechten
- [ ] Duidelijke foutmelding als iets mislukt

---

### U2 — Weten wat ik installeer `🔲 Open`
> When I select a project from the list,
> I want to see a description of what it does and what it installs,
> so I can make an informed choice before running anything.

- [ ] Beschrijving zichtbaar in de script-list
- [ ] Detail view per script (vereisten, versie, auteur)

---

### U3 — Voortgang volgen `🔲 Open`
> When a script is running,
> I want to see what's happening in real time,
> so I know whether to wait or whether something went wrong.

- [ ] Live output in runner screen
- [ ] Scrollable output view
- [ ] Spinner tijdens laden
- [ ] Exit code tonen na afloop

---

### U4 — Herstellen na fout `🔲 Open`
> When a script fails,
> I want to understand what went wrong and be able to try again,
> so I don't have to restart the whole tool.

- [ ] Error screen met duidelijke melding
- [ ] Optie om terug te gaan naar het menu
- [ ] Log opslaan naar `debug.log`

---

### U5 — Overzicht van wat geïnstalleerd is `🔲 Open`
> When I open the tool,
> I want to see which projects I already have and which I don't,
> so I don't install something twice.

- [ ] Status per script (geïnstalleerd / niet / verouderd)
- [ ] Visuele indicator in de lijst

---

### U6 — Updaten wat verouderd is `🔲 Open`
> When a newer version of a project is available,
> I want to update it without reinstalling everything,
> so I stay up to date without effort.

- [ ] Versievergelijking bij startup
- [ ] Update flow vanuit het menu

---

### U7 — Verwijderen wat ik niet meer nodig heb `🔲 Open`
> When I no longer use a project,
> I want to cleanly uninstall it,
> so my laptop stays clean.

- [ ] Uninstall flow per script
- [ ] Bevestigingsstap voor verwijderen

---

### U8 — Weten of het werkt `🔲 Open`
> When an installation completes,
> I want to see a quick check that everything is in order,
> so I don't have to open a terminal myself to verify.

- [ ] Post-install healthcheck per script
- [ ] Resultaat zichtbaar in result screen

---

### U9 — Tool werkt gewoon op mijn OS `🔲 Open`
> When I download the tool on Windows, Mac, or Linux,
> I want to start it without an admin prompt or security warning,
> so I don't have to contact IT just to open the tool.

- [ ] Gesigneerde binary (GoReleaser + cosign)
- [ ] Windows: geen UAC prompt bij normaal gebruik
- [ ] Mac: notarisatie zodat Gatekeeper niet blokkeert
- [ ] Linux: uitvoerbaar zonder sudo

---

### U10 — Begrijpelijke foutmeldingen `🔲 Open`
> When something goes wrong,
> I want to see a clear message in plain language,
> so I know what to do next (or who to ask).

- [ ] Foutmeldingen in begrijpelijke taal (niet raw stderr)
- [ ] Suggestie van vervolgstap in de foutmelding
- [ ] Onderscheid tussen gebruikersfout en systeemfout

---

### U11 — Prettige interface `🔲 Open`
> When I use the tool,
> I want to see a clean and attractive TUI,
> so I don't feel like I'm using a developer tool.

- [ ] Consistente stijl via Lip Gloss styles
- [ ] Toetsenbord- én muisondersteuning
- [ ] Responsive layout (past zich aan terminalgrootte aan)

---

## Development team

### D1 — Diagnose bij een gebruiker `🔲 Open`
> When a user reports a problem,
> I want to quickly see their environment and logs through the tool itself,
> so I can diagnose without manually accessing their laptop.

- [ ] Diagnose screen of commando (omgeving, OS, Go versie, script pad)
- [ ] `debug.log` eenvoudig te delen vanuit de tool

---

### D2 — Nieuw script toevoegen `🔲 Open`
> When a new CEDA project is ready,
> I want to register it in the store without modifying the TUI,
> so it's immediately available to users.

- [ ] Script discovery via bestandsconventie (geen code aanpassen)
- [ ] Metadata in script zelf (naam, beschrijving, versie)

---

### D3 — Script testen op alle platforms `🔲 Open`
> When I update a script,
> I want to test it on Windows, Mac, and Linux,
> so I don't discover breakage in production.

- [ ] Cross-platform test matrix in CI
- [ ] Lokaal testen mogelijk via runner interface

---

### D4 — Inzicht in gebruik `🔲 Open`
> When something goes wrong in the field,
> I want to be able to see logs or telemetry,
> so I can reproduce and fix the issue.

- [ ] Gestructureerde logging naar `debug.log`
- [ ] Optioneel: telemetrie (opt-in, privacybewust)

---

### D5 — Tool zelf updaten bij gebruikers `🔲 Open`
> When a new version of the tool is released,
> I want users to be able to install it easily,
> so I don't have to visit every laptop.

- [ ] Auto-update check bij startup (GitHub Releases API)
- [ ] Of: distribueren via Homebrew / Winget / Scoop

---

## Maintainer

### M1 — Codebase begrijpen `🚧 Gedeeltelijk`
> When I want to modify the tool,
> I want to quickly understand the architecture,
> so I don't spend hours navigating the code.

- [x] CLAUDE.md met architectuur en conventies
- [x] Folderstructuur gedocumenteerd
- [ ] ADRs voor belangrijke beslissingen (in `docs/adr/`)
- [ ] Inline documentatie op publieke interfaces

---

### M2 — Snel een release maken `🚧 Gedeeltelijk`
> When a version is done,
> I want to build, tag, and publish with a single command,
> so I don't forget manual steps.

- [x] GoReleaser geconfigureerd
- [x] git-cliff voor changelog
- [x] Conventional Commits workflow
- [ ] GitHub Actions workflow die alles automatiseert

---

### M3 — Zeker weten dat niets breekt `🔲 Open`
> When I make a change,
> I want tests and linting to run automatically,
> so I can ship with confidence.

- [ ] CI pipeline met `go test`, `go vet`, `golangci-lint`
- [ ] Pre-commit hooks voor lokale feedback

---

### M4 — Cross-platform builds valideren `🔲 Open`
> When I release,
> I want binaries for Windows, Mac, and Linux to be built and tested automatically,
> so I don't have to build per platform manually.

- [ ] GoReleaser cross-compile matrix in CI
- [ ] Smoke test per platform na build

---

## Versies

| Versie | Thema | Bijdrage aan jobs | Status |
|---|---|---|---|
| v0.1.0 | Project scaffolding | M1, M2 | ✅ Klaar |
| v0.2.0 | TBD | TBD | 🔲 Gepland |

*Versie-bestanden met details: `docs/roadmap/`*

---

## Out of scope (voorlopig)

- Webinterface
- Script marketplace / community scripts
- Authenticatie / gebruikersbeheer
- Telemetrie (D4) — privacy-overwegingen eerst uitwerken
