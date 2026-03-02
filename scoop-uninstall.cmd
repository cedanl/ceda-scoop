@echo off
:: scoop-first-run.bat
set "psScript=%~dp0ps\scoop-uninstall.ps1"
if not exist "%psScript%" (
    echo Error: scoop-uninstall.ps1 not found in ps folder!
    pause
    exit /b 1
)
powershell -ExecutionPolicy Bypass -File "%psScript%"
pause