# Cert Manager Webhook Infoblox Wapi Release Notes

## v1.7.0

### Release Summary

-  Upgraded all dependencies to latest.

### Dependency Changes

- Upgraded Go 1.21 to 1.23.2
- Upgraded github.com/infobloxopen/infoblox-go-client/v2 v2.0.0 to v2.7.0
- Upgraded Alpine 3.19 to 3.20

## v1.6.0

### Release Summary

-  First release of this form of https://github.com/luisico/cert-manager-webhook-infoblox-wapi
-  Add ability to pass Infoblox username and password via a Volume Mount from the host OS file system.  
   In some use cases this can be more secure or preferred than using a secret from Kubernetes.
   See the [README.md](README.md#hostpath-volume-mount) for instructions.
-  Update many package dependencies.
-  Added OpenSSF Scorecard, CodeQL and Dependabot to repo.
-  Improved [README.md](README.md) substantially.

### Dependency Changes

- Upgraded Go 1.16 to 1.21
- Upgraded Alpine 3.14 to 3.19
- Upgraded github.com/jetstack/cert-manager v1.5.4 to github.com/cert-manager/cert-manager v1.13.3
- Upgraded k8s.io/apiextensions-apiserver v0.21.3 to v0.28.1
- Upgraded k8s.io/apimachinery v0.21.3 to v0.28.1
- Upgraded k8s.io/client-go v0.21.3 to v0.28.1