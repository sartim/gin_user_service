#!/usr/bin/env bash
set -euo pipefail

readonly SERVICE_NAME="gin-user-service"
readonly RELEASE_DIR="/home/ubuntu/user-service"
readonly ENV_FILE="/etc/gin-user-service.env"
readonly UNIT_FILE="/etc/systemd/system/${SERVICE_NAME}.service"

before_install() {
    systemctl stop "${SERVICE_NAME}" 2>/dev/null || true
    rm -rf -- "${RELEASE_DIR}"
}

after_install() {
    test -f "${ENV_FILE}" || {
        echo "Required environment file ${ENV_FILE} is missing" >&2
        return 1
    }
    chmod 0755 "${RELEASE_DIR}/user-service"
    install -m 0644 "${RELEASE_DIR}/aws-codedeploy-helpers/gin-user-service.service" "${UNIT_FILE}"
    systemctl daemon-reload
}

application_start() {
    systemctl enable "${SERVICE_NAME}"
    systemctl restart "${SERVICE_NAME}"
}

application_stop() {
    systemctl stop "${SERVICE_NAME}" 2>/dev/null || true
}

validate_service() {
    curl --fail --silent --show-error --retry 10 --retry-delay 2 \
        "http://127.0.0.1:8000/health/ready" >/dev/null
}
