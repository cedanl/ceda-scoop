---
name: new-adr
description: Maak een nieuwe Architecture Decision Record aan in docs/adr/.
invocation: explicit
---

Voer uit in volgorde:

1. Bepaal volgnummer: tel bestaande bestanden in `docs/adr/` en gebruik N+1
2. Vraag naar de titel als niet meegegeven
3. Maak `docs/adr/00N-<kebab-titel>.md` aan:

```markdown
# ADR-00N: <Titel>

**Datum:** YYYY-MM-DD
**Status:** Proposed

## Context
Waarom is deze beslissing nodig? Wat is de situatie of het probleem?

## Beslissing
Wat hebben we gekozen?

## Consequenties
**Voordelen:**
- 

**Nadelen / trade-offs:**
- 

**Gerelateerd:** ADR-00X
```

4. Update de lijst in `CLAUDE.md` onder "Architecture Decisions"
5. Conventional commit suggestie: `docs(adr): add ADR-00N <titel>`
