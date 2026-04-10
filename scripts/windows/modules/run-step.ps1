# run-step.ps1
# Wrapper die door ceda-scoop (Go) wordt aangeroepen per stap.
# Gebruik: powershell.exe -File run-step.ps1 -Step <naam> -Root <pad>
#
# Returnt exit 0 bij succes, exit 1 bij fout (met Write-Error voor context).

param(
    [string]$Step,
    [string]$Root
)

$ErrorActionPreference = "Stop"
$ModulesPath = $PSScriptRoot

# ── Modules laden ──
. "$ModulesPath\scoop-install.ps1"
. "$ModulesPath\uv-install.ps1"
. "$ModulesPath\uv-config.ps1"
. "$ModulesPath\uv-sync.ps1"
. "$ModulesPath\uv-run.ps1"
. "$ModulesPath\r-install.ps1"
. "$ModulesPath\r-config.ps1"
. "$ModulesPath\r-sync.ps1"
. "$ModulesPath\r-run.ps1"

function Invoke-Step {
    param([string]$Label, [scriptblock]$Action)
    $result = & $Action
    if ($result -and ($result.Success -eq $false)) {
        Write-Error $result.Message
        exit 1
    }
    if ($result -and $result.Message) {
        Write-Host $result.Message
    }
}

switch ($Step) {

    # ── Scoop ──
    "scoop-check" {
        Invoke-Step "Scoop controleren" { Install-Scoop }
    }
    "core-deps" {
        Invoke-Step "Core dependencies" { Install-CoreDeps }
    }
    "buckets" {
        # Add-Buckets geeft exit 1 als bucket al bestaat — vang dat op
        try {
            scoop bucket add extras  2>&1 | Out-Null
        } catch {}
        try {
            scoop bucket add versions 2>&1 | Out-Null
        } catch {}
        Write-Host "Buckets OK"
    }

    # ── uv ──
    "uv-install" {
        Invoke-Step "uv installeren" { Install-UV }
    }
    "uv-config" {
        Invoke-Step "pyproject.toml configureren" { Set-UvCacheDir -Root $Root }
    }
    "uv-sync" {
        $extras   = Get-UvExtras -Root $Root
        $syncArgs = @()
        foreach ($e in $extras) { $syncArgs += "--extra"; $syncArgs += $e }
        Invoke-Step "uv sync" { Invoke-UvSync -Root $Root -SyncArgs $syncArgs }
    }
    "uv-run" {
        $isStreamlit = Test-IsStreamlitProject -Root $Root
        Invoke-Step "project starten" { Start-UvProject -Root $Root -IsStreamlit $isStreamlit }
    }

    # ── R ──
    "r-install" {
        Invoke-Step "R dependencies installeren" { Install-RDeps }
    }
    "r-paths" {
        Invoke-Step "Rtools paden instellen" { Set-RtoolsPaths }
    }
    "r-positron" {
        Invoke-Step "Positron configureren" { Set-PositronRConfig }
    }
    "r-sync" {
        Invoke-Step "renv restore" { Invoke-RenvRestore -Root $Root }
    }
    "r-run" {
        Invoke-Step "Positron openen" { Start-RProject -Root $Root }
    }

    default {
        Write-Error "Onbekende stap: $Step"
        exit 1
    }
}

exit 0
