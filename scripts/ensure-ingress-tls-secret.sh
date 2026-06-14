#!/usr/bin/env bash
set -euo pipefail

usage() {
  local status="${1:-1}"
  cat <<EOF
Usage: $0 --namespace NAME [--environment-file FILE] [--secret-name NAME]

Create or update the Kubernetes TLS secret used by local ingress.

Certificate source order:
  1. INGRESS_TLS_CERT_B64 and INGRESS_TLS_KEY_B64 environment variables.
  2. mkcert, when available on the runner.
  3. openssl self-signed certificate fallback.

Options:
  --namespace        namespace to install the TLS secret into
  --environment-file path under config/ (default: hp-prodesk-homelab.txt)
  --secret-name      Kubernetes TLS secret name (default: config value)
  -h, --help         show this help message
EOF
  exit "$status"
}

ENVIRONMENT_FILE="hp-prodesk-homelab.txt"
NAMESPACE=""
SECRET_NAME=""

while [[ $# -gt 0 ]]; do
  case "$1" in
    --namespace)
      NAMESPACE="$2"; shift 2
      ;;
    --environment-file)
      ENVIRONMENT_FILE="$2"; shift 2
      ;;
    --secret-name)
      SECRET_NAME="$2"; shift 2
      ;;
    -h|--help)
      usage 0
      ;;
    *)
      echo "Unknown option: $1" >&2
      usage
      ;;
  esac
done

if [[ -z "$NAMESPACE" ]]; then
  echo "--namespace is required" >&2
  usage
fi

CONFIG_FILE="config/$ENVIRONMENT_FILE"
if [[ ! -f "$CONFIG_FILE" ]]; then
  echo "Missing environment config: $CONFIG_FILE" >&2
  exit 1
fi

set -a
. "$CONFIG_FILE"
set +a

BASE_DOMAIN="${BASE_DOMAIN:-eco.test}"
SECRET_NAME="${SECRET_NAME:-${INGRESS_TLS_SECRET_NAME:-eco-local-ingress-tls}}"
CERT_DIR="${RUNNER_TEMP:-.local}/ingress-tls/$NAMESPACE"
CERT_FILE="$CERT_DIR/tls.crt"
KEY_FILE="$CERT_DIR/tls.key"

ensure_command() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "Required command not found: $1" >&2
    exit 1
  fi
}

ensure_command kubectl
mkdir -p "$CERT_DIR"

hosts=(
  "*.$BASE_DOMAIN"
  "$NAMESPACE.admin-web-app.$BASE_DOMAIN"
  "$NAMESPACE.user-gateway.$BASE_DOMAIN"
  "$NAMESPACE.user-service.$BASE_DOMAIN"
  "$NAMESPACE.product-admin-web-app.$BASE_DOMAIN"
  "$NAMESPACE.product-gateway.$BASE_DOMAIN"
  "$NAMESPACE.product-service.$BASE_DOMAIN"
  "$NAMESPACE.order-admin-web-app.$BASE_DOMAIN"
  "$NAMESPACE.order-gateway.$BASE_DOMAIN"
  "$NAMESPACE.order-service.$BASE_DOMAIN"
  "$NAMESPACE.storefront-web-app.$BASE_DOMAIN"
  "$NAMESPACE.storybook.$BASE_DOMAIN"
)

if [[ -n "${INGRESS_TLS_CERT_B64:-}" && -n "${INGRESS_TLS_KEY_B64:-}" ]]; then
  echo "Creating ingress TLS secret $SECRET_NAME in $NAMESPACE from configured certificate material"
  printf '%s' "$INGRESS_TLS_CERT_B64" | base64 --decode > "$CERT_FILE"
  printf '%s' "$INGRESS_TLS_KEY_B64" | base64 --decode > "$KEY_FILE"
elif command -v mkcert >/dev/null 2>&1; then
  echo "Creating ingress TLS secret $SECRET_NAME in $NAMESPACE with mkcert"
  mkcert -cert-file "$CERT_FILE" -key-file "$KEY_FILE" "${hosts[@]}"
else
  ensure_command openssl
  echo "Creating ingress TLS secret $SECRET_NAME in $NAMESPACE with openssl self-signed certificate"
  san=""
  for host in "${hosts[@]}"; do
    if [[ -n "$san" ]]; then
      san+=","
    fi
    san+="DNS:$host"
  done

  openssl req \
    -x509 \
    -newkey rsa:2048 \
    -sha256 \
    -days "${INGRESS_TLS_DAYS:-365}" \
    -nodes \
    -keyout "$KEY_FILE" \
    -out "$CERT_FILE" \
    -subj "/CN=$NAMESPACE.$BASE_DOMAIN" \
    -addext "subjectAltName=$san"
fi

kubectl create namespace "$NAMESPACE" --dry-run=client -o yaml | kubectl apply -f -
kubectl create secret tls "$SECRET_NAME" \
  --namespace "$NAMESPACE" \
  --cert "$CERT_FILE" \
  --key "$KEY_FILE" \
  --dry-run=client \
  -o yaml | kubectl apply -f -

echo "Installed ingress TLS secret $SECRET_NAME in namespace $NAMESPACE"
