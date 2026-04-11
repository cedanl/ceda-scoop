# ceda-scoop
✅ [TEMPLATE] - Centrale repository voor het beheren en installeren van Python & R dependencies van CEDA-producten, zonder dat adminrechten vereist zijn.

## Gebruik

**Windows:** dubbelklik op `ceda-scoop.exe`

**macOS / Linux:** dubbelklik op `ceda-scoop`

## Plan van Aanpak

### Core
- [x] Boilerplate Scoop met installatie van `uv` + `r` + core dependencies (`7zip`, `aria2`, `git`) + gebruikerspaden voor rtools
- [x] Uitbreiden met automatische installatie van `renv` bij R-projecten + openen van RStudio / Positron met R-project
- [x] Uitbreiden met automatische opstart van `uv` bij uv-Streamlit-projecten
- [x] Script refactoreren naar modulaire deelscripts, opgeslagen in `/windows`
- [x] Simpele TUI cmd wrapper
- [x] r-sync vereenvoudigd: configuratie via `.Rprofile` in repo, r-sync triggert alleen `renv::restore()`

### Extra
- [x] (T)UI-wrapper bouwen voor betere gebruikersinterface en separation of concerns van frontend & backend (bijv. Textual, Bubble Tea of Electron, nog te bepalen)
- [x] Automatisch selecteren / downloaden van production-ready repo's
- [ ] Automatische update functionaliteit van CEDA 'Store' (ceda-scoop)
- [ ] Cross-platform Executable (CPE) maken
- [ ] Download knop in README
- [ ] Naamgeving ceda-scoop aanpassen naar CEDA-Store?
- [ ] Versies vastzetten in Scoop
- [ ] Maak deïnstallatie optie voor modules
- [ ] Bovenstaande functionaliteiten uitbreiden naar `/mac` en `/linux` met respectievelijk Homebrew en een nog te bepalen packagemanager
- [ ] Verbeteringen UI / UX
- [ ] FAQ + ADR (Architecture Decision Record)

#### Extra - r
- [ ] Instellen Standaard Positron Preferences
- [ ] Maken van renv (R) repo guidelines ten behoeve van ceda-scoop

#### Extra - uv
- [ ] Maken van uv (python) repo guidelines ten behoeve van ceda-scoop
