# CEDA Store (`ceda-scoop`)

TUI desktop application voor het installeren en uitvoeren van CEDA-gecureerde R en Python projecten — zonder adminrechten.

## Install

Download de juiste binary voor jouw platform via [Releases](https://github.com/cedanl/ceda-scoop/releases).

| Platform | Binary |
|---|---|
| macOS (Apple Silicon) | `ceda-scoop_Darwin_arm64.tar.gz` |
| macOS (Intel) | `ceda-scoop_Darwin_x86_64.tar.gz` |
| Windows | `ceda-scoop_Windows_x86_64.zip` |
| Linux | `ceda-scoop_Linux_x86_64.tar.gz` |

## Usage

```bash
./ceda-scoop
```

Navigeren met pijltjestoetsen, `q` om af te sluiten, `?` voor help.

Updates worden automatisch gedetecteerd. Druk `u` om een nieuwe versie te downloaden.

## Development

```bash
go run ./cmd/main.go   # starten
go build ./...         # builden
go test ./...          # testen
```

Zie [CLAUDE.md](CLAUDE.md) voor architectuur en conventies.

## License

MIT — zie [LICENSE](LICENSE).
