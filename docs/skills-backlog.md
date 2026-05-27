# Skills Backlog

Potentiële Claude Code skills (`.claude/skills/`) om later toe te voegen.
Gebaseerd op practices van goede developers — gegroepeerd op wat het oplevert.

---

## Design-first

### `/spike`
Tijdgebonden exploratie voordat je richting kiest.
- Schrijf een schets (max 1 pagina): wat onderzoek je, wat zijn de opties?
- Stel een tijdlimiet (bv. 2 uur)
- Beslissing: doorgaan of weggooien
- Populariseerd door: XP / Extreme Programming

### `/rfc`
Grote feature eerst uitschrijven voor je code aanraakt.
- Beschrijf het probleem, de oplossing, alternatieven, trade-offs
- Dwingt je na te denken over API-design voor implementatie
- Populariseerd door: React team, Rust team

### `/api-design`
Definieer de interface / types voordat je implementeert.
- Types zijn de spec — code volgt
- Populariseerd door: Matt Pocock (TypeScript), Rich Harris

### `/changelog-first`
Schrijf de changelog-entry *voordat* je codeert.
- Dwingt je te verwoorden wat de gebruiker merkt van de wijziging
- Gebaseerd op: Keep a Changelog-filosofie

---

## Testing

### `/tdd`
Red/green/refactor cyclus begeleiden bij Go code.
1. Schrijf een falende test
2. Implementeer minimale code om de test te laten slagen
3. Refactor zonder de test te breken
- Populariseerd door: Kent C. Dodds, Robert C. Martin

### `/property-test`
Property-based testing: geef random input, beschrijf eigenschappen die altijd waar zijn.
- Voor pure functies met veel inputvarianten
- Go: `testing/quick` of `pgregory.net/rapid`
- Populariseerd door: John Hughes (QuickCheck / Haskell)

---

## Iteratie & scope

### `/shape-up`
Scope afbakenen met grove schetsen, niet pixel-perfecte specs.
- Fat marker sketch: wat is de kern, wat is buiten scope?
- Appetit bepalen: hoeveel tijd is dit waard?
- Populariseerd door: Basecamp / DHH, Ryan Singer

---

## Code kwaliteit

### `/error-design`
Foutmeldingen reviewen op duidelijkheid.
- Wat ging er fout? (concreet, niet technisch jargon)
- Waarom? (oorzaak in begrijpelijke taal)
- Wat nu? (vervolgstap of contactpersoon)
- Populariseerd door: Matt Pocock (TypeScript errors), Elm team

### `/boy-scout`
Laat de code schoner achter dan je hem aantrof.
- Kleine refactors tussendoor, niet bewaren voor later
- Populariseerd door: Robert C. Martin

---

## Diagnose & onderhoud

### `/postmortem`
Na een bug of incident: wat ging fout, waarom, wat voorkomt herhaling?
- Timeline van wat er gebeurde
- Root cause (geen schuldvraag)
- Actiepunten met eigenaar
- Populariseerd door: Google SRE, Etsy

### `/diagnose`
`tool --diagnose` commando dat omgeving, versies, en config valideert.
- OS, Go versie, script paden, permissions
- Geeft ontwikkelaar info om probleem te reproduceren
- Populariseerd door: Mitchell Hashimoto (HashiCorp tooling)

---

## Referenties

- Matt Pocock — [Total TypeScript](https://www.totaltypescript.com/)
- Kent C. Dodds — [Testing Trophy](https://kentcdodds.com/blog/the-testing-trophy-and-testing-classifications)
- Ryan Singer — [Shape Up](https://basecamp.com/shapeup)
- Dave Cheney — [Practical Go](https://dave.cheney.net/practical-go)
- Google SRE — [Site Reliability Engineering](https://sre.google/sre-book/postmortem-culture/)
- Mitchell Hashimoto — [Advanced Testing in Go](https://mitchellh.com/writing/advanced-testing-in-go)
