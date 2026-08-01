#!/usr/bin/env bash
set -euo pipefail

PORT_FORWARD_PIDS=()
PORT_FORWARD_LOGS=()

stop_port_forwards() {
  if [[ "${#PORT_FORWARD_PIDS[@]}" -gt 0 ]]; then
    kill "${PORT_FORWARD_PIDS[@]}" 2>/dev/null || true
    wait "${PORT_FORWARD_PIDS[@]}" 2>/dev/null || true
  fi
  PORT_FORWARD_PIDS=()
  PORT_FORWARD_LOGS=()
}

start_port_forward() {
  local service="$1"
  local mapping="$2"
  local log="${RUNNER_TEMP:-/tmp}/port-forward-${service}-$$.log"

  kubectl port-forward -n "$NAMESPACE" "service/$service" "$mapping" >"$log" 2>&1 &
  PORT_FORWARD_PIDS+=("$!")
  PORT_FORWARD_LOGS+=("$log")
}

show_port_forward_logs() {
  local log
  for log in "${PORT_FORWARD_LOGS[@]}"; do
    if [[ -s "$log" ]]; then
      cat "$log" >&2
    fi
  done
}

wait_for_services() {
  local description="$1"
  shift

  for attempt in {1..30}; do
    local pid
    for pid in "${PORT_FORWARD_PIDS[@]}"; do
      if ! kill -0 "$pid" 2>/dev/null; then
        echo "Port-forward exited while waiting for $description" >&2
        show_port_forward_logs
        return 1
      fi
    done

    local ready=true
    local url
    for url in "$@"; do
      if ! curl --silent "$url" >/dev/null; then
        ready=false
      fi
    done
    if [[ "$ready" == "true" ]]; then
      return 0
    fi
    sleep 1
  done

  echo "Timed out waiting for port-forwarded $description" >&2
  show_port_forward_logs
  return 1
}

trap stop_port_forwards EXIT

if [[ "${DEPLOY_ONLINE_STOREFRONT:-false}" == "true" ]]; then
  start_port_forward storefront-gateway 18082:8080
  wait_for_services storefront-gateway http://127.0.0.1:18082/health

  STOREFRONT_GATEWAY_BASE_URL=http://127.0.0.1:18082 \
  STOREFRONT_FUNCTIONAL_TEST_EMAIL=admin@test.com \
  STOREFRONT_FUNCTIONAL_TEST_PASSWORD=password \
    go test ./online-storefront/backend/app/storefront-gateway/routes \
      -tags integration \
      -run TestDeployedStorefrontGatewayFunctional \
      -v
  stop_port_forwards

  start_port_forward storefront-cart-gateway 18084:8080
  start_port_forward storefront-gateway 18082:8080
  wait_for_services online-storefront-backends \
    http://127.0.0.1:18084/health \
    http://127.0.0.1:18082/health

  CART_GATEWAY_BASE_URL=http://127.0.0.1:18084 \
  STOREFRONT_GATEWAY_BASE_URL=http://127.0.0.1:18082 \
  CART_FUNCTIONAL_TEST_EMAIL=admin@test.com \
  CART_FUNCTIONAL_TEST_PASSWORD=password \
    go test ./online-storefront/backend/app/cart-gateway/routes \
      -tags integration \
      -run TestDeployedCartGatewayFunctional \
      -v
  stop_port_forwards
fi

if [[ "${DEPLOY_USER_MANAGEMENT:-false}" == "true" ]]; then
  start_port_forward user-gateway 18180:8080
  wait_for_services user-gateway http://127.0.0.1:18180/health

  GATEWAY_BASE_URL=http://127.0.0.1:18180 \
  FUNCTIONAL_TEST_EMAIL=admin@test.com \
  FUNCTIONAL_TEST_PASSWORD=password \
    go test ./user-management/backend/app/user-gateway/routes \
      -tags integration \
      -run TestDeployedGatewayFunctional \
      -v
  stop_port_forwards
fi

if [[ "${DEPLOY_PRODUCT_MANAGEMENT:-false}" == "true" ]]; then
  start_port_forward product-gateway 18081:8080
  wait_for_services product-gateway http://127.0.0.1:18081/health

  PRODUCT_GATEWAY_BASE_URL=http://127.0.0.1:18081 \
  PRODUCT_FUNCTIONAL_TEST_EMAIL=admin@test.com \
  PRODUCT_FUNCTIONAL_TEST_PASSWORD=password \
    go test ./product-management/backend/app/product-gateway/routes \
      -tags integration \
      -run TestDeployedProductGatewayFunctional \
      -v
  stop_port_forwards
fi

if [[ "${DEPLOY_ORDER_MANAGEMENT:-false}" == "true" ]]; then
  start_port_forward order-gateway 18085:8080
  wait_for_services order-gateway http://127.0.0.1:18085/health

  ORDER_GATEWAY_BASE_URL=http://127.0.0.1:18085 \
  ORDER_FUNCTIONAL_TEST_EMAIL=admin@test.com \
  ORDER_FUNCTIONAL_TEST_PASSWORD=password \
    go test ./order-management/backend/app/order-gateway/routes \
      -tags integration \
      -run TestDeployedOrderGatewayFunctional \
      -v
  stop_port_forwards
fi
