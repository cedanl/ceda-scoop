$ErrorActionPreference = "Continue"

# ============================================
#  FUNCTIONS
# ============================================
# Progress bar function
function Show-Progress {
    param(
        [int]$DurationMs = 1000
    )
    
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
$options = @("Verwijdering Starten", "Afsluiten")
$selected = 0

do {
    Clear-Host
    
    # Redraw header
    Write-Host ""
    Write-Host " ===================================" -ForegroundColor Red
    Write-Host " SCOOP UNINSTALL TOOL" -ForegroundColor Red
    Write-Host " ===================================" -ForegroundColor Red
    Write-Host ""
    Write-Host "LET OP: Gebruik dit hulpmiddel alleen voor het verwijderen van Scoop" -ForegroundColor yellow
    Write-Host ""
    Write-Host " Dit script zal:" -ForegroundColor Gray
    Write-Host "  * Scoop & buckets verwijderen" -ForegroundColor Gray
    Write-Host "  * Alle applicaties verwijderen" -ForegroundColor Gray
    Write-Host "  * Gebruikerspaden voor rtools verwijderen" -ForegroundColor Gray    
    Write-Host ""
    Write-Host " Developed by: " -NoNewline -ForegroundColor DarkGray
    Write-Host " CEDA " -ForegroundColor Cyan
    Write-Host ""
    Write-Host " ===============================================" -ForegroundColor DarkGray
    Write-Host ""
    
    # Draw menu
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
    
    if ($key.VirtualKeyCode -eq 38) { # Up arrow
        $selected = ($selected - 1 + $options.Count) % $options.Count
    } elseif ($key.VirtualKeyCode -eq 40) { # Down arrow
        $selected = ($selected + 1) % $options.Count
    }
    
} while ($key.VirtualKeyCode -ne 13) # Enter key

Write-Host ""

if ($selected -eq 1) {
    Write-Host " Afsluiten..." -ForegroundColor Yellow
    exit
}


# ============================================
#  SECTION 1
# ============================================
# Ask about cleanup
Write-Host "Wil je doorgaan met het verwijderen van Scoop? " -NoNewline -ForegroundColor Red
Write-Host -NoNewline "(J/n) " -ForegroundColor Cyan
$response = Read-Host

if ($response -ne "j") {
    Show-CompletionPrompt
    return  # or exit, depending on your script structure
}

Write-Host "Weet je het zeker? Dit verwijdert alle geïnstalleerde Scoop-applicaties en kan niet ongedaan worden gemaakt. " -NoNewline -ForegroundColor Red
Write-Host -NoNewline "(J/n) " -ForegroundColor Cyan
$response = Read-Host

Write-Host "[1] Scoop, buckets & applicaties verwijderen..." -ForegroundColor Yellow

if (Get-Command scoop -ErrorAction SilentlyContinue) {
    scoop uninstall scoop
} else {
    Write-Host "Scoop is al verwijderd, stap overgeslagen..." -ForegroundColor Red
}

Show-Progress -DurationMs 800

# ============================================
#  SECTION 2
# ============================================

Write-Host "[2] Gebruikerspaden voor rtools & msys2 mingw64 (weasyprint) verwijderen..." -ForegroundColor Yellow
$currentUserPath = [System.Environment]::GetEnvironmentVariable("PATH", [System.EnvironmentVariableTarget]::User)

# Remove first path
$path1 = "%env:%USERPROFILE%\scoop\apps\rtools\current\usr\bin"
if ($currentUserPath.Contains($path1)) {
    Write-Host "$path1 verwijderen..."
    $currentUserPath = $currentUserPath.Replace(";$path1", "").Replace("$path1;", "").Replace($path1, "")
    [System.Environment]::SetEnvironmentVariable("PATH", $currentUserPath, [System.EnvironmentVariableTarget]::User)
    Write-Host "Succesvol verwijderd." -ForegroundColor Green
} else {
    Write-Host "$path1 staat niet in PATH." -ForegroundColor Gray
}

# Remove second path
$path2 = "%env:%USERPROFILE%\scoop\apps\rtools\current\x86_64-w64-mingw32.static.posix\bin"
if ($currentUserPath.Contains($path2)) {
    Write-Host "$path2 verwijderen..."
    $currentUserPath = $currentUserPath.Replace(";$path2", "").Replace("$path2;", "").Replace($path2, "")
    [System.Environment]::SetEnvironmentVariable("PATH", $currentUserPath, [System.EnvironmentVariableTarget]::User)
    Write-Host "Succesvol verwijderd." -ForegroundColor Green
} else {
    Write-Host "$path2 staat niet in PATH." -ForegroundColor Gray
}


Write-Host "Verwijdering van rtools-padconfiguratie voltooid." -ForegroundColor Yellow

# ============================================
#  END APPLICATION
# ============================================
Show-CompletionPrompt