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
  --namespace eco-foundation \
  --include-crds \
  --output-dir build/kubernetes

helm template shared-components shared-components/deployment \
  --namespace abc123def4 \
  --include-crds \
  --output-dir build/kubernetes

helm template user-management user-management/deployment \
  --namespace abc123def4 \
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
YAML files directly. Helm tracks releases and makes later upgrades cleaner.

The deployment is split into two layers:

- `foundation` deploys Postgres and Liquibase into `eco-foundation`.
- `application` deploys user services, admin web app and the design-system Storybook into an application namespace.

Feature deployments use a namespace derived from the branch name. Branches must
match `feature/<10 lowercase alphanumeric chars>`, for example
`feature/abc123def4`, and the namespace becomes `abc123def4`.

### Prerequisites

- Go installed for backend binaries.
- Node 24 and pnpm 11.0.9 installed for frontend builds.
- Docker installed for local runtime images.
- Helm and kubectl installed.
- kubectl must be able to read the kubeconfig.
- Your user or GitHub runner user must be able to access the Docker daemon.

On k3s, this may require fixing `/etc/rancher/k3s/k3s.yaml` permissions or
copying it to your user kubeconfig:

```sh
sudo chmod 644 /etc/rancher/k3s/k3s.yaml
export KUBECONFIG=/etc/rancher/k3s/k3s.yaml
kubectl get nodes
```

After this, run `kubectl` without `sudo`. In piped commands, `sudo` only
applies to the command before the pipe unless both sides use it.

### Load Environment Values

The homelab values live in `config/hp-prodesk-homelab.txt`:

```sh
set -a
. config/hp-prodesk-homelab.txt
set +a

export IMAGE_PREFIX=localhost/eco-system
export IMAGE_TAG="$(git rev-parse --short HEAD)"
export KUBE_CONTEXT="$(kubectl config current-context)"
export K3D_CLUSTER="${K3D_CLUSTER:-${KUBE_CONTEXT#k3d-}}"
export FOUNDATION_NAMESPACE="${FOUNDATION_NAMESPACE:-eco-foundation}"
export NAMESPACE="${DEFAULT_APPLICATION_NAMESPACE:-eco-test}"
export BACKEND_API_URL=""
export FRONTEND_ORIGIN="http://${NAMESPACE}.admin-web-app.${BASE_DOMAIN:-com}"
export ADMIN_WEB_APP_HOST="${NAMESPACE}.admin-web-app.${BASE_DOMAIN:-com}"
export USER_GATEWAY_HOST="${NAMESPACE}.user-gateway.${BASE_DOMAIN:-com}"
export USER_SERVICE_HOST="${NAMESPACE}.user-service.${BASE_DOMAIN:-com}"
export STORYBOOK_HOST="${NAMESPACE}.storybook.${BASE_DOMAIN:-com}"
```

For a feature namespace, override `NAMESPACE` with the branch suffix:

```sh
export NAMESPACE=abc123def4
```

Set database passwords before deploying:

```sh
export POSTGRES_PASSWORD=very_secure_password
export LIQUIBASE_PASSWORD=very_secure_password
export APP_PASSWORD=very_secure_password
```

### Build Chart Dependencies

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

### Deploy Foundation Layer

```sh
docker build \
  -f foundation/liquibase/dockerfile \
  -t "$IMAGE_PREFIX/liquibase:$IMAGE_TAG" \
  foundation/liquibase

docker image inspect "$IMAGE_PREFIX/liquibase:$IMAGE_TAG" >/dev/null

if command -v k3d >/dev/null 2>&1 && [ "${KUBE_CONTEXT#k3d-}" != "$KUBE_CONTEXT" ]; then
  k3d image import "$IMAGE_PREFIX/liquibase:$IMAGE_TAG" -c "$K3D_CLUSTER"
elif command -v k3s >/dev/null 2>&1; then
  docker save "$IMAGE_PREFIX/liquibase:$IMAGE_TAG" | sudo k3s ctr -n k8s.io images import -
  sudo k3s ctr -n k8s.io images list | grep "$IMAGE_PREFIX/liquibase"
fi

kubectl create namespace "$FOUNDATION_NAMESPACE" --dry-run=client -o yaml | kubectl apply -f -

helm upgrade --install foundation foundation/deployment \
  --namespace "$FOUNDATION_NAMESPACE" \
  --set postgres.app.env.POSTGRES_DB="${POSTGRES_DB:-ecoDB}" \
  --set liquibase.app.image.repository="$IMAGE_PREFIX/liquibase" \
  --set liquibase.app.image.tag="$IMAGE_TAG" \
  --set liquibase.app.image.pullPolicy=Never \
  --set-string liquibase.app.env.LIQUIBASE_COMMAND_URL="$LIQUIBASE_JDBC_URL" \
  --set-string liquibase.app.env.LIQUIBASE_COMMAND_USERNAME="${POSTGRES_LIQUIBASE_USER:-liquibase}" \
  --set-string postgres.app.secret.stringData.POSTGRES_PASSWORD="$POSTGRES_PASSWORD" \
  --set-string postgres.app.secret.stringData.LIQUIBASE_PASSWORD="$LIQUIBASE_PASSWORD" \
  --set-string postgres.app.secret.stringData.APP_PASSWORD="$APP_PASSWORD" \
  --set-string liquibase.app.secret.stringData.LIQUIBASE_COMMAND_PASSWORD="$LIQUIBASE_PASSWORD"

kubectl get all -n "$FOUNDATION_NAMESPACE"
```

