# r-install.ps1
# Installeert R, Rtools en Positron via Scoop.
# Exporteert: Install-RDeps

# ── R + Rtools + Positron ──
function Install-RDeps {
    $deps      = @("main/r", "main/rtools", "extras/positron")
    $installed = @()
    $skipped   = @()

    try {
        $installedApps = scoop list 2>&1

        foreach ($dep in $deps) {
            $name = $dep.Split("/")[-1]
            if ($installedApps | Select-String "^\s*$name\s") {
                $skipped += $name
            } else {
                $output = scoop install $dep 2>&1
                if ($LASTEXITCODE -ne 0) {
                    return [PSCustomObject]@{ Success = $false; Message = "Fout bij installatie van $name`: $output" }
                }
                $installed += $name
            }
        }

        return [PSCustomObject]@{
            Success = $true
            Message = "Geinstalleerd: $($installed -join ', ') | Overgeslagen: $($skipped -join ', ')"
        }
    } catch {
        return [PSCustomObject]@{ Success = $false; Message = "Onverwachte fout bij R deps: $_" }
    }
}
