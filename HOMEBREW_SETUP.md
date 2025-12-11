# Homebrew Setup Guide

This guide explains how to set up the Homebrew tap for bc-cli to enable easy installation on macOS.

## Prerequisites

- GitHub organization: `ButlerCoffee`
- Repository access to create a new repository

## Step 1: Create the Homebrew Tap Repository

1. Create a new repository named `homebrew-tap` in the ButlerCoffee organization:
   ```
   https://github.com/ButlerCoffee/homebrew-tap
   ```

2. Initialize the repository with a README:
   ```bash
   # Clone the new repository
   git clone https://github.com/ButlerCoffee/homebrew-tap.git
   cd homebrew-tap

   # Create the Formula directory
   mkdir -p Formula

   # Copy the formula from bc-cli repo
   cp /path/to/bc-cli/homebrew/bc-cli.rb Formula/

   # Commit and push
   git add .
   git commit -m "Initial commit with bc-cli formula"
   git push origin main
   ```

## Step 2: Create a GitHub Personal Access Token (PAT)

The release workflow needs permission to push updates to the homebrew-tap repository.

1. Go to GitHub Settings → Developer settings → Personal access tokens → Tokens (classic)
2. Click "Generate new token (classic)"
3. Give it a descriptive name: `bc-cli-homebrew-tap-updater`
4. Set expiration as needed (recommend: 1 year)
5. Select the following scopes:
   - `repo` (Full control of private repositories)
6. Generate the token and copy it

## Step 3: Add the Token as a Repository Secret

1. Go to the bc-cli repository settings
2. Navigate to Secrets and variables → Actions
3. Click "New repository secret"
4. Name: `HOMEBREW_TAP_TOKEN`
5. Value: Paste the PAT from Step 2
6. Click "Add secret"

## Step 4: Test the Setup

When you push commits to main that trigger a semantic-release, the workflow will:

1. Build binaries for multiple platforms
2. Create a GitHub release with all binaries
3. Automatically update the Homebrew formula in `homebrew-tap` with:
   - New version number
   - Updated SHA256 checksums for both Intel and Apple Silicon
   - Updated download URLs

## How Users Install

Once set up, users can install bc-cli with:

```bash
# Add the tap
brew tap butlercoffee/tap

# Install bc-cli
brew install bc-cli

# Or in one command
brew install butlercoffee/tap/bc-cli
```

## Updating the Formula

The formula is automatically updated on each release. You don't need to manually update it.

If you need to make manual changes to the formula:

1. Clone the homebrew-tap repository
2. Edit `Formula/bc-cli.rb`
3. Test locally: `brew install --build-from-source ./Formula/bc-cli.rb`
4. Commit and push changes

## Troubleshooting

### Formula update fails in workflow

Check that:
- The `HOMEBREW_TAP_TOKEN` secret is set correctly
- The PAT has the `repo` scope
- The PAT hasn't expired
- The homebrew-tap repository exists and has a `main` branch

### Local testing

Test the formula locally before releasing:

```bash
# Audit the formula
brew audit --strict Formula/bc-cli.rb

# Test installation
brew install --build-from-source Formula/bc-cli.rb

# Test the installed binary
bc-cli --version

# Uninstall when done testing
brew uninstall bc-cli
```

## Repository Structure

Your homebrew-tap repository should look like this:

```
homebrew-tap/
├── README.md
└── Formula/
    └── bc-cli.rb
```

The `Formula/` directory must be at the root level for Homebrew to recognize it.
