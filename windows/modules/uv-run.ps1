# uv-run.ps1
# Zoekt het entry-point en start het uv project op.
# Exporteert: Start-UvProject

function Start-UvProject {
    param(
        [string]$Root = $PSScriptRoot,
        [bool]$IsStreamlit = $false,
        [string]$EntryPoint
    )

    try {
        # ── Zoek entry-point als niet meegegeven ──
        if (-not $EntryPoint) {
            if (Test-Path (Join-Path $Root "src\main.py")) {
                $EntryPoint = "src\main.py"
            } elseif (Test-Path (Join-Path $Root "main.py")) {
                $EntryPoint = "main.py"
            } else {
                $found = Get-ChildItem -Path (Join-Path $Root "src") -Filter "*.py" -ErrorAction SilentlyContinue | Select-Object -First 1
                if ($found) { $EntryPoint = "src\$($found.Name)" }
            }
        }

        if (-not $EntryPoint) {
            return [PSCustomObject]@{ Success = $false; Message = "Geen Python entry-point gevonden in src\ of root" }
        }

        # ── Starten: Streamlit in nieuw venster, gewone uv run blokkeert bewust ──
        Push-Location $Root
        if ($IsStreamlit) {
            Write-Host "Streamlit app starten: $EntryPoint" -ForegroundColor Cyan
            Start-Process powershell -ArgumentList "-NoExit", "-Command", "uv run streamlit run $EntryPoint"
        } else {
            Write-Host "Project starten: $EntryPoint" -ForegroundColor Cyan
            uv run $EntryPoint
        }
        Pop-Location

        return [PSCustomObject]@{ Success = $true; Message = "Project gestart: $EntryPoint (streamlit: $IsStreamlit)" }

    } catch {
        Pop-Location -ErrorAction SilentlyContinue
        return [PSCustomObject]@{ Success = $false; Message = "Onverwachte fout bij opstarten project: $_" }
    }
}
