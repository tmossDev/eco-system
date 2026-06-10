#!/usr/bin/env bash
set -euo pipefail

usage() {
  cat <<EOF
Usage: $0 [--layer foundation|application|all] [--namespace NAME] [--environment-file FILE] [--tag TAG]

Automate local k3d/k3s deployment of the foundation and application layers.

Options:
  --layer            foundation, application, or all (default: all)
  --namespace        application namespace to deploy into (default: value from config or eco-system)
  --environment-file path under config/ (default: hp-prodesk-homelab.txt)
  --tag              image tag to use (default: git SHA)
  -h, --help         show this help message

Required environment variables for foundation deploy:
  POSTGRES_PASSWORD
  LIQUIBASE_PASSWORD
  APP_PASSWORD
EOF
  exit 1
}

LAYER=all
ENVIRONMENT_FILE="hp-prodesk-homelab.txt"
IMAGE_PREFIX="localhost/eco-system"
IMAGE_TAG=""
NAMESPACE=""
FOUNDATION_NAMESPACE=""
K3D_CLUSTER="eco-cluster"

while [[ $# -gt 0 ]]; do
  case "$1" in
    --layer)
      LAYER="$2"; shift 2
      ;;
    --namespace)
      NAMESPACE="$2"; shift 2
      ;;
    --environment-file)
      ENVIRONMENT_FILE="$2"; shift 2
      ;;
    --tag)
      IMAGE_TAG="$2"; shift 2
      ;;
    -h|--help)
      usage
      ;;
    *)
      echo "Unknown option: $1" >&2
      usage
      ;;
  esac
done

if [[ -z "$IMAGE_TAG" ]]; then
  if command -v git >/dev/null 2>&1; then
    IMAGE_TAG="$(git rev-parse --short HEAD)"
  else
    IMAGE_TAG="latest"
  fi
fi

CONFIG_FILE="config/$ENVIRONMENT_FILE"
if [[ ! -f "$CONFIG_FILE" ]]; then
  echo "Missing environment config: $CONFIG_FILE" >&2
  exit 1
fi

set -a
. "$CONFIG_FILE"
set +a

NAMESPACE="${NAMESPACE:-${DEFAULT_APPLICATION_NAMESPACE:-eco-system}}"
FOUNDATION_NAMESPACE="${FOUNDATION_NAMESPACE:-eco-foundation}"

ADMIN_WEB_APP_HOST="${ADMIN_WEB_APP_HOST:-$NAMESPACE.admin-web-app.${BASE_DOMAIN:-com}}"
USER_GATEWAY_HOST="${USER_GATEWAY_HOST:-$NAMESPACE.user-gateway.${BASE_DOMAIN:-com}}"
USER_SERVICE_HOST="${USER_SERVICE_HOST:-$NAMESPACE.user-service.${BASE_DOMAIN:-com}}"
PRODUCT_ADMIN_WEB_APP_HOST="${PRODUCT_ADMIN_WEB_APP_HOST:-$NAMESPACE.product-admin-web-app.${BASE_DOMAIN:-com}}"
PRODUCT_GATEWAY_HOST="${PRODUCT_GATEWAY_HOST:-$NAMESPACE.product-gateway.${BASE_DOMAIN:-com}}"
PRODUCT_SERVICE_HOST="${PRODUCT_SERVICE_HOST:-$NAMESPACE.product-service.${BASE_DOMAIN:-com}}"
ORDER_ADMIN_WEB_APP_HOST="${ORDER_ADMIN_WEB_APP_HOST:-$NAMESPACE.order-admin-web-app.${BASE_DOMAIN:-com}}"
ORDER_GATEWAY_HOST="${ORDER_GATEWAY_HOST:-$NAMESPACE.order-gateway.${BASE_DOMAIN:-com}}"
ORDER_SERVICE_HOST="${ORDER_SERVICE_HOST:-$NAMESPACE.order-service.${BASE_DOMAIN:-com}}"
STOREFRONT_WEB_APP_HOST="${STOREFRONT_WEB_APP_HOST:-$NAMESPACE.storefront-web-app.${BASE_DOMAIN:-com}}"
STORYBOOK_HOST="${STORYBOOK_HOST:-$NAMESPACE.storybook.${BASE_DOMAIN:-com}}"
PRODUCT_BACKEND_API_URL="${PRODUCT_BACKEND_API_URL:-http://product-service:8080}"

