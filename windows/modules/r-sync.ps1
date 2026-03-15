# r-sync.ps1
# Herstelt renv dependencies via renv::restore().
# Configuratie (PPM, binary) via .Rprofile in de repo.
# Exporteert: Invoke-RenvRestore

function Invoke-RenvRestore {
    param([string]$Root = $PSScriptRoot)

    # ── Check renv aanwezig ──
    if (-not (Test-Path (Join-Path $Root "renv\activate.R"))) {
        return [PSCustomObject]@{ Success = $false; Message = "Geen renv project gevonden (renv/activate.R ontbreekt)" }
    }

    # ── Rtools paden activeren ──
    $env:PATH = [System.Environment]::GetEnvironmentVariable("PATH", "User") + ";" + $env:PATH

    # ── R.exe pad via Scoop ──
    $rExe = "$env:USERPROFILE\scoop\apps\r\current\bin\x64\R.exe"

    if (-not (Test-Path $rExe)) {
        return [PSCustomObject]@{ Success = $false; Message = "R.exe niet gevonden op: $rExe" }
    }

    # ── Tijdelijk R script zonder BOM ──
    $tempScript = Join-Path $env:TEMP "ceda-renv-restore.R"
    [System.IO.File]::WriteAllText($tempScript, "renv::restore(prompt = FALSE)", [System.Text.UTF8Encoding]::new($false))

    $process = Start-Process `
        -FilePath $rExe `
        -ArgumentList "--no-save", "--no-restore", "--file=$tempScript" `
        -WorkingDirectory $Root `
        -Wait `
        -PassThru `
        -NoNewWindow

    Remove-Item $tempScript -ErrorAction SilentlyContinue

    if ($process.ExitCode -ne 0) {
        return [PSCustomObject]@{ Success = $false; Message = "renv restore mislukt (exit code $($process.ExitCode))" }
    }

    return [PSCustomObject]@{ Success = $true; Message = "renv restore geslaagd" }
}
