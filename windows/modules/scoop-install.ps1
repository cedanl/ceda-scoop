# scoop-core.ps1
# Installeert Scoop, core dependencies (7zip, aria2, git) en buckets.
# Aanroepvolgorde in entry-point: Install-Scoop → Install-CoreDeps → Add-Buckets
# Exporteert: Install-Scoop, Install-CoreDeps, Add-Buckets

# ── Scoop ──
function Install-Scoop {
    try {
        if (-not (Get-Command scoop -ErrorAction SilentlyContinue)) {
            Invoke-RestMethod -Uri 'https://get.scoop.sh' | Invoke-Expression
            return [PSCustomObject]@{ Success = $true; Message = "Scoop geinstalleerd" }
        }
        return [PSCustomObject]@{ Success = $true; Message = "Scoop al aanwezig, overgeslagen" }
    } catch {
        return [PSCustomObject]@{ Success = $false; Message = "Fout bij Scoop installatie: $_" }
    }
}

# ── Core dependencies (git vereist voor Add-Buckets) ──
function Install-CoreDeps {
    $deps      = @("git", "7zip", "aria2")
    $installed = @()
    $skipped   = @()

    try {
        if (-not (Get-Command scoop -ErrorAction SilentlyContinue)) {
            return [PSCustomObject]@{ Success = $false; Message = "Scoop niet gevonden" }
        }

        $installedApps = scoop list 2>&1

        foreach ($dep in $deps) {
            if ($installedApps | Select-String "^\s*$dep\s") {
                $skipped += $dep
            } else {
                $output = scoop install $dep 2>&1
                if ($LASTEXITCODE -ne 0) {
                    return [PSCustomObject]@{ Success = $false; Message = "Fout bij installatie van $dep`: $output" }
                }
                $installed += $dep
            }
        }

        return [PSCustomObject]@{
            Success = $true
            Message = "Geinstalleerd: $($installed -join ', ') | Overgeslagen: $($skipped -join ', ')"
        }
    } catch {
        return [PSCustomObject]@{ Success = $false; Message = "Onverwachte fout bij core deps: $_" }
    }
}

# ── Buckets (aanroepen na Install-CoreDeps, git is vereist) ──
function Add-Buckets {
    try {
        scoop bucket add extras  2>&1 | Out-Null
        scoop bucket add versions 2>&1 | Out-Null
        return [PSCustomObject]@{ Success = $true; Message = "Buckets toegevoegd (extras, versions)" }
    } catch {
        return [PSCustomObject]@{ Success = $false; Message = "Fout bij toevoegen buckets: $_" }
    }
}
