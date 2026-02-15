#!/usr/bin/env pwsh
<#
.SYNOPSIS
    Install a PR build of the azd exec extension
.DESCRIPTION
    Uninstalls existing extension, downloads PR registry, and installs the PR build
.PARAMETER PrNumber
    The PR number (e.g., 123)
.PARAMETER Version
    The PR version (e.g., 0.1.0-pr123)
.EXAMPLE
    .\install-pr.ps1 -PrNumber 123 -Version 0.1.0-pr123
.EXAMPLE
    iex "& { $(irm https://raw.githubusercontent.com/jongio/azd-exec/main/cli/scripts/install-pr.ps1) } -PrNumber 123 -Version 0.1.0-pr123"
#>

param(
    [Parameter(Mandatory=$true)]
    [int]$PrNumber,
    
    [Parameter(Mandatory=$true)]
    [string]$Version
)

$ErrorActionPreference = 'Stop'
$repo = "jongio/azd-exec"
$extensionId = "jongio.azd.exec"
$tag = "azd-ext-jongio-azd-exec_${Version}"
$registryUrl = "https://github.com/$repo/releases/download/$tag/pr-registry.json"

Write-Host "🚀 Installing azd exec PR #$PrNumber (version $Version)" -ForegroundColor Cyan
Write-Host ""

# Step 1: Uninstall existing extension
Write-Host "🗑️  Uninstalling existing extension (if any)..." -ForegroundColor Gray
azd extension uninstall $extensionId 2>&1 | Out-Null
# Ignore errors - extension might not be installed
Write-Host "   ✓" -ForegroundColor DarkGray

# Step 2: Download PR registry
Write-Host "📥 Downloading PR registry..." -ForegroundColor Gray
$registryPath = Join-Path $PWD "pr-registry.json"
try {
    Invoke-WebRequest -Uri $registryUrl -OutFile $registryPath
    Write-Host "   ✓ Downloaded to: $registryPath" -ForegroundColor DarkGray
} catch {
    Write-Host "❌ Failed to download registry from $registryUrl" -ForegroundColor Red
    Write-Host "   Make sure the PR build exists and is accessible" -ForegroundColor Yellow
    exit 1
}

# Step 3: Add registry source
Write-Host "🔗 Adding PR registry source..." -ForegroundColor Gray
azd extension source remove "pr-$PrNumber" 2>$null  # Remove if exists
azd extension source add -n "pr-$PrNumber" -t file -l $registryPath
if ($LASTEXITCODE -ne 0) {
    Write-Host "❌ Failed to add registry source" -ForegroundColor Red
    exit 1
}

# Step 4: Install PR version
Write-Host "📦 Installing version $Version..." -ForegroundColor Gray
azd extension install $extensionId --version $Version
if ($LASTEXITCODE -ne 0) {
    Write-Host "❌ Failed to install extension" -ForegroundColor Red
    exit 1
}

# Step 5: Verify installation
Write-Host ""
Write-Host "✅ Installation complete!" -ForegroundColor Green
Write-Host ""
Write-Host "🔍 Verifying installation..." -ForegroundColor Gray
$installedVersion = azd exec version 2>&1
if ($LASTEXITCODE -eq 0) {
    Write-Host "   $installedVersion" -ForegroundColor DarkGray
    if ($installedVersion -match $Version) {
        Write-Host ""
        Write-Host "✨ Success! PR build is ready to test." -ForegroundColor Green
    } else {
        Write-Host ""
        Write-Host "⚠️  Version mismatch - expected $Version" -ForegroundColor Yellow
    }
} else {
    Write-Host "⚠️  Could not verify version" -ForegroundColor Yellow
}

Write-Host ""
Write-Host "Try these commands:" -ForegroundColor Cyan
Write-Host "  azd exec version" -ForegroundColor White
Write-Host "  azd exec ./my-script.sh" -ForegroundColor White
Write-Host ""
Write-Host "To restore stable version, run:" -ForegroundColor Gray
Write-Host "  iex `"& { `$(irm https://raw.githubusercontent.com/$repo/main/cli/scripts/restore-stable.ps1) }`"" -ForegroundColor White
