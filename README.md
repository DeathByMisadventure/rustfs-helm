# RustFS Helm Chart

A Helm chart for RustFS S3-compatible object storage on Kubernetes

![Version: 0.1.0](https://img.shields.io/badge/Version-0.1.0-informational?style=flat-square) ![Type: application](https://img.shields.io/badge/Type-application-informational?style=flat-square) ![AppVersion: latest](https://img.shields.io/badge/AppVersion-latest-informational?style=flat-square)

## About RustFS

RustFS is a high-performance, **open-source**, S3-compatible object storage system built in Rust. It serves as a modern alternative to MinIO, combining simplicity with the memory safety and raw performance of the Rust programming language.

### Key Features

- **Blazing Fast Performance** — Up to 2.3x faster than MinIO for small (4KB) object workloads, leveraging Rust's efficiency and zero-cost abstractions.
- **Full S3 Compatibility** — Seamless integration with existing S3 tools, clients, and ecosystems; supports migration and coexistence with MinIO, Ceph, and other S3-compatible systems.
- **Distributed & Scalable** — Fault-tolerant distributed architecture designed for large-scale deployments, data lakes, AI/ML workloads, and big data analytics.
- **Memory Safe & Secure** — Built in Rust for maximum speed, resource efficiency, and protection against common memory-related vulnerabilities.
- **Permissive Licensing** — Released under Apache 2.0, avoiding the restrictions of AGPL—ideal for commercial and enterprise use without lock-in.

### Benefits

- **Drop-in MinIO Replacement** → Easily migrate existing MinIO deployments while gaining performance improvements and a more permissive license.
- **Optimized for Modern Workloads** → Excels in AI, data lakes, high-concurrency I/O, and small-object scenarios common in analytics and machine learning.
- **Cloud-Native Ready** → Supports Docker/Kubernetes deployments, multi-site replication, encryption, versioning, and WORM compliance.
- **Community-Driven & Open** → Fully open-source with active development; no hidden telemetry or proprietary restrictions.
- **Cost-Effective** → High performance and efficiency reduce infrastructure costs while maintaining enterprise-grade reliability.

## Prerequisites

- A Kubernetes cluster (version >=1.19).
- Helm 3 installed.
- `kubectl` configured to interact with your cluster.
- A default StorageClass configured in the cluster for dynamic PVC provisioning.
- An ingress controller such as ingress-nginx which can be deployed with:

```bash
helm upgrade --install ingress-nginx ingress-nginx \
  --repo https://kubernetes.github.io/ingress-nginx \
  --namespace ingress-nginx --create-namespace
```

## Configure values.yaml or local-values.yaml

Customize the deployment by editing values.yaml or creating a local-values.yaml file with overrides. Key configurations include:

- image: Specify the RustFS image repository, tag, and pull policy.
- persistence: Configure the PVC for data storage (size, storage class).
- ingress: Enable and configure separate Ingress resources for the S3 API and web console.
- auth: Set access and secret keys for S3 authentication.
- resources and securityContext: Define resource limits and security settings.

## Enabling Ingress

The chart creates two separate Ingress resources:

- One for the S3 API (port 9000)
- One for the RustFS web console (port 9001)

If you lack an Ingress controller, install one separately (e.g., NGINX Ingress as shown above).
For TLS, configure ingress.tls in values.yaml:

```YAML
ingress:
  enabled: true
  className: "nginx"
  tls:
    - secretName: rustfs-tls-secret
      hosts:
        - rustfsapi.lvh.me
        - rustfsconsole.lvh.me
```

### Create the TLS secret beforehand:

```Bash
kubectl create secret tls rustfs-tls-secret --cert=tls.crt --key=tls.key -n <your-namespace>
```

## Deployment

1. Verify your cluster's default StorageClass: `kubectl get storageclass`
1. Ensure a default StorageClass is set, or the PVC will fail to bind.
1. Create a namespace if needed: `kubectl create ns rustfs`
1. Customize values by creating a local-values.yaml file with overrides from values.yaml.
1. Install the Helm chart: `helm install --namespace rustfs -f ./local-values.yaml rustfs .`

## Values

| Key | Type | Description |
|-----|------|-------------|
| affinity | object | Affinity rules |
| auth | object | Root authentication (used for init and console login) |
| auth.accessKey | string | Root access key (console login username, S3 root user) |
| auth.secretKey | string | Root secret key (console login password, S3 root secret) |
| containerSecurityContext | object | Container security context |
| image.pullPolicy | string | Image pull policy |
| image.repository | string | RustFS container image repository |
| image.tag | string | Image tag (defaults to Chart.appVersion) |
| ingress | object | Ingress configuration |
| ingress.annotations | object | Annotations to add to both API and console Ingress resources |
| ingress.apiHosts | list | List of hostnames for the S3 API Ingress Multiple hosts can be added for multi-domain setups |
| ingress.className | string | IngressClass resource to use (e.g., "nginx", "traefik") |
| ingress.consoleHosts | list | List of hostnames for the RustFS web console Ingress Multiple hosts can be added if needed |
| ingress.enabled | bool | Enable creation of Ingress resources for API and console |
| ingress.tls | list | TLS configuration for the Ingress resources Example: tls:   - secretName: rustfs-tls-secret     hosts:       - rustfsapi.lvh.me       - rustfsconsole.lvh.me |
| initialization | object | Initialization settings (runs as a one-time Job after deployment) |
| initialization.buckets | list | List of initial buckets to create |
| initialization.enabled | bool | Enable automatic creation of buckets |
| nodeSelector | object | Node selector |
| persistence | object | Persistence configuration |
| persistence.accessMode | string | Access mode for the PVC |
| persistence.enabled | bool | Enable persistent storage |
| persistence.size | string | Size of the persistent volume |
| persistence.storageClass | string | Storage class to use (leave empty for default) |
| probes | object | Probes configuration (enabled by default) |
| probes.enabled | bool | Enable liveness, readiness, and startup probes |
| replicaCount | int | Number of RustFS pods to run. For single-node mode set to 1. For distributed mode set >=4 (Currently distributed mode is not supported). |
| resources | object | Resource requests and limits |
| s3region | string | S3 region |
| securityContext | object | Pod security context |
| securityContext.fsGroup | int | File system group ID |
| securityContext.runAsGroup | int | Group ID to run as |
| securityContext.runAsNonRoot | bool | Run as non-root user |
| securityContext.runAsUser | int | User ID to run as |
| service | object | Service configuration |
| service.consolePort | int | Port for the RustFS console (web UI) |
| service.port | int | Port for the S3 API |
| service.type | string | Service type (ClusterIP, NodePort, LoadBalancer, etc.) |
| serviceAccount | object | Service account configuration |
| serviceAccount.annotations | object | Annotations to add to the service account |
| serviceAccount.create | bool | Whether to create a service account |
| serviceAccount.name | string | Name of the service account (if not created) |
| tolerations | list | Tolerations |
