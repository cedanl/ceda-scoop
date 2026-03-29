# r-run.ps1
# Opent het R-project in Positron.
# Exporteert: Start-RProject

function Start-RProject {
    param([string]$Root = $PSScriptRoot)

    try {
        Write-Host "Positron openen met project: $Root" -ForegroundColor Cyan
        Start-Process positron -ArgumentList "--disable-workspace-trust", $Root

        return [PSCustomObject]@{ Success = $true; Message = "Positron geopend met $Root" }
    } catch {
        return [PSCustomObject]@{ Success = $false; Message = "Onverwachte fout bij opstarten Positron: $_" }
    }
}
