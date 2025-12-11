# Release Process

This project uses [semantic-release](https://github.com/semantic-release/semantic-release) to automate versioning and releases based on conventional commits.

## How It Works

1. **Conventional Commits**: Use conventional commit messages to trigger releases:
   - `feat: description` - Creates a minor version bump (e.g., 1.0.0 → 1.1.0)
   - `fix: description` - Creates a patch version bump (e.g., 1.0.0 → 1.0.1)
   - `feat!: description` or `BREAKING CHANGE:` - Creates a major version bump (e.g., 1.0.0 → 2.0.0)
   - `chore:`, `docs:`, `style:`, `refactor:`, `test:` - No release

2. **Automatic Releases**: When you push to `main` branch:
   - GitHub Actions runs semantic-release
   - If there are commits that trigger a release:
     - Version number is determined based on commit types
     - CHANGELOG.md is updated
     - VERSION file is created
     - Git tag is created (e.g., `v1.0.0`)
     - GitHub release is created with:
       - Release notes generated from commits
       - Binaries for multiple platforms (macOS, Linux, Windows)
       - SHA256 checksums

## Commit Message Format

Follow the [Conventional Commits](https://www.conventionalcommits.org/) specification:

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

### Examples

```bash
# Patch release (1.0.0 → 1.0.1)
git commit -m "fix: resolve authentication token refresh issue"

# Minor release (1.0.0 → 1.1.0)
git commit -m "feat: add support for subscription management"

# Major release (1.0.0 → 2.0.0)
git commit -m "feat!: redesign API client interface"
# or
git commit -m "feat: change authentication flow

BREAKING CHANGE: authentication now requires email verification"

# No release
git commit -m "chore: update dependencies"
git commit -m "docs: improve README documentation"
git commit -m "test: add integration tests for subscriptions"
```

## Release Artifacts

Each release includes pre-built binaries for:
- macOS (Intel and Apple Silicon)
- Linux (amd64 and arm64)
- Windows (amd64)

All binaries include embedded version information:
```bash
./bc-cli --version
# Output:
# bc-cli version 1.0.0
#   git commit: abc1234
#   built: 2024-01-15T10:30:00Z
```

## Local Development

Build with version information locally:
```bash
# Uses git describe to determine version
make compile

# Or specify version manually
VERSION=1.2.3 make compile
```

## Manual Release (if needed)

If you need to trigger a release manually:

1. Ensure you're on the `main` branch
2. Push your conventional commits to GitHub
3. The workflow will automatically run and create the release

## Configuration Files

- `.releaserc.json` - Semantic release configuration
- `.github/workflows/release.yml` - GitHub Actions workflow
- `Makefile` - Build configuration with version injection

## First Release

For the first release, you can use:
```bash
git commit --allow-empty -m "feat: initial release"
git push origin main
```

This will create version `1.0.0` as the first release.
