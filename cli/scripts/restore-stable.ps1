#!/usr/bin/env pwsh
<#
.SYNOPSIS
    Restore the stable version of azd exec extension
.DESCRIPTION
    Uninstalls PR build, removes PR registries, and installs latest stable version
.EXAMPLE
    .\restore-stable.ps1
.EXAMPLE
    iex "& { $(irm https://raw.githubusercontent.com/jongio/azd-exec/main/cli/scripts/restore-stable.ps1) }"
#>

$ErrorActionPreference = 'Stop'
$repo = "jongio/azd-exec"
$extensionId = "jongio.azd.exec"
$stableRegistryUrl = "https://raw.githubusercontent.com/$repo/refs/heads/main/registry.json"

Write-Host "🔄 Restoring stable azd exec extension" -ForegroundColor Cyan
Write-Host ""

# Step 1: Uninstall current extension
Write-Host "🗑️  Uninstalling current extension..." -ForegroundColor Gray
azd extension uninstall $extensionId 2>$null
# Ignore errors - extension might not be installed

# Step 2: Remove all PR registry sources
Write-Host "🧹 Removing PR registry sources..." -ForegroundColor Gray
$sources = azd extension source list 2>$null | Select-String "^pr-\d+"
foreach ($source in $sources) {
    $sourceName = $source -replace '\s+.*$'
    Write-Host "   Removing: $sourceName" -ForegroundColor DarkGray
    azd extension source remove $sourceName 2>&1 | Out-Null
}

# Step 3: Clean up pr-registry.json files
Write-Host "🧹 Cleaning up pr-registry.json files..." -ForegroundColor Gray
@("./pr-registry.json", "~/pr-registry.json", "$HOME/pr-registry.json") | ForEach-Object {
    if (Test-Path $_) {
        Remove-Item -Force $_
    }
}

# Step 4: Add stable registry source
Write-Host "🔗 Adding stable registry source..." -ForegroundColor Gray
# Remove if exists
azd extension source remove "azd-exec" 2>&1 | Out-Null
azd extension source add -n "azd-exec" -t url -l $stableRegistryUrl
if ($LASTEXITCODE -ne 0) {
    Write-Host "❌ Failed to add stable registry source" -ForegroundColor Red
    exit 1
}

# Step 5: Install stable version
Write-Host "📦 Installing latest stable version..." -ForegroundColor Gray
azd extension install $extensionId
if ($LASTEXITCODE -ne 0) {
    Write-Host "❌ Failed to install stable version" -ForegroundColor Red
    exit 1
}

# Step 6: Verify installation
Write-Host ""
Write-Host "✅ Restoration complete!" -ForegroundColor Green
Write-Host ""
Write-Host "🔍 Verifying installation..." -ForegroundColor Gray
$installedVersion = azd exec version 2>&1
if ($LASTEXITCODE -eq 0) {
    Write-Host "   $installedVersion" -ForegroundColor DarkGray
    Write-Host ""
    Write-Host "✨ Success! Stable version restored." -ForegroundColor Green
} else {
    Write-Host "⚠️  Could not verify version" -ForegroundColor Yellow
}

Write-Host ""
Write-Host "Try these commands:" -ForegroundColor Cyan
Write-Host "  azd exec version" -ForegroundColor White
Write-Host "  azd exec ./my-script.sh" -ForegroundColor White
