# Homebrew Tap Setup Guide

This guide explains how to set up and maintain the Homebrew tap for `ipctl`.

## Overview

Homebrew support is configured in `.goreleaser.yaml` and automatically publishes formula updates to the `itential/homebrew-tap` repository when new releases are created.

## Initial Setup

### 1. Create the Homebrew Tap Repository

Create a new GitHub repository at `github.com/itential/homebrew-tap`:

```bash
# Create the repository on GitHub
gh repo create itential/homebrew-tap --public --description "Homebrew formulae for Itential tools"

# Clone it locally
git clone git@github.com:itential/homebrew-tap.git
cd homebrew-tap

# Create Formula directory
mkdir -p Formula

# Create initial README
cat > README.md << 'EOF'
# Itential Homebrew Tap

Homebrew formulae for Itential tools.

## Installation

```bash
brew tap itential/tap
brew install ipctl
```

## Available Formulae

- `ipctl` - Command-line interface for managing Itential Platform servers
EOF

# Commit and push
git add .
git commit -m "Initial commit"
git push origin main
```

### 2. Create GitHub Personal Access Token

GoReleaser needs a GitHub token with permission to push to the tap repository:

1. Go to GitHub Settings → Developer settings → Personal access tokens → Tokens (classic)
2. Click "Generate new token (classic)"
3. Name it: `HOMEBREW_TAP_GITHUB_TOKEN`
4. Select scopes:
   - `repo` (Full control of private repositories)
5. Click "Generate token"
6. Copy the token (you won't see it again)

### 3. Configure GitHub Actions Secret

Add the token as a secret to the `ipctl` repository:

```bash
# Using GitHub CLI
gh secret set HOMEBREW_TAP_GITHUB_TOKEN -b "ghp_your_token_here" -R itential/ipctl

# Or manually:
# Go to github.com/itential/ipctl → Settings → Secrets and variables → Actions
# Click "New repository secret"
# Name: HOMEBREW_TAP_GITHUB_TOKEN
# Value: paste your token
```

### 4. Create GitHub Actions Workflow

Create `.github/workflows/release.yml` in the `ipctl` repository:

```yaml
name: Release

on:
  push:
    tags:
      - 'v*'

permissions:
  contents: write
  packages: write

jobs:
  release:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout
        uses: actions/checkout@v4
        with:
          fetch-depth: 0

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.23'

      - name: Run GoReleaser
        uses: goreleaser/goreleaser-action@v6
        with:
          version: latest
          args: release --clean
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
          HOMEBREW_TAP_GITHUB_TOKEN: ${{ secrets.HOMEBREW_TAP_GITHUB_TOKEN }}
          BUILD: ${{ github.sha }}
```

## Release Process

When you're ready to release a new version:

```bash
# 1. Ensure all changes are committed
git add .
git commit -m "Prepare release vX.Y.Z"

# 2. Create and push a new tag
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0

# 3. GitHub Actions will automatically:
#    - Build binaries for all platforms
#    - Create GitHub release with artifacts
#    - Update Homebrew formula in itential/homebrew-tap
```

### Manual Release (if needed)

If you need to release manually without GitHub Actions:

```bash
# Set environment variables
export GITHUB_TOKEN="your_github_token"
export HOMEBREW_TAP_GITHUB_TOKEN="your_homebrew_tap_token"
export BUILD=$(git rev-parse --short HEAD)

# Run goreleaser
goreleaser release --clean
```

## Testing the Homebrew Formula

### Before Publishing

Test the formula locally before releasing:

```bash
# Create a local tap
brew tap-new itential/tap
cd $(brew --repository itential/tap)

# Copy the generated formula from homebrew-tap repo
cp /path/to/homebrew-tap/Formula/ipctl.rb Formula/ipctl.rb

# Install locally
brew install --build-from-source itential/tap/ipctl

# Test
ipctl --version
```

### After Publishing

Test the published formula:

```bash
# Remove local tap if you created one
brew untap itential/tap

# Add the tap
brew tap itential/tap

# Install ipctl
brew install ipctl

# Verify
ipctl --version
ipctl --help
```

## Updating the Formula

The formula is automatically updated by GoReleaser when you create a new release. However, if you need to make manual changes:

1. Clone the homebrew-tap repository
2. Edit `Formula/ipctl.rb`
3. Test the changes locally
4. Commit and push

```bash
cd homebrew-tap
vim Formula/ipctl.rb

# Test
brew uninstall ipctl
brew install --build-from-source ./Formula/ipctl.rb

# Commit
git add Formula/ipctl.rb
git commit -m "Update ipctl formula"
git push origin main
```

## Troubleshooting

### Formula Not Found

If users get "formula not found" errors:

```bash
# Update tap
brew update

# Ensure tap is added
brew tap itential/tap
```

### Installation Fails

Check the formula syntax:

```bash
brew audit --strict itential/tap/ipctl
brew style itential/tap/ipctl
```

### Token Permissions

If GoReleaser fails to push to the tap repository:

1. Verify the token has `repo` scope
2. Check the token hasn't expired
3. Ensure the secret is correctly set in GitHub Actions

### Testing Multiple Versions

To test installing different versions:

```bash
# Install specific version
brew install itential/tap/ipctl@1.0.0

# List available versions
brew search ipctl

# Switch versions
brew unlink ipctl
brew link ipctl@1.0.0
```

## Architecture Support

The formula supports:
- macOS (Intel): `darwin/amd64`
- macOS (Apple Silicon): `darwin/arm64`
- Linux (x86_64): `linux/amd64`
- Linux (ARM64): `linux/arm64`

Windows users should download binaries directly from the releases page.

## Resources

- [Homebrew Formula Cookbook](https://docs.brew.sh/Formula-Cookbook)
- [GoReleaser Homebrew Documentation](https://goreleaser.com/customization/homebrew/)
- [GitHub Actions Documentation](https://docs.github.com/en/actions)
