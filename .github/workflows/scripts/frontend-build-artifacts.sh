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
  ds-tokens
  ds-button
  ds-icon
  ds-navigation-bar
  ds-login-form
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
PNPM_CMD="${PNPM_CMD:-corepack pnpm}"

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
  [ds-tokens]="shared-components/frontend/package/design-system/ds/tokens"
  [ds-button]="shared-components/frontend/package/design-system/ds/button"
  [ds-icon]="shared-components/frontend/package/design-system/ds/icon"
  [ds-navigation-bar]="shared-components/frontend/package/design-system/ds/navigation-bar"
  [ds-login-form]="shared-components/frontend/package/design-system/ds/login-form"
)

declare -A BUILD_COMMANDS=(
  [admin-web-app]="$PNPM_CMD --filter admin-web-app build && $PNPM_CMD --filter admin-web-app build-storybook"
  [product-admin-web-app]="$PNPM_CMD --filter product-admin-web-app build && $PNPM_CMD --filter product-admin-web-app build-storybook"
  [order-admin-web-app]="$PNPM_CMD --filter order-admin-web-app build && $PNPM_CMD --filter order-admin-web-app build-storybook"
  [storefront-web-app]="$PNPM_CMD --filter storefront-web-app build"
  [storybook]="$PNPM_CMD --filter storybook build-storybook"
  [ds-tokens]="$PNPM_CMD --filter @ds/tokens build"
  [ds-button]="$PNPM_CMD --filter design-system build:button"
  [ds-icon]="$PNPM_CMD --filter design-system build:icon"
  [ds-navigation-bar]="$PNPM_CMD --filter design-system build:navigation-bar"
  [ds-login-form]="$PNPM_CMD --filter design-system build:login-form"
)

declare -A OUTPUTS=(
  [admin-web-app]="user-management/frontend/app/admin-web-app/dist user-management/frontend/app/admin-web-app/storybook-static"
  [product-admin-web-app]="product-management/frontend/app/product-admin-web-app/dist product-management/frontend/app/product-admin-web-app/storybook-static"
  [order-admin-web-app]="order-management/frontend/app/order-admin-web-app/dist order-management/frontend/app/order-admin-web-app/storybook-static"
  [storefront-web-app]="online-storefront/frontend/app/storefront-web-app/dist"
  [storybook]="shared-components/frontend/app/storybook/storybook-static"
  [ds-tokens]="shared-components/frontend/package/design-system/ds/tokens/package.json shared-components/frontend/package/design-system/ds/tokens/style-dictionary.config.json shared-components/frontend/package/design-system/ds/tokens/tokens shared-components/frontend/package/design-system/ds/tokens/src/scss"
  [ds-button]="shared-components/frontend/package/design-system/dist/button"
  [ds-icon]="shared-components/frontend/package/design-system/dist/icon"
  [ds-navigation-bar]="shared-components/frontend/package/design-system/dist/navigation-bar"
  [ds-login-form]="shared-components/frontend/package/design-system/dist/login-form"
)

declare -A CLEAN_OUTPUTS=(
  [ds-tokens]="shared-components/frontend/package/design-system/ds/tokens/src/scss/_variables.scss shared-components/frontend/package/design-system/ds/tokens/src/scss/_map.scss shared-components/frontend/package/design-system/ds/tokens/src/scss/_css-variables.scss"
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

  case "$artifact" in
    admin-web-app|product-admin-web-app|order-admin-web-app)
      paths+=(
        shared-components/frontend/package/admin-features
        shared-components/frontend/package/auth-features
      )
      ;;
  esac

  case "$artifact" in
    ds-button|ds-navigation-bar|ds-login-form)
      paths+=(shared-components/frontend/package/design-system/ds/tokens)
      ;;
  esac

  git ls-files -z -- "${paths[@]}" \
    | while IFS= read -r -d '' file; do
        [[ -f "$file" ]] && printf '%s\0' "$file"
      done \
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
  for path in ${CLEAN_OUTPUTS[$artifact]:-${OUTPUTS[$artifact]}}; do
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
  if [[ "${FORCE_FRONTEND_ARTIFACT_BUILD:-false}" == "true" ]]; then
    build_and_publish "$artifact" "$hash"
  elif ! restore_artifact "$artifact" "$hash"; then
    build_and_publish "$artifact" "$hash"
  fi
done