If the Liquibase pod reports `ErrImageNeverPull`, the image exists in Docker
but is missing from the container runtime for the node running the pod.
For k3d, the node name starts with `k3d-`; import the exact tag into the k3d
cluster and restart the deployment:

```sh
k3d image import "$IMAGE_PREFIX/liquibase:$IMAGE_TAG" -c "$K3D_CLUSTER"
kubectl rollout restart deployment/liquibase -n "$FOUNDATION_NAMESPACE"
kubectl get pods -n "$FOUNDATION_NAMESPACE" -o wide
```

For plain k3s, re-import the exact tag and restart the deployment:

```sh
docker save "$IMAGE_PREFIX/liquibase:$IMAGE_TAG" | sudo k3s ctr -n k8s.io images import -
sudo k3s ctr -n k8s.io images list | grep "$IMAGE_PREFIX/liquibase"
kubectl rollout restart deployment/liquibase -n "$FOUNDATION_NAMESPACE"
kubectl get pods -n "$FOUNDATION_NAMESPACE" -o wide
```

For a multi-node k3s cluster, import the image on every node that may run the
Liquibase pod, or push the image to a registry and use a pull policy other than
`Never`.

### Build Application Layer

Build the application artifacts first. The Dockerfiles copy these built files
into lightweight runtime images:

```sh
command -v go >/dev/null || {
  echo "Go is required to build user-service and user-gateway"
  exit 1
}

mkdir -p build/bin
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/bin/user-service ./user-management/backend/app/user-service/cmd
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/bin/user-gateway ./user-management/backend/app/user-gateway/cmd
for binary in build/bin/user-service build/bin/user-gateway; do
  test -x "$binary" || {
    echo "Missing built binary: $binary"
    exit 1
  }
done

pnpm install --frozen-lockfile
pnpm --filter admin-web-app build
pnpm --filter admin-web-app build-storybook
pnpm --filter storybook build-storybook
for directory in \
  user-management/frontend/app/admin-web-app/dist/admin-web-app/browser \
  user-management/frontend/app/admin-web-app/storybook-static \
  shared-components/frontend/app/storybook/storybook-static; do
  test -d "$directory" || {
    echo "Missing frontend build output: $directory"
    exit 1
  }
done

docker build \
  -f dockerfile \
  --build-arg APP_NAME=user-service \
  -t "$IMAGE_PREFIX/user-service:$IMAGE_TAG" \
  .

docker build \
  -f dockerfile \
  --build-arg APP_NAME=user-gateway \
  -t "$IMAGE_PREFIX/user-gateway:$IMAGE_TAG" \
  .

docker build \
  -f user-management/frontend/app/admin-web-app/Dockerfile \
  -t "$IMAGE_PREFIX/admin-web-app:$IMAGE_TAG" \
  .

docker build \
  -f shared-components/frontend/app/storybook/Dockerfile \
  -t "$IMAGE_PREFIX/storybook:$IMAGE_TAG" \
  .

if command -v k3d >/dev/null 2>&1 && [ "${KUBE_CONTEXT#k3d-}" != "$KUBE_CONTEXT" ]; then
  k3d image import \
    "$IMAGE_PREFIX/user-service:$IMAGE_TAG" \
    "$IMAGE_PREFIX/user-gateway:$IMAGE_TAG" \
    "$IMAGE_PREFIX/admin-web-app:$IMAGE_TAG" \
    "$IMAGE_PREFIX/storybook:$IMAGE_TAG" \
    -c "$K3D_CLUSTER"
elif command -v k3s >/dev/null 2>&1; then
  for image in \
    "$IMAGE_PREFIX/user-service:$IMAGE_TAG" \
    "$IMAGE_PREFIX/user-gateway:$IMAGE_TAG" \
    "$IMAGE_PREFIX/admin-web-app:$IMAGE_TAG" \
    "$IMAGE_PREFIX/storybook:$IMAGE_TAG"; do
    docker image inspect "$image" >/dev/null || {
      echo "Missing Docker image: $image"
      exit 1
    }
  done

  docker save \
    "$IMAGE_PREFIX/user-service:$IMAGE_TAG" \
    "$IMAGE_PREFIX/user-gateway:$IMAGE_TAG" \
    "$IMAGE_PREFIX/admin-web-app:$IMAGE_TAG" \
    "$IMAGE_PREFIX/storybook:$IMAGE_TAG" \
    | sudo k3s ctr -n k8s.io images import -

  sudo k3s ctr -n k8s.io images list | grep "$IMAGE_PREFIX"
fi
```

