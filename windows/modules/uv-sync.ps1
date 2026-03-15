# uv-sync.ps1
# Synchroniseert uv project dependencies met retry-flow en kill-prompt.
# Exporteert: Invoke-UvSync
# Vereist: uv-config.ps1 (Get-UvExtras)

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

# ── uv sync met retry-flow ──
function Invoke-UvSync {
    param(
        [string]$Root,
        [string[]]$SyncArgs = @()
    )

    Push-Location $Root

    Write-Host ""
    Write-Host "Packages installeren via uv sync..." -ForegroundColor Yellow
    uv sync @SyncArgs
    $syncExit = $LASTEXITCODE

    # Poging 2: .venv verwijderen
    if ($syncExit -ne 0) {
        Write-Host "Sync mislukt, .venv verwijderen en opnieuw proberen..." -ForegroundColor Yellow
        $removed = Remove-DirectoryWithFallback -Path (Join-Path $Root ".venv")
        if (-not $removed) {
            Pop-Location
            return [PSCustomObject]@{ Success = $false; Message = "Kon .venv niet verwijderen, sync gestopt." }
        }
        uv sync @SyncArgs
        $syncExit = $LASTEXITCODE
    }

    # Poging 3: uv cache clean + .uv_cache verwijderen
    if ($syncExit -ne 0) {
        Write-Host "Sync mislukt, uv cache legen en opnieuw proberen..." -ForegroundColor Yellow
        uv cache clean 2>&1 | Out-Null
        $removed = Remove-DirectoryWithFallback -Path (Join-Path $Root ".uv_cache")
        if (-not $removed) {
            Pop-Location
            return [PSCustomObject]@{ Success = $false; Message = "Kon .uv_cache niet verwijderen, sync gestopt." }
        }
        uv sync @SyncArgs
        $syncExit = $LASTEXITCODE
    }

    Pop-Location

    if ($syncExit -ne 0) {
        return [PSCustomObject]@{
            Success = $false
            Message = "uv sync mislukt na 3 pogingen. Probeer 'uv cache purge' handmatig en start opnieuw."
        }
    }

    return [PSCustomObject]@{ Success = $true; Message = "uv sync geslaagd" }
}