if [[ "$LAYER" != "foundation" && "$LAYER" != "application" && "$LAYER" != "all" ]]; then
  echo "Invalid layer: $LAYER" >&2
  usage
fi

if [[ -z "${POSTGRES_PASSWORD:-}" || -z "${LIQUIBASE_PASSWORD:-}" || -z "${APP_PASSWORD:-}" ]]; then
  echo "Environment variables POSTGRES_PASSWORD, LIQUIBASE_PASSWORD, and APP_PASSWORD are required." >&2
  echo "Example: POSTGRES_PASSWORD=... LIQUIBASE_PASSWORD=... APP_PASSWORD=... ./.github/workflows/scripts/local-deploy.sh" >&2
  exit 1
fi

ensure_command() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "Required command not found: $1" >&2
    exit 1
  fi
}

ensure_command helm
ensure_command kubectl
ensure_command docker

import_image() {
  local image="$1"
  if command -v k3d >/dev/null 2>&1; then
    k3d image import "$image" -c "$K3D_CLUSTER"
  elif command -v k3s >/dev/null 2>&1; then
    docker save "$image" | sudo k3s ctr -n k8s.io images import -
  else
    echo "No supported local Kubernetes runtime found. Install k3d or k3s." >&2
    exit 1
  fi
}

create_namespace() {
  local namespace="$1"
  if ! kubectl get namespace "$namespace" >/dev/null 2>&1; then
    kubectl create namespace "$namespace"
  else
    echo "Namespace $namespace already exists"
  fi
}

build_liquibase_image() {
  echo "Building Liquibase image $IMAGE_PREFIX/liquibase:$IMAGE_TAG"
  docker build -f foundation/liquibase/dockerfile -t "$IMAGE_PREFIX/liquibase:$IMAGE_TAG" foundation/liquibase
  import_image "$IMAGE_PREFIX/liquibase:$IMAGE_TAG"
}

deploy_foundation() {
  echo "Deploying foundation to namespace $FOUNDATION_NAMESPACE"
  helm dependency build foundation/postgres/deployment
  helm dependency build foundation/liquibase/deployment
  helm dependency build foundation/deployment

  create_namespace "$FOUNDATION_NAMESPACE"

  helm upgrade --install foundation foundation/deployment \
    --namespace "$FOUNDATION_NAMESPACE" \
    --wait --wait-for-jobs \
    --set postgres.app.env.POSTGRES_DB="${POSTGRES_DB:-ecoDB}" \
    --set liquibase.app.image.repository="$IMAGE_PREFIX/liquibase" \
    --set liquibase.app.image.tag="$IMAGE_TAG" \
    --set liquibase.app.image.pullPolicy=Never \
    --set-string liquibase.app.env.LIQUIBASE_COMMAND_URL="${LIQUIBASE_JDBC_URL:-jdbc:postgresql://postgres:5432/${POSTGRES_DB:-ecoDB}}" \
    --set-string liquibase.app.env.LIQUIBASE_COMMAND_USERNAME="${POSTGRES_LIQUIBASE_USER:-liquibase}" \
    --set-string postgres.app.secret.stringData.POSTGRES_PASSWORD="$POSTGRES_PASSWORD" \
    --set-string postgres.app.secret.stringData.LIQUIBASE_PASSWORD="$LIQUIBASE_PASSWORD" \
    --set-string postgres.app.secret.stringData.APP_PASSWORD="$APP_PASSWORD" \
    --set-string liquibase.app.secret.stringData.LIQUIBASE_COMMAND_PASSWORD="$LIQUIBASE_PASSWORD" \
    --set-string liquibase.app.secret.stringData.LIQUIBASE_PASSWORD="$LIQUIBASE_PASSWORD"
}

