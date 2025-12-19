# GitHub Actions Workflows

This directory contains all GitHub Actions workflows for the cert-manager-webhook-infoblox-wapi project.

## Workflows Overview

### CI/CD Workflows

#### [ci.yml](ci.yml)
**Purpose:** Continuous Integration for Go code  
**Triggers:** Push to main, Pull Requests  
**Jobs:**
- **Lint:** Runs golangci-lint, go vet, gofmt, and go mod verification
- **Test:** Executes unit tests with race detector and generates coverage reports
- **Build:** Builds Docker image to verify build process

**Secrets Required:** None (uses `GITHUB_TOKEN`)

#### [helm.yml](helm.yml)
**Purpose:** Lint and test Helm charts  
**Triggers:** Push/PR affecting `charts/**`  
**Jobs:**
- **Lint:** Runs `helm lint` on all charts
- **Test:** Deploys charts to a Kind cluster to verify installation

**Secrets Required:** None

#### [shellcheck.yml](shellcheck.yml)
**Purpose:** Validate shell scripts  
**Triggers:** Push/PR affecting `**.sh` files  
**Jobs:** Runs ShellCheck on all bash scripts

**Secrets Required:** None

### Release Workflows

#### [release.yml](release.yml)
**Purpose:** Build and publish releases from tags  
**Triggers:** Git tags (e.g., `v1.0.0`)  
**Jobs:**
- **Build and Push:** Builds and pushes Docker image with semantic versioning tags
- **Helm Version Updater:** Creates PRs to update Helm chart versions

**Secrets Required:**
- `GITHUB_TOKEN` (automatic)
- `USER_PAT` (for creating PRs that trigger workflows)

#### [release-monthly.yml](release-monthly.yml)
**Purpose:** Automated monthly releases for dependency updates  
**Triggers:** Monthly cron (1st of month at 00:00 UTC), manual  
**Jobs:**
- **Release Build and Push:** Builds new image only if content differs from previous
- **Helm Version Updater:** Updates charts if image was released

**Secrets Required:**
- `GITHUB_TOKEN`
- `USER_PAT`

#### [deploy-chart-on-pr-close.yml](deploy-chart-on-pr-close.yml)
**Purpose:** Deploy Helm charts to gh-pages after PR merge  
**Triggers:** PR close events affecting `charts/**`  
**Jobs:** Runs chart-releaser to publish charts

**Secrets Required:**
- `USER_PAT`

### Security Workflows

#### [codeql.yml](codeql.yml)
**Purpose:** Static code analysis for security vulnerabilities  
**Triggers:** Push to main, PRs, weekly schedule, manual  
**Jobs:** Runs CodeQL analysis on Go code

**Secrets Required:** None

#### [trivy.yml](trivy.yml)
**Purpose:** Container vulnerability scanning  
**Triggers:** Push to main, PRs affecting Docker/Go files, weekly schedule  
**Jobs:** Scans Docker images for CVEs and uploads results to Security tab

**Secrets Required:** None

#### [scorecard.yml](scorecard.yml)
**Purpose:** OpenSSF Scorecard security analysis  
**Triggers:** Push to main, branch protection changes, weekly schedule  
**Jobs:** Analyzes repository security posture

**Secrets Required:** None (uses `GITHUB_TOKEN`)

#### [dependency-review.yml](dependency-review.yml)
**Purpose:** Review dependencies in PRs for vulnerabilities  
**Triggers:** Pull Requests  
**Jobs:** Scans dependency changes for known vulnerabilities

**Secrets Required:** None

## Concurrency Control

All workflows use concurrency groups to prevent overlapping runs:
- Most workflows cancel in-progress runs when new commits are pushed
- Release workflows do NOT cancel to ensure releases complete

## Caching Strategy

- **Go modules:** Cached via `actions/setup-go` with `cache: true`
- **Docker layers:** Cached via GitHub Actions cache (`type=gha`)
- **Kubebuilder test assets:** Downloaded and cached per workflow run

## Security Hardening

All workflows use:
- **step-security/harden-runner:** Network egress control and auditing
- **SHA-pinned actions:** All actions pinned to specific commit SHAs
- **Minimal permissions:** Each job declares only required permissions

## Maintenance

- **Dependabot:** Automatically updates action versions monthly (see `.github/dependabot.yml`)
- **Version comments:** All pinned actions include version comments for readability

## Local Testing

To test workflows locally, you can use [act](https://github.com/nektos/act):

```bash
# Install act
brew install act  # macOS
# or download from https://github.com/nektos/act/releases

# Test CI workflow
act pull_request -W .github/workflows/ci.yml

# Test with specific secrets
act -s GITHUB_TOKEN=your_token
```

## Troubleshooting

### Workflow fails on go.mod verification
Run locally: `go mod tidy` and commit changes

### golangci-lint fails
Run locally: `golangci-lint run`  
Fix issues or update `.golangci.yml` to exclude false positives

### Docker build fails
Check Dockerfile syntax and ensure all dependencies are available

### Helm lint fails
Run locally: `helm lint charts/cert-manager-webhook-infoblox-wapi/`

## Adding New Workflows

When creating new workflows:

1. Add concurrency control
2. Use step-security/harden-runner
3. Pin all actions to SHA with version comments
4. Set minimal permissions
5. Add caching where applicable
6. Update this README
7. Test with `act` if possible
