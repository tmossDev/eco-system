#!/usr/bin/env bash
set -euo pipefail

if [[ "${DEPLOY_ONLINE_STOREFRONT:-false}" == "true" ]]; then
  kubectl port-forward -n "$NAMESPACE" service/storefront-gateway 18082:8080 &
  GATEWAY_PORT_FORWARD_PID="$!"
  kubectl port-forward -n "$NAMESPACE" service/storefront-web-app 18083:80 &
  WEB_PORT_FORWARD_PID="$!"
  trap 'kill "$GATEWAY_PORT_FORWARD_PID" "$WEB_PORT_FORWARD_PID" 2>/dev/null || true' EXIT

  for attempt in {1..30}; do
    if curl --silent --fail http://127.0.0.1:18082/health >/dev/null \
      && curl --silent --fail http://127.0.0.1:18083/ >/dev/null; then
      break
    fi
    if [[ "$attempt" == "30" ]]; then
      echo "Timed out waiting for online storefront services" >&2
      exit 1
    fi
    sleep 1
  done

  curl --silent --fail http://127.0.0.1:18082/api/products >/dev/null
  kill "$GATEWAY_PORT_FORWARD_PID" "$WEB_PORT_FORWARD_PID" 2>/dev/null || true
  trap - EXIT
fi

if [[ "${DEPLOY_ONLINE_STOREFRONT:-false}" == "true" ]]; then
  show_gateway_logs() {
    echo "::group::Recent storefront-gateway logs"
    kubectl logs -n "$NAMESPACE" deployment/storefront-gateway --tail=120 || true
    echo "::endgroup::"
  }
  kubectl port-forward -n "$NAMESPACE" service/storefront-gateway 18082:8080 &
  PORT_FORWARD_PID="$!"
  trap 'kill "$PORT_FORWARD_PID" 2>/dev/null || true; show_gateway_logs' EXIT

  for attempt in {1..30}; do
    if curl --silent --fail http://127.0.0.1:18082/health >/dev/null; then
      break
    fi
    if [[ "$attempt" == "30" ]]; then
      echo "Timed out waiting for port-forwarded storefront-gateway service" >&2
      exit 1
    fi
    sleep 1
  done

  STOREFRONT_GATEWAY_BASE_URL=http://127.0.0.1:18082 \
  STOREFRONT_FUNCTIONAL_TEST_EMAIL=admin@test.com \
  STOREFRONT_FUNCTIONAL_TEST_PASSWORD=password \
    go test ./online-storefront/backend/app/storefront-gateway/routes \
      -tags integration \
      -run TestDeployedStorefrontGatewayFunctional \
      -v
  kill "$PORT_FORWARD_PID" 2>/dev/null || true
  trap - EXIT

  show_gateway_logs() {
    echo "::group::Recent storefront-cart-gateway logs"
    kubectl logs -n "$NAMESPACE" deployment/storefront-cart-gateway --tail=120 || true
    echo "::endgroup::"
    echo "::group::Recent storefront-gateway logs"
    kubectl logs -n "$NAMESPACE" deployment/storefront-gateway --tail=120 || true
    echo "::endgroup::"
  }
  kubectl port-forward -n "$NAMESPACE" service/storefront-cart-gateway 18084:8080 &
  CART_PORT_FORWARD_PID="$!"
  kubectl port-forward -n "$NAMESPACE" service/storefront-gateway 18082:8080 &
  STOREFRONT_PORT_FORWARD_PID="$!"
  trap 'kill "$CART_PORT_FORWARD_PID" "$STOREFRONT_PORT_FORWARD_PID" 2>/dev/null || true; show_gateway_logs' EXIT

  for attempt in {1..30}; do
    if curl --silent --fail http://127.0.0.1:18084/health >/dev/null \
      && curl --silent --fail http://127.0.0.1:18082/health >/dev/null; then
      break
    fi
    if [[ "$attempt" == "30" ]]; then
      echo "Timed out waiting for port-forwarded online storefront backend services" >&2
      exit 1
    fi
    sleep 1
  done

  CART_GATEWAY_BASE_URL=http://127.0.0.1:18084 \
  STOREFRONT_GATEWAY_BASE_URL=http://127.0.0.1:18082 \
  CART_FUNCTIONAL_TEST_EMAIL=admin@test.com \
  CART_FUNCTIONAL_TEST_PASSWORD=password \
    go test ./online-storefront/backend/app/cart-gateway/routes \
      -tags integration \
      -run TestDeployedCartGatewayFunctional \
      -v
  kill "$CART_PORT_FORWARD_PID" "$STOREFRONT_PORT_FORWARD_PID" 2>/dev/null || true
  trap - EXIT
