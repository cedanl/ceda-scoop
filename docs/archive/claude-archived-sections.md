# CLAUDE.md — Gearchiveerde secties

> Inhoud die verwijderd is uit CLAUDE.md omdat het te lang werd of geen directe instructiewaarde had voor Claude. Bewaard hier als referentie.

---

## Development Workflow — Modern alternatives (reference)

Verwijderd uit: `## Development Workflow > Modern alternatives (reference)`

| Pattern | Flow | Best for |
|---|---|---|
| **GitHub Flow** | issue → branch → PR → merge main → auto-release | Small teams, continuous delivery |
| **Conventional Commits + semantic-release** | commit with type prefix → CI bumps version + changelog + release | Fully automated release pipelines |
| **ADR-first** | brainstorm → ADR → plan → build → commits → release | Projects with significant architecture decisions |
| **Ship/It loop** | feature flag off → build → merge → flag on → observe | Large teams, continuous deployment |

---

## Development Workflow — Automation gap (next step)

Verwijderd uit: `## Development Workflow > Automation gap (next step)`

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

Verwijderd uit: `## Future Recommendations (Out of Scope)`

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
- Auto-update mechanism in the TUI (check GitHub Releases API on startup) ← **gedaan in v0.2.0**

### Quality
- `govulncheck` in CI for dependency vulnerability scanning
- Integration tests that run the actual TUI against a mock script registry
- E2E smoke test: install the binary, run it, verify it exits cleanly

### Developer experience
- `devcontainer.json` for reproducible dev environment (Go + tools pre-installed)
- `pre-commit` hooks for `go vet`, `golangci-lint`, and commit message linting
- VS Code / Zed workspace settings committed in `.vscode/` / `.zed/`

---

## Subagents — volledige frontmatter (v0.2.x)

Verwijderd uit: `## Subagents`

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
