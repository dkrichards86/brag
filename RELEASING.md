# Release Process

This document describes how to create a new release of brag.

## Quick Start

To create a new release:

```bash
make release VERSION=v1.0.0
```

This will:
1. Create a git tag
2. Push the tag to GitHub
3. Trigger the release workflow to build and publish binaries

## Detailed Process

### 1. Prepare for Release

Before creating a release:

- Ensure all tests pass: `make test`
- Run the linter: `make lint`
- Update version information if needed
- Commit all changes

### 2. Create the Release

Use semantic versioning (vMAJOR.MINOR.PATCH):

```bash
# For a new feature
make release VERSION=v1.1.0

# For a bug fix
make release VERSION=v1.0.1

# For breaking changes
make release VERSION=v2.0.0
```

### 3. Automated Build Process

Once you push the tag, GitHub Actions will automatically:

1. **Run tests** - Ensures everything works
2. **Build binaries** for:
   - Linux (AMD64, ARM64)
   - macOS (AMD64, ARM64)
   - Windows (AMD64)
3. **Generate checksums** - SHA256 sums for verification
4. **Create GitHub Release** with:
   - All binaries attached
   - Checksums file
   - Auto-generated changelog from commits
   - Installation instructions

### 4. Verify the Release

After the workflow completes:

1. Visit https://github.com/YOUR_USERNAME/brag/releases
2. Verify all binaries are attached
3. Check the checksums file
4. Test installation instructions

## Manual Release (Alternative)

If you need to create a release manually without the Makefile:

```bash
# Create and push tag
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0

# The GitHub Actions workflow will handle the rest
```

## Release Workflow Details

The release workflow ([.github/workflows/release.yml](.github/workflows/release.yml)) includes:

- **Trigger**: Runs on tags matching `v*.*.*` pattern
- **Tests**: Runs full test suite before building
- **Builds**: Cross-compiles for all supported platforms
- **Checksums**: Generates SHA256 checksums for verification
- **Release Notes**: Auto-generates changelog from commit messages

## Binary Naming Convention

Binaries follow this naming pattern:
```
brag-{OS}-{ARCH}[.exe]
```

Examples:
- `brag-linux-amd64`
- `brag-darwin-arm64`
- `brag-windows-amd64.exe`

## Troubleshooting

### Tag Already Exists

If you get "Tag already exists" error:

```bash
# List existing tags
git tag -l

# Delete local tag
git tag -d v1.0.0

# Delete remote tag (use carefully!)
git push origin :refs/tags/v1.0.0
```

### Workflow Failed

If the GitHub Actions workflow fails:

1. Check the Actions tab on GitHub
2. Review the error logs
3. Fix the issue
4. Delete and recreate the tag:
   ```bash
   git tag -d v1.0.0
   git push origin :refs/tags/v1.0.0
   make release VERSION=v1.0.0
   ```

### Missing Binaries

If binaries are missing from the release:

1. Check that the workflow completed successfully
2. Verify the build matrix in [.github/workflows/release.yml](.github/workflows/release.yml)
3. Check build logs for errors

## Best Practices

1. **Test Before Release**: Always run `make ci` locally first
2. **Meaningful Versions**: Follow semantic versioning
3. **Clean Commits**: Write clear commit messages (they appear in release notes)
4. **Tag Messages**: Include a summary in the tag message
5. **Verify After Release**: Download and test at least one binary

## Version Numbering

Follow [Semantic Versioning](https://semver.org/):

- **MAJOR** (v2.0.0): Breaking changes
- **MINOR** (v1.1.0): New features, backward compatible
- **PATCH** (v1.0.1): Bug fixes, backward compatible

Examples:
- `v1.0.0` - Initial stable release
- `v1.1.0` - Added search command
- `v1.1.1` - Fixed search bug
- `v2.0.0` - Changed file format (breaking)

## Pre-releases

For beta versions, use this format:

```bash
git tag -a v1.0.0-beta.1 -m "Release v1.0.0-beta.1"
git push origin v1.0.0-beta.1
```

The workflow will automatically mark these as "pre-release" on GitHub.
