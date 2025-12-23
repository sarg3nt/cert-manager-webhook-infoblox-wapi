# Commit Message Guidelines for AI Assistants

This repository uses [Conventional Commits](https://www.conventionalcommits.org/) for automated semantic versioning and changelog generation.

## Format

All commit messages MUST follow this structure:

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

## Commit Types (REQUIRED)

Choose the appropriate type based on the change:

- **feat**: A new feature (minor version bump, e.g., 1.0.0 → 1.1.0)
- **fix**: A bug fix (patch version bump, e.g., 1.0.0 → 1.0.1)
- **docs**: Documentation only changes (no version bump)
- **style**: Code style/formatting changes (no version bump)
- **refactor**: Code restructuring without changing behavior (no version bump)
- **perf**: Performance improvements (patch version bump)
- **test**: Adding or updating tests (no version bump)
- **chore**: Maintenance, dependency updates (no version bump)
- **ci**: CI/CD configuration changes (no version bump)
- **build**: Build system or external dependency changes (no version bump)
- **revert**: Reverting a previous commit (no version bump)
- **security**: Security fixes (patch version bump)

## Scope (OPTIONAL)

Scope provides additional context about what part of the codebase is affected:

- `helm` - Helm chart changes
- `webhook` - Webhook core logic
- `auth` - Authentication/authorization
- `api` - API changes
- `deps` - Dependency updates
- `ci` - CI/CD changes
- `docs` - Documentation

Examples: `feat(helm):`, `fix(webhook):`, `chore(deps):`

## Breaking Changes (MAJOR VERSION)

For breaking changes that require a major version bump (e.g., 1.0.0 → 2.0.0):

**Option 1:** Add `!` after type/scope:
```
feat!: change webhook API endpoint structure
```

**Option 2:** Add `BREAKING CHANGE:` in the footer:
```
feat: update authentication method

BREAKING CHANGE: The old token-based auth is no longer supported.
Users must migrate to certificate-based authentication.
```

## Description (REQUIRED)

- Use imperative mood ("add" not "added" or "adds")
- Start with lowercase
- No period at the end
- Be concise but descriptive (50-72 characters)

## Body (OPTIONAL)

- Provide more detailed explanation if needed
- Wrap at 72 characters
- Explain what and why, not how
- Separate from description with blank line

## Footer (OPTIONAL)

- Reference GitHub issues: `Closes #123`, `Fixes #456`, `Refs #789`
- Note breaking changes: `BREAKING CHANGE: <description>`
- Multiple footers allowed

## Examples

### Good Examples

```
feat: add support for custom DNS zone configuration

Allows users to specify custom DNS zones in the webhook configuration.
This enables more flexible DNS management across different Infoblox instances.

Closes #42
```

```
fix(webhook): resolve timeout issue during DNS validation

The webhook was timing out for large zone files. Increased the timeout
from 30s to 60s and added retry logic with exponential backoff.

Fixes #156
```

```
feat(helm)!: change default service port from 443 to 8443

BREAKING CHANGE: The default service port has changed. Users must update
their firewall rules and ingress configurations to use port 8443.
```

```
docs: improve installation instructions in README

Added clarification about Kubernetes version requirements and
troubleshooting steps for common installation issues.
```

```
chore(deps): update Go dependencies to latest versions

- Bumped cert-manager to v1.14.0
- Updated infoblox-go-client to v2.7.1
- Updated k8s.io dependencies to v0.29.0
```

```
ci: add automated security scanning with Trivy
```

```
refactor(webhook): simplify DNS record creation logic

Extracted common functionality into helper functions to improve
readability and maintainability. No functional changes.
```

### Bad Examples (DO NOT USE)

```
❌ Updated files
   (Too vague, no type)

❌ fix: Fixed the bug.
   (Not imperative, has period, not descriptive)

❌ Added new feature for DNS
   (No type prefix)

❌ feat: added support for custom zones.
   (Past tense "added", has period)

❌ FIX: Resolve timeout
   (Type should be lowercase)

❌ feat : add feature
   (Extra space before colon)
```

## Special Cases

### Dependency Updates
```
chore(deps): update Go to 1.25.5
```

### Multiple Changes
Create separate commits for each logical change:
```
feat: add DNS zone validation
fix: resolve webhook timeout issue
docs: update README with new examples
```

### Merge Commits
For merge commits, ensure the PR title follows conventional commits as it becomes the merge commit message.

## AI Assistant Instructions

When generating commit messages:

1. **Analyze the changes** to determine the appropriate type
2. **Use scope** when changes are localized to a specific component
3. **Keep description concise** but meaningful
4. **Add body** for complex changes that need explanation
5. **Reference issues** in footer when applicable
6. **Mark breaking changes** explicitly with `!` or `BREAKING CHANGE:`
7. **Use imperative mood** consistently
8. **Validate format** before committing

## Validation

Commit messages are validated in CI. Invalid messages will fail the build.

## Questions?

See the [README.md](../README.md#contributions) for more details on the contribution process.