build_application_images() {
  echo "Building application images with tag $IMAGE_TAG"
  docker build -f dockerfile --build-arg APP_NAME=user-service -t "$IMAGE_PREFIX/user-service:$IMAGE_TAG" .
  docker build -f dockerfile --build-arg APP_NAME=user-gateway -t "$IMAGE_PREFIX/user-gateway:$IMAGE_TAG" .
  docker build -f dockerfile --build-arg APP_NAME=product-service -t "$IMAGE_PREFIX/product-service:$IMAGE_TAG" .
  docker build -f dockerfile --build-arg APP_NAME=product-gateway -t "$IMAGE_PREFIX/product-gateway:$IMAGE_TAG" .
  docker build -f dockerfile --build-arg APP_NAME=order-service -t "$IMAGE_PREFIX/order-service:$IMAGE_TAG" .
  docker build -f dockerfile --build-arg APP_NAME=order-gateway -t "$IMAGE_PREFIX/order-gateway:$IMAGE_TAG" .
  docker build -f dockerfile --build-arg APP_NAME=cart-gateway -t "$IMAGE_PREFIX/cart-gateway:$IMAGE_TAG" .
  docker build -f dockerfile --build-arg APP_NAME=storefront-gateway -t "$IMAGE_PREFIX/storefront-gateway:$IMAGE_TAG" .
  docker build -f user-management/frontend/app/admin-web-app/Dockerfile -t "$IMAGE_PREFIX/admin-web-app:$IMAGE_TAG" .
  docker build -f product-management/frontend/app/product-admin-web-app/Dockerfile -t "$IMAGE_PREFIX/product-admin-web-app:$IMAGE_TAG" .
  docker build -f order-management/frontend/app/order-admin-web-app/Dockerfile -t "$IMAGE_PREFIX/order-admin-web-app:$IMAGE_TAG" .
  docker build -f online-storefront/frontend/app/storefront-web-app/Dockerfile -t "$IMAGE_PREFIX/storefront-web-app:$IMAGE_TAG" .
  docker build -f shared-components/frontend/app/storybook/Dockerfile -t "$IMAGE_PREFIX/storybook:$IMAGE_TAG" .

  import_image "$IMAGE_PREFIX/user-service:$IMAGE_TAG"
  import_image "$IMAGE_PREFIX/user-gateway:$IMAGE_TAG"
  import_image "$IMAGE_PREFIX/product-service:$IMAGE_TAG"
  import_image "$IMAGE_PREFIX/product-gateway:$IMAGE_TAG"
  import_image "$IMAGE_PREFIX/order-service:$IMAGE_TAG"
  import_image "$IMAGE_PREFIX/order-gateway:$IMAGE_TAG"
  import_image "$IMAGE_PREFIX/cart-gateway:$IMAGE_TAG"
  import_image "$IMAGE_PREFIX/storefront-gateway:$IMAGE_TAG"
  import_image "$IMAGE_PREFIX/admin-web-app:$IMAGE_TAG"
  import_image "$IMAGE_PREFIX/product-admin-web-app:$IMAGE_TAG"
  import_image "$IMAGE_PREFIX/order-admin-web-app:$IMAGE_TAG"
  import_image "$IMAGE_PREFIX/storefront-web-app:$IMAGE_TAG"
  import_image "$IMAGE_PREFIX/storybook:$IMAGE_TAG"
}

