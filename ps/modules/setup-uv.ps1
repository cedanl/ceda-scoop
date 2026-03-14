# setup-uv.ps1
# Installeert uv (via Scoop) en start een uv-project op.
# Exporteert: Install-UV, Start-UvProject

# ── uv installeren ──
function Install-UV {
    try {
        if (Get-Command uv -ErrorAction SilentlyContinue) {
            return [PSCustomObject]@{ Success = $true; Message = "uv al aanwezig, overgeslagen" }
        }

        $output = scoop install uv 2>&1
        if ($LASTEXITCODE -ne 0) {
            return [PSCustomObject]@{ Success = $false; Message = "Fout bij installatie van uv: $output" }
        }

        return [PSCustomObject]@{ Success = $true; Message = "uv geinstalleerd via Scoop" }
    } catch {
        return [PSCustomObject]@{ Success = $false; Message = "Onverwachte fout bij uv installatie: $_" }
    }
}

# ── Streamlit check via pyproject.toml ──
function Test-IsStreamlitProject {
    param([string]$Root)

    $pyproject = Join-Path $Root "pyproject.toml"
    if (-not (Test-Path $pyproject)) { return $false }

    $content = Get-Content $pyproject -Raw
    return $content -match "streamlit"
}

# ── Entry-point zoeken en project starten ──
function Start-UvProject {
    param([string]$Root = $PSScriptRoot)

    try {
        # Zoek entry-point in volgorde
        $entryPoint = $null

        if (Test-Path (Join-Path $Root "src\main.py")) {
            $entryPoint = "src\main.py"
        } elseif (Test-Path (Join-Path $Root "main.py")) {
            $entryPoint = "main.py"
        } else {
            $found = Get-ChildItem -Path (Join-Path $Root "src") -Filter "*.py" -ErrorAction SilentlyContinue | Select-Object -First 1
            if ($found) { $entryPoint = "src\$($found.Name)" }
        }

        if (-not $entryPoint) {
            return [PSCustomObject]@{ Success = $false; Message = "Geen Python entry-point gevonden in src\ of root" }
        }

        # uv sync vanuit de project root
        Push-Location $Root
        $syncOutput = uv sync 2>&1
        $syncExit = $LASTEXITCODE
        Pop-Location

        if ($syncExit -ne 0) {
            return [PSCustomObject]@{ Success = $false; Message = "uv sync mislukt: $syncOutput" }
        }

        # Streamlit of gewone uv run
        $isStreamlit = Test-IsStreamlitProject -Root $Root
        Push-Location $Root
        if ($isStreamlit) {
            uv run streamlit run $entryPoint
        } else {
            uv run $entryPoint
        }
        Pop-Location

        return [PSCustomObject]@{ Success = $true; Message = "Project gestart: $entryPoint (streamlit: $isStreamlit)" }

    } catch {
        return [PSCustomObject]@{ Success = $false; Message = "Onverwachte fout bij opstarten project: $_" }
    }
}
