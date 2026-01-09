---
title: GitHub Actions Workflows Setup
description: Required permissions and settings for GitHub Actions workflows
lastUpdated: 2026-01-09
tags: [github-actions, ci-cd, automation]
---

# GitHub Actions Workflows - Setup Guide

This document describes the required permissions and settings for the GitHub Actions workflows in this repository.

## Repository Settings Required

### 1. Enable GitHub Actions to Create Pull Requests

The release workflow creates a PR for version bumps. This requires enabling a repository setting:

1. Go to **Settings** → **Actions** → **General**
2. Scroll to **Workflow permissions**
3. Select **Read and write permissions**
4. ✅ Check **"Allow GitHub Actions to create and approve pull requests"**
5. Click **Save**

> ⚠️ Without this setting, the release workflow will fail with:
> `GitHub Actions is not permitted to create or approve pull requests`

### 2. Enable Auto-Merge (Optional but Recommended)

The release workflow uses auto-merge to automatically merge the version bump PR once CI passes:

1. Go to **Settings** → **General**
2. Scroll to **Pull Requests**
3. ✅ Check **"Allow auto-merge"**
4. Click **Save**

### 3. Branch Protection Rules (If Enabled)

If you have branch protection rules on `main`, ensure:

- **Require status checks to pass before merging**: The release workflow waits for CI to pass
- **Require pull request reviews before merging**: If enabled, you'll need to use a PAT (see Alternative Setup below)

---

## Secrets Required

| Secret | Required For | How to Get |
|--------|--------------|------------|
| `GITHUB_TOKEN` | All workflows | Automatically provided by GitHub Actions |
| `CODECOV_TOKEN` | Code coverage uploads | Get from [codecov.io](https://codecov.io) after linking your repo |

### Note on GITHUB_TOKEN

The built-in `GITHUB_TOKEN` is used for:
- Creating release PRs
- Creating GitHub releases
- Updating the registry.json
- Enabling auto-merge

No additional PAT (Personal Access Token) is needed if repository settings are configured correctly.

---

## Alternative Setup: Using a Personal Access Token (PAT)

If your repository has stricter branch protection rules (like required reviews), you'll need a PAT:

### Create a Fine-Grained PAT

1. Go to **GitHub** → **Settings** → **Developer settings** → **Personal access tokens** → **Fine-grained tokens**
2. Click **Generate new token**
3. Set:
   - **Token name**: `azd-exec-release`
   - **Expiration**: 90 days (or custom)
   - **Repository access**: Only select repositories → `jongio/azd-exec`
   - **Permissions**:
     - **Contents**: Read and write
     - **Pull requests**: Read and write
     - **Metadata**: Read-only (automatically selected)
4. Click **Generate token** and copy it

### Add as Repository Secret

1. Go to your repo → **Settings** → **Secrets and variables** → **Actions**
2. Click **New repository secret**
3. Name: `RELEASE_PAT`
4. Value: Paste your PAT
5. Click **Add secret**

### Update Workflow

If using a PAT, change the release workflow:

```yaml
- name: Update version files and create PR
  env:
    GH_TOKEN: ${{ secrets.RELEASE_PAT }}  # Changed from GITHUB_TOKEN
```

---

## Workflow Overview

### `ci.yml` - Continuous Integration
- **Triggers**: Pull requests to main, manual dispatch
- **What it does**: Runs preflight checks, tests (all OS), lint, and build
- **Secrets needed**: `CODECOV_TOKEN` (optional, for coverage)

### `release.yml` - Release Workflow
- **Triggers**: Manual dispatch only (workflow_dispatch)
- **What it does**:
  1. Runs preflight, test, lint, build
  2. Creates version bump PR
  3. Waits for CI, auto-merges
  4. Creates git tag
  5. Builds and packages binaries
  6. Creates GitHub release
  7. Updates registry.json
  8. Deploys website
- **Secrets needed**: `GITHUB_TOKEN` (or `RELEASE_PAT` if using PAT method)

### `pr-build.yml` - PR Preview Builds
- **Triggers**: After CI passes, PR labeled, manual dispatch
- **What it does**: Creates pre-release builds for testing PRs
- **Secrets needed**: `GITHUB_TOKEN`

### `website.yml` - Website Deployment
- **Triggers**: Push to main (web/** changes), PR events, manual/called
- **What it does**: Builds and deploys to GitHub Pages, creates PR previews
- **Secrets needed**: `GITHUB_TOKEN`

### `codeql.yml` - Security Scanning
- **Triggers**: Push/PR to main (Go files), weekly schedule
- **What it does**: Runs CodeQL security analysis
- **Secrets needed**: None (uses built-in permissions)

---

## Troubleshooting

### "GitHub Actions is not permitted to create or approve pull requests"

**Solution**: Enable the setting in repository Settings → Actions → General → Workflow permissions → Check "Allow GitHub Actions to create and approve pull requests"

### "Resource not accessible by integration"

**Solution**: Ensure workflow has correct permissions and repository settings allow the action.

### Auto-merge not working

**Solutions**:
1. Enable auto-merge in repository settings
2. Ensure branch protection rules allow merging after checks pass
3. If reviews are required, use a PAT from a user with merge permissions

### PR created but CI doesn't run

**Solution**: The CI workflow triggers on `pull_request` events. Ensure the PR is targeting `main` and paths match.

