---
name: plan-versie
description: Brainstorm en schrijf het versie-bestand voor de volgende release. Gebruik aan het begin van een nieuwe iteratie.
invocation: explicit
---

Voer uit in volgorde:

1. **Context ophalen**
   - Lees `docs/ROADMAP.md` voor JTBD-context: welke user jobs staan centraal?
   - Draai `git tag --sort=-version:refname | head -5` voor recente versies
   - Draai `git log $(git describe --tags --abbrev=0)..HEAD --oneline` voor commits sinds laatste tag
   - Lees bestaande versie-bestanden in `docs/roadmap/` voor patroon en stijl

2. **Brainstorm met de gebruiker**
   Stel gerichte vragen — één tegelijk, wacht op antwoord:
   - Wat is het thema of doel van deze versie? (één zin)
   - Welke user job uit de roadmap staat centraal?
   - Wat zijn de absolute must-haves (geen release zonder)?
   - Wat zou mooi zijn maar kan ook later?
   - Wat schuiven we bewust door naar de volgende versie?

3. **Versie bepalen**
   - Bepaal volgende semver op basis van commit types: `feat` → minor, `fix` → patch, `feat!` → major
   - Controleer of er al een versie-bestand bestaat voor die versie

4. **Versie-bestand schrijven**
   Maak `docs/roadmap/00N-vX.Y.Z.md` aan (N = volgnummer):

   ```markdown
   # vX.Y.Z — <thema>

   **JTBD:** When <situatie>, I want to <motivatie>, so I can <uitkomst>

   ## Must
   - [ ] <item>

   ## Nice-to-have
   - [ ] <item>

   ## Out of scope → vX.Y+1.Z
   - <item>
   ```

5. **Afsluiten**
   - Conventional commit suggestie: `docs: add vX.Y.Z roadmap`
   - Zeg: "Start bouwen met `/bouw`"
