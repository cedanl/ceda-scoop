# CEDA Store — Roadmap

Georganiseerd rond **Jobs to be Done (JTBD)**: wat probeert de gebruiker bereiken, in welke situatie, en waarom?

Format: *When [situatie], I want to [actie], so I can [uitkomst].*

---

## Jobs to be Done

### Job 1: Software installeren zonder IT
> When I need to run an R or Python project from CEDA on my work laptop,
> I want to install it without needing admin rights or IT support,
> so I can start working immediately without bureaucratic delay.

**Features die dit job vervullen:**
- [ ] Script discovery per platform (Windows/Mac/Linux)
- [ ] Uitvoeren van scripts zonder admin rechten
- [ ] Duidelijke foutmelding als iets mislukt

---

### Job 2: Weten wat ik installeer
> When I select a project from the list,
> I want to see a description of what it does and what it installs,
> so I can make an informed choice before running anything.

**Features die dit job vervullen:**
- [ ] Beschrijving tonen in de script-list
- [ ] Detail view per script (vereisten, versie, auteur)

---

### Job 3: Voortgang volgen
> When a script is running,
> I want to see what's happening in real time,
> so I know whether to wait or whether something went wrong.

**Features die dit job vervullen:**
- [ ] Live output in runner screen
- [ ] Scrollable output view
- [ ] Exit code tonen na afloop
- [ ] Spinner tijdens laden

---

### Job 4: Terugkeren na fout
> When a script fails,
> I want to understand what went wrong and be able to try again,
> so I don't have to restart the whole tool.

**Features die dit job vervullen:**
- [ ] Error screen met duidelijke melding
- [ ] Optie om terug te gaan naar het menu
- [ ] Log opslaan naar `debug.log`

---

## Versies

| Versie | Thema | Status |
|---|---|---|
| v0.1.0 | Project scaffolding | ✅ Klaar |
| v0.2.0 | TBD | 🔲 Gepland |

Versie-bestanden: zie `docs/roadmap/`

---

## Out of scope (voorlopig)

- Auto-update mechanisme
- Webinterface
- Script marketplace / community scripts
- Authenticatie / gebruikersbeheer
