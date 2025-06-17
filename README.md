<!--cspell:ignore mycompany dvcm Gracia sarg  -->
# Cert Manager Webhook for InfoBlox WAPI

[![OpenSSF Scorecard](https://api.scorecard.dev/projects/github.com/sarg3nt/cert-manager-webhook-infoblox-wapi/badge)](https://scorecard.dev/viewer/?uri=github.com/sarg3nt/cert-manager-webhook-infoblox-wapi)
[![CodeQL](https://github.com/sarg3nt/cert-manager-webhook-infoblox-wapi/actions/workflows/codeql.yml/badge.svg)](https://github.com/sarg3nt/cert-manager-webhook-infoblox-wapi/actions/workflows/codeql.yml)
[![trivy](https://github.com/sarg3nt/cert-manager-webhook-infoblox-wapi/actions/workflows/trivy.yml/badge.svg)](https://github.com/sarg3nt/cert-manager-webhook-infoblox-wapi/actions/workflows/trivy.yml)
[![Release](https://github.com/sarg3nt/cert-manager-webhook-infoblox-wapi/actions/workflows/release.yml/badge.svg)](https://github.com/sarg3nt/cert-manager-webhook-infoblox-wapi/actions/workflows/release.yml)
[![Weekly Release](https://github.com/sarg3nt/cert-manager-webhook-infoblox-wapi/actions/workflows/release-weekly.yml/badge.svg)](https://github.com/sarg3nt/cert-manager-webhook-infoblox-wapi/actions/workflows/release-weekly.yml)  
[![Scorecard Analyzer](https://github.com/sarg3nt/cert-manager-webhook-infoblox-wapi/actions/workflows/scorecard.yml/badge.svg)](https://github.com/sarg3nt/cert-manager-webhook-infoblox-wapi/actions/workflows/scorecard.yml)
[![Dependabot Updates](https://github.com/sarg3nt/cert-manager-webhook-infoblox-wapi/actions/workflows/dependabot/dependabot-updates/badge.svg)](https://github.com/sarg3nt/cert-manager-webhook-infoblox-wapi/actions/workflows/dependabot/dependabot-updates)
[![Dependency Review](https://github.com/sarg3nt/cert-manager-webhook-infoblox-wapi/actions/workflows/dependency-review.yml/badge.svg)](https://github.com/sarg3nt/cert-manager-webhook-infoblox-wapi/actions/workflows/dependency-review.yml)
****

An InfoBlox WAPI webhook for cert-manager.

This project provides a custom [ACME DNS01 Challenge Provider](https://cert-manager.io/docs/configuration/acme/dns01) as a webhook for [Cert Manager](https://cert-manager.io/). This webhook integrates Cert Manager with InfoBlox WAPI via its REST API. You can learn more about WAPI in this [PDF](https://www.infoblox.com/wp-content/uploads/infoblox-deployment-infoblox-rest-api.pdf).

This implementation is based on the [infoblox-go-client](https://github.com/infobloxopen/infoblox-go-client) library.

This project is a fork of https://github.com/luisico/cert-manager-webhook-infoblox-wapi, which was forked from 
https://github.com/cert-manager/webhook-example.

- [Requirements](#requirements)
- [Installation](#installation)
  - [Install Cert Manager](#install-cert-manager)
  - [Install Infoblox Wapi Webhook](#install-infoblox-wapi-webhook)
    - [Using the Public Helm Chart](#using-the-public-helm-chart)
    - [From Source](#from-source)
    - [Values](#values)
  - [Infoblox User Account](#infoblox-user-account)
    - [Kubernetes Secret](#kubernetes-secret)
    - [Hostpath Volume Mount](#hostpath-volume-mount)
  - [Create Issuers](#create-issuers)
    - [Cluster Issuer for Let's Encrypt Staging using Secrets For the Infoblox Account](#cluster-issuer-for-lets-encrypt-staging-using-secrets-for-the-infoblox-account)
    - [Cluster Issuer for Let's Encrypt Production using Volume Mount For the Infoblox Account](#cluster-issuer-for-lets-encrypt-production-using-volume-mount-for-the-infoblox-account)
    - [Issuer for Let's Encrypt Production using Volume Mount For the Infoblox Account](#issuer-for-lets-encrypt-production-using-volume-mount-for-the-infoblox-account)
    - [Issuer Webhook Configuration Options](#issuer-webhook-configuration-options)
  - [Creating Certificates](#creating-certificates)
    - [Manually](#manually)
    - [Ingress Annotations](#ingress-annotations)
    - [Setting Default Issuer in Let's Encrypt](#setting-default-issuer-in-lets-encrypt)
- [Building](#building)
- [Contributions](#contributions)
- [License](#license)
- [Author](#author)

## Requirements

- InfoBlox GRID installation with WAPI 2.5 or above
- helm v3
- kubernetes 1.21+
- cert-manager 1.5+

> [!NOTE]
> Other versions might work, but have not been tested.

## Installation

There are three steps needed to make this work.

1. [Install Cert Manager](#install-cert-manager)
2. [Install Infoblox Wapi Webhook](#install-infoblox-wapi-webhook) (this plugin)
3. [Create Issuers](#create-issuers)

### Install Cert Manager

Follow the [instructions](https://cert-manager.io/docs/installation/) to install Cert Manager.

### Install Infoblox Wapi Webhook

At a minimum you will need to customize `groupName` with your own group name. See [charts/cert-manager-webhook-infoblox-wapi/values.yaml](./charts/cert-manager-webhook-infoblox-wapi/values.yaml) for an in-depth explanation and other values that might require tweaking. With either method below, follow [helm instructions](https://helm.sh/docs/intro/using_helm/#customizing-the-chart-before-installing) to customize your deployment.

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
| Name                           | Description                                                                                                                                                                                                                                                                                                                                                                       | Value                                              |
|--------------------------------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|----------------------------------------------------|
| nameOverride                   | String to partially override chart name.                                                                                                                                                                                                                                                                                                                                          | ""                                                 |
| fullNameOverride               | String to fully override chart fullname.                                                                                                                                                                                                                                                                                                                                          | ""                                                 |
| groupName                      | The GroupName here is used to identify your company or business unit that created this webhook. This name will need to be referenced in each Issuer's `webhook` stanza to inform cert-manager of where to send ChallengePayload resources in order to solve the DNS01 challenge. This group name should be **unique**, hence using your own company's domain here is recommended. | acme.mycompany.com                                 |
| certManager.namespace          | Namespace where cert-manager is deployed.                                                                                                                                                                                                                                                                                                                                         | cert-manager                                       |
| certManager.serviceAccountName | Service account name of cert-manager.                                                                                                                                                                                                                                                                                                                                             | cert-manager                                       |
| rootCACertificate.duration     | Duration of root CA certificate                                                                                                                                                                                                                                                                                                                                                   | 43800h                                             |
| servingCertificate.duration    | Duration of serving certificate                                                                                                                                                                                                                                                                                                                                                   | 8760h                                              |
| image.repository               | Deployment image repository                                                                                                                                                                                                                                                                                                                                                       | ghcr.io/sarg3nt/cert-manager-webhook-infoblox-wapi |
| image.tag                      | Deployment image tag                                                                                                                                                                                                                                                                                                                                                              | 1.5                                                |
| image.pullPolicy               | Image pull policy                                                                                                                                                                                                                                                                                                                                                                 | IfNotPresent                                       |
| secretVolume.hostPath          | Location of a secrets file on the host file system to use instead of a Kubernetes secret                                                                                                                                                                                                                                                                                          | /etc/secrets/secrets.json                          |
| service.type                   | Service type to expose                                                                                                                                                                                                                                                                                                                                                            | ClusterIP                                          |
| service.port                   | Service port to expose                                                                                                                                                                                                                                                                                                                                                            | 443                                                |
| resources                      | Deployment resource limits                                                                                                                                                                                                                                                                                                                                                        | {}                                                 |
| nodeSelector                   | Deployment node selector object                                                                                                                                                                                                                                                                                                                                                   | {}                                                 |
| tolerations                    | Deployment tolerations                                                                                                                                                                                                                                                                                                                                                            | []                                                 |
| affinity                       | Deployment affinity                                                                                                                                                                                                                                                                                                                                                               | {}                                                 |

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
See [Issuer Examples](#issuer-examples)

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
See [Create Issuers](#create-issuers)

```yaml
getUserFromVolume: true
```

### Create Issuers

An issuer is the method that Cert Manager will use to request a certificate and the configuration Let's Encrypt will use to validate that the requester (you) owns the domain the certificate request is for.

The part of an issuer that defines the use of this webhook plugin starts in the `webhook` section as shown in the examples below.

All settings under `config` are specific to this plugin.  See the list of [Issuer Webhook Configuration Options](#issuer-webhook-configuration-options) below.

See: [Cert Manager Issuers](https://cert-manager.io/docs/concepts/issuer/) in the official Cert Manager documentation for more information.

There are two different kinds of issuers:
- `Issuer` is for a specific namespace.
- `ClusterIssuer` is for an entire cluster and is most often used.

#### Cluster Issuer for Let's Encrypt Staging using Secrets For the Infoblox Account

```yaml 
apiVersion: cert-manager.io/v1
kind: ClusterIssuer
metadata:
  name: letsencrypt-staging
spec:
  acme:
    # The email matching the account that your Lets Encrypt account key was created in.
    email: your.email@example.com
    # What Lets Encrypt server to use. This one is for staging certificates during development.
    server: https://acme-staging-v02.api.letsencrypt.org/directory
    privateKeySecretRef:
      name: letsencrypt-account-key
    solvers:
    - dns01:
        webhook:
          # groupName must match the groupName you set while installing this plugin via the Helm chart.
          groupName: acme.mycompany.com
          solverName: infoblox-wapi
          config:
            host: my-infoblox.company.com # required
            view: "InfoBlox View" # required
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
            host: my-infoblox.company.com
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
            host: my-infoblox.company.com
            view: "InfoBlox View"
            getUserFromVolume: true
```
> [!NOTE] 
> You can create more than one `ClusterIssuer`.  For example, one for Let's Encrypt staging and one for Let's Encrypt production.  You can then reference which one you want to use when creating a cert or annotating an ingress.  See below for examples.

#### Issuer Webhook Configuration Options

This is the full list of webhook configuration options:

- `groupName`: This must match the `groupName` you specified in the Helm chart config during install.
- `host`: FQDN or IP address of the InfoBlox server.
- `view`: DNS View in the InfoBlox server to manipulate TXT records in.
- `usernameSecretRef`: Reference to the secret name holding the username for the InfoBlox server (optional if getUserFromVolume is true)
- `passwordSecretRef`: Reference to the secret name holding the password for the InfoBlox server (optional if getUserFromVolume is true)
- `getUserFromVolume: true`: Get the Infoblox user from the host file system. (default: false)
- `port`: Port of the InfoBlox server (default: 443).
- `version`: Version of the InfoBlox server (default: 2.10).
- `sslVerify`: Verify SSL connection (default: false).
- `httpRequestTimeout`: Timeout for HTTP request to the InfoBlox server, in seconds (default: 60).
- `httpPoolConnections`: Maximum number of connections to the InfoBlox server (default: 10).
- `ttl`: The time to live of the TXT record. (default: 90)
- `useTtl`: Whether or not to use the ttl.  (default: true)

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

> [!NOTE] 
> To use `kubernetes.io/tls-acme: "true"`, a `defaultIssuerName` must be set.  
> See: [Setting Default Issuer in Let's Encrypt](#setting-default-issuer-in-lets-encrypt)

#### Setting Default Issuer in Let's Encrypt

When deploying the Let's Encrypt Helm chart you can set a default issuer with the following config.

```yaml
ingressShim:
  defaultIssuerName: "letsencrypt-production"
```

Once this is done you can then use the `kubernetes.io/tls-acme: "true"` annotation and the default issuer will be used.

<!-- ## Running the test suite

Requirements:

- go >= 1.21

First create you own `config.json` and `credentials.yaml` inside `testdata/infoblox-wapi/` based on the corresponding `.sample` files. The values in `config.json` correspond to the webhook `config` section in the example `ClusterIssuer` above, while `credentials.yaml` will create a secret. Ensure that you fill in the values for the test to connect to an InfoBlox instance.

You can then run the test suite with:

```bash
TEST_ZONE_NAME=example.com. make test
``` -->

## Building

1. If you've made any changes to `go.mod`, run `go mod tidy`
1. Update the `Makefile` with a new `IMAGE_TAG` if necessary.
1. Run `make build`.  This will use `go` to build the project.
1. Run `make build-container`.  A new Docker container will be generated and tagged as:  
   `$(IMAGE_NAME):$(IMAGE_TAG)-$(GIT_BRANCH)` as given in the `Makefile`
1. Run `make push-container`. This will push the above tagged image to the repo defined in the `IMAGE_NAME`

## Contributions

If you would like to contribute to this project, please, open a PR via GitHub. Thanks.

## License

This project inherits the Apache 2.0 license from https://github.com/cert-manager/webhook-example.

Modifications to files are listed in [NOTICE](./NOTICE).

## Author

Luis Gracia while at [The Rockefeller University](http://www.rockefeller.edu), taken over by Dave Sargent:
- dave [at] sarg3.net
- GitHub at [sarg3nt](https://github.com/sarg3nt)
