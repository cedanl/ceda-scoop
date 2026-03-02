# ceda-scoop
✅ [TEMPLATE] - Centrale repository voor het beheren en installeren van Python & R dependencies van CEDA-producten, zonder dat adminrechten vereist zijn.

## Gebruik
1. Voer één van de onderstaande scripts uit:

### Installeren (en Terminal UI openen)

```cmd
scoop-install.cmd
```

### Deïnstalleren
```cmd
scoop-uninstall.cmd
```

## Plan van Aanpak

### Core
- [x] Boilerplate Scoop met installatie van `uv` + `r` + core dependencies (`7zip`, `aria2`, `git`) + gebruikerspaden voor rtools
- [ ] Uitbreiden met automatische installatie van `renv` bij R-projecten + openen van RStudio / Positron met R-project
- [ ] Uitbreiden met automatische opstart van `uv` bij uv-Streamlit-projecten
- [ ] Versies vastzetten in Scoop

### Extra
- [ ] Script refactoreren naar modulaire deelscripts, opgeslagen in `/windows`
- [ ] (T)UI-wrapper bouwen voor betere gebruikersinterface en separation of concerns van frontend & backend  (bijv. Textual, Bubble Tea of Electron, nog te bepalen)
- [ ] Bovenstaande functionaliteiten uitbreiden naar `/mac` en `/linux` met respectievelijk Homebrew en een nog te bepalen packagemanager
