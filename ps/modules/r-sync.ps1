# r-sync.ps1
# Herstelt renv dependencies via pak met PPM binaries.
# Exporteert: Invoke-RenvRestore

function Invoke-RenvRestore {
    param([string]$Root = $PSScriptRoot)

    # ── Check renv aanwezig ──
    if (-not (Test-Path (Join-Path $Root "renv\activate.R"))) {
        return [PSCustomObject]@{ Success = $false; Message = "Geen renv project gevonden (renv/activate.R ontbreekt)" }
    }

    # ── Rtools paden activeren in huidige sessie ──
    $env:PATH = [System.Environment]::GetEnvironmentVariable("PATH", "User") + ";" + $env:PATH

    # ── R restore script via pak + PPM ──
    $rScript = @"
# ── Repo configuratie: PPM binaries + CRAN fallback ──
options(repos = c(
  PPM  = 'https://packagemanager.posit.co/cran/latest',
  CRAN = 'https://cloud.r-project.org'
))
options(renv.config.ppm.enabled = TRUE)
options(pkgType = 'binary')

# ── Parallel cores: totaal - 2, minimaal 1 ──
ncpus <- max(1, parallel::detectCores() - 2)
options(Ncpus = ncpus)

# ── renv cache optimalisaties ──
options(renv.config.cache.symlinks  = TRUE)
options(renv.config.install.shortcuts = TRUE)

# ── Info ──
cat('\n=== Environment ===\n')
cat('R versie     :', as.character(getRversion()), '\n')
cat('Cores        :', parallel::detectCores(), '| Gebruikt:', ncpus, '\n')

# ── renv activeren ──
source('renv/activate.R')

# ── pak installeren als het er niet is ──
if (!requireNamespace('pak', quietly = TRUE)) {
  cat('\npak installeren...\n')
  renv::install('pak', prompt = FALSE)
  cat('pak geinstalleerd\n')
}

# ── Pak exact restore ──
if (!file.exists('renv.lock')) stop('renv.lock niet gevonden')

lock <- renv::lockfile_read()
pkgs <- names(lock[['Packages']])

cat('\n=== Restore ===\n')
cat('Packages     :', length(pkgs), '\n')
cat('Cache        :', renv::paths[['cache']](), '\n\n')

start <- Sys.time()

pkg_specs <- sapply(pkgs, function(pkg) {
  info    <- lock[['Packages']][[pkg]]
  version <- info[['Version']]
  source  <- info[['Source']]

  if (source == 'GitHub') {
    ref <- if (is.null(info[['RemoteRef']])) info[['RemoteSha']] else info[['RemoteRef']]
    sprintf('%s/%s@%s', info[['RemoteUsername']], info[['RemoteRepo']], ref)
  } else if (source == 'Bioconductor') {
    sprintf('bioc::%s@%s', pkg, version)
  } else {
    sprintf('%s@%s', pkg, version)
  }
})

pak::pkg_install(pkg_specs, ask = FALSE)

elapsed <- round(difftime(Sys.time(), start, units = 'secs'), 1)
cat('\nRestore klaar in', elapsed, 'seconden\n')
"@

    # ── Uitvoeren via Rscript ──
    Push-Location $Root
    Rscript --vanilla -e $rScript
    $exitCode = $LASTEXITCODE
    Pop-Location

    if ($exitCode -ne 0) {
        return [PSCustomObject]@{ Success = $false; Message = "renv restore mislukt (exit code $exitCode)" }
    }

    return [PSCustomObject]@{ Success = $true; Message = "renv restore geslaagd via pak" }
}
