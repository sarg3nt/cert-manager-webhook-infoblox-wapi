# Cert Manager Webhook Infoblox Wapi Release Notes

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
