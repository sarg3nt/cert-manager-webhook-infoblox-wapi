# Cert Manager Webhook Infoblox Wapi Release Notes

## [1.13.1](https://github.com/sarg3nt/cert-manager-webhook-infoblox-wapi/compare/v1.13.0...v1.13.1) (2026-06-22)


### Bug Fixes

* **ci:** allow api.deps.dev egress and submit base-branch snapshot ([889273b](https://github.com/sarg3nt/cert-manager-webhook-infoblox-wapi/commit/889273b663668f77be18439bf425c79e0c858199))
* **ci:** allow githubusercontent egress in Dependency Review job ([de6df83](https://github.com/sarg3nt/cert-manager-webhook-infoblox-wapi/commit/de6df83eaac9ebe6e47848247ea2b8e3f5def9ec))
* **ci:** allow release-assets.githubusercontent.com egress in Test job ([0e05594](https://github.com/sarg3nt/cert-manager-webhook-infoblox-wapi/commit/0e05594a9d18f5a30829e2c068c6c2c0e5ebe167))
* **ci:** retry dependency-review on missing snapshot warnings ([68d5495](https://github.com/sarg3nt/cert-manager-webhook-infoblox-wapi/commit/68d5495a76c720fd9d15d4516d61fbf1b9d60182))
* **ci:** submit Go dependency snapshot for PR head commit ([d35dd04](https://github.com/sarg3nt/cert-manager-webhook-infoblox-wapi/commit/d35dd04b31e0ffbc3c4ac060fed1e188ca7d44aa))
* **deps:** bump controller-runtime to v0.24.1 for k8s 0.36 compat ([8582b24](https://github.com/sarg3nt/cert-manager-webhook-infoblox-wapi/commit/8582b245bb9a18e16eb1f5be2bdd50587c5ebb81))
* **deps:** bump controller-runtime to v0.24.1 for k8s 0.36 compatibility ([991d1b1](https://github.com/sarg3nt/cert-manager-webhook-infoblox-wapi/commit/991d1b1978a0001edcb277f9479b3b3994af5206))
* **dns:** prevent unnecessary deletion of existing TXT records ([0e6b67e](https://github.com/sarg3nt/cert-manager-webhook-infoblox-wapi/commit/0e6b67ea7ee82e5f0d39e03c41f89eae8c22d119))
* **dns:** prevent unnecessary deletion of existing TXT records ([6ce2c04](https://github.com/sarg3nt/cert-manager-webhook-infoblox-wapi/commit/6ce2c048ebc25c077faba75ee30c3a3f94f9c21f))

## [1.13.0](https://github.com/sarg3nt/cert-manager-webhook-infoblox-wapi/compare/v1.12.0...v1.13.0) (2025-12-23)


### Features

* Add CI workflow for Go project with linting, testing, and Docker build ([8a3d0ed](https://github.com/sarg3nt/cert-manager-webhook-infoblox-wapi/commit/8a3d0ed4235bfdea23883828689041fb17af4852))
* add commit message guidelines and release process documentation ([985e46b](https://github.com/sarg3nt/cert-manager-webhook-infoblox-wapi/commit/985e46bfa84190d707da0e03353d34330491a4d1))
* Add golangci-lint configuration file for enhanced linting ([8a3d0ed](https://github.com/sarg3nt/cert-manager-webhook-infoblox-wapi/commit/8a3d0ed4235bfdea23883828689041fb17af4852))
* Add Helm workflow for linting and testing Helm charts ([8a3d0ed](https://github.com/sarg3nt/cert-manager-webhook-infoblox-wapi/commit/8a3d0ed4235bfdea23883828689041fb17af4852))
* Add ShellCheck workflow for shell script linting ([8a3d0ed](https://github.com/sarg3nt/cert-manager-webhook-infoblox-wapi/commit/8a3d0ed4235bfdea23883828689041fb17af4852))
* Add workflow to deploy Helm charts on PR close ([8a3d0ed](https://github.com/sarg3nt/cert-manager-webhook-infoblox-wapi/commit/8a3d0ed4235bfdea23883828689041fb17af4852))


### Bug Fixes

* Update main_test.go to skip integration tests if TEST_ZONE_NAME is not set and enable conformance tests ([8a3d0ed](https://github.com/sarg3nt/cert-manager-webhook-infoblox-wapi/commit/8a3d0ed4235bfdea23883828689041fb17af4852))

## [1.12.0](https://github.com/sarg3nt/cert-manager-webhook-infoblox-wapi/compare/v1.11.0...v1.12.0) (2025-12-23)


### Features

* Add CI workflow for Go project with linting, testing, and Docker build ([8a3d0ed](https://github.com/sarg3nt/cert-manager-webhook-infoblox-wapi/commit/8a3d0ed4235bfdea23883828689041fb17af4852))
* add commit message guidelines and release process documentation ([985e46b](https://github.com/sarg3nt/cert-manager-webhook-infoblox-wapi/commit/985e46bfa84190d707da0e03353d34330491a4d1))
* Add golangci-lint configuration file for enhanced linting ([8a3d0ed](https://github.com/sarg3nt/cert-manager-webhook-infoblox-wapi/commit/8a3d0ed4235bfdea23883828689041fb17af4852))
* Add Helm workflow for linting and testing Helm charts ([8a3d0ed](https://github.com/sarg3nt/cert-manager-webhook-infoblox-wapi/commit/8a3d0ed4235bfdea23883828689041fb17af4852))
* Add ShellCheck workflow for shell script linting ([8a3d0ed](https://github.com/sarg3nt/cert-manager-webhook-infoblox-wapi/commit/8a3d0ed4235bfdea23883828689041fb17af4852))
* Add workflow to deploy Helm charts on PR close ([8a3d0ed](https://github.com/sarg3nt/cert-manager-webhook-infoblox-wapi/commit/8a3d0ed4235bfdea23883828689041fb17af4852))


### Bug Fixes

* Update main_test.go to skip integration tests if TEST_ZONE_NAME is not set and enable conformance tests ([8a3d0ed](https://github.com/sarg3nt/cert-manager-webhook-infoblox-wapi/commit/8a3d0ed4235bfdea23883828689041fb17af4852))

## v1.6.0

### Release Summary

-  First release of this fork of https://github.com/luisico/cert-manager-webhook-infoblox-wapi
-  Add ability to pass Infoblox username and password via a Volume Mount from the host OS file system.  
   In some use cases this can be more secure or preferred than using a secret from Kubernetes.
   See the [README.md](README.md#hostpath-volume-mount) for instructions.
-  Update many package dependencies.
-  Added OpenSSF Scorecard, CodeQL, Trivy, Weekly builds and Dependabot as Github Actions.
-  Improved [README.md](README.md) substantially.

### Dependency Changes

- Upgraded Go 1.16 to 1.23.2
- Upgraded Alpine 3.14 to 3.20
- Upgraded github.com/infobloxopen/infoblox-go-client/v2 v2.0.0 to v2.7.0
- Upgraded github.com/jetstack/cert-manager v1.5.4 to github.com/cert-manager/cert-manager v1.13.3
- Upgraded k8s.io/apiextensions-apiserver v0.21.3 to v0.28.1
- Upgraded k8s.io/apimachinery v0.21.3 to v0.28.1
- Upgraded k8s.io/client-go v0.21.3 to v0.28.1

## v1.7.0 

- Add zone, ttl and useTtl as options to the plugin by @sarg3nt in #3
  NOTE: zone is not currently documented as it may be read only in Infoblox. Not sure if this is a problem with our environment or truly not writable in the Infoblox API. We will test more and update as needed.
- Improve Dockerfile to reduce image size.  No longer contains Alpine OS, just the Go binary and ca-certificates.


## V1.8.0

- Improved logging.
- Troubleshooting for init failures.

## v1.9.0

- Update package dependencies.
- Upgraded Go 1.23.2 to 1.24.2
- Various updates to k8s API dependencies.
