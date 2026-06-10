#!/usr/bin/env bash
set -euo pipefail

usage() {
  cat >&2 <<'USAGE'
Usage: frontend-build-artifacts.sh <artifact> [<artifact>...]

Artifacts:
  admin-web-app
  product-admin-web-app
  order-admin-web-app
  storefront-web-app
  storybook
USAGE
}

if [[ "$#" -eq 0 ]]; then
  usage
  exit 1
fi

NEXUS_ARTIFACT_REPOSITORY="${NEXUS_ARTIFACT_REPOSITORY:-eco-node-builds}"
NEXUS_ARTIFACT_URL="${NEXUS_ARTIFACT_URL:-http://127.0.0.1:18081/repository/$NEXUS_ARTIFACT_REPOSITORY}"
NEXUS_USERNAME="${NEXUS_USERNAME:-admin}"
NEXUS_PASSWORD="${NEXUS_PASSWORD:-${NEXUS_ADMIN_PASSWORD:-}}"
CI_CD_NAMESPACE="${CI_CD_NAMESPACE:-eco-cicd}"

if [[ -z "$NEXUS_PASSWORD" ]]; then
  echo "Missing NEXUS_PASSWORD or NEXUS_ADMIN_PASSWORD for Nexus artifact access." >&2
  exit 1
fi

declare -A PACKAGE_DIRS=(
  [admin-web-app]="user-management/frontend/app/admin-web-app"
  [product-admin-web-app]="product-management/frontend/app/product-admin-web-app"
  [order-admin-web-app]="order-management/frontend/app/order-admin-web-app"
  [storefront-web-app]="online-storefront/frontend/app/storefront-web-app"
  [storybook]="shared-components/frontend/app/storybook"
)

declare -A BUILD_COMMANDS=(
  [admin-web-app]="pnpm --filter admin-web-app build && pnpm --filter admin-web-app build-storybook"
  [product-admin-web-app]="pnpm --filter product-admin-web-app build && pnpm --filter product-admin-web-app build-storybook"
  [order-admin-web-app]="pnpm --filter order-admin-web-app build && pnpm --filter order-admin-web-app build-storybook"
  [storefront-web-app]="pnpm --filter storefront-web-app build"
  [storybook]="pnpm --filter storybook build-storybook"
)

declare -A OUTPUTS=(
  [admin-web-app]="user-management/frontend/app/admin-web-app/dist user-management/frontend/app/admin-web-app/storybook-static"
  [product-admin-web-app]="product-management/frontend/app/product-admin-web-app/dist product-management/frontend/app/product-admin-web-app/storybook-static"
  [order-admin-web-app]="order-management/frontend/app/order-admin-web-app/dist order-management/frontend/app/order-admin-web-app/storybook-static"
  [storefront-web-app]="online-storefront/frontend/app/storefront-web-app/dist"
  [storybook]="shared-components/frontend/app/storybook/storybook-static"
)

start_nexus_tunnel() {
  if [[ "$NEXUS_ARTIFACT_URL" != http://127.0.0.1:* && "$NEXUS_ARTIFACT_URL" != http://localhost:* ]]; then
    return
  fi

  kubectl rollout status deployment/nexus -n "$CI_CD_NAMESPACE" --timeout=600s
  kubectl port-forward -n "$CI_CD_NAMESPACE" deployment/nexus 18081:8081 &
  NEXUS_PORT_FORWARD_PID="$!"
  trap 'kill "$NEXUS_PORT_FORWARD_PID" 2>/dev/null || true' EXIT

  for attempt in {1..90}; do
    if curl --silent --fail http://127.0.0.1:18081/service/rest/v1/status >/dev/null; then
      return
    fi
    sleep 2
  done

  echo "Timed out waiting for Nexus artifact API." >&2
  exit 1
}

ensure_raw_repository() {
  local repository_api="http://127.0.0.1:18081/service/rest/v1/repositories/raw/hosted/$NEXUS_ARTIFACT_REPOSITORY"
  local repositories_api="http://127.0.0.1:18081/service/rest/v1/repositories/raw/hosted"

  if [[ "$NEXUS_ARTIFACT_URL" != http://127.0.0.1:* && "$NEXUS_ARTIFACT_URL" != http://localhost:* ]]; then
    return
  fi

  if curl --silent --fail --user "$NEXUS_USERNAME:$NEXUS_PASSWORD" "$repository_api" >/dev/null; then
    return
  fi

  curl --silent --fail --user "$NEXUS_USERNAME:$NEXUS_PASSWORD" \
    --request POST \
    --header "Content-Type: application/json" \
    --data "{
      \"name\": \"$NEXUS_ARTIFACT_REPOSITORY\",
      \"online\": true,
      \"storage\": {
        \"blobStoreName\": \"default\",
        \"strictContentTypeValidation\": false,
        \"writePolicy\": \"ALLOW\"
      },
      \"cleanup\": {
        \"policyNames\": []
      },
      \"component\": {
        \"proprietaryComponents\": false
      },
      \"raw\": {
        \"contentDisposition\": \"ATTACHMENT\"
      }
    }" \
    "$repositories_api"
}

