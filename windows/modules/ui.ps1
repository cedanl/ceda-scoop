# ui.ps1
# Gedeelde presentatiefuncties voor CEDA Quickrun.
# Geen business logica — alleen output naar de terminal.

# ── Voortgangsbalk ──
function Show-Progress {
    param([int]$DurationMs = 1000)

    $steps = 20
    $delay = $DurationMs / $steps

    Write-Host ""
    Write-Host "[" -NoNewline -ForegroundColor DarkGray
    for ($i = 0; $i -lt 30; $i++) {
        Write-Host "=" -NoNewline -ForegroundColor Yellow
        Start-Sleep -Milliseconds $delay
    }
    Write-Host "]" -ForegroundColor DarkGray
    Write-Host ""
}

# ── Titelbanner ──
function Show-Header {
    param(
        [string]$Title = "CEDA QUICKRUN",
        [string]$Color = "Cyan"
    )

    Clear-Host
    Write-Host ""
    Write-Host " ===================================" -ForegroundColor $Color
    Write-Host " $Title"                             -ForegroundColor $Color
    Write-Host " ===================================" -ForegroundColor $Color
    Write-Host ""
}

# ── Pijltjesmenu — compatibel met cmd en conhost ──
function Show-Menu {
    param(
        [string[]]$Options,
        [string]$Title = "",
        [string]$Color = "Green"
    )

    $selected = 0

    do {
        Clear-Host
        if ($Title) {
            Write-Host " $Title" -ForegroundColor Cyan
            Write-Host ""
        }

        for ($i = 0; $i -lt $Options.Count; $i++) {
            if ($i -eq $selected) {
                Write-Host " > $($Options[$i])" -ForegroundColor $Color
            } else {
                Write-Host "   $($Options[$i])" -ForegroundColor Gray
            }
        }

        Write-Host ""
        Write-Host " Pijltjestoetsen om te selecteren, ENTER om te bevestigen" -ForegroundColor DarkGray

        $key = $Host.UI.RawUI.ReadKey("NoEcho,IncludeKeyDown")

        if ($key.VirtualKeyCode -eq 38) { $selected = ($selected - 1 + $Options.Count) % $Options.Count }
        if ($key.VirtualKeyCode -eq 40) { $selected = ($selected + 1) % $Options.Count }

    } while ($key.VirtualKeyCode -ne 13)

    return $selected
}

# ── Afsluiten / opnieuw starten prompt ──
function Show-CompletionPrompt {
    param([string]$ScriptPath = $PSCommandPath)

    Write-Host ""
    Write-Host " =============================" -ForegroundColor Green
    Write-Host " >> ALLE STAPPEN VOLTOOID <<" -ForegroundColor Green
    Write-Host " =============================" -ForegroundColor Green
    Write-Host ""
    Write-Host "Wil je de tool opnieuw uitvoeren? " -NoNewline -ForegroundColor Yellow
    Write-Host "(J/n) " -NoNewline -ForegroundColor Cyan
    $response = Read-Host

    if ($response -eq "j") {
        Write-Host "`nTool opnieuw starten..." -ForegroundColor Cyan
        Start-Sleep -Milliseconds 500
        & $ScriptPath
    } else {
        Write-Host "`nAfsluiten. Tot ziens!" -ForegroundColor Green
    }
}
