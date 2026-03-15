# r-config.ps1
# Configureert de omgeving voor R projecten.
# Exporteert: Set-RtoolsPaths, Set-PositronRConfig

# ── Rtools PATH-entries toevoegen aan gebruikers-PATH ──
function Set-RtoolsPaths {
    $paths = @(
        "$env:USERPROFILE\scoop\apps\rtools\current\usr\bin",
        "$env:USERPROFILE\scoop\apps\rtools\current\x86_64-w64-mingw32.static.posix\bin"
    )

    try {
        $currentPath = [System.Environment]::GetEnvironmentVariable("PATH", [System.EnvironmentVariableTarget]::User)
        $added   = @()
        $skipped = @()

        foreach ($path in $paths) {
            if ($currentPath -like "*$path*") {
                $skipped += $path
            } else {
                $currentPath = "$currentPath;$path"
                $added += $path
            }
        }

        [System.Environment]::SetEnvironmentVariable("PATH", $currentPath, [System.EnvironmentVariableTarget]::User)

        return [PSCustomObject]@{
            Success = $true
            Message = "Toegevoegd: $($added.Count) pad(en) | Overgeslagen: $($skipped.Count) pad(en)"
        }
    } catch {
        return [PSCustomObject]@{ Success = $false; Message = "Fout bij instellen Rtools PATH: $_" }
    }
}

# ── Positron R interpreter pad instellen via settings.json (raw tekst, JSONC-safe) ──
function Set-PositronRConfig {
    $settingsDir  = "$env:APPDATA\Positron\User"
    $settingsPath = "$settingsDir\settings.json"

    # Enkelvoudige backslash escaping voor JSON: \ wordt \\
    $rScoopPath = "$env:USERPROFILE\scoop\apps\r" -replace '\\', '\\'

    try {
        # Map aanmaken als die nog niet bestaat
        if (-not (Test-Path $settingsDir)) {
            New-Item -ItemType Directory -Path $settingsDir -Force | Out-Null
        }

        # File aanmaken als die nog niet bestaat
        if (-not (Test-Path $settingsPath)) {
            Set-Content $settingsPath '{}' -Encoding UTF8
        }

        $content = Get-Content $settingsPath -Raw

        # ── Geval 1: key bestaat al met lege array [] ──
        if ($content -match '"positron\.r\.customRootFolders"\s*:\s*\[\s*\]') {
            $content = $content -replace '"positron\.r\.customRootFolders"\s*:\s*\[\s*\]', """positron.r.customRootFolders"": [""$rScoopPath""]"
            Set-Content $settingsPath $content -Encoding UTF8
            return [PSCustomObject]@{ Success = $true; Message = "Positron R pad ingesteld in bestaande lege array" }
        }

        # ── Geval 2: key bestaat al met waarden — check of pad er al in zit ──
        if ($content -match '"positron\.r\.customRootFolders"') {
            if ($content -match [regex]::Escape($rScoopPath)) {
                return [PSCustomObject]@{ Success = $true; Message = "Positron R config al correct ingesteld, overgeslagen" }
            }
            $content = $content -replace '("positron\.r\.customRootFolders"\s*:\s*\[)', "`$1`n`t`t""$rScoopPath"","
            Set-Content $settingsPath $content -Encoding UTF8
            return [PSCustomObject]@{ Success = $true; Message = "Positron R pad toegevoegd aan bestaande array" }
        }

        # ── Geval 3: key bestaat niet — injecteren na eerste { ──
        $content = $content -replace '(\{)', "`$1`n`t""positron.r.customRootFolders"": [""$rScoopPath""],"
        Set-Content $settingsPath $content -Encoding UTF8
        return [PSCustomObject]@{ Success = $true; Message = "Positron R pad toegevoegd aan settings.json" }

    } catch {
        return [PSCustomObject]@{ Success = $false; Message = "Fout bij instellen Positron R config: $_" }
    }
}