deploy_application() {
  echo "Deploying application to namespace $NAMESPACE"
  helm dependency build shared-components/frontend/app/storybook/deployment
  helm dependency build user-management/backend/app/user-service/deployment
  helm dependency build user-management/backend/app/user-gateway/deployment
  helm dependency build user-management/frontend/app/admin-web-app/deployment
  helm dependency build product-management/backend/app/product-service/deployment
  helm dependency build product-management/backend/app/product-gateway/deployment
  helm dependency build product-management/frontend/app/product-admin-web-app/deployment
  helm dependency build order-management/backend/app/order-service/deployment
  helm dependency build order-management/backend/app/order-gateway/deployment
  helm dependency build order-management/frontend/app/order-admin-web-app/deployment
  helm dependency build online-storefront/backend/app/cart-gateway/deployment
  helm dependency build online-storefront/backend/app/storefront-gateway/deployment
  helm dependency build online-storefront/frontend/app/storefront-web-app/deployment
  helm dependency build shared-components/deployment
  helm dependency build user-management/deployment
  helm dependency build product-management/deployment
  helm dependency build order-management/deployment
  helm dependency build online-storefront/deployment

  create_namespace "$NAMESPACE"

  printf '{"BACKEND_API_URL":"","ENABLE_MOCK_API":false}\n' > runtime-config.json

  helm upgrade --install shared-components shared-components/deployment \
    --namespace "$NAMESPACE" \
    --set storybook.app.image.repository="$IMAGE_PREFIX/storybook" \
    --set storybook.app.image.tag="$IMAGE_TAG" \
    --set storybook.app.image.pullPolicy=Never \
    --set storybook.app.ingress.enabled=true \
    --set storybook.app.ingress.className="${INGRESS_CLASS_NAME:-traefik}" \
    --set storybook.app.ingress.hosts[0].host="${STORYBOOK_HOST:-$NAMESPACE.storybook.${BASE_DOMAIN:-com}}" \
    --set storybook.app.ingress.hosts[0].paths[0].path=/ \
    --set storybook.app.ingress.hosts[0].paths[0].pathType=Prefix

  helm upgrade --install user-management user-management/deployment \
    --namespace "$NAMESPACE" \
    --set user-service.app.image.repository="$IMAGE_PREFIX/user-service" \
    --set user-service.app.image.tag="$IMAGE_TAG" \
    --set user-service.app.image.pullPolicy=Never \
    --set user-service.app.configMap.enabled=true \
    --set-string user-service.app.configMap.data.DB_DIALECT=postgresql \
    --set-string user-service.app.configMap.data.DB_NAME="${POSTGRES_DB:-ecoDB}" \
    --set-string user-service.app.configMap.data.DB_HOST="${POSTGRES_HOST:-postgres.${FOUNDATION_NAMESPACE}.svc.cluster.local}" \
    --set-string user-service.app.configMap.data.DB_PORT="${POSTGRES_PORT:-5432}" \
    --set-string user-service.app.configMap.data.DB_USER="${POSTGRES_APP_USER:-app_user}" \
    --set user-service.app.secret.enabled=true \
    --set-string user-service.app.secret.stringData.DB_PASSWORD="$APP_PASSWORD" \
    --set user-service.app.ingress.enabled=true \
    --set user-service.app.ingress.className="${INGRESS_CLASS_NAME:-traefik}" \
    --set user-service.app.ingress.hosts[0].host="${USER_SERVICE_HOST:-$NAMESPACE.user-service.${BASE_DOMAIN:-com}}" \
    --set user-service.app.ingress.hosts[0].paths[0].path=/ \
    --set user-service.app.ingress.hosts[0].paths[0].pathType=Prefix \
    --set user-gateway.app.image.repository="$IMAGE_PREFIX/user-gateway" \
    --set user-gateway.app.image.tag="$IMAGE_TAG" \
    --set user-gateway.app.image.pullPolicy=Never \
    --set user-gateway.app.configMap.enabled=true \
    --set-string user-gateway.app.configMap.data.DB_DIALECT=postgresql \
    --set-string user-gateway.app.configMap.data.DB_NAME="${POSTGRES_DB:-ecoDB}" \
    --set-string user-gateway.app.configMap.data.DB_HOST="${POSTGRES_HOST:-postgres.${FOUNDATION_NAMESPACE}.svc.cluster.local}" \
    --set-string user-gateway.app.configMap.data.DB_PORT="${POSTGRES_PORT:-5432}" \
    --set-string user-gateway.app.configMap.data.DB_USER="${POSTGRES_APP_USER:-app_user}" \
    --set-string user-gateway.app.configMap.data.USER_SERVICE_URL="${USER_SERVICE_INTERNAL_URL:-http://user-service:8080}" \
    --set-string user-gateway.app.configMap.data.FRONTEND_ORIGIN="http://${ADMIN_WEB_APP_HOST:-$NAMESPACE.admin-web-app.${BASE_DOMAIN:-com}}" \
    --set user-gateway.app.secret.enabled=true \
    --set-string user-gateway.app.secret.stringData.DB_PASSWORD="$APP_PASSWORD" \
    --set user-gateway.app.ingress.enabled=true \
    --set user-gateway.app.ingress.className="${INGRESS_CLASS_NAME:-traefik}" \
    --set user-gateway.app.ingress.hosts[0].host="${USER_GATEWAY_HOST:-$NAMESPACE.user-gateway.${BASE_DOMAIN:-com}}" \
    --set user-gateway.app.ingress.hosts[0].paths[0].path=/ \
    --set user-gateway.app.ingress.hosts[0].paths[0].pathType=Prefix \
    --set admin-web-app.app.image.repository="$IMAGE_PREFIX/admin-web-app" \
    --set admin-web-app.app.image.tag="$IMAGE_TAG" \
    --set admin-web-app.app.image.pullPolicy=Never \
    --set admin-web-app.app.configMap.enabled=true \
    --set admin-web-app.app.configMap.envFrom=false \
    --set-file admin-web-app.app.configMap.data.config\\.json=runtime-config.json \
    --set admin-web-app.app.ingress.enabled=true \
    --set admin-web-app.app.ingress.className="${INGRESS_CLASS_NAME:-traefik}" \
    --set admin-web-app.app.ingress.hosts[0].host="${ADMIN_WEB_APP_HOST:-$NAMESPACE.admin-web-app.${BASE_DOMAIN:-com}}" \
    --set admin-web-app.app.ingress.hosts[0].paths[0].path=/ \
    --set admin-web-app.app.ingress.hosts[0].paths[0].pathType=Prefix \
    --set adminWebAppApiIngress.enabled=true \
    --set adminWebAppApiIngress.className="${INGRESS_CLASS_NAME:-traefik}" \
    --set adminWebAppApiIngress.host="${ADMIN_WEB_APP_HOST:-$NAMESPACE.admin-web-app.${BASE_DOMAIN:-com}}" \
    --set adminWebAppApiIngress.path=/api \
    --set adminWebAppApiIngress.pathType=Prefix \
    --set adminWebAppApiIngress.serviceName=user-gateway \
    --set adminWebAppApiIngress.servicePortName=http

  printf '{"BACKEND_API_URL":"%s","ENABLE_MOCK_API":false}\n' "http://${PRODUCT_SERVICE_INTERNAL_URL:-product-service:8080}" > product-runtime-config.json

  helm upgrade --install product-management product-management/deployment \
    --namespace "$NAMESPACE" \
    --set product-service.app.image.repository="$IMAGE_PREFIX/product-service" \
    --set product-service.app.image.tag="$IMAGE_TAG" \
    --set product-service.app.image.pullPolicy=Never \
    --set product-service.app.configMap.enabled=true \
    --set-string product-service.app.configMap.data.DB_DIALECT=postgresql \
    --set-string product-service.app.configMap.data.DB_NAME="${POSTGRES_DB:-ecoDB}" \
    --set-string product-service.app.configMap.data.DB_HOST="${POSTGRES_HOST:-postgres.${FOUNDATION_NAMESPACE}.svc.cluster.local}" \
    --set-string product-service.app.configMap.data.DB_PORT="${POSTGRES_PORT:-5432}" \
    --set-string product-service.app.configMap.data.DB_USER="${POSTGRES_APP_USER:-app_user}" \
    --set product-service.app.secret.enabled=true \
    --set-string product-service.app.secret.stringData.DB_PASSWORD="$APP_PASSWORD" \
    --set product-service.app.ingress.enabled=true \
    --set product-service.app.ingress.className="${INGRESS_CLASS_NAME:-traefik}" \
    --set product-service.app.ingress.hosts[0].host="${PRODUCT_SERVICE_HOST:-$NAMESPACE.product-service.${BASE_DOMAIN:-com}}" \
    --set product-service.app.ingress.hosts[0].paths[0].path=/ \
    --set product-service.app.ingress.hosts[0].paths[0].pathType=Prefix \
    --set product-gateway.app.image.repository="$IMAGE_PREFIX/product-gateway" \
    --set product-gateway.app.image.tag="$IMAGE_TAG" \
    --set product-gateway.app.image.pullPolicy=Never \
    --set product-gateway.app.configMap.enabled=true \
    --set-string product-gateway.app.configMap.data.DB_DIALECT=postgresql \
    --set-string product-gateway.app.configMap.data.DB_NAME="${POSTGRES_DB:-ecoDB}" \
    --set-string product-gateway.app.configMap.data.DB_HOST="${POSTGRES_HOST:-postgres.${FOUNDATION_NAMESPACE}.svc.cluster.local}" \
    --set-string product-gateway.app.configMap.data.DB_PORT="${POSTGRES_PORT:-5432}" \
    --set-string product-gateway.app.configMap.data.DB_USER="${POSTGRES_APP_USER:-app_user}" \
    --set-string product-gateway.app.configMap.data.PRODUCT_SERVICE_URL="${PRODUCT_SERVICE_INTERNAL_URL:-http://product-service:8080}" \
    --set-string product-gateway.app.configMap.data.USER_SERVICE_URL="${USER_SERVICE_INTERNAL_URL:-http://user-service:8080}" \
    --set-string product-gateway.app.configMap.data.FRONTEND_ORIGIN="http://${PRODUCT_ADMIN_WEB_APP_HOST}" \
    --set product-gateway.app.secret.enabled=true \
    --set-string product-gateway.app.secret.stringData.DB_PASSWORD="$APP_PASSWORD" \
    --set product-gateway.app.ingress.enabled=true \
    --set product-gateway.app.ingress.className="${INGRESS_CLASS_NAME:-traefik}" \
    --set product-gateway.app.ingress.hosts[0].host="${PRODUCT_GATEWAY_HOST:-$NAMESPACE.product-gateway.${BASE_DOMAIN:-com}}" \
    --set product-gateway.app.ingress.hosts[0].paths[0].path=/ \
    --set product-gateway.app.ingress.hosts[0].paths[0].pathType=Prefix \
    --set product-admin-web-app.app.image.repository="$IMAGE_PREFIX/product-admin-web-app" \
    --set product-admin-web-app.app.image.tag="$IMAGE_TAG" \
    --set product-admin-web-app.app.image.pullPolicy=Never \
    --set product-admin-web-app.app.configMap.enabled=true \
    --set product-admin-web-app.app.configMap.envFrom=false \
    --set-file product-admin-web-app.app.configMap.data.config\\.json=product-runtime-config.json \
    --set product-admin-web-app.app.ingress.enabled=true \
    --set product-admin-web-app.app.ingress.className="${INGRESS_CLASS_NAME:-traefik}" \
    --set product-admin-web-app.app.ingress.hosts[0].host="${PRODUCT_ADMIN_WEB_APP_HOST:-$NAMESPACE.product-admin-web-app.${BASE_DOMAIN:-com}}" \
    --set product-admin-web-app.app.ingress.hosts[0].paths[0].path=/ \
    --set product-admin-web-app.app.ingress.hosts[0].paths[0].pathType=Prefix \
    --set productAdminWebAppApiIngress.enabled=true \
    --set productAdminWebAppApiIngress.className="${INGRESS_CLASS_NAME:-traefik}" \
    --set productAdminWebAppApiIngress.host="${PRODUCT_ADMIN_WEB_APP_HOST:-$NAMESPACE.product-admin-web-app.${BASE_DOMAIN:-com}}" \
    --set productAdminWebAppApiIngress.path=/api \
    --set productAdminWebAppApiIngress.pathType=Prefix \
    --set productAdminWebAppApiIngress.serviceName=product-gateway \
    --set productAdminWebAppApiIngress.servicePortName=http

  printf '{"BACKEND_API_URL":"","ENABLE_MOCK_API":false}\n' > order-runtime-config.json

  helm upgrade --install order-management order-management/deployment \
    --namespace "$NAMESPACE" \
    --set order-service.app.image.repository="$IMAGE_PREFIX/order-service" \
    --set order-service.app.image.tag="$IMAGE_TAG" \
    --set order-service.app.image.pullPolicy=Never \
    --set order-service.app.configMap.enabled=true \
    --set-string order-service.app.configMap.data.DB_DIALECT=postgresql \
    --set-string order-service.app.configMap.data.DB_NAME="${POSTGRES_DB:-ecoDB}" \
    --set-string order-service.app.configMap.data.DB_HOST="${POSTGRES_HOST:-postgres.${FOUNDATION_NAMESPACE}.svc.cluster.local}" \
    --set-string order-service.app.configMap.data.DB_PORT="${POSTGRES_PORT:-5432}" \
    --set-string order-service.app.configMap.data.DB_USER="${POSTGRES_APP_USER:-app_user}" \
    --set order-service.app.secret.enabled=true \
    --set-string order-service.app.secret.stringData.DB_PASSWORD="$APP_PASSWORD" \
    --set order-service.app.ingress.enabled=true \
    --set order-service.app.ingress.className="${INGRESS_CLASS_NAME:-traefik}" \
    --set order-service.app.ingress.hosts[0].host="$ORDER_SERVICE_HOST" \
    --set order-service.app.ingress.hosts[0].paths[0].path=/ \
    --set order-service.app.ingress.hosts[0].paths[0].pathType=Prefix \
    --set order-gateway.app.image.repository="$IMAGE_PREFIX/order-gateway" \
    --set order-gateway.app.image.tag="$IMAGE_TAG" \
    --set order-gateway.app.image.pullPolicy=Never \
    --set order-gateway.app.configMap.enabled=true \
    --set-string order-gateway.app.configMap.data.DB_DIALECT=postgresql \
    --set-string order-gateway.app.configMap.data.DB_NAME="${POSTGRES_DB:-ecoDB}" \
    --set-string order-gateway.app.configMap.data.DB_HOST="${POSTGRES_HOST:-postgres.${FOUNDATION_NAMESPACE}.svc.cluster.local}" \
    --set-string order-gateway.app.configMap.data.DB_PORT="${POSTGRES_PORT:-5432}" \
    --set-string order-gateway.app.configMap.data.DB_USER="${POSTGRES_APP_USER:-app_user}" \
    --set-string order-gateway.app.configMap.data.USER_SERVICE_URL="${USER_SERVICE_INTERNAL_URL:-http://user-service:8080}" \
    --set-string order-gateway.app.configMap.data.FRONTEND_ORIGIN="http://${ORDER_ADMIN_WEB_APP_HOST}" \
    --set order-gateway.app.secret.enabled=true \
    --set-string order-gateway.app.secret.stringData.DB_PASSWORD="$APP_PASSWORD" \
    --set order-gateway.app.ingress.enabled=true \
    --set order-gateway.app.ingress.className="${INGRESS_CLASS_NAME:-traefik}" \
    --set order-gateway.app.ingress.hosts[0].host="$ORDER_GATEWAY_HOST" \
    --set order-gateway.app.ingress.hosts[0].paths[0].path=/ \
    --set order-gateway.app.ingress.hosts[0].paths[0].pathType=Prefix \
    --set order-admin-web-app.app.image.repository="$IMAGE_PREFIX/order-admin-web-app" \
    --set order-admin-web-app.app.image.tag="$IMAGE_TAG" \
    --set order-admin-web-app.app.image.pullPolicy=Never \
    --set order-admin-web-app.app.configMap.enabled=true \
    --set order-admin-web-app.app.configMap.envFrom=false \
    --set-file order-admin-web-app.app.configMap.data.config\\.json=order-runtime-config.json \
    --set order-admin-web-app.app.ingress.enabled=true \
    --set order-admin-web-app.app.ingress.className="${INGRESS_CLASS_NAME:-traefik}" \
    --set order-admin-web-app.app.ingress.hosts[0].host="$ORDER_ADMIN_WEB_APP_HOST" \
    --set order-admin-web-app.app.ingress.hosts[0].paths[0].path=/ \
    --set order-admin-web-app.app.ingress.hosts[0].paths[0].pathType=Prefix \
    --set orderAdminWebAppApiIngress.enabled=true \
    --set orderAdminWebAppApiIngress.className="${INGRESS_CLASS_NAME:-traefik}" \
    --set orderAdminWebAppApiIngress.host="$ORDER_ADMIN_WEB_APP_HOST" \
    --set orderAdminWebAppApiIngress.path=/api \
    --set orderAdminWebAppApiIngress.pathType=Prefix \
    --set orderAdminWebAppApiIngress.serviceName=order-gateway \
    --set orderAdminWebAppApiIngress.servicePortName=http

  helm upgrade --install online-storefront online-storefront/deployment \
    --namespace "$NAMESPACE" \
    --set cart-gateway.app.image.repository="$IMAGE_PREFIX/cart-gateway" \
    --set cart-gateway.app.image.tag="$IMAGE_TAG" \
    --set cart-gateway.app.image.pullPolicy=Never \
    --set cart-gateway.app.configMap.enabled=true \
    --set-string cart-gateway.app.configMap.data.DB_DIALECT=postgresql \
    --set-string cart-gateway.app.configMap.data.DB_NAME="${POSTGRES_DB:-ecoDB}" \
    --set-string cart-gateway.app.configMap.data.DB_HOST="${POSTGRES_HOST:-postgres.${FOUNDATION_NAMESPACE}.svc.cluster.local}" \
    --set-string cart-gateway.app.configMap.data.DB_PORT="${POSTGRES_PORT:-5432}" \
    --set-string cart-gateway.app.configMap.data.DB_USER="${POSTGRES_APP_USER:-app_user}" \
    --set cart-gateway.app.secret.enabled=true \
    --set-string cart-gateway.app.secret.stringData.DB_PASSWORD="$APP_PASSWORD" \
    --set storefront-gateway.app.image.repository="$IMAGE_PREFIX/storefront-gateway" \
    --set storefront-gateway.app.image.tag="$IMAGE_TAG" \
    --set storefront-gateway.app.image.pullPolicy=Never \
    --set-string storefront-gateway.app.env.USER_SERVICE_URL="${USER_SERVICE_INTERNAL_URL:-http://user-service:8080}" \
    --set-string storefront-gateway.app.env.PRODUCT_SERVICE_URL="${PRODUCT_SERVICE_INTERNAL_URL:-http://product-service:8080}" \
    --set-string storefront-gateway.app.env.PRODUCT_GATEWAY_URL="${PRODUCT_GATEWAY_INTERNAL_URL:-http://product-gateway:8080}" \
    --set-string storefront-gateway.app.env.FRONTEND_ORIGIN="http://${STOREFRONT_WEB_APP_HOST}" \
    --set storefront-web-app.app.image.repository="$IMAGE_PREFIX/storefront-web-app" \
    --set storefront-web-app.app.image.tag="$IMAGE_TAG" \
    --set storefront-web-app.app.image.pullPolicy=Never \
    --set storefront-web-app.app.ingress.enabled=true \
    --set storefront-web-app.app.ingress.className="${INGRESS_CLASS_NAME:-traefik}" \
    --set storefront-web-app.app.ingress.hosts[0].host="$STOREFRONT_WEB_APP_HOST" \
    --set storefront-web-app.app.ingress.hosts[0].paths[0].path=/ \
    --set storefront-web-app.app.ingress.hosts[0].paths[0].pathType=Prefix \
    --set cartGatewayApiIngress.enabled=true \
    --set cartGatewayApiIngress.className="${INGRESS_CLASS_NAME:-traefik}" \
    --set cartGatewayApiIngress.host="$STOREFRONT_WEB_APP_HOST" \
    --set cartGatewayApiIngress.paths[0].path=/api/cart \
    --set cartGatewayApiIngress.paths[0].pathType=Prefix \
    --set cartGatewayApiIngress.paths[1].path=/api/orders \
    --set cartGatewayApiIngress.paths[1].pathType=Prefix \
    --set cartGatewayApiIngress.serviceName=storefront-cart-gateway \
    --set cartGatewayApiIngress.servicePortName=http \
    --set storefrontWebAppApiIngress.enabled=true \
    --set storefrontWebAppApiIngress.className="${INGRESS_CLASS_NAME:-traefik}" \
    --set storefrontWebAppApiIngress.host="$STOREFRONT_WEB_APP_HOST" \
    --set storefrontWebAppApiIngress.path=/api \
    --set storefrontWebAppApiIngress.pathType=Prefix \
    --set storefrontWebAppApiIngress.serviceName=storefront-gateway \
    --set storefrontWebAppApiIngress.servicePortName=http
}

main() {
  if [[ "$LAYER" == "foundation" || "$LAYER" == "all" ]]; then
    build_liquibase_image
    deploy_foundation
  fi

  if [[ "$LAYER" == "application" || "$LAYER" == "all" ]]; then
    build_application_images
    deploy_application
  fi
}

main
