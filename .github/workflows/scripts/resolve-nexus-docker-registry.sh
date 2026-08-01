#!/usr/bin/env bash
set -euo pipefail

mode="${1:-pull}"

docker_port="$(kubectl get svc nexus -n "$CI_CD_NAMESPACE" -o jsonpath='{.spec.ports[?(@.name=="docker")].port}')"
cluster_ip="$(kubectl get svc nexus -n "$CI_CD_NAMESPACE" -o jsonpath='{.spec.clusterIP}')"
pull_registry="${NEXUS_DOCKER_PULL_REGISTRY:-$cluster_ip:$docker_port}"

echo "NEXUS_DOCKER_PULL_REGISTRY=$pull_registry" >> "$GITHUB_ENV"

if [[ "$mode" == "push" ]]; then
  push_port="${NEXUS_DOCKER_PUSH_PORT:-${NEXUS_DOCKER_NODE_PORT:-30500}}"
  push_registry="${NEXUS_DOCKER_PUSH_REGISTRY:-127.0.0.1:$push_port}"
  echo "NEXUS_DOCKER_PUSH_REGISTRY=$push_registry" >> "$GITHUB_ENV"
  echo "NEXUS_DOCKER_REGISTRY=$push_registry" >> "$GITHUB_ENV"
  echo "IMAGE_PREFIX=$push_registry/$IMAGE_NAMESPACE" >> "$GITHUB_ENV"
  echo "Docker push registry: $push_registry"
else
  echo "NEXUS_DOCKER_REGISTRY=$pull_registry" >> "$GITHUB_ENV"
  echo "IMAGE_PREFIX=$pull_registry/$IMAGE_NAMESPACE" >> "$GITHUB_ENV"
fi

echo "Kubernetes pull registry: $pull_registry"
