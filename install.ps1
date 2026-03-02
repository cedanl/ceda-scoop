<#
.SYNOPSIS
    Install script for r-scoop-template using Scoop.

.DESCRIPTION
    This script installs all required dependencies for the R project using Scoop,
    a command-line installer for Windows. It will:
      1. Install Scoop if not already present
      2. Install R via Scoop
      3. Install Rtools via Scoop (needed to build R packages from source)
      4. Restore R package dependencies using renv

.NOTES
    Run this script once after cloning the repository.
    Requires Windows PowerShell 5.1 or PowerShell 7+.

.EXAMPLE
    .\install.ps1
#>

#Requires -Version 5.1

Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"

# ── Helper functions ──────────────────────────────────────────────────────────

function Write-Step {
    param([string]$Message)
    Write-Host ""
    Write-Host "==> $Message" -ForegroundColor Cyan
}

function Write-Success {
    param([string]$Message)
    Write-Host "  ✓ $Message" -ForegroundColor Green
}

function Write-Warning {
    param([string]$Message)
    Write-Host "  ! $Message" -ForegroundColor Yellow
}

# ── 1. Ensure Scoop is installed ──────────────────────────────────────────────

Write-Step "Checking Scoop installation"

if (-not (Get-Command scoop -ErrorAction SilentlyContinue)) {
    Write-Host "  Scoop not found. Installing Scoop..." -ForegroundColor Yellow

    # Allow running scripts for the current user
    Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser -Force

    Invoke-RestMethod -Uri "https://get.scoop.sh" | Invoke-Expression

    # Reload PATH so scoop is available in this session
    $env:PATH = [System.Environment]::GetEnvironmentVariable("PATH", "User") +
                ";" + [System.Environment]::GetEnvironmentVariable("PATH", "Machine")

    if (-not (Get-Command scoop -ErrorAction SilentlyContinue)) {
        Write-Error "Scoop installation failed. Please install it manually: https://scoop.sh"
        exit 1
    }
    Write-Success "Scoop installed"
} else {
    Write-Success "Scoop is already installed"
}

# ── 2. Add required Scoop buckets ─────────────────────────────────────────────

Write-Step "Adding Scoop buckets"

$buckets = scoop bucket list 2>$null | ForEach-Object { $_.Trim() }

foreach ($bucket in @("main", "extras")) {
    if ($buckets -notcontains $bucket) {
        scoop bucket add $bucket
        Write-Success "Added bucket: $bucket"
    } else {
        Write-Success "Bucket already present: $bucket"
    }
}

# ── 3. Install R ──────────────────────────────────────────────────────────────

Write-Step "Installing R"

if (Get-Command Rscript -ErrorAction SilentlyContinue) {
    $rVersion = (Rscript --version 2>&1) -replace ".*version\s+", ""
    Write-Success "R is already installed (version $rVersion)"
} else {
    scoop install r
    Write-Success "R installed"
}

# ── 4. Install Rtools ─────────────────────────────────────────────────────────

Write-Step "Installing Rtools (required to build R packages from source)"

$rtoolsInstalled = scoop list 2>$null | Where-Object { $_ -match "^rtools" }

if ($rtoolsInstalled) {
    Write-Success "Rtools is already installed"
} else {
    scoop install rtools
    Write-Success "Rtools installed"
}

# Reload PATH to pick up R and Rtools
$env:PATH = [System.Environment]::GetEnvironmentVariable("PATH", "User") +
            ";" + [System.Environment]::GetEnvironmentVariable("PATH", "Machine")

# ── 5. Restore R packages with renv ──────────────────────────────────────────

Write-Step "Restoring R packages via renv"

if (-not (Test-Path "renv.lock")) {
    Write-Warning "No renv.lock file found. Skipping package restoration."
} else {
    Rscript -e "if (!requireNamespace('renv', quietly = TRUE)) install.packages('renv'); renv::restore()"
    Write-Success "R packages restored"
}

# ── Done ──────────────────────────────────────────────────────────────────────

Write-Host ""
Write-Host "Installation complete!" -ForegroundColor Green
Write-Host "Open 'r-scoop-template.Rproj' in RStudio to get started." -ForegroundColor Cyan
