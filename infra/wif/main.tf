data "google_project" "current" {
  project_id = var.project_id
}

data "google_iam_workload_identity_pool" "github" {
  workload_identity_pool_id = var.workload_identity_pool_id
}

data "google_iam_workload_identity_pool_provider" "github" {
  workload_identity_pool_id          = var.workload_identity_pool_id
  workload_identity_pool_provider_id = var.workload_identity_provider_id
}

resource "google_service_account" "adcplus_ci" {
  account_id   = var.service_account_id
  display_name = var.service_account_display_name
  description  = "GitHub Actions WIF identity for apstndb/adcplus integration tests."
}

locals {
  wif_pool_resource = data.google_iam_workload_identity_pool.github.name
  wif_member        = "principalSet://iam.googleapis.com/${local.wif_pool_resource}/attribute.repository/${var.github_repository}"
}

# Allow GitHub Actions (via WIF) to impersonate adcplus-ci.
resource "google_service_account_iam_member" "adcplus_ci_wif_user" {
  service_account_id = google_service_account.adcplus_ci.name
  role               = "roles/iam.workloadIdentityUser"
  member             = local.wif_member
}

# Allow adcplus-ci to mint tokens for itself (Credentials API / impersonation self-tests).
resource "google_service_account_iam_member" "adcplus_ci_token_creator_self" {
  service_account_id = google_service_account.adcplus_ci.name
  role               = "roles/iam.serviceAccountTokenCreator"
  member             = "serviceAccount:${google_service_account.adcplus_ci.email}"
}

# Optional: allow adcplus-ci to impersonate another SA used as an impersonation target in tests.
resource "google_service_account_iam_member" "adcplus_ci_token_creator_target" {
  count = var.impersonation_target_service_account != "" ? 1 : 0

  service_account_id = "projects/${var.project_id}/serviceAccounts/${var.impersonation_target_service_account}"
  role               = "roles/iam.serviceAccountTokenCreator"
  member             = "serviceAccount:${google_service_account.adcplus_ci.email}"
}

# Minimal project role so integration tests can call IAM Credentials API as the SA.
resource "google_project_iam_member" "adcplus_ci_service_account_user" {
  project = var.project_id
  role    = "roles/iam.serviceAccountUser"
  member  = "serviceAccount:${google_service_account.adcplus_ci.email}"
}
