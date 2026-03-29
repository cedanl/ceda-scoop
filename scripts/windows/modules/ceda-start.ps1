# ceda-start.ps1
# Entry-point voor CEDA Quickrun. Detecteert projecttype en voert de juiste flow uit.
# Aanroepvolgorde: detect → scoop → install → config → sync → run
# Structuur na kopiëren: repo-root/modules/ceda-start.ps1 → Root is één niveau omhoog

$ErrorActionPreference = "Continue"

# ── Paden: Root = repo root, modules naast dit script ──
$Root        = Split-Path -Parent $PSScriptRoot
$ModulesPath = $PSScriptRoot

# ── Modules laden ──
. "$ModulesPath\ui.ps1"
. "$ModulesPath\detect-project.ps1"
. "$ModulesPath\scoop-install.ps1"
. "$ModulesPath\uv-install.ps1"
. "$ModulesPath\uv-config.ps1"
. "$ModulesPath\uv-sync.ps1"
. "$ModulesPath\uv-run.ps1"
. "$ModulesPath\r-install.ps1"
. "$ModulesPath\r-config.ps1"
. "$ModulesPath\r-sync.ps1"
. "$ModulesPath\r-run.ps1"

# ── Helper: stap uitvoeren en bij fout stoppen ──
function Invoke-Step {
    param(
        [string]$Label,
        [scriptblock]$Action
    )

    Write-Host ""
    Write-Host $Label -ForegroundColor Yellow
    $result = & $Action
    if ($result -and -not $result.Success) {
        Write-Host "FOUT: $($result.Message)" -ForegroundColor Red
        Read-Host "`nDruk op Enter om af te sluiten"
        exit 1
    }
    if ($result) {
        Write-Host $result.Message -ForegroundColor Gray
    }
}

# ── Menu ──
Show-Header -Title "CEDA QUICKRUN"
Write-Host " LET OP: Dit hulpmiddel vereist geen beheerdersrechten." -ForegroundColor Yellow
Write-Host ""

$keuze = Show-Menu -Options @("Starten", "Afsluiten") -Title "CEDA QUICKRUN"
if ($keuze -eq 1) { exit }

# ── Projecttype detecteren ──
Show-Header -Title "CEDA QUICKRUN"
Write-Host ""
Write-Host "Projecttype detecteren..." -ForegroundColor Yellow
$project = Get-ProjectType -Root $Root

if ($project.Type -eq "unknown") {
    Write-Host "FOUT: $($project.Message)" -ForegroundColor Red
    Read-Host "`nDruk op Enter om af te sluiten"
    exit 1
}

Write-Host $project.Message -ForegroundColor Gray

# ── Scoop + core deps (altijd) ──
Invoke-Step "[1] Scoop installeren..."              { Install-Scoop }
Invoke-Step "[2] Core dependencies installeren..."  { Install-CoreDeps }
Invoke-Step "[3] Buckets toevoegen..."              { Add-Buckets }

# ── Projecttype-specifieke flow ──
if ($project.Type -eq "uv") {

    Invoke-Step "[4] uv installeren..."              { Install-UV }
    Invoke-Step "[5] pyproject.toml configureren..." { Set-UvCacheDir -Root $Root }

    $extras   = Get-UvExtras -Root $Root
    $syncArgs = @()
    foreach ($extra in $extras) { $syncArgs += "--extra"; $syncArgs += $extra }
    if ($extras.Count -gt 0) {
        Write-Host "Optionele dependencies: $($extras -join ', ')" -ForegroundColor Gray
    }

    Invoke-Step "[6] Packages installeren..."        { Invoke-UvSync -Root $Root -SyncArgs $syncArgs }
    Invoke-Step "[7] Project starten..."             { Start-UvProject -Root $Root -IsStreamlit (Test-IsStreamlitProject -Root $Root) }

} elseif ($project.Type -eq "r") {

    Invoke-Step "[4] R, Rtools en Positron installeren..." { Install-RDeps }
    Invoke-Step "[5] Rtools paden instellen..."            { Set-RtoolsPaths }
    Invoke-Step "[6] Positron R interpreter instellen..."  { Set-PositronRConfig }
    Invoke-Step "[7] Packages installeren (renv)..."       { Invoke-RenvRestore -Root $Root }
    Invoke-Step "[8] Project openen in Positron..."        { Start-RProject -Root $Root }

}

Show-CompletionPrompt -ScriptPath $PSCommandPath