The workspace explicitly approves the install-time build scripts needed by
Angular and Storybook in `pnpm-workspace.yaml` under `allowBuilds`. If pnpm asks
you to run `pnpm approve-builds`, check that this file still includes
`@parcel/watcher`, `esbuild`, `lmdb` and `msgpackr-extract` with `true` values.

### Deploy Application Layer

```sh
kubectl create namespace "$NAMESPACE" --dry-run=client -o yaml | kubectl apply -f -

mkdir -p build/deploy

cat > build/deploy/shared-components-values-runtime.yaml <<EOF
storybook:
  app:
    image:
      repository: "$IMAGE_PREFIX/storybook"
      tag: "$IMAGE_TAG"
    ingress:
      hosts:
        - host: "$STORYBOOK_HOST"
          paths:
            - path: /
              pathType: Prefix
      tls:
        - secretName: eco-test-ingress-tls
          hosts:
            - "$STORYBOOK_HOST"
EOF

helm upgrade --install shared-components shared-components/deployment \
  --namespace "$NAMESPACE" \
  --values shared-components/deployment/values-test.yaml \
  --values build/deploy/shared-components-values-runtime.yaml

kubectl apply -n "$NAMESPACE" \
  -f user-management/frontend/app/admin-web-app/deployment/runtime-config.test.yaml

cat > build/deploy/user-management-values-runtime.yaml <<EOF
user-service:
  app:
    image:
      repository: "$IMAGE_PREFIX/user-service"
      tag: "$IMAGE_TAG"
    configMap:
      data:
        DB_NAME: "${POSTGRES_DB:-ecoDB}"
        DB_HOST: "$POSTGRES_HOST"
        DB_USER: "${POSTGRES_APP_USER:-app_user}"
    secret:
      stringData:
        DB_PASSWORD: "$APP_PASSWORD"
    ingress:
      hosts:
        - host: "$USER_SERVICE_HOST"
          paths:
            - path: /
              pathType: Prefix
      tls:
        - secretName: eco-test-ingress-tls
          hosts:
            - "$USER_SERVICE_HOST"

user-gateway:
  app:
    image:
      repository: "$IMAGE_PREFIX/user-gateway"
      tag: "$IMAGE_TAG"
    configMap:
      data:
        DB_NAME: "${POSTGRES_DB:-ecoDB}"
        DB_HOST: "$POSTGRES_HOST"
        DB_USER: "${POSTGRES_APP_USER:-app_user}"
        FRONTEND_ORIGIN: "https://$ADMIN_WEB_APP_HOST"
    secret:
      stringData:
        DB_PASSWORD: "$APP_PASSWORD"
    ingress:
      hosts:
        - host: "$USER_GATEWAY_HOST"
          paths:
            - path: /
              pathType: Prefix
      tls:
        - secretName: eco-test-ingress-tls
          hosts:
            - "$USER_GATEWAY_HOST"

admin-web-app:
  app:
    image:
      repository: "$IMAGE_PREFIX/admin-web-app"
      tag: "$IMAGE_TAG"
    ingress:
      hosts:
        - host: "$ADMIN_WEB_APP_HOST"
          paths:
            - path: /
              pathType: Prefix
      tls:
        - secretName: eco-test-ingress-tls
          hosts:
            - "$ADMIN_WEB_APP_HOST"

adminWebAppApiIngress:
  host: "$ADMIN_WEB_APP_HOST"
  tls:
    - secretName: eco-test-ingress-tls
      hosts:
        - "$ADMIN_WEB_APP_HOST"

adminWebAppStoriesIngress:
  host: "$STORYBOOK_HOST"
  tls:
    - secretName: eco-test-ingress-tls
      hosts:
        - "$STORYBOOK_HOST"
EOF

helm upgrade --install user-management user-management/deployment \
  --namespace "$NAMESPACE" \
  --values user-management/deployment/values-test.yaml \
  --values build/deploy/user-management-values-runtime.yaml
```

