# ceda-scoop
✅ [TEMPLATE] - Centrale repository voor het beheren en installeren van Python & R dependencies van CEDA-producten, zonder dat adminrechten vereist zijn.

## Gebruik

Kopieer de inhoud van de `/windows` map naar de **root van je CEDA repo**:

```
jouw-repo/
├── CEDA-Start.cmd        ← gekopieerd uit /windows
├── CEDA-Remove.cmd       ← gekopieerd uit /windows
└── modules/              ← gekopieerd uit /windows
```

Voer daarna vanuit de root van je repo uit:

```cmd
CEDA-Start.cmd
```

## Plan van Aanpak

### Core
- [x] Boilerplate Scoop met installatie van `uv` + `r` + core dependencies (`7zip`, `aria2`, `git`) + gebruikerspaden voor rtools
- [x] Uitbreiden met automatische installatie van `renv` bij R-projecten + openen van RStudio / Positron met R-project
- [x] Uitbreiden met automatische opstart van `uv` bij uv-Streamlit-projecten
- [x] Script refactoreren naar modulaire deelscripts, opgeslagen in `/windows`
- [x] Simpele TUI cmd wrapper
- [x] r-sync vereenvoudigd: configuratie via `.Rprofile` in repo, r-sync triggert alleen `renv::restore()`
- [ ] Maak deïnstallatie optie voor modules

### Extra
- [ ] (T)UI-wrapper bouwen voor betere gebruikersinterface en separation of concerns van frontend & backend (bijv. Textual, Bubble Tea of Electron, nog te bepalen)
- [ ] Automatisch selecteren / downloaden van production-ready repo's
- [ ] Versies vastzetten in Scoop
- [ ] Bovenstaande functionaliteiten uitbreiden naar `/mac` en `/linux` met respectievelijk Homebrew en een nog te bepalen packagemanager
- [ ] FAQ + ADR (Architecture Decision Record)

#### Extra - r
- [ ] Instellen Standaard Positron Preferences
- [ ] Maken van renv (R) repo guidelines ten behoeve van ceda-scoop

#### Extra - uv
- [ ] Maken van uv (python) repo guidelines ten behoeve van ceda-scoop
