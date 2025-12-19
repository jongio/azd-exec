#!/bin/bash
# Multi-line bash script with env vars
echo "=== Environment Test ==="
echo "PATH exists: ${PATH:+yes}"
echo "HOME exists: ${HOME:+yes}"
echo "=== Script Complete ==="
exit 0
