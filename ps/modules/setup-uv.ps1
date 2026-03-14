# setup-uv.ps1
# Installeert uv (via Scoop) en start een uv-project op.
# Exporteert: Install-UV, Set-UvCacheDir, Get-UvExtras, Start-UvProject

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
        return [PSCustomObject]@{ Success = $true; Message = "cache-dir al correct ingesteld in pyproject.toml, overgeslagen" }
    }

    # Geval 2: [tool.uv] + andere/incomplete cache-dir → vervangen
    if ($content -match '\[tool\.uv\]' -and $content -match 'cache-dir') {
        $content = $content -replace 'cache-dir\s*=\s*"[^"]*"', $correctValue
        Set-Content $pyproject $content -NoNewline
        return [PSCustomObject]@{ Success = $true; Message = "Bestaande cache-dir in pyproject.toml vervangen door correcte waarde" }
    }

    # Geval 3: [tool.uv] zonder cache-dir → toevoegen onder de sectie
    if ($content -match '\[tool\.uv\]') {
        $content = $content -replace '(\[tool\.uv\])', "`$1`n$correctValue"
        Set-Content $pyproject $content -NoNewline
        return [PSCustomObject]@{ Success = $true; Message = "cache-dir toegevoegd aan [tool.uv] sectie in pyproject.toml" }
    }

    # Geval 4: [tool.uv] bestaat niet → bovenaan toevoegen
    $insert = "[tool.uv]`n$correctValue`n`n"
    $content = $insert + $content
    Set-Content $pyproject $content -NoNewline
    return [PSCustomObject]@{ Success = $true; Message = "[tool.uv] sectie met cache-dir toegevoegd aan pyproject.toml" }
}

# ── Lees optional-dependencies uit pyproject.toml, filter 'dev' eruit ──
function Get-UvExtras {
    param([string]$Root)

    $pyproject = Join-Path $Root "pyproject.toml"
    if (-not (Test-Path $pyproject)) { return @() }

    $content = Get-Content $pyproject -Raw

    if ($content -notmatch '\[project\.optional-dependencies\]') { return @() }

    $extras = @()
    $inSection = $false

    foreach ($line in ($content -split "`n")) {
        if ($line -match '^\[project\.optional-dependencies\]') {
            $inSection = $true
            continue
        }
        if ($inSection -and $line -match '^\[') { break }
        if ($inSection -and $line -match '^\s*(\w+)\s*=\s*\[') {
            $name = $matches[1]
            if ($name -ne "dev") { $extras += $name }
        }
    }

    return $extras
}

# ── Vraag gebruiker of Python processen gestopt mogen worden ──
function Invoke-KillPythonPrompt {
    Write-Host ""
    Write-Host "Kon bestanden niet verwijderen - een Python process heeft ze nog open." -ForegroundColor Red
    Write-Host "Wil je alle actieve Python processen stoppen en opnieuw proberen?" -ForegroundColor Yellow
    Write-Host "  [J] Ja, stop processen" -ForegroundColor Cyan
    Write-Host "  [N] Nee, annuleren" -ForegroundColor Gray
    $keuze = Read-Host "Keuze (J/N)"

    if ($keuze -eq "j") {
        Write-Host "Python processen stoppen..." -ForegroundColor Yellow
        Get-Process -Name "python*", "pythonw*" -ErrorAction SilentlyContinue | Stop-Process -Force
        Start-Sleep -Milliseconds 500
        return $true
    }

    return $false
}

# ── Verwijder map met fallback naar kill-prompt ──
function Remove-DirectoryWithFallback {
    param([string]$Path)

    if (-not (Test-Path $Path)) { return $true }

    try {
        Remove-Item -Recurse -Force $Path -ErrorAction Stop
        return $true
    } catch {
        $kill = Invoke-KillPythonPrompt
        if (-not $kill) { return $false }

        try {
            Remove-Item -Recurse -Force $Path -ErrorAction Stop
            return $true
        } catch {
            Write-Host "Verwijderen mislukt ook na stoppen van processen: $_" -ForegroundColor Red
            return $false
        }
    }
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

        # Cache-dir instellen en feedback tonen
        $cacheResult = Set-UvCacheDir -Root $Root
        Write-Host $cacheResult.Message -ForegroundColor Gray
        if (-not $cacheResult.Success) { return $cacheResult }

        # Streamlit check + feedback
        $isStreamlit = Test-IsStreamlitProject -Root $Root
        if ($isStreamlit) {
            Write-Host "Streamlit app gedetecteerd, app opent automatisch in de browser..." -ForegroundColor Cyan
        } else {
            Write-Host "Python project gedetecteerd, entry-point: $entryPoint" -ForegroundColor Cyan
        }

        # Extras ophalen en sync-argumenten opbouwen
        $extras = Get-UvExtras -Root $Root
        $syncArgs = @()
        foreach ($extra in $extras) { $syncArgs += "--extra"; $syncArgs += $extra }

        if ($extras.Count -gt 0) {
            Write-Host "Optionele dependencies gevonden: $($extras -join ', ')" -ForegroundColor Gray
        }

        # ── uv sync (live output) met retry-flow ──
        Push-Location $Root

        Write-Host ""
        Write-Host "Packages installeren via uv sync..." -ForegroundColor Yellow
        uv sync @syncArgs
        $syncExit = $LASTEXITCODE

        # Poging 2: .venv verwijderen met kill-prompt als fallback
        if ($syncExit -ne 0) {
            Write-Host "Sync mislukt, .venv verwijderen en opnieuw proberen..." -ForegroundColor Yellow
            $venvPath = Join-Path $Root ".venv"
            $removed = Remove-DirectoryWithFallback -Path $venvPath
            if (-not $removed) {
                Pop-Location
                return [PSCustomObject]@{ Success = $false; Message = "Kon .venv niet verwijderen, sync gestopt." }
            }
            uv sync @syncArgs
            $syncExit = $LASTEXITCODE
        }

        # Poging 3: uv cache clean + .uv_cache verwijderen met kill-prompt als fallback
        if ($syncExit -ne 0) {
            Write-Host "Sync mislukt, uv cache legen en opnieuw proberen..." -ForegroundColor Yellow
            uv cache clean 2>&1 | Out-Null
            $cachePath = Join-Path $Root ".uv_cache"
            $removed = Remove-DirectoryWithFallback -Path $cachePath
            if (-not $removed) {
                Pop-Location
                return [PSCustomObject]@{ Success = $false; Message = "Kon .uv_cache niet verwijderen, sync gestopt." }
            }
            uv sync @syncArgs
            $syncExit = $LASTEXITCODE
        }

        Pop-Location

        if ($syncExit -ne 0) {
            return [PSCustomObject]@{
                Success = $false
                Message = "uv sync mislukt na 3 pogingen. Probeer 'uv cache prune' handmatig en start opnieuw."
            }
        }

        # Streamlit in nieuw venster, gewone uv run blokkeert bewust
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