fi

if [[ "${DEPLOY_USER_MANAGEMENT:-false}" == "true" ]]; then
  show_gateway_logs() {
    echo "::group::Recent user-gateway logs"
    kubectl logs -n "$NAMESPACE" deployment/user-gateway --tail=120 || true
    echo "::endgroup::"
  }
  kubectl port-forward -n "$NAMESPACE" service/user-gateway 18180:8080 &
  PORT_FORWARD_PID="$!"
  trap 'kill "$PORT_FORWARD_PID" 2>/dev/null || true; show_gateway_logs' EXIT

  for attempt in {1..30}; do
    if curl --silent --fail --request OPTIONS http://127.0.0.1:18180/api/auth/login >/dev/null; then
      break
    fi
    if [[ "$attempt" == "30" ]]; then
      echo "Timed out waiting for port-forwarded user-gateway service" >&2
      exit 1
    fi
    sleep 1
  done

  GATEWAY_BASE_URL=http://127.0.0.1:18180 \
  FUNCTIONAL_TEST_EMAIL=admin@test.com \
  FUNCTIONAL_TEST_PASSWORD=password \
    go test ./user-management/backend/app/user-gateway/routes \
      -tags integration \
      -run TestDeployedGatewayFunctional \
      -v
  kill "$PORT_FORWARD_PID" 2>/dev/null || true
  trap - EXIT
fi

if [[ "${DEPLOY_PRODUCT_MANAGEMENT:-false}" == "true" ]]; then
  show_gateway_logs() {
    echo "::group::Recent product-gateway logs"
    kubectl logs -n "$NAMESPACE" deployment/product-gateway --tail=120 || true
    echo "::endgroup::"
  }
  kubectl port-forward -n "$NAMESPACE" service/product-gateway 18081:8080 &
  PORT_FORWARD_PID="$!"
  trap 'kill "$PORT_FORWARD_PID" 2>/dev/null || true; show_gateway_logs' EXIT

  for attempt in {1..30}; do
    if curl --silent --fail --request OPTIONS http://127.0.0.1:18081/api/auth/login >/dev/null; then
      break
    fi
    if [[ "$attempt" == "30" ]]; then
      echo "Timed out waiting for port-forwarded product-gateway service" >&2
      exit 1
    fi
    sleep 1
  done

  PRODUCT_GATEWAY_BASE_URL=http://127.0.0.1:18081 \
  PRODUCT_FUNCTIONAL_TEST_EMAIL=admin@test.com \
  PRODUCT_FUNCTIONAL_TEST_PASSWORD=password \
    go test ./product-management/backend/app/product-gateway/routes \
      -tags integration \
      -run TestDeployedProductGatewayFunctional \
      -v
  kill "$PORT_FORWARD_PID" 2>/dev/null || true
  trap - EXIT
fi

if [[ "${DEPLOY_ORDER_MANAGEMENT:-false}" == "true" ]]; then
  show_gateway_logs() {
    echo "::group::Recent order-gateway logs"
    kubectl logs -n "$NAMESPACE" deployment/order-gateway --tail=120 || true
    echo "::endgroup::"
  }
  kubectl port-forward -n "$NAMESPACE" service/order-gateway 18085:8080 &
  PORT_FORWARD_PID="$!"
  trap 'kill "$PORT_FORWARD_PID" 2>/dev/null || true; show_gateway_logs' EXIT

  for attempt in {1..30}; do
    if curl --silent --fail http://127.0.0.1:18085/health >/dev/null; then
      break
    fi
    if [[ "$attempt" == "30" ]]; then
      echo "Timed out waiting for port-forwarded order-gateway service" >&2
      exit 1
    fi
    sleep 1
  done

  ORDER_GATEWAY_BASE_URL=http://127.0.0.1:18085 \
  ORDER_FUNCTIONAL_TEST_EMAIL=admin@test.com \
  ORDER_FUNCTIONAL_TEST_PASSWORD=password \
    go test ./order-management/backend/app/order-gateway/routes \
      -tags integration \
      -run TestDeployedOrderGatewayFunctional \
      -v
  kill "$PORT_FORWARD_PID" 2>/dev/null || true
  trap - EXIT
fi
