#!/usr/bin/env pwsh
# PowerShell script with parameters
param(
    [string]$Name = "World",
    [string]$Greeting = "Hello"
)

Write-Host "$Greeting, $Name!"
Write-Host "Total args: $($args.Count)"
exit 0
