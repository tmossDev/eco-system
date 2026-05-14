# Eco System

See [CHANGELOG.md](CHANGELOG.md) for the project change log.

## Render Kubernetes Resources From Helm

Use this workflow to build the Helm chart dependencies and render the Kubernetes
resource YAML files without applying anything to a cluster.

### Prerequisites

- Helm installed locally.
- Run commands from the repository root.

Check Helm is available:

```sh
helm version --short
```

### Build Chart Dependencies

The service charts depend on the shared local `app` chart, and the umbrella
charts depend on those service charts. Build the leaf chart dependencies first,
then the umbrella chart dependencies:

```sh
helm dependency build foundation/postgres/deployment
helm dependency build foundation/liquibase/deployment
helm dependency build shared-components/frontend/app/storybook/deployment
helm dependency build user-management/backend/app/user-service/deployment
helm dependency build user-management/backend/app/user-gateway/deployment
helm dependency build user-management/frontend/app/admin-web-app/deployment

helm dependency build foundation/deployment
helm dependency build shared-components/deployment
helm dependency build user-management/deployment
```

This writes packaged chart dependencies into each chart's `charts/` directory.

### Render Kubernetes Resource Files

Render all umbrella charts into `build/kubernetes/`:

```sh
mkdir -p build/kubernetes

helm template foundation foundation/deployment \
  --namespace eco-system \
  --include-crds \
  --output-dir build/kubernetes

helm template shared-components shared-components/deployment \
  --namespace eco-system \
  --include-crds \
  --output-dir build/kubernetes

helm template user-management user-management/deployment \
  --namespace eco-system \
  --include-crds \
  --output-dir build/kubernetes
```

The rendered Kubernetes resource files are written under:

```text
build/kubernetes/foundation/
build/kubernetes/shared-components/
build/kubernetes/user-management/
```

### Render One Combined YAML File

If you want a single YAML file instead of one file per template, omit
`--output-dir` and redirect the output:

```sh
helm template user-management user-management/deployment \
  --namespace eco-system \
  --include-crds \
  > build/user-management.yaml
```

### Validate Before Applying

Use `kubectl apply --dry-run=client` to check rendered manifests locally:

```sh
kubectl apply --dry-run=client --recursive -f build/kubernetes/
```

To apply the rendered resources to the current Kubernetes context:

```sh
kubectl apply --recursive -f build/kubernetes/
```
