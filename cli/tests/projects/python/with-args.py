#!/usr/bin/env python3
# Python script with arguments
import sys

print("Script arguments test")
print(f"Script name: {sys.argv[0]}")
if len(sys.argv) > 1:
    print(f"Arg 1: {sys.argv[1]}")
if len(sys.argv) > 2:
    print(f"Arg 2: {sys.argv[2]}")
print(f"Total args: {len(sys.argv) - 1}")
sys.exit(0)
