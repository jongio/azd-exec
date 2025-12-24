#!/usr/bin/env pwsh
# Local test script to simulate the release workflow

$ErrorActionPreference = "Stop"
Set-Location $PSScriptRoot/cli

Write-Host "🧹 Cleaning up..." -ForegroundColor Cyan
Remove-Item -Recurse -Force bin -ErrorAction SilentlyContinue
Remove-Item -Recurse -Force "$HOME/.azd/registry/jongio.azd.exec" -ErrorAction SilentlyContinue

Write-Host "`n🏗️  Building binaries..." -ForegroundColor Cyan
$platforms = @(
    @{OS='windows'; Arch='amd64'; Ext='.exe'},
    @{OS='windows'; Arch='arm64'; Ext='.exe'},
    @{OS='linux'; Arch='amd64'; Ext=''},
    @{OS='linux'; Arch='arm64'; Ext=''},
    @{OS='darwin'; Arch='amd64'; Ext=''},
    @{OS='darwin'; Arch='arm64'; Ext=''}
)

foreach ($platform in $platforms) {
    $binaryName = "jongio-azd-exec-$($platform.OS)-$($platform.Arch)$($platform.Ext)"
    Write-Host "  Building $binaryName..."
    $env:GOOS = $platform.OS
    $env:GOARCH = $platform.Arch
    $env:CGO_ENABLED = '0'
    & go build -o "bin/$binaryName" ./src/cmd/script
    if ($LASTEXITCODE -ne 0) {
        Write-Host "❌ Build failed for $binaryName" -ForegroundColor Red
        exit 1
    }
}

Remove-Item env:GOOS -ErrorAction SilentlyContinue
Remove-Item env:GOARCH -ErrorAction SilentlyContinue
Remove-Item env:CGO_ENABLED -ErrorAction SilentlyContinue

Write-Host "`n📦 Built binaries:" -ForegroundColor Green
Get-ChildItem bin | ForEach-Object {
    $size = [math]::Round($_.Length / 1MB, 2)
    Write-Host "  ✅ $($_.Name) ($size MB)"
}

Write-Host "`n📦 Packaging with azd x pack..." -ForegroundColor Cyan
$packOutput = azd x pack 2>&1 | Out-String
Write-Host $packOutput

if ($packOutput -match "ERROR|error" -and $packOutput -notmatch "SUCCESS") {
    Write-Host "❌ azd x pack failed" -ForegroundColor Red
    exit 1
}

Write-Host "`n🔍 Checking registry..." -ForegroundColor Cyan
$registryPath = "$HOME/.azd/registry/jongio.azd.exec"
if (Test-Path $registryPath) {
    $versions = Get-ChildItem $registryPath -Directory
    foreach ($version in $versions) {
        Write-Host "`n  Version: $($version.Name)" -ForegroundColor Yellow
        $versionPath = $version.FullName
        $files = Get-ChildItem $versionPath -Recurse
        if ($files.Count -eq 0) {
            Write-Host "  ❌ EMPTY DIRECTORY - THIS IS THE BUG!" -ForegroundColor Red
            exit 1
        } else {
            Write-Host "  ✅ Found $($files.Count) files:" -ForegroundColor Green
            $files | ForEach-Object {
                $size = if ($_.PSIsContainer) { "DIR" } else { "$([math]::Round($_.Length / 1KB, 1)) KB" }
                Write-Host "    - $($_.Name) ($size)"
            }
        }
    }
} else {
    Write-Host "  ❌ Registry path not found: $registryPath" -ForegroundColor Red
    exit 1
}

Write-Host "`n✅ All checks passed! Release workflow should work." -ForegroundColor Green
