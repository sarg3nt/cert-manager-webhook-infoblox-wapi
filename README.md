# Cert Manager Webhook for InfoBlox WAPI

[![OpenSSF Scorecard](https://api.scorecard.dev/projects/github.com/sarg3nt/cert-manager-webhook-infoblox-wapi/badge)](https://scorecard.dev/viewer/?uri=github.com/sarg3nt/cert-manager-webhook-infoblox-wapi)
[![CodeQL](https://github.com/sarg3nt/cert-manager-webhook-infoblox-wapi/actions/workflows/github-code-scanning/codeql/badge.svg)](https://github.com/sarg3nt/cert-manager-webhook-infoblox-wapi/actions/workflows/github-code-scanning/codeql)

An InfoBlox WAPI webhook for cert-manager.

This project provides a custom [ACME DNS01 Challenge Provider](https://cert-manager.io/docs/configuration/acme/dns01) as a webhook for [cert-manager](https://cert-manager.io/). This webhook integrates cert-manager with InfoBlox WAPI is a REST API. You can learn more about WAPI in this [PDF](https://www.infoblox.com/wp-content/uploads/infoblox-deployment-infoblox-rest-api.pdf).

This implementation is based on [infoblox-go-client](https://github.com/infobloxopen/infoblox-go-client) library.

This project is a fork of https://github.com/luisico/cert-manager-webhook-infoblox-wapi, which was forked from 
https://github.com/cert-manager/webhook-example.

- [Requirements](#requirements)
- [Installation](#installation)
  - [Install Cert-manager](#install-cert-manager)
  - [Install infoblox-wapi webhook](#install-infoblox-wapi-webhook)
    - [Using the Public Helm Chart](#using-the-public-helm-chart)
    - [From Source](#from-source)
    - [Values](#values)
  - [Infoblox User Account](#infoblox-user-account)
    - [Kubernetes Secret](#kubernetes-secret)
    - [Hostpath Volume Mount](#hostpath-volume-mount)
  - [Issuer Examples](#issuer-examples)
    - [Cluster Issuer for Let's Encrypt Staging using Secrets For the Infoblox Account](#cluster-issuer-for-lets-encrypt-staging-using-secrets-for-the-infoblox-account)
    - [Cluster Issuer for Let's Encrypt Production using Volume Mount For the Infoblox Account](#cluster-issuer-for-lets-encrypt-production-using-volume-mount-for-the-infoblox-account)
    - [Issuer for Let's Encrypt Production using Volume Mount For the Infoblox Account](#issuer-for-lets-encrypt-production-using-volume-mount-for-the-infoblox-account)
    - [Issuer Webhook Configuration Options](#issuer-webhook-configuration-options)
  - [Creating Certificates](#creating-certificates)
    - [Manually](#manually)
    - [Ingress Annotations](#ingress-annotations)
    - [Setting Default Issuer in Let's Encrypt](#setting-default-issuer-in-lets-encrypt)
- [Running the test suite](#running-the-test-suite)
- [Building](#building)
- [Contributions](#contributions)
- [License](#license)
- [Author](#author)

## Requirements

- InfoBlox GRID installation with WAPI 2.5 or above
- helm v3
- kubernetes 1.21+
- cert-manager 1.5+

Note that other versions might work, but have not been tested.

## Installation

1. Cert-manager
2. Infoblox-wapi webhook
3. Issuer

### Install Cert-manager

Follow [instructions](https://cert-manager.io/docs/installation/) to install cert-manager.

### Install infoblox-wapi webhook

At a minimum you will need to customize `groupName` with your own group name. See [deploy/cert-manager-webhook-infoblox-wapi/values.yaml](./deploy/cert-manager-webhook-infoblox-wapi/values.yaml) for an in-depth explanation and other values that might require tweaking. With either method below, follow [helm instructions](https://helm.sh/docs/intro/using_helm/#customizing-the-chart-before-installing) to customize your deployment.


Docker images are stored in GitHub's [ghcr.io](ghcr.io) registry, specifically at [ghcr.io/sarg3nt/cert-manager-webhook-infoblox-wapi](ghcr.io/sarg3nt/cert-manager-webhook-infoblox-wapi).

#### Using the Public Helm Chart

```sh
helm repo add cert-manager-webhook-infoblox-wapi https://sarg3nt.github.io/cert-manager-webhook-infoblox-wapi

# The values file below is optional, if you don't need it you can remove that line.
helm -n cert-manager install \
  cert-manager-webhook \
  cert-manager-webhook-infoblox-wapi/cert-manager-webhook-infoblox-wapi \
  -f cert-manager-infoblox-values.yaml
```

#### From Source

Check out this repository and run the following command:

```sh
helm -n cert-manager install webhook-infoblox-wapi deploy/cert-manager-webhook-infoblox-wapi
```

#### Values
| Name                           | Description                                                                                                                                                                                                                                                                                                                                                                         | Value                                              |
|--------------------------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|----------------------------------------------------|
| nameOverride                   | String to partially override chart name.                                                                                                                                                                                                                                                                                                                                            | ""                                                 |
| fullNameOverride               | String to fully override chart fullname.                                                                                                                                                                                                                                                                                                                                            | ""                                                 |
| groupName                      | The GroupName here is used to identify your company or business unit that created this webhook. This name will need to be referenced in each Issuer's `webhook` stanza to inform cert-manager of where to send ChallengePayload resources in order to solve the DNS01 challenge. This group name should be **unique**, hence using your own company's domain # here is recommended. | acme.mycompany.com                                 |
| certManager.namespace          | Namespace where cert-manager is deployed.                                                                                                                                                                                                                                                                                                                                           | cert-manager                                       |
| certManager.serviceAccountName | Service account name of cert-manager.                                                                                                                                                                                                                                                                                                                                               | cert-manager                                       |
| rootCACertificate.duration     | Duration of root CA certificate                                                                                                                                                                                                                                                                                                                                                     | 43800h                                             |
| servingCertificate.duration    | Duration of serving certificate                                                                                                                                                                                                                                                                                                                                                     | 8760h                                              |
| image.repository               | Deployment image repository                                                                                                                                                                                                                                                                                                                                                         | ghcr.io/sarg3nt/cert-manager-webhook-infoblox-wapi |
| image.tag                      | Deployment image tag                                                                                                                                                                                                                                                                                                                                                                | 1.5                                                |
| image.pullPolicy               | Image pull policy                                                                                                                                                                                                                                                                                                                                                                   | IfNotPresent                                       |
| secretVolume.hostPath          | Location of a secrets file on the host file system to use instead of a Kubernetes secret                                                                                                                                                                                                                                                                                            | /etc/secrets/secrets.json                          |
| service.type                   | Service type to expose                                                                                                                                                                                                                                                                                                                                                              | ClusterIP                                          |
| service.port                   | Service port to expose                                                                                                                                                                                                                                                                                                                                                              | 443                                                |
| resources                      | Deployment resource limits                                                                                                                                                                                                                                                                                                                                                          | {}                                                 |
| nodeSelector                   | Deployment node selector object                                                                                                                                                                                                                                                                                                                                                     | {}                                                 |
| tolerations                    | Deployment tolerations                                                                                                                                                                                                                                                                                                                                                              | []                                                 |
| affinity                       | Deployment affinity                                                                                                                                                                                                                                                                                                                                                                 | {}                                                 |

### Infoblox User Account

A user account with the ability to create TXT records in the required domain is needed.  
We support two ways of loading this service account.

#### Kubernetes Secret

The first method is to create a Kubernetes secret that include the Infoblox users `username` and `password`.

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: infoblox-credentials
  namespace: cert-manager
type: Opaque
data:
  username: dXNlcm5hbWUK      # base64 encoded: "username"
  password: cGFzc3dvcmQK      # base64 encoded: "password"

---
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: webhook-infoblox-wapi:secret-reader
  namespace: cert-manager
rules:
  - apiGroups: [""]
    resources:
      - secrets
    resourceNames:
      - infoblox-credentials
    verbs:
      - get
      - watch

---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: webhook-infoblox-wapi:secret-reader
  namespace: cert-manager
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: webhook-infoblox-wapi:secret-reader
subjects:
  - apiGroup: ""
    kind: ServiceAccount
    name: cert-manager-webhook-cert-manager-webhook-infoblox-wapi
    namespace: cert-manager
```

Then create a `ClusterIssuer` with the following in the `config` section.  
See full examples below.

```yaml
usernameSecretRef:
  name: infoblox-credentials
  key: username
passwordSecretRef:
  name: infoblox-credentials
  key: password
```

#### Hostpath Volume Mount

The second method is to create a file on the hosts file system that contains the `username` and `password`.  
This file must be created in the path given in `secretVolume.hostPath` in the Helm chart's `values.yaml` file.  Default location is `/etc/secrets/secrets.json`.

**Example:**
The values must be base64 encoded.
```json
{
  "username": "dXNlcm5hbWUK",
  "password": "cGFzc3dvcmQK"
}
```

Then create a `ClusterIssuer` with the following in the `config` section.  
See full examples below.

```yaml
getUserFromVolume: true
```

### Issuer Examples

See: [Cert Manager Issuers](https://cert-manager.io/docs/concepts/issuer/) for more information.

There are two different kind of issuers:
- `Issuer` is for a specific namespace.
- `ClusterIssuer` is for an entire cluster.

#### Cluster Issuer for Let's Encrypt Staging using Secrets For the Infoblox Account

```yaml 
apiVersion: cert-manager.io/v1
kind: ClusterIssuer
metadata:
  name: letsencrypt-staging
spec:
  acme:
    email: your.email@example.com
    server: https://acme-staging-v02.api.letsencrypt.org/directory
    privateKeySecretRef:
      name: letsencrypt-account-key
    solvers:
    - dns01:
        webhook:
          groupName: acme.mycompany.com
          solverName: infoblox-wapi
          config:
            host: infoblox.fqdn
            view: "InfoBlox View"
            usernameSecretRef:
              name: infoblox-credentials
              key: username
            passwordSecretRef:
              name: infoblox-credentials
              key: password
```

#### Cluster Issuer for Let's Encrypt Production using Volume Mount For the Infoblox Account

```yaml 
apiVersion: cert-manager.io/v1
kind: ClusterIssuer
metadata:
  name: letsencrypt-production
spec:
  acme:
    email: your.email@example.com
    server: https://acme-v02.api.letsencrypt.org/directory
    privateKeySecretRef:
      name: letsencrypt-account-key
    solvers:
    - dns01:
        webhook:
          groupName: acme.mycompany.com
          solverName: infoblox-wapi
          config:
            host: infoblox.fqdn
            view: "InfoBlox View"
            getUserFromVolume: true
```
#### Issuer for Let's Encrypt Production using Volume Mount For the Infoblox Account

```yaml 
apiVersion: cert-manager.io/v1
kind: Issuer
metadata:
  name: letsencrypt-production
  namespace: mesh-system
spec:
  acme:
    email: your.email@example.com
    server: https://acme-v02.api.letsencrypt.org/directory
    privateKeySecretRef:
      name: letsencrypt-account-key
    solvers:
    - dns01:
        webhook:
          groupName: acme.mycompany.com
          solverName: infoblox-wapi
          config:
            host: infoblox.fqdn
            view: "InfoBlox View"
            getUserFromVolume: true
```
> **NOTE:** You can create more than one `ClusterIssuer`.  For example, one for Let's Encrypt staging and one for Let's Encrypt production.  You can then reference which one you want to use when creating a cert or annotating an ingress.  See below for examples.

#### Issuer Webhook Configuration Options

This is the full list of webhook configuration options:

- `host`: FQDN or IP address of the InfoBlox server.
- `view`: DNS View in the InfoBlox server to manipulate TXT records in.
- `usernameSecretRef`: Reference to the secret name holding the username for the InfoBlox server (optional if getUserFromVolume is used)
- `passwordSecretRef`: Reference to the secret name holding the password for the InfoBlox server (optional if getUserFromVolume is used)
- `getUserFromVolume: true`: Get the Infoblox user from the host file system. (default: false)
- `port`: Port of the InfoBlox server (default: 443).
- `version`: Version of the InfoBlox server (default: 2.5).
- `sslVerify`: Verify SSL connection (default: false).
- `httpRequestTimeout`: Timeout for HTTP request to the InfoBlox server, in seconds (default: 60).
- `httpPoolConnections`: Maximum number of connections to the InfoBlox server (default: 10).

### Creating Certificates

You can create certificates either manually or via Ingress Annotations.

#### Manually

Now you can create a certificate, for example:

```yaml
apiVersion: cert-manager.io/v1
kind: Certificate
metadata:
  name: infoblox-wapi-test
  namespace: cert-manager
spec:
  commonName: example.com
  dnsNames:
    - example.com
  issuerRef:
    # The name of the issuer created above.
    name: letsencrypt-production
    kind: ClusterIssuer
  secretName: infoblox-wapi-test-tls
```

#### Ingress Annotations

If you are using Nginx Ingress you can add an annotation and Let's Encrypt will automatically create a certificate for you.  
See: [Cert Manager Annotated Ingress resource](https://cert-manager.io/docs/usage/ingress/)

```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  annotations:
    # Use the default issuer: 
    kubernetes.io/tls-acme: "true"
    # OR Use a specific issuer:
    cert-manager.io/cluster-issuer: letsencrypt-staging
# Rest of normal ingress config goes here.
```

> **NOTE:** To use `kubernetes.io/tls-acme: "true"`, a `defaultIssuerName` must be set.  
See: [Setting Default Issuer in Let's Encrypt](#setting-default-issuer-in-lets-encrypt)

#### Setting Default Issuer in Let's Encrypt

When deploying the Let's Encrypt Helm chart you can set a default issuer with the following config.

```yaml
ingressShim:
  defaultIssuerName: "letsencrypt-production"
```

Once this is done you can then use the `kubernetes.io/tls-acme: "true"` annotation and the default issuer will be used.

## Running the test suite

Requirements:

- go >= 1.21

First create you own `config.json` and `credentials.yaml` inside `testdata/infoblox-wapi/` based on the corresponding `.sample` files. The values in `config.json` correspond to the webhook `config` section in the example `ClusterIssuer` above, while `credentials.yaml` will create a secret. Ensure that you fill in the values for the test to connect to an InfoBlox instance.

You can then run the test suite with:

```bash
TEST_ZONE_NAME=example.com. make test
```

## Building

1. If you've made any changes to `go.mod`, run `go mod tidy`
1. Update the `Makefile` with a new `IMAGE_TAG` if necessary.
1. Run `make build`.  A new Docker container will be generated with the `IMAGE_NAME` and `IMAGE_TAG` given in the `Makefile`
1. Run `make push`. This will tag the version given to latest and push both images to the repo in the `IMAGE_NAME`

## Contributions

If you would like to contribute to this projects, please, open a PR via GitHub. Thanks.

## License

This project inherits the Apache 2.0 license from https://github.com/cert-manager/webhook-example.

Modifications to files are listed in [NOTICE](./NOTICE).

## Author

Luis Gracia while at [The Rockefeller University](http://www.rockefeller.edu), taken over by Dave Sargent:
- dave [at] sarg3.net
- GitHub at [sarg3nt](https://github.com/sarg3nt)
