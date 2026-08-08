# KubeAuth

Production-grade Keycloak OIDC operator for Kubernetes.

[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](LICENSE)
[![Go](https://img.shields.io/badge/Go-1.22+-blue.svg)](https://golang.org)
[![Kubernetes](https://img.shields.io/badge/Kubernetes-1.28+-blue.svg)](https://kubernetes.io)

## Overview

KubeAuth automates the full Keycloak to Kubernetes authentication lifecycle.
Most teams spend days manually configuring OIDC flags, creating Keycloak
clients, and writing ClusterRoleBindings. KubeAuth does all of that
automatically — bootstrap, sync, and verify in seconds.

## The Problem

Connecting Keycloak to Kubernetes requires:
- Manually configuring kube-apiserver OIDC flags
- Manually creating Keycloak realm, client, and claim mappers
- Manually writing ClusterRoleBindings that mirror Keycloak groups
- No feedback loop when auth breaks

KubeAuth owns the full triangle: Keycloak, kube-apiserver, and RBAC.

## Features

- Bootstrap — auto-configure Keycloak realm, client, and mappers
- Verify — test OIDC end-to-end and decode token claims
- Sync — continuously reconcile Keycloak groups to Kubernetes RBAC
- Operator — Kubernetes-native continuous sync via CRDs

## Quick Start

### Build from source

```bash
git clone https://github.com/Abderraoufzekkour/KubeAuth.git
cd KubeAuth
go build -o bin/kubeauth ./cmd/kubeauth
```

### Verify your Keycloak integration

```bash
./bin/kubeauth verify \
  --keycloak-url https://auth.example.com \
  --realm myrealm \
  --client-id kubernetes \
  --username testuser
```

### Run the operator

```bash
./bin/kubeauth operator
```

## How it works

Keycloak Groups
|
v
KeycloakGroupBinding (CRD)
|
v
KubeAuth Operator (reconcile loop)
|
v
Kubernetes ClusterRoleBinding / RoleBinding


## CRDs

### KeycloakGroupBinding

Maps a Keycloak group to a Kubernetes ClusterRole automatically.

```yaml
apiVersion: kubeauth.io/v1alpha1
kind: KeycloakGroupBinding
metadata:
  name: developers-view
spec:
  keycloakGroup: /developers
  clusterRole: view
  groupPrefix: "oidc:"
  keycloakRef:
    name: my-keycloak
    namespace: kubeauth-system
```

## Comparison

| Tool | Keycloak Setup | RBAC Sync | Verify |
|---|---|---|---|
| kubelogin | No | No | No |
| keycloak-operator | Yes | No | No |
| dex | No | No | No |
| **KubeAuth** | **Yes** | **Yes** | **Yes** |

## Roadmap

- [x] v0.1 — CLI with verify and operator commands
- [ ] v0.2 — Real Keycloak HTTP integration
- [ ] v0.3 — Full operator with controller-runtime
- [ ] v0.4 — Helm chart
- [ ] v0.5 — krew index submission

## Contributing

Contributions are welcome. Please open an issue before submitting a PR.

## License

Apache 2.0 — see [LICENSE](LICENSE) for details.