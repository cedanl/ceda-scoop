---
name: new-feature
description: Scaffold een nieuwe TUI feature volgens projectconventies. Maakt model, registreert Msg types, en koppelt aan root.
invocation: explicit
---

Voer uit in volgorde:

1. Vraag naar de naam en beschrijving van de feature (als niet meegegeven)

2. Maak `internal/ui/<naam>/model.go` aan met:
   - Package declaratie
   - Lip Gloss styles als package-level vars bovenaan (nooit inline)
   - `Model` struct met benodigde state
   - `Init() tea.Cmd`
   - `Update(msg tea.Msg) (tea.Model, tea.Cmd)`
   - `View() string`

3. Voeg benodigde Msg types toe in `internal/messages.go`

4. Registreer handler in de `Update()` switch in `internal/ui/root.go`

5. Voeg view toe in de `View()` conditional in `internal/ui/root.go`

6. Controleer: geen `fmt.Println` in de nieuwe code, geen inline styles, geen direct `exec.Command`

7. Conventional commit suggestie: `feat(<naam>): <beschrijving>`
