# uv-install.ps1
# Installeert uv via Scoop als het nog niet aanwezig is.
# Exporteert: Install-UV

function Install-UV {
    try {
        if (Get-Command uv -ErrorAction SilentlyContinue) {
            return [PSCustomObject]@{ Success = $true; Message = "uv al aanwezig, overgeslagen" }
        }

        $output = scoop install uv 2>&1
        if ($LASTEXITCODE -ne 0) {
            return [PSCustomObject]@{ Success = $false; Message = "Fout bij installatie van uv: $output" }
        }

        return [PSCustomObject]@{ Success = $true; Message = "uv geinstalleerd via Scoop" }
    } catch {
        return [PSCustomObject]@{ Success = $false; Message = "Onverwachte fout bij uv installatie: $_" }
    }
}
