$ErrorActionPreference = "Continue"

# ============================================
#  FUNCTIONS
# ============================================
function Show-Progress {
    param([int]$DurationMs = 1000)
    
    $barLength = 30
    $steps = 20
    $delay = $DurationMs / $steps
    
    Write-Host ""
    Write-Host "[" -NoNewline -ForegroundColor DarkGray
    
    for ($i = 0; $i -lt $barLength; $i++) {
        Write-Host "=" -NoNewline -ForegroundColor Yellow
        Start-Sleep -Milliseconds $delay
    }
    
    Write-Host "]" -ForegroundColor DarkGray
    Write-Host ""
}

function Show-CompletionPrompt {
    Write-Host ""
    Write-Host ""
    Write-Host " =============================" -ForegroundColor Green
    Write-Host " >> ALLE STAPPEN VOLTOOID <<" -ForegroundColor Green
    Write-Host " =============================" -ForegroundColor Green
    Write-Host ""
    Write-Host "Wil je de tool opnieuw uitvoeren? " -NoNewline -ForegroundColor Yellow
    Write-Host -NoNewline "(J/n) " -ForegroundColor Cyan
    $response = Read-Host
    
    if ($response -eq "j") {
        Write-Host "`nTool opnieuw starten..." -ForegroundColor Cyan
        Start-Sleep -Milliseconds 500
        & $PSCommandPath
    } else {
        Write-Host "`nAfsluiten. Tot ziens!" -ForegroundColor Green
    }
}

# ============================================
#  [0] MENU
# ============================================
$options = @("Setup Starten", "Afsluiten")
$selected = 0

do {
    Clear-Host
    
    Write-Host ""
    Write-Host " ===================================" -ForegroundColor Cyan
    Write-Host " CEDA QUICKRUN - MVP" -ForegroundColor Cyan
    Write-Host " ===================================" -ForegroundColor Cyan
    Write-Host ""
    Write-Host "LET OP: Dit hulpmiddel heeft geen invloed op geïnstalleerde applicaties en vereist geen beheerdersrechten." -ForegroundColor yellow
    Write-Host ""
    Write-Host " Dit script zal:" -ForegroundColor Gray
    Write-Host "  * Scoop & buckets installeren" -ForegroundColor Gray
    Write-Host "  * 7zip, aria2 & git installeren voor geoptimaliseerde installatie" -ForegroundColor Gray
    Write-Host "  * R, RStudio & Rtools installeren" -ForegroundColor Gray
    Write-Host "  * Gebruikerspaden instellen die nodig zijn voor rtools" -ForegroundColor Gray  
    Write-Host "  * RStudio openen met huidig project & renv::restore()" -ForegroundColor Gray  
    Write-Host ""
    Write-Host " Developed by: " -NoNewline -ForegroundColor DarkGray
    Write-Host " CEDA " -ForegroundColor Cyan
    Write-Host ""
    Write-Host " ===============================================" -ForegroundColor DarkGray
    Write-Host ""
    
    for ($i = 0; $i -lt $options.Count; $i++) {
        if ($i -eq $selected) {
            Write-Host " > " -NoNewline -ForegroundColor Green
            Write-Host $options[$i] -ForegroundColor Green
        } else {
            Write-Host "   $($options[$i])" -ForegroundColor Gray
        }
    }
    
    Write-Host ""
    Write-Host " Gebruik de pijltjestoetsen om te selecteren, druk op ENTER om te bevestigen" -ForegroundColor DarkGray
    
    $key = $Host.UI.RawUI.ReadKey("NoEcho,IncludeKeyDown")
    
    if ($key.VirtualKeyCode -eq 38) { 
        $selected = ($selected - 1 + $options.Count) % $options.Count
    } elseif ($key.VirtualKeyCode -eq 40) { 
        $selected = ($selected + 1) % $options.Count
    }
    
} while ($key.VirtualKeyCode -ne 13)

Write-Host ""

if ($selected -eq 1) {
    Write-Host " Afsluiten..." -ForegroundColor Yellow
    exit
}

# ============================================
#  SECTION 1 - Scoop / Rtools Config
# ============================================
Write-Host "[1] Scoop en buckets installeren..." -ForegroundColor Yellow

if (-not (Get-Command scoop -ErrorAction SilentlyContinue)) {
    Invoke-RestMethod -Uri 'https://get.scoop.sh' | Invoke-Expression
} else {
    Write-Host "Scoop is al geïnstalleerd, stap overgeslagen..." -ForegroundColor Red
}

scoop bucket add extras
scoop bucket add versions

Show-Progress -DurationMs 800

Write-Host "[2] Apps installeren voor geoptimaliseerde installatie..." -ForegroundColor Yellow

scoop install 7zip
scoop install git
scoop install aria2

Show-Progress -DurationMs 800

Write-Host "[3] R, RStudio & Rtools installeren..." -ForegroundColor Yellow

scoop install main/r
scoop install extras/rstudio
scoop install main/rtools
Show-Progress -DurationMs 800

Write-Host "[4] Gebruikerspaden instellen voor rtools..." -ForegroundColor Yellow
$currentUserPath = [System.Environment]::GetEnvironmentVariable("PATH", [System.EnvironmentVariableTarget]::User)

$path1 = [System.Environment]::ExpandEnvironmentVariables("%USERPROFILE%\scoop\apps\rtools\current\usr\bin")
if ($currentUserPath.Contains($path1)) {
    Write-Host "$path1 staat al in PATH." -ForegroundColor Gray
} else {
    Write-Host "$path1 toevoegen..."
    $currentUserPath = "$currentUserPath;$path1"
    [System.Environment]::SetEnvironmentVariable("PATH", $currentUserPath, [System.EnvironmentVariableTarget]::User)
    Write-Host "Succesvol toegevoegd." -ForegroundColor Green
}

$path2 = [System.Environment]::ExpandEnvironmentVariables("%USERPROFILE%\scoop\apps\rtools\current\x86_64-w64-mingw32.static.posix\bin")
if ($currentUserPath.Contains($path2)) {
    Write-Host "$path2 staat al in PATH." -ForegroundColor Gray
} else {
    Write-Host "$path2 toevoegen..."
    $currentUserPath = "$currentUserPath;$path2"
    [System.Environment]::SetEnvironmentVariable("PATH", $currentUserPath, [System.EnvironmentVariableTarget]::User)
    Write-Host "Succesvol toegevoegd." -ForegroundColor Green
}
Write-Host "Padconfiguratie voltooid." -ForegroundColor Yellow

Show-CompletionPrompt