#!/usr/bin/env bash
set -euo pipefail

PORT_FORWARD_PID=""
PORT_FORWARD_LOG=""

stop_port_forward() {
  if [[ -n "$PORT_FORWARD_PID" ]]; then
    kill "$PORT_FORWARD_PID" 2>/dev/null || true
    wait "$PORT_FORWARD_PID" 2>/dev/null || true
    PORT_FORWARD_PID=""
  fi
}

start_port_forward() {
  local service="$1"
  local mapping="$2"

  PORT_FORWARD_LOG="${RUNNER_TEMP:-/tmp}/port-forward-${service//\//-}-$$.log"
  kubectl port-forward -n "$NAMESPACE" "$service" "$mapping" >"$PORT_FORWARD_LOG" 2>&1 &
  PORT_FORWARD_PID="$!"
}

wait_for_port_forward() {
  local service="$1"
  local mapping="$2"
  local description="$3"
  shift 3

  for attempt in {1..30}; do
    if [[ -z "$PORT_FORWARD_PID" ]] || ! kill -0 "$PORT_FORWARD_PID" 2>/dev/null; then
      if [[ -n "$PORT_FORWARD_PID" ]]; then
        wait "$PORT_FORWARD_PID" 2>/dev/null || true
        echo "Port-forward for $description exited; retrying ($attempt/30)" >&2
        if [[ -s "$PORT_FORWARD_LOG" ]]; then
          cat "$PORT_FORWARD_LOG" >&2
        fi
      fi
      start_port_forward "$service" "$mapping"
    fi

    if curl --silent --fail "$@" >/dev/null; then
      return 0
    fi
    sleep 1
  done

  echo "Timed out waiting for port-forwarded $description service" >&2
  if [[ -s "$PORT_FORWARD_LOG" ]]; then
    cat "$PORT_FORWARD_LOG" >&2
  fi
  return 1
}

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
  trap 'stop_port_forward; show_gateway_logs' EXIT
  wait_for_port_forward \
    service/storefront-gateway 18082:8080 storefront-gateway \
    http://127.0.0.1:18082/health

  STOREFRONT_GATEWAY_BASE_URL=http://127.0.0.1:18082 \
  STOREFRONT_FUNCTIONAL_TEST_EMAIL=admin@test.com \
  STOREFRONT_FUNCTIONAL_TEST_PASSWORD=password \
    go test ./online-storefront/backend/app/storefront-gateway/routes \
      -tags integration \
      -run TestDeployedStorefrontGatewayFunctional \
      -v
  stop_port_forward
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
  trap 'stop_port_forward; show_gateway_logs' EXIT
  wait_for_port_forward \
    service/user-gateway 18180:8080 user-gateway \
    --request OPTIONS http://127.0.0.1:18180/api/auth/login

  GATEWAY_BASE_URL=http://127.0.0.1:18180 \
  FUNCTIONAL_TEST_EMAIL=admin@test.com \
  FUNCTIONAL_TEST_PASSWORD=password \
    go test ./user-management/backend/app/user-gateway/routes \
      -tags integration \
      -run TestDeployedGatewayFunctional \
      -v
  stop_port_forward
  trap - EXIT
fi

if [[ "${DEPLOY_PRODUCT_MANAGEMENT:-false}" == "true" ]]; then
  show_gateway_logs() {
    echo "::group::Recent product-gateway logs"
    kubectl logs -n "$NAMESPACE" deployment/product-gateway --tail=120 || true
    echo "::endgroup::"
  }
  trap 'stop_port_forward; show_gateway_logs' EXIT
  wait_for_port_forward \
    service/product-gateway 18081:8080 product-gateway \
    --request OPTIONS http://127.0.0.1:18081/api/auth/login

  PRODUCT_GATEWAY_BASE_URL=http://127.0.0.1:18081 \
  PRODUCT_FUNCTIONAL_TEST_EMAIL=admin@test.com \
  PRODUCT_FUNCTIONAL_TEST_PASSWORD=password \
    go test ./product-management/backend/app/product-gateway/routes \
      -tags integration \
      -run TestDeployedProductGatewayFunctional \
      -v
  stop_port_forward
  trap - EXIT
fi

if [[ "${DEPLOY_ORDER_MANAGEMENT:-false}" == "true" ]]; then
  show_gateway_logs() {
    echo "::group::Recent order-gateway logs"
    kubectl logs -n "$NAMESPACE" deployment/order-gateway --tail=120 || true
    echo "::endgroup::"
  }
  trap 'stop_port_forward; show_gateway_logs' EXIT
  wait_for_port_forward \
    service/order-gateway 18085:8080 order-gateway \
    http://127.0.0.1:18085/health

  ORDER_GATEWAY_BASE_URL=http://127.0.0.1:18085 \
  ORDER_FUNCTIONAL_TEST_EMAIL=admin@test.com \
  ORDER_FUNCTIONAL_TEST_PASSWORD=password \
    go test ./order-management/backend/app/order-gateway/routes \
      -tags integration \
      -run TestDeployedOrderGatewayFunctional \
      -v
  stop_port_forward
  trap - EXIT
fi
