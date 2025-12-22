#!/usr/bin/env pwsh
# Build script called by azd x build
# This script is invoked by `azd x build` with specific environment variables:
# - EXTENSION_ID: The extension identifier (e.g., jongio.azd.exec)
# - EXTENSION_VERSION: The extension version (e.g., 0.1.0)
# - GOOS: Target operating system
# - GOARCH: Target architecture
# - OUTPUT_PATH: Where to write the binary

$ErrorActionPreference = 'Stop'

# Get environment variables set by azd x build
$extensionId = $env:EXTENSION_ID
$extensionVersion = $env:EXTENSION_VERSION
$targetOS = $env:GOOS
$targetArch = $env:GOARCH
$outputPath = $env:OUTPUT_PATH

Write-Host "DEBUG: GOOS=$targetOS GOARCH=$targetArch OUTPUT_PATH=$outputPath" -ForegroundColor Yellow

# Write debug to file for troubleshooting
"GOOS=$targetOS GOARCH=$targetArch OUTPUT_PATH=$outputPath" | Out-File -FilePath "build-debug.txt" -Append

if (-not $extensionId) {
    Write-Host "ERROR: EXTENSION_ID environment variable not set" -ForegroundColor Red
    exit 1
}

if (-not $extensionVersion) {
    Write-Host "ERROR: EXTENSION_VERSION environment variable not set" -ForegroundColor Red
    exit 1
}

# Get the directory of the script (cli folder)
$extensionDir = Split-Path -Parent $MyInvocation.MyCommand.Path

# Change to the script directory
Set-Location -Path $extensionDir

# Build metadata
$buildDate = (Get-Date).ToUniversalTime().ToString("yyyy-MM-ddTHH:mm:ssZ")
$gitCommit = try {
    git rev-parse --short HEAD 2>$null
} catch {
    "unknown"
}

# Build flags with version info
$ldflags = "-X main.version=$extensionVersion -X main.buildDate=$buildDate -X main.gitCommit=$gitCommit"

Write-Host "Building $extensionId v$extensionVersion" -ForegroundColor Cyan

# If OUTPUT_PATH is set, this is a targeted build for azd x build
if ($outputPath) {
    # Default OS/Arch if not explicitly set by azd x build
    if (-not $targetOS) {
        $targetOS = if ($IsWindows) { "windows" } elseif ($IsLinux) { "linux" } elseif ($IsMacOS) { "darwin" } else { "windows" }
    }
    if (-not $targetArch) {
        $targetArch = if ([System.Environment]::Is64BitOperatingSystem) { "amd64" } else { "386" }
    }
    
    Write-Host "  OS/Arch: $targetOS/$targetArch" -ForegroundColor Gray
    Write-Host "  Output: $outputPath" -ForegroundColor Gray
    
    $env:GOOS = $targetOS
    $env:GOARCH = $targetArch
    
    go build -ldflags $ldflags -o $outputPath ./src/cmd/script
    
    if ($LASTEXITCODE -ne 0) {
        Write-Host "ERROR: Build failed" -ForegroundColor Red
        exit $LASTEXITCODE
    }
    
    # IMPORTANT: azd x pack expects binaries with platform-specific names
    # Copy the binary with the correct naming convention for pack to find it
    $extensionIdSafe = $extensionId -replace '\.', '-'

    $binaryExt = if ($targetOS -eq "windows") { ".exe" } else { "" }
    $platformSpecificName = "$extensionIdSafe-$targetOS-$targetArch$binaryExt"
    
    $binDir = Split-Path -Parent $outputPath
    $platformSpecificPath = Join-Path $binDir $platformSpecificName
    
    Write-Host "  Creating platform-specific copy: $platformSpecificName" -ForegroundColor Magenta
    Write-Host "    From: $outputPath" -ForegroundColor DarkGray
    Write-Host "    To: $platformSpecificPath" -ForegroundColor DarkGray
    
    if (Test-Path $outputPath) {
        Copy-Item $outputPath $platformSpecificPath -Force
        Write-Host "    ✓ Copy successful" -ForegroundColor Green
    } else {
        Write-Host "    ✗ Source file not found!" -ForegroundColor Red
    }
    
    Write-Host "✅ Build successful!" -ForegroundColor Green
} else {
    # Fallback: build for current platform to bin/
    Write-Host "  Building for current platform..." -ForegroundColor Gray
    
    $binDir = Join-Path $extensionDir "bin"
    New-Item -ItemType Directory -Force -Path $binDir | Out-Null
    
    $binaryName = "exec"
    if ($env:GOOS -eq "windows" -or $IsWindows) {
        $binaryName += ".exe"
    }
    
    $outputPath = Join-Path $binDir $binaryName
    
    go build -ldflags $ldflags -o $outputPath ./src/cmd/script
    
    if ($LASTEXITCODE -ne 0) {
        Write-Host "ERROR: Build failed" -ForegroundColor Red
        exit $LASTEXITCODE
    }
    
    Write-Host "✅ Build successful: $outputPath" -ForegroundColor Green
}
