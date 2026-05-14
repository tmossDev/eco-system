# Eco System

See [CHANGELOG.md](CHANGELOG.md) for the project change log.

## Render Kubernetes Resources From Helm

Use this workflow to build the Helm chart dependencies and render the Kubernetes
resource YAML files without applying anything to a cluster.

### Prerequisites

- Helm installed locally.
- kubectl installed locally if you want to validate or deploy to a cluster.
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

## Deploy To Kubernetes Manually

For a real cluster deploy, prefer `helm upgrade --install` over applying rendered
YAML files directly. Helm tracks the release and makes later upgrades cleaner.

### 1. Point kubectl At The Cluster

Confirm your local kube context is targeting the right cluster:

```sh
kubectl config current-context
kubectl get nodes
```

Set the namespace used by the examples:

```sh
export NAMESPACE=eco-system
kubectl create namespace "$NAMESPACE" --dry-run=client -o yaml | kubectl apply -f -
```

### 2. Build And Push Images

The Helm charts deploy these application images:

- `user-service`
- `user-gateway`
- `admin-web-app`
- `storybook`
- `liquibase`

Choose a container registry and tag. This example uses GitHub Container Registry:

```sh
export REGISTRY=ghcr.io
export OWNER=<your-github-user-or-org>
export IMAGE_PREFIX="$REGISTRY/$OWNER/eco-system"
export IMAGE_TAG="$(git rev-parse --short HEAD)"
```

Log in to the registry:

```sh
docker login ghcr.io
```

Build and push the images:

```sh
docker buildx build --push \
  -f dockerfile \
  --build-arg SERVICE_PATH=./user-management/backend/app/user-service/cmd \
  -t "$IMAGE_PREFIX/user-service:$IMAGE_TAG" \
  .

docker buildx build --push \
  -f dockerfile \
  --build-arg SERVICE_PATH=./user-management/backend/app/user-gateway/cmd \
  -t "$IMAGE_PREFIX/user-gateway:$IMAGE_TAG" \
  .

docker buildx build --push \
  -f user-management/frontend/app/admin-web-app/Dockerfile \
  -t "$IMAGE_PREFIX/admin-web-app:$IMAGE_TAG" \
  .

docker buildx build --push \
  -f shared-components/frontend/app/storybook/Dockerfile \
  -t "$IMAGE_PREFIX/storybook:$IMAGE_TAG" \
  .

docker buildx build --push \
  -f foundation/liquibase/dockerfile \
  -t "$IMAGE_PREFIX/liquibase:$IMAGE_TAG" \
  foundation/liquibase
```

If your registry images are private, create an image pull secret in the cluster:

```sh
kubectl create secret docker-registry ghcr-pull-secret \
  --namespace "$NAMESPACE" \
  --docker-server=ghcr.io \
  --docker-username=<github-user> \
  --docker-password=<github-token-with-read-packages>
```

If your registry images are public, remove the
`--set <chart>.app.imagePullSecrets[0].name=ghcr-pull-secret` arguments from
the Helm commands below.

### 3. Deploy With Helm

Build chart dependencies:

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

Set database passwords for the first deploy:

```sh
export POSTGRES_PASSWORD=<postgres-password>
export LIQUIBASE_PASSWORD=<liquibase-password>
export APP_PASSWORD=<app-password>
```

Deploy the foundation chart first because it provides postgres and liquibase:

```sh
helm upgrade --install foundation foundation/deployment \
  --namespace "$NAMESPACE" \
  --create-namespace \
  --set liquibase.app.image.repository="$IMAGE_PREFIX/liquibase" \
  --set liquibase.app.image.tag="$IMAGE_TAG" \
  --set 'liquibase.app.imagePullSecrets[0].name=ghcr-pull-secret' \
  --set-string postgres.app.secret.stringData.POSTGRES_PASSWORD="$POSTGRES_PASSWORD" \
  --set-string postgres.app.secret.stringData.LIQUIBASE_PASSWORD="$LIQUIBASE_PASSWORD" \
  --set-string postgres.app.secret.stringData.APP_PASSWORD="$APP_PASSWORD" \
  --set-string liquibase.app.secret.stringData.LIQUIBASE_COMMAND_PASSWORD="$LIQUIBASE_PASSWORD"
```

Deploy shared components:

```sh
helm upgrade --install shared-components shared-components/deployment \
  --namespace "$NAMESPACE" \
  --create-namespace \
  --set storybook.app.image.repository="$IMAGE_PREFIX/storybook" \
  --set storybook.app.image.tag="$IMAGE_TAG" \
  --set 'storybook.app.imagePullSecrets[0].name=ghcr-pull-secret'
```

Deploy user management:

```sh
helm upgrade --install user-management user-management/deployment \
  --namespace "$NAMESPACE" \
  --create-namespace \
  --set user-service.app.image.repository="$IMAGE_PREFIX/user-service" \
  --set user-service.app.image.tag="$IMAGE_TAG" \
  --set 'user-service.app.imagePullSecrets[0].name=ghcr-pull-secret' \
  --set user-gateway.app.image.repository="$IMAGE_PREFIX/user-gateway" \
  --set user-gateway.app.image.tag="$IMAGE_TAG" \
  --set 'user-gateway.app.imagePullSecrets[0].name=ghcr-pull-secret' \
  --set admin-web-app.app.image.repository="$IMAGE_PREFIX/admin-web-app" \
  --set admin-web-app.app.image.tag="$IMAGE_TAG" \
  --set 'admin-web-app.app.imagePullSecrets[0].name=ghcr-pull-secret'
```

Check the result:

```sh
helm list -n "$NAMESPACE"
kubectl get pods,svc,pvc -n "$NAMESPACE"
```

If you need to remove the releases:

```sh
helm uninstall user-management shared-components foundation -n "$NAMESPACE"
```

## Deploy With GitHub Actions

The workflow at `.github/workflows/deploy.yml` builds the deployable images,
pushes them to GitHub Container Registry, and deploys the Helm releases to the
cluster.

### Required GitHub Secrets

Add these in GitHub under `Settings > Secrets and variables > Actions`:

```text
KUBE_CONFIG_B64
GHCR_USERNAME
GHCR_TOKEN
POSTGRES_PASSWORD
LIQUIBASE_PASSWORD
APP_PASSWORD
```

`KUBE_CONFIG_B64` is a base64-encoded kubeconfig for a service account that can
deploy into the target namespace:

```sh
base64 -i ~/.kube/config | pbcopy
```

`GHCR_TOKEN` should be a GitHub token with `read:packages` permission so the
cluster can pull private GHCR images. The workflow uses the built-in
`GITHUB_TOKEN` to push images during the GitHub Actions run.

### Run The Workflow

The workflow runs automatically on pushes to `main`. You can also run it
manually from the GitHub Actions tab and choose the namespace.

The deploy order is:

```text
foundation -> shared-components -> user-management
```
