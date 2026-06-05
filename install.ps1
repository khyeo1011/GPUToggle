#Requires -Version 5.1
<#
.SYNOPSIS
  Builds the GPU Toggle native helper and registers it with Chrome.

.DESCRIPTION
  Run this script once after cloning the repo, and again whenever you
  reload the extension in Chrome (which changes its extension ID).

.EXAMPLE
  .\install.ps1
#>

Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"

$RepoRoot  = $PSScriptRoot
$HelperDir = Join-Path $RepoRoot "helper"
$HelperExe = Join-Path $HelperDir "helper.exe"
$ManifestTemplate = Join-Path $RepoRoot "com.gputoggle.helper.json"

# ── 1. Build the helper ──────────────────────────────────────────────────────
Write-Host "Building helper.exe..." -ForegroundColor Cyan
Push-Location $HelperDir
try {
    go build -o $HelperExe .
    if (-not $?) { throw "go build failed" }
} finally {
    Pop-Location
}
Write-Host "  -> $HelperExe" -ForegroundColor Green

# ── 2. Get extension ID ───────────────────────────────────────────────────────
Write-Host ""
Write-Host "Open Chrome and go to chrome://extensions"
Write-Host "Enable 'Developer mode', click 'Load unpacked', and select:"
Write-Host "  $(Join-Path $RepoRoot 'extension')"
Write-Host ""
$ExtensionId = Read-Host "Paste the Extension ID shown on the card"
$ExtensionId = $ExtensionId.Trim()
if ($ExtensionId -notmatch '^[a-z]{32}$') {
    Write-Warning "That doesn't look like a valid extension ID (32 lowercase letters). Continuing anyway."
}

# ── 3. Write the native messaging manifest ───────────────────────────────────
$InstallDir  = Join-Path $env:LOCALAPPDATA "GPUToggle"
$ManifestDst = Join-Path $InstallDir "com.gputoggle.helper.json"

if (-not (Test-Path $InstallDir)) {
    New-Item -ItemType Directory -Path $InstallDir | Out-Null
}

$manifest = Get-Content $ManifestTemplate -Raw
$manifest = $manifest -replace "PLACEHOLDER_HELPER_PATH", ($HelperExe -replace "\\", "\\")
$manifest = $manifest -replace "PLACEHOLDER_EXTENSION_ID", $ExtensionId

Set-Content -Path $ManifestDst -Value $manifest -Encoding utf8
Write-Host "Manifest written to: $ManifestDst" -ForegroundColor Green

# ── 4. Register in Windows registry ──────────────────────────────────────────
$RegKey = "HKCU:\Software\Google\Chrome\NativeMessagingHosts\com.gputoggle.helper"
New-Item -Path $RegKey -Force | Out-Null
Set-ItemProperty -Path $RegKey -Name "(default)" -Value $ManifestDst
Write-Host "Registry key set: $RegKey" -ForegroundColor Green

# ── 5. Done ───────────────────────────────────────────────────────────────────
Write-Host ""
Write-Host "Done! Now:" -ForegroundColor Cyan
Write-Host "  1. If Chrome was open during install, restart it once (the registry is read at startup)."
Write-Host "  2. Click the GPU Toggle icon in the toolbar."
Write-Host "  3. The popup should show the current hardware acceleration state."