Check the result:

```sh
helm list -n "$NAMESPACE"
kubectl get all,ingress,configmap -n "$NAMESPACE"
kubectl get all,pvc -n "$FOUNDATION_NAMESPACE"
```

For local browser access, map the ingress address to the configured hostnames.
On the hp-prodesk homelab this is usually Traefik's external IP:

```sh
kubectl get ingress -n "$NAMESPACE"
INGRESS_ADDRESS="$(kubectl get ingress admin-web-app -n "$NAMESPACE" -o jsonpath='{.status.loadBalancer.ingress[0].ip}')"
echo "$INGRESS_ADDRESS $ADMIN_WEB_APP_HOST $STORYBOOK_HOST $USER_GATEWAY_HOST $USER_SERVICE_HOST" | sudo tee -a /etc/hosts
getent hosts "$STORYBOOK_HOST"
```

Then open the deployed apps over HTTP:

```text
http://eco-test.admin-web-app.com
http://eco-test.storybook.com
```

As an alternative to editing `/etc/hosts`, run the local ingress proxy:

```sh
scripts/local-ingress-proxy.py --namespace "$NAMESPACE"
```

Then set your browser or operating system HTTP proxy to `127.0.0.1:18080`.
The proxy discovers the namespace's ingress hosts and forwards requests to the
ingress address while preserving the original `Host` header. It refreshes the
ingress list while it runs, so newly deployed ingress hosts are picked up
without editing `/etc/hosts`.

For feature namespaces, point the proxy at that namespace as well:

```sh
scripts/local-ingress-proxy.py --namespace eco-test --namespace ssoadm2026
```

Or discover ingress hosts from every namespace:

```sh
scripts/local-ingress-proxy.py --all-namespaces
```

If `getent hosts "$STORYBOOK_HOST"` returns nothing, the browser cannot resolve
the hostname; re-run the `/etc/hosts` command above. If the hostname resolves
but the browser still cannot connect, the k3d cluster may not expose Traefik on
the host. Use a local port-forward as a fallback:

```sh
echo "127.0.0.1 $ADMIN_WEB_APP_HOST $STORYBOOK_HOST $USER_GATEWAY_HOST $USER_SERVICE_HOST" | sudo tee -a /etc/hosts
kubectl port-forward -n kube-system svc/traefik 8080:80
```

For the proxy with a port-forward, run:

```sh
scripts/local-ingress-proxy.py --namespace "$NAMESPACE" --target-host 127.0.0.1 --http-port 8080
```

With the proxy configured, keep using the normal ingress URLs. Without the
proxy, open the port-forwarded URL directly:

```text
http://eco-test.storybook.com:8080
```

If you need to remove the releases:

```sh
helm uninstall user-management shared-components -n "$NAMESPACE"
helm uninstall foundation -n "$FOUNDATION_NAMESPACE"
```

## Deploy With GitHub Actions

The workflow at `.github/workflows/deploy.yml` builds deployable images locally,
imports them into k3s when available, and deploys the Helm releases to the
cluster.

### Required GitHub Secrets

Add these in GitHub under `Settings > Secrets and variables > Actions`:

```text
KUBE_CONFIG_B64
POSTGRES_PASSWORD
LIQUIBASE_PASSWORD
APP_PASSWORD
```

`KUBE_CONFIG_B64` is a base64-encoded kubeconfig for a service account that can
deploy into the target namespace:

```sh
base64 -i ~/.kube/config | pbcopy
```

The workflow expects a self-hosted Linux runner with Go, Node 24, pnpm, Docker,
Helm and kubectl installed. If the runner deploys to k3s, it also expects
passwordless `sudo k3s ctr -n k8s.io images import -` so locally built Docker
images are available to the cluster runtime.

### Run The Workflow

The workflow runs automatically on pushes to `main` and on pull requests.

Manual runs support three layer choices:

```text
foundation
application
all
```

Pull request feature deploys only deploy the application layer and require a
source branch named `feature/<10 lowercase alphanumeric chars>`. That suffix
becomes the namespace and ingress prefix.

The deploy order is:

```text
foundation -> shared-components storybook -> user-management
```
