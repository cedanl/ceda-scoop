# uv-config.ps1
# Leest en configureert pyproject.toml voor uv projecten.
# Exporteert: Set-UvCacheDir, Get-UvExtras, Test-IsStreamlitProject

# ── Streamlit check via pyproject.toml ──
function Test-IsStreamlitProject {
    param([string]$Root)

    $pyproject = Join-Path $Root "pyproject.toml"
    if (-not (Test-Path $pyproject)) { return $false }

    $content = Get-Content $pyproject -Raw
    return $content -match "streamlit"
}

# ── Zorg dat [tool.uv] cache-dir correct aanwezig is in pyproject.toml ──
function Set-UvCacheDir {
    param([string]$Root)

    $pyproject = Join-Path $Root "pyproject.toml"

    if (-not (Test-Path $pyproject)) {
        return [PSCustomObject]@{ Success = $false; Message = "pyproject.toml niet gevonden" }
    }

    $content = Get-Content $pyproject -Raw
    $correctValue = 'cache-dir = "./.uv_cache"'
    $correctPattern = 'cache-dir\s*=\s*"\.\/\.uv_cache"'

    # Geval 1: [tool.uv] + correcte cache-dir → niets doen
    if ($content -match '\[tool\.uv\]' -and $content -match $correctPattern) {
        return [PSCustomObject]@{ Success = $true; Message = "cache-dir al correct ingesteld in pyproject.toml, overgeslagen" }
    }

    # Geval 2: [tool.uv] + andere/incomplete cache-dir → vervangen
    if ($content -match '\[tool\.uv\]' -and $content -match 'cache-dir') {
        $content = $content -replace 'cache-dir\s*=\s*"[^"]*"', $correctValue
        Set-Content $pyproject $content -NoNewline
        return [PSCustomObject]@{ Success = $true; Message = "Bestaande cache-dir in pyproject.toml vervangen door correcte waarde" }
    }

    # Geval 3: [tool.uv] zonder cache-dir → toevoegen onder de sectie
    if ($content -match '\[tool\.uv\]') {
        $content = $content -replace '(\[tool\.uv\])', "`$1`n$correctValue"
        Set-Content $pyproject $content -NoNewline
        return [PSCustomObject]@{ Success = $true; Message = "cache-dir toegevoegd aan [tool.uv] sectie in pyproject.toml" }
    }

    # Geval 4: [tool.uv] bestaat niet → bovenaan toevoegen
    $insert = "[tool.uv]`n$correctValue`n`n"
    $content = $insert + $content
    Set-Content $pyproject $content -NoNewline
    return [PSCustomObject]@{ Success = $true; Message = "[tool.uv] sectie met cache-dir toegevoegd aan pyproject.toml" }
}

# ── Lees optional-dependencies uit pyproject.toml, filter 'dev' eruit ──
function Get-UvExtras {
    param([string]$Root)

    $pyproject = Join-Path $Root "pyproject.toml"
    if (-not (Test-Path $pyproject)) { return @() }

    $content = Get-Content $pyproject -Raw

    if ($content -notmatch '\[project\.optional-dependencies\]') { return @() }

    $extras = @()
    $inSection = $false

    foreach ($line in ($content -split "`n")) {
        if ($line -match '^\[project\.optional-dependencies\]') {
            $inSection = $true
            continue
        }
        if ($inSection -and $line -match '^\[') { break }
        if ($inSection -and $line -match '^\s*(\w+)\s*=\s*\[') {
            $name = $matches[1]
            if ($name -ne "dev") { $extras += $name }
        }
    }

    return $extras
}
