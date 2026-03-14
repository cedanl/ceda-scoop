# uv-run.ps1
# Zoekt het entry-point en start het uv project op.
# Exporteert: Start-UvProject
# Vereist: uv-config.ps1, uv-sync.ps1

. "$PSScriptRoot\uv-config.ps1"
. "$PSScriptRoot\uv-sync.ps1"

function Start-UvProject {
    param([string]$Root = $PSScriptRoot)

    try {
        # ── Zoek entry-point in volgorde ──
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

        # ── Cache-dir instellen ──
        $cacheResult = Set-UvCacheDir -Root $Root
        Write-Host $cacheResult.Message -ForegroundColor Gray
        if (-not $cacheResult.Success) { return $cacheResult }

        # ── Streamlit check + feedback ──
        $isStreamlit = Test-IsStreamlitProject -Root $Root
        if ($isStreamlit) {
            Write-Host "Streamlit app gedetecteerd, app opent automatisch in de browser..." -ForegroundColor Cyan
        } else {
            Write-Host "Python project gedetecteerd, entry-point: $entryPoint" -ForegroundColor Cyan
        }

        # ── Extras ophalen en sync-argumenten opbouwen ──
        $extras = Get-UvExtras -Root $Root
        $syncArgs = @()
        foreach ($extra in $extras) { $syncArgs += "--extra"; $syncArgs += $extra }

        if ($extras.Count -gt 0) {
            Write-Host "Optionele dependencies gevonden: $($extras -join ', ')" -ForegroundColor Gray
        }

        # ── Sync ──
        $syncResult = Invoke-UvSync -Root $Root -SyncArgs $syncArgs
        if (-not $syncResult.Success) { return $syncResult }

        # ── Starten: Streamlit in nieuw venster, gewone uv run blokkeert bewust ──
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
