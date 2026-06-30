#!/usr/bin/env bash
# Import adcplus WIF resources created outside Terraform into local state.
#
# Usage: ./scripts/import-wif-state.sh [project_id]

set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
TF_DIR="${ROOT}/infra/wif"
PROJECT_ID="${1:-apstndb-sandbox}"
PROJECT_NUMBER="$(gcloud projects describe "$PROJECT_ID" --format='value(projectNumber)')"
SA="adcplus-ci@${PROJECT_ID}.iam.gserviceaccount.com"
WIF_MEMBER="principalSet://iam.googleapis.com/projects/${PROJECT_NUMBER}/locations/global/workloadIdentityPools/github/attribute.repository/apstndb/adcplus"

cd "${TF_DIR}"

terraform init -upgrade

terraform import -var="project_id=${PROJECT_ID}" google_service_account.adcplus_ci \
  "projects/${PROJECT_ID}/serviceAccounts/${SA}"

terraform import -var="project_id=${PROJECT_ID}" google_service_account_iam_member.adcplus_ci_wif_user \
  "projects/${PROJECT_ID}/serviceAccounts/${SA} roles/iam.workloadIdentityUser ${WIF_MEMBER}"

terraform import -var="project_id=${PROJECT_ID}" google_service_account_iam_member.adcplus_ci_token_creator_self \
  "projects/${PROJECT_ID}/serviceAccounts/${SA} roles/iam.serviceAccountTokenCreator serviceAccount:${SA}"

terraform import -var="project_id=${PROJECT_ID}" google_project_iam_member.adcplus_ci_service_account_user \
  "${PROJECT_ID} roles/iam.serviceAccountUser serviceAccount:${SA}"

echo "Import complete. Run: cd infra/wif && terraform plan -var=project_id=${PROJECT_ID}"
