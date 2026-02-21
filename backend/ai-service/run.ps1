# Run AI service (Windows). Create venv first: python -m venv venv
Set-Location $PSScriptRoot
if (Test-Path venv\Scripts\Activate.ps1) {
    & venv\Scripts\Activate.ps1
}
python main.py
