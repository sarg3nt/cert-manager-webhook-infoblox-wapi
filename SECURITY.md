# Security Policy

## Supported Versions

- OS patches are applied Weekly, typically at 00:00 on Sunday, which include applying security patches during the `docker build`, a release is automatically cut adding one to the patch version and a new Docker container is built.
- A new Helm chart is released when new OS patches are applied, the `appVersion` is auto updated. 
- The user is responsible for applying the new Helm chart.
Dependabot is used to maintain dependency versions.

## Reporting a Vulnerability

- Create an Issue for low impact vulnerabilities.
- Report high impact vulnerabilities by contacting us at [security@sarg3.net](mailto:security@sarg3.net)