artifact_hash() {
  local artifact="$1"
  local package_dir="${PACKAGE_DIRS[$artifact]}"
  local paths=(
    package.json
    pnpm-lock.yaml
    pnpm-workspace.yaml
    "$package_dir"
  )

  case "$artifact" in
    admin-web-app|product-admin-web-app|order-admin-web-app|storybook)
      paths+=(shared-components/frontend/package/design-system)
      ;;
  esac

  git ls-files -z -- "${paths[@]}" \
    | sort -z \
    | xargs -0 sha256sum \
    | sha256sum \
    | awk '{print $1}'
}

artifact_url() {
  local artifact="$1"
  local hash="$2"
  printf '%s/%s/%s.tgz' "$NEXUS_ARTIFACT_URL" "$artifact" "$hash"
}

clean_outputs() {
  local artifact="$1"
  local path
  for path in ${OUTPUTS[$artifact]}; do
    rm -rf "$path"
  done
}

restore_artifact() {
  local artifact="$1"
  local hash="$2"
  local url
  url="$(artifact_url "$artifact" "$hash")"

  mkdir -p .frontend-artifacts
  if curl --silent --fail --user "$NEXUS_USERNAME:$NEXUS_PASSWORD" \
    --output ".frontend-artifacts/$artifact.tgz" \
    "$url"; then
    echo "Restored $artifact build artifact from Nexus ($hash)."
    clean_outputs "$artifact"
    tar -xzf ".frontend-artifacts/$artifact.tgz"
    return 0
  fi

  echo "No reusable $artifact artifact found in Nexus ($hash)."
  return 1
}

publish_artifact() {
  local artifact="$1"
  local hash="$2"
  local archive=".frontend-artifacts/$artifact.tgz"
  local url
  url="$(artifact_url "$artifact" "$hash")"

  # Nexus raw repositories accept PUT to create intermediate paths.
  curl --silent --fail --user "$NEXUS_USERNAME:$NEXUS_PASSWORD" \
    --upload-file "$archive" \
    "$url"
  echo "Published $artifact build artifact to Nexus ($hash)."
}

build_and_publish() {
  local artifact="$1"
  local hash="$2"

  echo "Building $artifact because no Nexus artifact matched its input hash."
  clean_outputs "$artifact"
  eval "${BUILD_COMMANDS[$artifact]}"
  tar -czf ".frontend-artifacts/$artifact.tgz" ${OUTPUTS[$artifact]}
  publish_artifact "$artifact" "$hash"
}

start_nexus_tunnel
ensure_raw_repository

for artifact in "$@"; do
  if [[ -z "${PACKAGE_DIRS[$artifact]:-}" ]]; then
    echo "Unknown frontend artifact: $artifact" >&2
    usage
    exit 1
  fi

  hash="$(artifact_hash "$artifact")"
  if ! restore_artifact "$artifact" "$hash"; then
    build_and_publish "$artifact" "$hash"
  fi
done
