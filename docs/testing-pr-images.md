# Testing PR Images on a Cluster

Every pull request opened from a branch in this repository gets a multi-arch (`linux/amd64`, `linux/arm64`) container image built from its head commit. Use it to try a change on a real cluster before it merges. The [pr-image.yml](../.github/workflows/pr-image.yml) workflow builds the image, pushes it, and comments the image reference on the PR.

For unit and integration tests that run locally or in CI, see [TESTING.md](TESTING.md).

## Where the Images Live

PR images are published to a separate, public package so they never mix with release images:

```text
ghcr.io/sarg3nt/cert-manager-webhook-infoblox-wapi-pr
```

Each PR gets two tags:

| Tag             | Meaning                                                          |
| --------------- | ---------------------------------------------------------------- |
| `pr-<N>`        | Moves to the newest push on PR `<N>`                             |
| `pr-<N>-<sha7>` | Immutable; one per pushed head commit. Prefer this for deploying |

The workflow keeps a single "PR test image" comment on the PR up to date with the current tags, digest, and a ready-to-run `helm` command.

## Deploying a PR Image

1. Add the chart repository if you haven't already:

   ```bash
   helm repo add cert-manager-webhook-infoblox-wapi https://sarg3nt.github.io/cert-manager-webhook-infoblox-wapi
   helm repo update
   ```

2. Point an existing release at the PR image, keeping your other values:

   ```bash
   helm upgrade --install webhook-infoblox-wapi \
     cert-manager-webhook-infoblox-wapi/cert-manager-webhook-infoblox-wapi \
     -n cert-manager --reuse-values \
     --set image.repository=ghcr.io/sarg3nt/cert-manager-webhook-infoblox-wapi-pr \
     --set image.tag=pr-123-abc1234
   ```

3. Confirm the pod is running the PR image:

   ```bash
   kubectl -n cert-manager get pods -l app=cert-manager-webhook-infoblox-wapi \
     -o jsonpath='{.items[*].spec.containers[*].image}'
   ```

If the PR changes the Helm chart itself, install from `charts/cert-manager-webhook-infoblox-wapi` in a checkout of the PR branch instead of the published chart, with the same `--set` flags.

If you deploy the moving `pr-<N>` tag instead of the immutable one, also pass `--set image.pullPolicy=Always` and restart the deployment after each push so the node pulls the new image.

## Cleanup

When the PR is merged or closed, the workflow's cleanup job deletes the `pr-<N>` and `pr-<N>-*` tags along with their per-platform manifests. Switch any cluster still running a PR image back to a release image before (or soon after) merging, or its pods will fail to pull on the next reschedule.

## Pull Requests That Don't Get Images

GitHub gives pull requests from forks and from Dependabot a read-only `GITHUB_TOKEN`, so they can't push packages and the build job is skipped. The workflow deliberately doesn't use `pull_request_target` to get around this, because that would build untrusted code with a write token.

To build an image for a Dependabot PR (or any open PR whose branch lives in this repository):

1. Open **Actions → PR Image → Run workflow**.
2. Select the PR's head branch under **Use workflow from**.
3. Enter the PR number and run.

The run checks that the PR is open and that its head branch matches the selected branch. The images are tagged `pr-<N>` as usual and deleted when the PR closes.

## One-Time Package Setup

The `-pr` package is created by the first PR build. Because the workflow publishes it with the repository's `GITHUB_TOKEN`, it is linked to this public repository and inherits its public visibility, so clusters can pull without an `imagePullSecret`. To confirm, or if the package is ever recreated:

1. Open the package at **github.com/sarg3nt?tab=packages → cert-manager-webhook-infoblox-wapi-pr → Package settings**.
2. Check that visibility is **Public**. If it isn't, change it under **Danger Zone → Change visibility**; GitHub has no API for this.
3. Under **Manage Actions access**, confirm `cert-manager-webhook-infoblox-wapi` has the **Admin** role. The cleanup job needs Admin to delete package versions.
