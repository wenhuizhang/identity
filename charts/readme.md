# Identity Node Helm Chart

This Helm chart deploys the Identity Node service with a PostgreSQL database on Kubernetes.

## Overview

The Identity Node service provides identity management capabilities including Verifiable Credential verification and management. This chart includes:

- Identity Node deployment
- PostgreSQL database (using Bitnami's PostgreSQL chart)
- HTTP and gRPC service endpoints
- Ingress configurations for both HTTP and gRPC

## Prerequisites

- Kubernetes 1.19+
- Helm 3.2.0+
- Ingress controller (like NGINX)
- Cert-manager (for TLS certificates)

## Installation

```bash
# Add the required repositories
helm repo add bitnami https://charts.bitnami.com/bitnami

# Update dependencies
helm dependency update ./charts
```

Create a `values-custom.yaml` file with your domain configuration:

```yaml
ingress:
  domainPrefixHttp: api.yourdomain
  domainPrefixGrpc: api.grpc.yourdomain
  apiDomainName: yourdomain.com
```

Then install with:

```bash
# Create your own values-override.yaml file
helm install identity-node ./charts -f values-custom.yaml
```

## Configuration

The following table lists the configurable parameters of the Identity Node chart and their default values.

| Parameter | Description | Default |
|-----------|-------------|---------|
| `replicaCount` | Number of Identity Node replicas | `1` |
| `image.repository` | Identity Node container image repository | `ghcr.io/agntcy/identity/node` |
| `image.tag` | Identity Node container image tag | `latest` |
| `image.pullPolicy` | Image pull policy | `IfNotPresent` |
| `service.type` | Kubernetes service type | `ClusterIP` |
| `service.http.port` | HTTP service port | `4000` |
| `service.grpc.port` | gRPC service port | `5006` |
| `ingress.enabled` | Enable ingress | `true` |
| `ingress.domainPrefixHttp` | Prefix for HTTP domain | `api.example` |
| `ingress.domainPrefixGrpc` | Prefix for gRPC domain | `api.grpc.example` |
| `ingress.apiDomainName` | Base domain name | `example.com` |
| `postgresql.enabled` | Enable PostgreSQL | `true` |
| `postgresql.auth.postgresPassword` | PostgreSQL password | `change-me` |
| `resources` | CPU/Memory resource requests/limits | See values.yaml |

## Security Notes

**IMPORTANT**: The default values in this chart are for demonstration purposes only. For production deployments:

1. Change all default passwords in PostgreSQL
2. Use a specific image tag instead of `latest`
3. Consider using external secrets management
4. Apply appropriate resource limits
5. Use a values override file for sensitive information

## License

This chart is licensed under Apache 2.0. See the LICENSE file for details.
