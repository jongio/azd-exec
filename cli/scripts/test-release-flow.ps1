#!/usr/bin/env pwsh
# Test script for azd x release flow
# This simulates what the GitHub Actions workflow will do

param(
    [string]$Version = "0.2.0-test",
    [switch]$SkipBuild,
    [switch]$SkipPack
)

$ErrorActionPreference = 'Stop'

Write-Host "🧪 Testing azd x release flow" -ForegroundColor Cyan
Write-Host "Version: $Version" -ForegroundColor Yellow
Write-Host ""

# Check if azd is installed
Write-Host "Checking azd installation..." -ForegroundColor Gray
try {
    $azdVersion = azd version 2>&1
    Write-Host "✅ azd installed: $azdVersion" -ForegroundColor Green
} catch {
    Write-Host "❌ azd not installed. Install from https://aka.ms/azd" -ForegroundColor Red
    exit 1
}

# Check if microsoft.azd.extensions is installed
Write-Host "Checking azd extensions..." -ForegroundColor Gray
try {
    $extensions = azd extension list --output json | ConvertFrom-Json
    $hasExtension = $extensions | Where-Object { $_.id -eq "microsoft.azd.extensions" }
    if ($hasExtension) {
        Write-Host "✅ microsoft.azd.extensions installed" -ForegroundColor Green
    } else {
        Write-Host "Installing microsoft.azd.extensions..." -ForegroundColor Yellow
        azd extension install microsoft.azd.extensions
    }
} catch {
    Write-Host "⚠️  Could not verify extensions" -ForegroundColor Yellow
}

# Navigate to cli directory
$cliDir = Join-Path $PSScriptRoot ".."
Set-Location $cliDir

Write-Host ""
Write-Host "=== Step 1: Build Website ===" -ForegroundColor Cyan
if (Test-Path "../web") {
    Set-Location "../web"
    if (Test-Path "package.json") {
        Write-Host "Running pnpm install..." -ForegroundColor Gray
        pnpm install
        Write-Host "Running pnpm build..." -ForegroundColor Gray
        pnpm run build
        Write-Host "✅ Website built" -ForegroundColor Green
    } else {
        Write-Host "⚠️  No package.json found, skipping website build" -ForegroundColor Yellow
    }
    Set-Location "../cli"
} else {
    Write-Host "⚠️  No web directory found" -ForegroundColor Yellow
}

Write-Host ""
Write-Host "=== Step 2: Build Extension Binaries ===" -ForegroundColor Cyan
if (-not $SkipBuild) {
    $env:EXTENSION_ID = "jongio.azd.exec"
    $env:EXTENSION_VERSION = $Version
    $env:BUILD_ALL = "true"
    
    Write-Host "Building for all platforms..." -ForegroundColor Gray
    Write-Host "  EXTENSION_ID=$env:EXTENSION_ID" -ForegroundColor DarkGray
    Write-Host "  EXTENSION_VERSION=$env:EXTENSION_VERSION" -ForegroundColor DarkGray
    Write-Host "  BUILD_ALL=$env:BUILD_ALL" -ForegroundColor DarkGray
    
    azd x build --all
    
    if (Test-Path "bin") {
        $binaries = Get-ChildItem "bin" -Recurse -Filter "exec*"
        Write-Host "✅ Built $($binaries.Count) binaries:" -ForegroundColor Green
        $binaries | ForEach-Object { 
            $relativePath = $_.FullName.Replace((Get-Location).Path, "").TrimStart('\', '/')
            Write-Host "   - $relativePath" -ForegroundColor DarkGray 
        }
    } else {
        Write-Host "❌ No binaries found in bin/" -ForegroundColor Red
        exit 1
    }
} else {
    Write-Host "⏭️  Skipped (--SkipBuild)" -ForegroundColor Yellow
}

Write-Host ""
Write-Host "=== Step 3: Package Extension ===" -ForegroundColor Cyan
if (-not $SkipPack) {
    Write-Host "Packaging extension..." -ForegroundColor Gray
    
    # Temporarily update extension.yaml version for testing
    $extensionYaml = Get-Content "extension.yaml" -Raw
    $originalVersion = ($extensionYaml -match 'version:\s*(.+)') ? $Matches[1] : "0.2.0"
    Write-Host "  Original version in extension.yaml: $originalVersion" -ForegroundColor DarkGray
    Write-Host "  Using test version: $Version" -ForegroundColor DarkGray
    
    # Backup and update
    Copy-Item "extension.yaml" "extension.yaml.bak"
    (Get-Content "extension.yaml") -replace "version:\s*.+", "version: $Version" | Set-Content "extension.yaml"
    
    try {
        azd x pack
    } finally {
        # Restore original version
        Move-Item "extension.yaml.bak" "extension.yaml" -Force
        Write-Host "  Restored extension.yaml to version $originalVersion" -ForegroundColor DarkGray
    }
    
    $registryPath = Join-Path $env:USERPROFILE ".azd\registry\jongio.azd.exec\$Version"
    if (Test-Path $registryPath) {
        $archives = Get-ChildItem $registryPath -Include "*.zip","*.tar.gz" -File
        Write-Host "✅ Packaged $($archives.Count) archives:" -ForegroundColor Green
        $archives | ForEach-Object { Write-Host "   - $($_.Name)" -ForegroundColor DarkGray }
        
        Write-Host ""
        Write-Host "Registry location:" -ForegroundColor Gray
        Write-Host "  $registryPath" -ForegroundColor DarkGray
    } else {
        Write-Host "❌ No packages found at $registryPath" -ForegroundColor Red
        Write-Host ""
        Write-Host "Expected location:" -ForegroundColor Yellow
        Write-Host "  $registryPath" -ForegroundColor DarkGray
        exit 1
    }
} else {
    Write-Host "⏭️  Skipped (--SkipPack)" -ForegroundColor Yellow
}

Write-Host ""
Write-Host "=== Test Summary ===" -ForegroundColor Cyan
Write-Host "✅ All steps completed successfully!" -ForegroundColor Green
Write-Host ""
Write-Host "Artifacts location:" -ForegroundColor Yellow
Write-Host "  $registryPath" -ForegroundColor Gray
Write-Host ""
Write-Host "To clean up this test:" -ForegroundColor Yellow
Write-Host "  Remove-Item -Recurse -Force '$registryPath'" -ForegroundColor Gray
