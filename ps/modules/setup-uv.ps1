# setup-uv.ps1
# Installeert uv (via Scoop) en start een uv-project op.
# Exporteert: Install-UV, Set-UvCacheDir, Start-UvProject

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

# ── Zorg dat [tool.uv] cache-dir correct aanwezig is in pyproject.toml ──
function Set-UvCacheDir {
    param([string]$Root)

    $pyproject = Join-Path $Root "pyproject.toml"

    if (-not (Test-Path $pyproject)) {
        return [PSCustomObject]@{ Success = $false; Message = "pyproject.toml niet gevonden" }
    }

    $content = Get-Content $pyproject -Raw
    $correctValue = 'cache-dir = "./.uv_cache"'
    $correctPattern = 'cache-dir\s*=\s*"\.\/\.uv_cache"'

    # Geval 1: [tool.uv] + correcte cache-dir → niets doen
    if ($content -match '\[tool\.uv\]' -and $content -match $correctPattern) {
        return [PSCustomObject]@{ Success = $true; Message = "uv cache-dir al correct geconfigureerd, overgeslagen" }
    }

    # Geval 2: [tool.uv] + andere/incomplete cache-dir → vervangen
    if ($content -match '\[tool\.uv\]' -and $content -match 'cache-dir') {
        $content = $content -replace 'cache-dir\s*=\s*"[^"]*"', $correctValue
        Set-Content $pyproject $content -NoNewline
        return [PSCustomObject]@{ Success = $true; Message = "Bestaande cache-dir vervangen door correcte waarde" }
    }

    # Geval 3: [tool.uv] zonder cache-dir → toevoegen onder de sectie
    if ($content -match '\[tool\.uv\]') {
        $content = $content -replace '(\[tool\.uv\])', "`$1`n$correctValue"
        Set-Content $pyproject $content -NoNewline
        return [PSCustomObject]@{ Success = $true; Message = "cache-dir toegevoegd aan bestaande [tool.uv] sectie" }
    }

    # Geval 4: [tool.uv] bestaat niet → bovenaan toevoegen
    $insert = "[tool.uv]`n$correctValue`n`n"
    $content = $insert + $content
    Set-Content $pyproject $content -NoNewline
    return [PSCustomObject]@{ Success = $true; Message = "[tool.uv] sectie met cache-dir bovenaan toegevoegd" }
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

        # Cache-dir instellen voor uv sync
        $cacheResult = Set-UvCacheDir -Root $Root
        if (-not $cacheResult.Success) {
            return $cacheResult
        }

        # ── uv sync met retry-flow ──
        Push-Location $Root

        $syncOutput = uv sync 2>&1
        $syncExit = $LASTEXITCODE

        # Poging 2: .venv verwijderen en opnieuw proberen
        if ($syncExit -ne 0) {
            $venvPath = Join-Path $Root ".venv"
            if (Test-Path $venvPath) {
                Remove-Item -Recurse -Force $venvPath
            }
            $syncOutput = uv sync 2>&1
            $syncExit = $LASTEXITCODE
        }

        # Poging 3: uv cache clean en opnieuw proberen
        if ($syncExit -ne 0) {
            uv cache clean 2>&1 | Out-Null
            $syncOutput = uv sync 2>&1
            $syncExit = $LASTEXITCODE
        }

        Pop-Location

        if ($syncExit -ne 0) {
            return [PSCustomObject]@{
                Success = $false
                Message = "uv sync mislukt na 3 pogingen. Probeer 'uv cache purge' handmatig en start opnieuw. Fout: $syncOutput"
            }
        }

        # Streamlit in nieuw venster, gewone uv run blokkeert bewust
        $isStreamlit = Test-IsStreamlitProject -Root $Root
        Push-Location $Root
        if ($isStreamlit) {
            Start-Process powershell -ArgumentList "-NoExit", "-Command", "uv run streamlit run $entryPoint"
        } else {
            uv run $entryPoint
        }
        Pop-Location

        return [PSCustomObject]@{ Success = $true; Message = "Project gestart: $entryPoint (streamlit: $isStreamlit)" }

    } catch {
        Pop-Location -ErrorAction SilentlyContinue
        return [PSCustomObject]@{ Success = $false; Message = "Onverwachte fout bij opstarten project: $_" }
    }
}
