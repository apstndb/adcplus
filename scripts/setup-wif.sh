#!/usr/bin/env bash
# Provision adcplus GitHub Actions Workload Identity Federation bindings.
#
# Prerequisites:
#   - gcloud CLI authenticated with IAM admin on the target project
#   - terraform >= 1.5
#
# Usage:
#   ./scripts/setup-wif.sh [project_id]
#
# Default project_id: value from `gcloud config get-value project`.

set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
TF_DIR="${ROOT}/infra/wif"
PROJECT_ID="${1:-$(gcloud config get-value project 2>/dev/null)}"

if [[ -z "${PROJECT_ID}" || "${PROJECT_ID}" == "(unset)" ]]; then
  echo "error: set a GCP project via argument or gcloud config" >&2
  exit 1
fi

echo "Using project: ${PROJECT_ID}"
echo "Terraform directory: ${TF_DIR}"

cd "${TF_DIR}"

if [[ ! -f terraform.tfvars ]]; then
  cp terraform.tfvars.example terraform.tfvars
  sed -i.bak "s/apstndb-sandbox/${PROJECT_ID}/" terraform.tfvars
  rm -f terraform.tfvars.bak
  echo "Created terraform.tfvars from example (project_id=${PROJECT_ID})"
fi

terraform init -upgrade
terraform plan -var="project_id=${PROJECT_ID}"
terraform apply -var="project_id=${PROJECT_ID}" -auto-approve

echo
echo "=== GitHub Actions values ==="
terraform output -json github_actions_auth_snippet
