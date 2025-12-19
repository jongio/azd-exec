#!/usr/bin/env pwsh
# PowerShell env test
Write-Host "=== Environment Test ==="
Write-Host "PATH exists: $($null -ne $env:PATH)"
Write-Host "HOME exists: $($null -ne $env:HOME)"
Write-Host "=== Script Complete ==="
exit 0
