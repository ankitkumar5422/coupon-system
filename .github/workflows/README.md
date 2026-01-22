# GitHub Actions Workflows

This directory contains GitHub Actions workflows for the Coupon System project. These workflows automate various aspects of continuous integration, deployment, and code quality.

## Workflows Overview

### 1. Go CI (`go-ci.yml`)

**Trigger:** Push or Pull Request to `main` or `develop` branches

**Purpose:** Builds, tests, and lints the Go application across multiple Go versions.

**Jobs:**
- **Build**: 
  - Runs on multiple Go versions (1.18, 1.19, 1.20)
  - Installs dependencies
  - Verifies dependencies
  - Builds the application
  - Runs `go vet`
  - Runs tests with race detection and coverage
  - Uploads coverage to Codecov (for Go 1.20 only)

- **Lint**:
  - Runs golangci-lint for code quality checks
  - Uses the latest version of golangci-lint

**What it checks:**
- Code compiles successfully
- All tests pass
- No race conditions
- Code follows Go best practices
- Code coverage is tracked

### 2. Docker Build and Push (`docker-image.yml`)

**Trigger:** 
- Push to `main` branch
- Pull Request to `main` branch
- Tags matching `v*`

**Purpose:** Builds and pushes Docker images to GitHub Container Registry.

**Features:**
- Uses Docker Buildx for multi-platform support
- Implements build caching for faster builds
- Automatically tags images based on:
  - Branch name
  - Pull request number
  - Semantic versioning (for tags)
  - Git SHA
  - `latest` tag for default branch
- Only pushes images on non-PR events
- Stores images in GitHub Container Registry (ghcr.io)

**What it does:**
- Builds optimized Docker image with multi-stage build
- Tags images appropriately
- Pushes to GitHub Container Registry (when not a PR)

### 3. Security and Code Quality (`security.yml`)

**Trigger:** 
- Push or Pull Request to `main` or `develop` branches
- Weekly schedule (Sundays at midnight)

**Purpose:** Performs security scanning and code quality checks.

**Jobs:**
- **Security Scan**:
  - Runs Gosec security scanner
  - Uploads results to GitHub Security tab (SARIF format)
  - Checks for common security vulnerabilities

- **Dependency Check**:
  - Runs govulncheck to detect vulnerable dependencies
  - Checks against Go vulnerability database

- **CodeQL Analysis**:
  - Runs GitHub's CodeQL analysis
  - Identifies security vulnerabilities and code quality issues
  - Results available in Security > Code scanning alerts

**What it checks:**
- Security vulnerabilities in code
- Vulnerable dependencies
- Code quality issues
- Common security mistakes

### 4. PR Validation (`pr-validation.yml`)

**Trigger:** Pull Request to `main` or `develop` branches

**Purpose:** Quick validation checks for pull requests.

**Jobs:**
- **Validate**:
  - Checks code formatting with `go fmt`
  - Runs `go vet` for static analysis
  - Builds the application
  - Checks for TODO/FIXME comments (informational)
  - Checks for debug statements (informational)

- **Size Check**:
  - Reports repository statistics
  - Counts Go files and lines of code
  - Shows repository size

**What it does:**
- Provides quick feedback on PRs
- Ensures code is formatted correctly
- Identifies potential issues early

## Setting Up Workflows

All workflows are automatically triggered based on their configuration. No manual setup is required beyond ensuring the `.github/workflows/` directory is committed to the repository.

## Required Secrets

The workflows use GitHub's built-in secrets:
- `GITHUB_TOKEN` - Automatically provided by GitHub Actions
  - Used for pushing Docker images
  - Used for uploading security scan results

## Viewing Results

### Build Status
- Go to the "Actions" tab in your GitHub repository
- Click on any workflow to see its runs
- Click on a specific run to see detailed logs

### Test Coverage
- Coverage reports are uploaded to Codecov (if configured)
- View coverage badges in your README

### Security Alerts
- Go to "Security" > "Code scanning alerts"
- View issues found by Gosec, CodeQL, and govulncheck

### Docker Images
- Go to repository "Packages" section
- Find the `coupon-system` package
- View all available tags and versions

## Best Practices

1. **Before Merging PRs:**
   - Ensure all CI checks pass
   - Review security scan results
   - Check test coverage

2. **For Releases:**
   - Tag your release with semantic versioning (e.g., `v1.0.0`)
   - Docker image will be automatically built and tagged

3. **Security:**
   - Review weekly security scan results
   - Update dependencies when vulnerabilities are found

## Customization

To customize these workflows:

1. Edit the workflow files in `.github/workflows/`
2. Adjust triggers, job steps, or environment variables
3. Commit and push changes
4. Workflows will automatically use the new configuration

## Troubleshooting

### Build Failures
- Check the workflow logs in the Actions tab
- Look for specific error messages in failed steps
- Ensure all dependencies are correctly specified in `go.mod`

### Docker Build Issues
- Verify Dockerfile is correctly configured
- Check that all required files are in the build context
- Review Docker build logs for specific errors

### Security Alerts
- Review the specific vulnerability or issue
- Update code or dependencies as needed
- Re-run security scans to verify fixes

## Additional Resources

- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [Go CI Best Practices](https://github.com/mvdan/github-actions-golang)
- [Docker Build and Push Action](https://github.com/docker/build-push-action)
- [CodeQL Documentation](https://codeql.github.com/docs/)
