variable "project_id" {
  description = "GCP project ID for adcplus CI integration tests."
  type        = string
}

variable "region" {
  description = "Default GCP region (used by provider metadata only)."
  type        = string
  default     = "us-central1"
}

variable "github_repository" {
  description = "GitHub repository allowed to impersonate the CI service account (org/repo)."
  type        = string
  default     = "apstndb/adcplus"
}

variable "workload_identity_pool_id" {
  description = "Existing Workload Identity Pool ID (shared apstndb GitHub pool)."
  type        = string
  default     = "github"
}

variable "workload_identity_provider_id" {
  description = "Existing OIDC provider ID within the pool."
  type        = string
  default     = "github"
}

variable "service_account_id" {
  description = "Short account ID for the adcplus CI service account."
  type        = string
  default     = "adcplus-ci"
}

variable "service_account_display_name" {
  description = "Display name for the adcplus CI service account."
  type        = string
  default     = "adcplus GitHub Actions CI"
}

variable "impersonation_target_service_account" {
  description = "Optional second service account email that adcplus-ci may impersonate for integration tests."
  type        = string
  default     = ""
}
