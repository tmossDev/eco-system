#!/usr/bin/env bash
set -euo pipefail

DRY_RUN="${DRY_RUN:-true}"
KEEP_NAMESPACES="${KEEP_NAMESPACES:-eco-test eco-foundation eco-cicd default kube-system kube-public kube-node-lease}"
ACTIVE_FEATURE_NAMESPACES="${ACTIVE_FEATURE_NAMESPACES:-}"

if [[ -z "$ACTIVE_FEATURE_NAMESPACES" ]]; then
  ACTIVE_FEATURE_NAMESPACES="$(git ls-remote --heads origin 'refs/heads/feature/*' \
    | awk -F/ '{print $NF}' \
    | grep -E '^[a-z0-9]{1,10}$' \
    | sort -u \
    | tr '\n' ' ' || true)"
fi

is_listed() {
  local needle="$1"
  local item
  shift
  for item in "$@"; do
    if [[ "$item" == "$needle" ]]; then
      return 0
    fi
  done
  return 1
}

read -r -a keep_namespaces <<< "$KEEP_NAMESPACES"
read -r -a active_namespaces <<< "$ACTIVE_FEATURE_NAMESPACES"

mapfile -t namespaces < <(kubectl get namespaces -o jsonpath='{range .items[*]}{.metadata.name}{"\n"}{end}')

for namespace in "${namespaces[@]}"; do
  if is_listed "$namespace" "${keep_namespaces[@]}"; then
    echo "Keeping reserved namespace: $namespace"
    continue
  fi

  if [[ "$namespace" == kube-* ]]; then
    echo "Keeping Kubernetes namespace: $namespace"
    continue
  fi

  if [[ ! "$namespace" =~ ^[a-z0-9]{1,10}$ ]]; then
    echo "Skipping non-feature-shaped namespace: $namespace"
    continue
  fi

  deployment_label="$(kubectl get namespace "$namespace" -o jsonpath='{.metadata.labels.eco-system-deployment}' 2>/dev/null || true)"
  managed_by_label="$(kubectl get namespace "$namespace" -o jsonpath='{.metadata.labels.eco-system-managed-by}' 2>/dev/null || true)"
  if [[ "$deployment_label" != "feature" || "$managed_by_label" != "github-actions" ]]; then
    echo "Skipping unlabeled namespace: $namespace"
    continue
  fi

  if is_listed "$namespace" "${active_namespaces[@]}"; then
    echo "Keeping active feature namespace: $namespace"
    continue
  fi

  if [[ "$DRY_RUN" == "true" ]]; then
    echo "Would delete inactive feature namespace: $namespace"
  else
    echo "Deleting inactive feature namespace: $namespace"
    kubectl delete namespace "$namespace" --ignore-not-found --wait=false
  fi
done
