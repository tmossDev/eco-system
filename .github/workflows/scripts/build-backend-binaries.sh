#!/usr/bin/env bash
set -euo pipefail

mkdir -p build/bin
built_bins=()

if [[ "${BUILD_USER_MANAGEMENT:-false}" == "true" ]]; then
  go test ./user-management/backend/app/user-gateway/routes -run 'TestGatewayFunctional'
  CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/bin/user-service ./user-management/backend/app/user-service/cmd
  CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/bin/user-gateway ./user-management/backend/app/user-gateway/cmd
  built_bins+=(build/bin/user-service build/bin/user-gateway)
fi

if [[ "${BUILD_PRODUCT_MANAGEMENT:-false}" == "true" ]]; then
  go test ./product-management/backend/app/product-gateway/routes -run 'TestProductGatewayFunctional'
  CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/bin/product-service ./product-management/backend/app/product-service/cmd
  CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/bin/product-gateway ./product-management/backend/app/product-gateway/cmd
  built_bins+=(build/bin/product-service build/bin/product-gateway)
fi

if [[ "${BUILD_ORDER_MANAGEMENT:-false}" == "true" ]]; then
  go test ./order-management/backend/app/order-gateway/routes
  go test ./order-management/backend/app/order-service/routes
  CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/bin/order-service ./order-management/backend/app/order-service/cmd
  CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/bin/order-gateway ./order-management/backend/app/order-gateway/cmd
  built_bins+=(build/bin/order-service build/bin/order-gateway)
fi

if [[ "${BUILD_ONLINE_STOREFRONT:-false}" == "true" ]]; then
  go test ./online-storefront/backend/app/cart-gateway/routes
  go test ./online-storefront/backend/app/storefront-gateway/routes
  CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/bin/cart-gateway ./online-storefront/backend/app/cart-gateway/cmd
  CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/bin/storefront-gateway ./online-storefront/backend/app/storefront-gateway/cmd
  built_bins+=(build/bin/cart-gateway build/bin/storefront-gateway)
fi

if [[ "${#built_bins[@]}" -eq 0 ]]; then
  echo "No backend binaries selected."
  exit 0
fi

file "${built_bins[@]}"
