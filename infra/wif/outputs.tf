output "project_id" {
  description = "GCP project ID."
  value       = var.project_id
}

output "project_number" {
  description = "GCP project number."
  value       = data.google_project.current.number
}

output "service_account_email" {
  description = "Service account email for GitHub Actions WIF auth."
  value       = google_service_account.adcplus_ci.email
}

output "workload_identity_provider" {
  description = "Full resource name for google-github-actions/auth workload_identity_provider input."
  value       = data.google_iam_workload_identity_pool_provider.github.name
}

output "github_actions_auth_snippet" {
  description = "Copy-paste values for .github/workflows/integration-test.yml."
  value = {
    workload_identity_provider = data.google_iam_workload_identity_pool_provider.github.name
    service_account            = google_service_account.adcplus_ci.email
  }
}

output "wif_principal_set" {
  description = "Principal set bound as workloadIdentityUser on the CI service account."
  value       = local.wif_member
}
