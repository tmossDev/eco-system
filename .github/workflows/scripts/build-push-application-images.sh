#!/usr/bin/env bash
set -euo pipefail

require_var() {
  local name="$1"
  if [[ -z "${!name:-}" ]]; then
    echo "Missing required environment variable: $name" >&2
    exit 1
  fi
}

require_var IMAGE_PREFIX
require_var IMAGE_TAG
require_var NEXUS_DOCKER_PUSH_REGISTRY
require_var CI_CD_NAMESPACE
require_var NEXUS_ADMIN_PASSWORD

build_image() {
  local dockerfile="$1"
  local image="$2"
  shift 2
  docker build -f "$dockerfile" "$@" -t "$IMAGE_PREFIX/$image:$IMAGE_TAG" .
}

if [[ "${BUILD_USER_MANAGEMENT:-false}" == "true" ]]; then
  build_image dockerfile user-service --build-arg APP_NAME=user-service
  build_image dockerfile user-gateway --build-arg APP_NAME=user-gateway
  build_image user-management/frontend/app/admin-web-app/Dockerfile admin-web-app
fi

if [[ "${BUILD_PRODUCT_MANAGEMENT:-false}" == "true" ]]; then
  build_image dockerfile product-service --build-arg APP_NAME=product-service
  build_image dockerfile product-gateway --build-arg APP_NAME=product-gateway
  build_image product-management/frontend/app/product-admin-web-app/Dockerfile product-admin-web-app
fi

if [[ "${BUILD_ORDER_MANAGEMENT:-false}" == "true" ]]; then
  build_image dockerfile order-service --build-arg APP_NAME=order-service
  build_image dockerfile order-gateway --build-arg APP_NAME=order-gateway
  build_image order-management/frontend/app/order-admin-web-app/Dockerfile order-admin-web-app
fi

if [[ "${BUILD_ONLINE_STOREFRONT:-false}" == "true" ]]; then
  build_image dockerfile cart-gateway --build-arg APP_NAME=cart-gateway
  build_image dockerfile storefront-gateway --build-arg APP_NAME=storefront-gateway
  build_image online-storefront/frontend/app/storefront-web-app/Dockerfile storefront-web-app
fi

if [[ "${BUILD_STORYBOOK:-false}" == "true" ]]; then
  build_image shared-components/frontend/app/storybook/Dockerfile storybook
fi

for attempt in {1..60}; do
  endpoints="$(kubectl get endpoints nexus -n "$CI_CD_NAMESPACE" -o jsonpath='{.subsets[*].addresses[*].ip}' 2>/dev/null || true)"
  if [[ -n "$endpoints" ]]; then
    break
  fi
  if [[ "$attempt" == "60" ]]; then
    echo "Timed out waiting for existing Nexus service endpoints in namespace $CI_CD_NAMESPACE." >&2
    kubectl get deploy,pod,svc,endpoints -n "$CI_CD_NAMESPACE" -l app.kubernetes.io/name=nexus -o wide || true
    exit 1
  fi
  sleep 2
done

kubectl port-forward -n "$CI_CD_NAMESPACE" svc/nexus "${NEXUS_DOCKER_PUSH_REGISTRY##*:}:5000" &
NEXUS_DOCKER_PORT_FORWARD_PID="$!"
trap 'kill "$NEXUS_DOCKER_PORT_FORWARD_PID" 2>/dev/null || true' EXIT

for attempt in {1..60}; do
  status="$(curl --silent --output /dev/null --write-out '%{http_code}' "http://$NEXUS_DOCKER_PUSH_REGISTRY/v2/" || true)"
  if [[ "$status" == "200" || "$status" == "401" ]]; then
    break
  fi
  if [[ "$attempt" == "60" ]]; then
    echo "Timed out waiting for local Nexus Docker registry forward at $NEXUS_DOCKER_PUSH_REGISTRY" >&2
    exit 1
  fi
  sleep 2
done

echo "$NEXUS_ADMIN_PASSWORD" | docker login "$NEXUS_DOCKER_PUSH_REGISTRY" --username admin --password-stdin

push_image() {
  local name="$1"
  docker push "$IMAGE_PREFIX/$name:$IMAGE_TAG"
  if [[ "${TAG_LATEST:-false}" == "true" ]]; then
    docker tag "$IMAGE_PREFIX/$name:$IMAGE_TAG" "$IMAGE_PREFIX/$name:latest"
    docker push "$IMAGE_PREFIX/$name:latest"
    docker image rm "$IMAGE_PREFIX/$name:latest" >/dev/null 2>&1 || true
  fi
  docker image rm "$IMAGE_PREFIX/$name:$IMAGE_TAG" >/dev/null 2>&1 || true
}

if [[ "${BUILD_USER_MANAGEMENT:-false}" == "true" ]]; then
  push_image user-service
  push_image user-gateway
  push_image admin-web-app
fi

if [[ "${BUILD_PRODUCT_MANAGEMENT:-false}" == "true" ]]; then
  push_image product-service
  push_image product-gateway
  push_image product-admin-web-app
fi

if [[ "${BUILD_ORDER_MANAGEMENT:-false}" == "true" ]]; then
  push_image order-service
  push_image order-gateway
  push_image order-admin-web-app
fi

if [[ "${BUILD_ONLINE_STOREFRONT:-false}" == "true" ]]; then
  push_image cart-gateway
  push_image storefront-gateway
  push_image storefront-web-app
fi

if [[ "${BUILD_STORYBOOK:-false}" == "true" ]]; then
  push_image storybook
fi
