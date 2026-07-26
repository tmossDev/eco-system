#!/usr/bin/env bash
set -euo pipefail

require_var() {
  local name="$1"
  if [[ -z "${!name:-}" ]]; then
    echo "Missing required environment variable: $name" >&2
    exit 1
  fi
}

require_var NEXUS_DOCKER_PUSH_REGISTRY
require_var CI_CD_NAMESPACE
require_var NEXUS_ADMIN_PASSWORD

DRY_RUN="${DRY_RUN:-true}"
IMAGE_NAMESPACE="${IMAGE_NAMESPACE:-eco-system}"
ACTIVE_FEATURE_NAMESPACES="${ACTIVE_FEATURE_NAMESPACES:-}"
IMAGES="${IMAGES:-user-service user-gateway admin-web-app product-service product-gateway product-admin-web-app order-service order-gateway order-admin-web-app cart-gateway storefront-gateway storefront-web-app storybook liquibase}"

if [[ -z "$ACTIVE_FEATURE_NAMESPACES" ]]; then
  ACTIVE_FEATURE_NAMESPACES="$(git ls-remote --heads origin 'refs/heads/feature/*' \
    | awk -F/ '{print $NF}' \
    | grep -E '^[a-z0-9]{1,10}$' \
    | sort -u \
    | tr '\n' ' ' || true)"
fi

is_active_feature() {
  local namespace="$1"
  local active
  for active in $ACTIVE_FEATURE_NAMESPACES; do
    if [[ "$active" == "$namespace" ]]; then
      return 0
    fi
  done
  return 1
}

for attempt in {1..60}; do
  endpoints="$(kubectl get endpoints nexus -n "$CI_CD_NAMESPACE" -o jsonpath='{.subsets[*].addresses[*].ip}' 2>/dev/null || true)"
  if [[ -n "$endpoints" ]]; then
    break
  fi
  if [[ "$attempt" == "60" ]]; then
    echo "Timed out waiting for Nexus service endpoints in namespace $CI_CD_NAMESPACE." >&2
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

for image in $IMAGES; do
  repository="$IMAGE_NAMESPACE/$image"
  tags_json="$(curl --silent --fail --user "admin:$NEXUS_ADMIN_PASSWORD" "http://$NEXUS_DOCKER_PUSH_REGISTRY/v2/$repository/tags/list" || true)"
  if [[ -z "$tags_json" ]]; then
    echo "No tags found for $repository"
    continue
  fi

  mapfile -t tags < <(python3 -c 'import json, sys; print("\n".join(json.load(sys.stdin).get("tags") or []))' <<< "$tags_json")
  for tag in "${tags[@]}"; do
    if [[ ! "$tag" =~ ^feature-([a-z0-9]{1,10})-[0-9a-f]{40}$ ]]; then
      continue
    fi

    feature_namespace="${BASH_REMATCH[1]}"
    if is_active_feature "$feature_namespace"; then
      echo "Keeping active feature image: $repository:$tag"
      continue
    fi

    if [[ "$DRY_RUN" == "true" ]]; then
      echo "Would delete inactive feature image: $repository:$tag"
      continue
    fi

    digest="$(curl --silent --head --fail \
      --user "admin:$NEXUS_ADMIN_PASSWORD" \
      -H 'Accept: application/vnd.docker.distribution.manifest.v2+json' \
      "http://$NEXUS_DOCKER_PUSH_REGISTRY/v2/$repository/manifests/$tag" \
      | awk -F': ' 'tolower($1) == "docker-content-digest" {gsub("\r", "", $2); print $2}')"

    if [[ -z "$digest" ]]; then
      echo "Could not resolve digest for $repository:$tag" >&2
      continue
    fi

    echo "Deleting inactive feature image: $repository:$tag"
    curl --silent --fail --request DELETE \
      --user "admin:$NEXUS_ADMIN_PASSWORD" \
      "http://$NEXUS_DOCKER_PUSH_REGISTRY/v2/$repository/manifests/$digest" >/dev/null
  done
done
