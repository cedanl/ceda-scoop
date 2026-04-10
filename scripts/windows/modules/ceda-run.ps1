# ceda-run.ps1
# Voert de volledige CEDA run-flow uit voor een geïnstalleerd project.
# Aanroep vanuit Go: powershell.exe -File ceda-run.ps1 -Root <installatiepad> -Step <stap>
#
# Werkt per stap zodat Go de voortgang kan tonen.
# Stappen: scoop-check | core-deps | buckets | detect | [r of uv stappen] | open

param(
    [string]$Root,
    [string]$Step
)

$ErrorActionPreference = "Continue"
$ModulesPath = $PSScriptRoot

# ── Modules laden ──────────────────────────────────────────────────────────────
. "$ModulesPath\scoop-install.ps1"
. "$ModulesPath\detect-project.ps1"
. "$ModulesPath\uv-install.ps1"
. "$ModulesPath\uv-config.ps1"
. "$ModulesPath\uv-sync.ps1"
. "$ModulesPath\uv-run.ps1"
. "$ModulesPath\r-install.ps1"
. "$ModulesPath\r-config.ps1"
. "$ModulesPath\r-sync.ps1"
. "$ModulesPath\r-run.ps1"

# ── Helper: resultaat evalueren ────────────────────────────────────────────────
function Invoke-Step {
    param([scriptblock]$Action)
    $result = & $Action
    if ($result -and ($result.Success -eq $false)) {
        Write-Error $result.Message
        exit 1
    }
    if ($result -and $result.Message) {
        Write-Host $result.Message
    }
}

# ── Stap dispatcher ────────────────────────────────────────────────────────────
switch ($Step) {

    "scoop-check" {
        Invoke-Step { Install-Scoop }
    }

    "core-deps" {
        Invoke-Step { Install-CoreDeps }
    }

    "buckets" {
        # Vang "al bestaat" fout op per bucket
        try { scoop bucket add extras  2>&1 | Out-Null } catch {}
        try { scoop bucket add versions 2>&1 | Out-Null } catch {}
        Write-Host "Buckets OK"
    }

    "detect" {
        # Detecteer projecttype en schrijf naar stdout zodat Go het kan lezen
        $project = Get-ProjectType -Root $Root
        if ($project.Type -eq "unknown") {
            Write-Error $project.Message
            exit 1
        }
        Write-Host "TYPE:$($project.Type)"
        Write-Host $project.Message
    }

    # ── uv flow ────────────────────────────────────────────────────────────────

    "uv-install" {
        Invoke-Step { Install-UV }
    }

    "uv-config" {
        Invoke-Step { Set-UvCacheDir -Root $Root }
    }

    "uv-sync" {
        $extras   = Get-UvExtras -Root $Root
        $syncArgs = @()
        foreach ($e in $extras) { $syncArgs += "--extra"; $syncArgs += $e }
        Invoke-Step { Invoke-UvSync -Root $Root -SyncArgs $syncArgs }
    }

    "uv-run" {
        $isStreamlit = Test-IsStreamlitProject -Root $Root
        Invoke-Step { Start-UvProject -Root $Root -IsStreamlit $isStreamlit }
    }

    # ── R flow ─────────────────────────────────────────────────────────────────

    "r-install" {
        Invoke-Step { Install-RDeps }
    }

    "r-paths" {
        Invoke-Step { Set-RtoolsPaths }
    }

    "r-positron" {
        Invoke-Step { Set-PositronRConfig }
    }

    "r-sync" {
        Invoke-Step { Invoke-RenvRestore -Root $Root }
    }

    "r-run" {
        Invoke-Step { Start-RProject -Root $Root }
    }

    default {
        Write-Error "Onbekende stap: $Step"
        exit 1
    }
}

exit 0
