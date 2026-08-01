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
KEEP_FEATURE_IMAGES="${KEEP_FEATURE_IMAGES:-1}"
KEEP_MAIN_IMAGES="${KEEP_MAIN_IMAGES:-5}"
IMAGES="${IMAGES:-user-service user-gateway admin-web-app product-service product-gateway product-admin-web-app order-service order-gateway order-admin-web-app cart-gateway storefront-gateway storefront-web-app storybook liquibase}"

if [[ ! "$KEEP_FEATURE_IMAGES" =~ ^[0-9]+$ || ! "$KEEP_MAIN_IMAGES" =~ ^[0-9]+$ ]]; then
  echo "KEEP_FEATURE_IMAGES and KEEP_MAIN_IMAGES must be non-negative integers." >&2
  exit 1
fi

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

tag_epoch() {
  local tag="$1"
  local last_modified
  last_modified="$(curl --silent --head --fail \
    --user "admin:$NEXUS_ADMIN_PASSWORD" \
    -H 'Accept: application/vnd.docker.distribution.manifest.v2+json' \
    "http://$NEXUS_DOCKER_PUSH_REGISTRY/v2/$repository/manifests/$tag" \
    | awk -F': ' 'tolower($1) == "last-modified" {gsub("\r", "", $2); print $2}' \
    || true)"
  date --date="$last_modified" +'%s' 2>/dev/null || printf '0\n'
}

sort_tags_by_registry_date() {
  local tag
  local epoch
  for tag in "$@"; do
    epoch="$(tag_epoch "$tag")"
    printf '%s\t%s\n' "$epoch" "$tag"
  done | sort -rn | cut -f2-
}

is_running_image() {
  local repository="$1"
  local tag="$2"
  grep -Fxq "$repository:$tag" <<< "$RUNNING_IMAGES"
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

RUNNING_IMAGES="$(
  kubectl get pods --all-namespaces \
    -o jsonpath='{range .items[*].spec.initContainers[*]}{.image}{"\n"}{end}{range .items[*].spec.containers[*]}{.image}{"\n"}{end}' \
    2>/dev/null \
    | sed -E "s#^([^/]+/)?${IMAGE_NAMESPACE}/#${IMAGE_NAMESPACE}/#" \
    | sort -u \
    || true
)"

for image in $IMAGES; do
  repository="$IMAGE_NAMESPACE/$image"
  tags_json="$(curl --silent --fail --user "admin:$NEXUS_ADMIN_PASSWORD" "http://$NEXUS_DOCKER_PUSH_REGISTRY/v2/$repository/tags/list" || true)"
  if [[ -z "$tags_json" ]]; then
    echo "No tags found for $repository"
    continue
  fi

  mapfile -t tags < <(python3 -c 'import json, sys; print("\n".join(json.load(sys.stdin).get("tags") or []))' <<< "$tags_json")
  declare -a protected_tags=()
  declare -a deletion_candidates=()
  declare -a main_tags=()
  declare -A active_feature_tags=()

  for tag in "${tags[@]}"; do
    if [[ "$tag" == "latest" ]] || is_running_image "$repository" "$tag"; then
      protected_tags+=("$tag")
      continue
    fi

    if [[ "$tag" =~ ^feature-([a-z0-9]{1,10})-([0-9a-f]{40})$ ]]; then
      feature_namespace="${BASH_REMATCH[1]}"
      if is_active_feature "$feature_namespace"; then
        active_feature_tags["$feature_namespace"]+="$tag "
      else
        deletion_candidates+=("$tag")
      fi
      continue
    fi

    if [[ "$tag" =~ ^[0-9a-f]{40}$ ]]; then
      main_tags+=("$tag")
      continue
    fi

    echo "Keeping unrecognized image tag: $repository:$tag"
    protected_tags+=("$tag")
  done

  mapfile -t sorted_main_tags < <(sort_tags_by_registry_date "${main_tags[@]}")
  for index in "${!sorted_main_tags[@]}"; do
    tag="${sorted_main_tags[$index]}"
    if (( index < KEEP_MAIN_IMAGES )); then
      protected_tags+=("$tag")
    else
      deletion_candidates+=("$tag")
    fi
  done

  for feature_namespace in "${!active_feature_tags[@]}"; do
    read -r -a feature_tags <<< "${active_feature_tags[$feature_namespace]}"
    mapfile -t sorted_feature_tags < <(sort_tags_by_registry_date "${feature_tags[@]}")
    for index in "${!sorted_feature_tags[@]}"; do
      tag="${sorted_feature_tags[$index]}"
      if (( index < KEEP_FEATURE_IMAGES )); then
        protected_tags+=("$tag")
      else
        deletion_candidates+=("$tag")
      fi
    done
  done

  declare -A protected_digests=()
  for tag in "${protected_tags[@]}"; do
    digest="$(curl --silent --head --fail \
      --user "admin:$NEXUS_ADMIN_PASSWORD" \
      -H 'Accept: application/vnd.docker.distribution.manifest.v2+json' \
      "http://$NEXUS_DOCKER_PUSH_REGISTRY/v2/$repository/manifests/$tag" \
      | awk -F': ' 'tolower($1) == "docker-content-digest" {gsub("\r", "", $2); print $2}' \
      || true)"
    if [[ -n "$digest" ]]; then
      protected_digests["$digest"]=1
    fi
  done

  echo "$repository: keeping ${#protected_tags[@]} tag(s), pruning ${#deletion_candidates[@]} tag(s)"
  for tag in "${deletion_candidates[@]}"; do
    if is_running_image "$repository" "$tag"; then
      echo "Keeping running image: $repository:$tag"
      continue
    fi

    digest="$(curl --silent --head --fail \
      --user "admin:$NEXUS_ADMIN_PASSWORD" \
      -H 'Accept: application/vnd.docker.distribution.manifest.v2+json' \
      "http://$NEXUS_DOCKER_PUSH_REGISTRY/v2/$repository/manifests/$tag" \
      | awk -F': ' 'tolower($1) == "docker-content-digest" {gsub("\r", "", $2); print $2}' \
      || true)"

    if [[ -z "$digest" ]]; then
      echo "Could not resolve digest for $repository:$tag" >&2
      continue
    fi

    if [[ -n "${protected_digests[$digest]:-}" ]]; then
      echo "Keeping shared manifest: $repository:$tag"
      continue
    fi

    if [[ "$DRY_RUN" == "true" ]]; then
      echo "Would delete image outside retention: $repository:$tag"
      continue
    fi

    echo "Deleting image outside retention: $repository:$tag"
    curl --silent --fail --request DELETE \
      --user "admin:$NEXUS_ADMIN_PASSWORD" \
      "http://$NEXUS_DOCKER_PUSH_REGISTRY/v2/$repository/manifests/$digest" >/dev/null
  done

  unset active_feature_tags protected_digests
done
