---
name: bouw
description: Begeleidt het bouwen van de huidige versie. Leest het versie-bestand, toont Must-checklist, stelt commits voor tussendoor, en triggert release als alles klaar is.
invocation: explicit
---

Voer uit in volgorde:

1. **Huidig versie-bestand ophalen**
   - Zoek het meest recente bestand in `docs/roadmap/` (hoogste volgnummer)
   - Lees het en toon de Must-lijst als genummerde checklist
   - Toon ook Nice-to-have items apart, maar markeer ze als optioneel

2. **Checklist presenteren**
   Voorbeeld output:
   ```
   📋 v0.2.0 — Script runner output

   Must (verplicht voor release):
   [ ] 1. Voeg scrollable output view toe aan runner screen
   [ ] 2. Toon exit code na script voltooiing
   [ ] 3. Foutmelding tonen als script niet gevonden

   Nice-to-have (als er tijd is):
   [ ] A. Syntax highlighting in output
   ```

3. **Per item begeleiden**
   - Vraag: "Welk Must-item pak je als eerste op?"
   - Nadat de gebruiker aangeeft dat een item klaar is:
     - Vink het af in de mentale checklist
     - Stel een conventional commit voor: `feat(<scope>): <beschrijving>`
     - Vraag of de gebruiker wil committen voor verder gaan
   - Herhaal voor elk item

4. **Klaar-check**
   Als alle Must-items afgevinkt:
   - Vraag: "Wil je Nice-to-have items oppakken of direct releasen?"
   - Bij release: geef instructies:
     ```
     1. git-cliff -o CHANGELOG.md
     2. git add CHANGELOG.md
     3. git commit -m "chore: release vX.Y.Z"
     4. git tag vX.Y.Z
     5. git push origin main --tags
     6. goreleaser release --clean
     ```
   - Of gebruik de `/github-release` skill

5. **Tussendoor**
   - Als de gebruiker vastloopt: help debuggen zonder de checklist te vergeten
   - Houd altijd bij hoeveel Must-items nog open staan
