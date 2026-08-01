#!/usr/bin/env bash
set -euo pipefail

if ! command -v docker >/dev/null 2>&1; then
  echo "Docker is not available; skipping k3d registry mirror verification"
  exit 0
fi

K3D_CLUSTER_NAME="${K3D_CLUSTER:-eco-cluster}"
mapfile -t K3D_NODES < <(docker ps \
  --filter "name=k3d-${K3D_CLUSTER_NAME}-" \
  --format '{{.Names}}' \
  | grep -E '^k3d-.+-(server|agent)-[0-9]+$' || true)

if [[ "${#K3D_NODES[@]}" -eq 0 ]]; then
  echo "No k3d node containers found; skipping registry mirror verification"
  exit 0
fi

missing_mirrors=()
for node in "${K3D_NODES[@]}"; do
  if docker exec "$node" sh -c "test -f /etc/rancher/k3s/registries.yaml && grep -F '\"${NEXUS_DOCKER_PULL_REGISTRY}\":' /etc/rancher/k3s/registries.yaml >/dev/null && grep -F '\"http://${NEXUS_DOCKER_PULL_REGISTRY}\"' /etc/rancher/k3s/registries.yaml >/dev/null"; then
    echo "Registry mirror for $NEXUS_DOCKER_PULL_REGISTRY already configured on $node"
    continue
  fi

  missing_mirrors+=("$node")
done

if [[ "${#missing_mirrors[@]}" -gt 0 ]]; then
  printf 'Missing Nexus registry mirror %s on k3d node(s): %s\n' "$NEXUS_DOCKER_PULL_REGISTRY" "${missing_mirrors[*]}" >&2
  echo "Configure k3d registry mirrors outside the application deploy before running this workflow." >&2
  exit 1
fi

echo "All k3d registry mirrors are configured; no node restarts needed."
