# detect-project.ps1
# Detecteert of een repo een R- of uv-project is op basis van lockfiles.
# Exporteert: Get-ProjectType

function Get-ProjectType {
    param(
        [string]$Root = $PSScriptRoot
    )

    $hasR  = Test-Path (Join-Path $Root "renv.lock")
    $hasUv = Test-Path (Join-Path $Root "uv.lock")

    # ── Ambiguiteit: beide gevonden, vraag gebruiker ──
    if ($hasR -and $hasUv) {
        Write-Host "Zowel een R- als uv-project gevonden in deze map." -ForegroundColor Yellow
        Write-Host "Welk projecttype wil je gebruiken?" -ForegroundColor Yellow
        Write-Host "  [1] R (renv.lock)"
        Write-Host "  [2] Python/uv (uv.lock)"
        $keuze = Read-Host "Keuze (1/2)"

        $type = if ($keuze -eq "1") { "r" } elseif ($keuze -eq "2") { "uv" } else {
            Write-Host "Ongeldige keuze, afsluiten." -ForegroundColor Red
            "unknown"
        }

        return [PSCustomObject]@{
            Type    = $type
            Root    = $Root
            Message = "Handmatig gekozen: $type"
        }
    }

    # ── Enkelvoudige detectie ──
    if ($hasR) {
        return [PSCustomObject]@{
            Type    = "r"
            Root    = $Root
            Message = "R-project gedetecteerd (renv.lock gevonden)"
        }
    }

    if ($hasUv) {
        return [PSCustomObject]@{
            Type    = "uv"
            Root    = $Root
            Message = "uv-project gedetecteerd (uv.lock gevonden)"
        }
    }

    # ── Onbekend ──
    return [PSCustomObject]@{
        Type    = "unknown"
        Root    = $Root
        Message = "Geen bekend projecttype gevonden in $Root"
    }
}
