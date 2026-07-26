#!/usr/bin/env bash
set -euo pipefail

mkdir -p "$HOME/.kube"
if [[ -f "$HOME/.kube/config" ]] && kubectl config current-context >/dev/null 2>&1; then
  echo "Using existing runner kubeconfig at $HOME/.kube/config"
elif [[ -f "$HOME/.config/k3d/kubeconfig-${K3D_CLUSTER:-eco-cluster}.yaml" ]]; then
  cp "$HOME/.config/k3d/kubeconfig-${K3D_CLUSTER:-eco-cluster}.yaml" "$HOME/.kube/config"
elif command -v k3d >/dev/null 2>&1 && k3d cluster list "${K3D_CLUSTER:-eco-cluster}" >/dev/null 2>&1; then
  k3d kubeconfig get "${K3D_CLUSTER:-eco-cluster}" > "$HOME/.kube/config"
elif [[ -f /etc/rancher/k3s/k3s.yaml ]]; then
  sudo cat /etc/rancher/k3s/k3s.yaml > "$HOME/.kube/config"
elif [[ -n "${KUBE_CONFIG_B64:-}" ]]; then
  echo "$KUBE_CONFIG_B64" | base64 --decode > "$HOME/.kube/config"
else
  echo "No usable kubeconfig found. Add one at $HOME/.kube/config on the runner, install k3s config at /etc/rancher/k3s/k3s.yaml, or set KUBE_CONFIG_B64." >&2
  exit 1
fi
chmod 600 "$HOME/.kube/config"
