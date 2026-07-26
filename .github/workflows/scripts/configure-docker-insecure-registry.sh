#!/usr/bin/env bash
set -euo pipefail

if [[ "$NEXUS_DOCKER_REGISTRY" =~ ^(127\.|localhost:|0\.0\.0\.0:) ]]; then
  echo "Using local Nexus registry forward $NEXUS_DOCKER_REGISTRY; Docker daemon insecure-registry configuration is not required."
  exit 0
fi

if docker info 2>/dev/null | grep -F "$NEXUS_DOCKER_REGISTRY" >/dev/null; then
  echo "Docker already trusts insecure registry $NEXUS_DOCKER_REGISTRY"
  exit 0
fi

if ! sudo -n true 2>/dev/null; then
  echo "Docker must be configured to trust HTTP registry $NEXUS_DOCKER_REGISTRY before images can be pushed." >&2
  echo "Either preconfigure /etc/docker/daemon.json on the self-hosted runner or allow passwordless sudo for the runner user." >&2
  exit 1
fi

sudo mkdir -p /etc/docker
if [[ -f /etc/docker/daemon.json ]]; then
  sudo cp /etc/docker/daemon.json /tmp/docker-daemon.json
else
  printf '{}\n' > /tmp/docker-daemon.json
fi

python3 - "$NEXUS_DOCKER_REGISTRY" /tmp/docker-daemon.json <<'PY'
import json
import sys

registry = sys.argv[1]
path = sys.argv[2]

with open(path, "r", encoding="utf-8") as fh:
    config = json.load(fh)

registries = config.setdefault("insecure-registries", [])
if registry not in registries:
    registries.append(registry)

with open(path, "w", encoding="utf-8") as fh:
    json.dump(config, fh, indent=2)
    fh.write("\n")
PY

sudo cp /tmp/docker-daemon.json /etc/docker/daemon.json
sudo systemctl restart docker || sudo service docker restart
docker info
