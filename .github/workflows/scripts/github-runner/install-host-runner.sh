#!/usr/bin/env bash
set -euo pipefail

RUNNER_VERSION="${RUNNER_VERSION:-2.334.0}"
RUNNER_SHA256="${RUNNER_SHA256:-048024cd2c848eb6f14d5646d56c13a4def2ae7ee3ad12122bee960c56f3d271}"
RUNNER_REPO_URL="${RUNNER_REPO_URL:-https://github.com/tmossDev/eco-system}"
RUNNER_DIR="${RUNNER_DIR:-$HOME/actions-runner/eco-system}"
RUNNER_NAME="${RUNNER_NAME:-$(hostname)-eco-system}"
RUNNER_LABELS="${RUNNER_LABELS:-self-hosted,linux}"
RUNNER_WORK_DIR="${RUNNER_WORK_DIR:-_work}"

if [[ -z "${RUNNER_TOKEN:-}" ]]; then
  echo "RUNNER_TOKEN is required. Generate a fresh repo runner token in GitHub and export it before running this script." >&2
  exit 1
fi

mkdir -p "$RUNNER_DIR"
cd "$RUNNER_DIR"

archive="actions-runner-linux-x64-${RUNNER_VERSION}.tar.gz"
url="https://github.com/actions/runner/releases/download/v${RUNNER_VERSION}/${archive}"

if [[ ! -f "$archive" ]]; then
  curl -fsSL "$url" -o "$archive"
fi

echo "${RUNNER_SHA256}  ${archive}" | shasum -a 256 -c

if [[ ! -x ./config.sh ]]; then
  tar xzf "$archive"
fi

if [[ "${INSTALL_RUNNER_DEPENDENCIES:-false}" == "true" ]]; then
  sudo ./bin/installdependencies.sh
fi

if [[ ! -f .runner ]]; then
  ./config.sh \
    --unattended \
    --url "$RUNNER_REPO_URL" \
    --token "$RUNNER_TOKEN" \
    --name "$RUNNER_NAME" \
    --labels "$RUNNER_LABELS" \
    --work "$RUNNER_WORK_DIR" \
    --replace
fi

if [[ "${INSTALL_SERVICE:-true}" == "true" ]]; then
  sudo ./svc.sh install
  sudo ./svc.sh start
  sudo ./svc.sh status
else
  ./run.sh
fi
