package runner

import (
	"fmt"
	"os"
	"path/filepath"
)

// ProjectType geeft aan of een repo R of uv is.
type ProjectType string

const (
	ProjectTypeR       ProjectType = "r"
	ProjectTypeUV      ProjectType = "uv"
	ProjectTypeUnknown ProjectType = "unknown"
)

// RunStep beschrijft één stap in de run-flow.
type RunStep struct {
	Label string
	Cmd   string // PowerShell commando
}

// DetectProjectType kijkt naar lockfiles in de installmap.
// Geeft (type, ambiguous) terug. ambiguous = true als beide aanwezig.
func DetectProjectType(installPath string) (ProjectType, bool) {
	hasR := fileExists(filepath.Join(installPath, "renv.lock"))
	hasUV := fileExists(filepath.Join(installPath, "uv.lock"))
	if hasR && hasUV {
		return ProjectTypeUnknown, true
	}
	if hasR {
		return ProjectTypeR, false
	}
	if hasUV {
		return ProjectTypeUV, false
	}
	return ProjectTypeUnknown, false
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// BuildRunSteps geeft de geordende lijst van stappen terug voor een projecttype.
func BuildRunSteps(installPath string, projectType ProjectType) []RunStep {
	root := psPath(installPath)

	common := []RunStep{
		{
			Label: "Scoop controleren",
			Cmd:   `if (-not (Get-Command scoop -ErrorAction SilentlyContinue)) { Invoke-RestMethod -Uri 'https://get.scoop.sh' | Invoke-Expression }`,
		},
		{
			Label: "Core dependencies installeren",
			Cmd:   `@("git","7zip","aria2") | ForEach-Object { if (-not (scoop list 2>&1 | Select-String "^\s*$_\s")) { scoop install $_ 2>&1 | Out-Null } }`,
		},
		{
			Label: "Scoop buckets toevoegen",
			Cmd:   `scoop bucket add extras 2>&1 | Out-Null; scoop bucket add versions 2>&1 | Out-Null`,
		},
	}

	var specific []RunStep

	switch projectType {
	case ProjectTypeR:
		specific = []RunStep{
			{
				Label: "R, Rtools en Positron installeren",
				Cmd:   `@("main/r","main/rtools","extras/positron") | ForEach-Object { $n=$_.Split("/")[-1]; if (-not (scoop list 2>&1 | Select-String "^\s*$n\s")) { scoop install $_ 2>&1 | Out-Null } }`,
			},
			{
				Label: "Rtools paden instellen",
				Cmd: `$paths=@("$env:USERPROFILE\scoop\apps\rtools\current\usr\bin","$env:USERPROFILE\scoop\apps\rtools\current\x86_64-w64-mingw32.static.posix\bin");` +
					`$cur=[System.Environment]::GetEnvironmentVariable("PATH","User");` +
					`foreach($p in $paths){if($cur -notlike "*$p*"){$cur="$cur;$p"}};` +
					`[System.Environment]::SetEnvironmentVariable("PATH",$cur,"User")`,
			},
			{
				Label: "Positron R interpreter instellen",
				Cmd: `$dir="$env:APPDATA\Positron\User";$file="$dir\settings.json";` +
					`$rPath=("$env:USERPROFILE\scoop\apps\r" -replace '\\','\\');` +
					`if(-not(Test-Path $dir)){New-Item -ItemType Directory $dir -Force|Out-Null};` +
					`if(-not(Test-Path $file)){Set-Content $file '{}' -Encoding UTF8};` +
					`$c=Get-Content $file -Raw;` +
					`if($c -notmatch 'customRootFolders'){$c=$c -replace '(\{)',"`$1`n`t""positron.r.customRootFolders"": [""$rPath""],";Set-Content $file $c -Encoding UTF8}`,
			},
			{
				Label: "R packages installeren via renv",
				Cmd: fmt.Sprintf(
					`$env:PATH=[System.Environment]::GetEnvironmentVariable("PATH","User")+";"+$env:PATH;`+
						`$r="$env:USERPROFILE\scoop\apps\r\current\bin\x64\R.exe";`+
						`if(-not(Test-Path $r)){Write-Error "R.exe niet gevonden";exit 1};`+
						`$tmp=Join-Path $env:TEMP "ceda-renv.R";`+
						`[System.IO.File]::WriteAllText($tmp,"renv::restore(prompt=FALSE)",[System.Text.UTF8Encoding]::new($false));`+
						`$p=Start-Process $r -ArgumentList "--no-save","--no-restore","--file=$tmp" -WorkingDirectory "%s" -Wait -PassThru -NoNewWindow;`+
						`Remove-Item $tmp -ErrorAction SilentlyContinue;`+
						`if($p.ExitCode -ne 0){Write-Error "renv restore mislukt";exit 1}`,
					root),
			},
			{
				Label: "Project openen in Positron",
				Cmd:   fmt.Sprintf(`Start-Process positron -ArgumentList "--disable-workspace-trust","%s"`, root),
			},
		}

	case ProjectTypeUV:
		specific = []RunStep{
			{
				Label: "uv installeren",
				Cmd:   `if(-not(Get-Command uv -ErrorAction SilentlyContinue)){scoop install uv 2>&1|Out-Null}`,
			},
			{
				Label: "pyproject.toml configureren",
				Cmd: fmt.Sprintf(
					`$f="%s\pyproject.toml";if(Test-Path $f){`+
						`$c=Get-Content $f -Raw;`+
						`if($c -notmatch 'cache-dir'){`+
						`$c="[tool.uv]`ncache-dir = ""./.uv_cache""`n`n"+$c;`+
						`Set-Content $f $c -NoNewline}}`,
					root),
			},
			{
				Label: "Packages installeren via uv sync",
				Cmd:   fmt.Sprintf(`Push-Location "%s"; uv sync; $e=$LASTEXITCODE; Pop-Location; if($e -ne 0){Write-Error "uv sync mislukt";exit 1}`, root),
			},
			{
				Label: "Project starten",
				Cmd: fmt.Sprintf(
					`Push-Location "%s";`+
						`$entry=$null;`+
						`if(Test-Path "src\main.py"){$entry="src\main.py"}`+
						`elseif(Test-Path "main.py"){$entry="main.py"}`+
						`else{$f=Get-ChildItem -Path "src" -Filter "*.py" -EA SilentlyContinue|Select-Object -First 1;if($f){$entry="src\$($f.Name)"}};`+
						`if(-not $entry){Pop-Location;Write-Error "Geen entry-point gevonden";exit 1};`+
						`$sl=(Get-Content "pyproject.toml" -Raw -EA SilentlyContinue) -match "streamlit";`+
						`if($sl){Start-Process powershell -ArgumentList "-NoExit","-Command","uv run streamlit run $entry"}`+
						`else{Start-Process powershell -ArgumentList "-NoExit","-Command","cd '%s'; uv run $entry"};`+
						`Pop-Location`,
					root, root),
			},
		}
	}

	return append(common, specific...)
}

// ExecuteStep voert één stap synchroon uit en geeft een fout terug of nil.
func ExecuteStep(step RunStep) error {
	_, err := runPSCommand(step.Cmd)
	return err
}

// psPath escaped backslashes voor gebruik in PS strings.
func psPath(path string) string {
	out := ""
	for _, c := range path {
		if c == '\\' {
			out += `\\`
		} else {
			out += string(c)
		}
	}
	return out
}
