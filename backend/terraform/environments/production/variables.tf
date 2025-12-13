variable "project_id" {
  description = "GCP Project ID"
  type        = string
}

variable "region" {
  description = "GCP region for Cloud Run"
  type        = string
  default     = "asia-northeast1"
}

variable "service_name" {
  description = "Cloud Run service name"
  type        = string
  default     = "backend-api"
}

variable "image" {
  description = "Container image URL (e.g., asia-northeast1-docker.pkg.dev/PROJECT/REPO/IMAGE:TAG)"
  type        = string
}

variable "domain" {
  description = "Custom domain for the service"
  type        = string
}

variable "cloudflare_origin_ca_cert" {
  description = "Cloudflare Origin CA certificate (PEM format)"
  type        = string
  sensitive   = true
}

variable "cloudflare_origin_ca_key" {
  description = "Cloudflare Origin CA private key (PEM format)"
  type        = string
  sensitive   = true
}

variable "env_vars" {
  description = "Environment variables for Cloud Run"
  type        = map(string)
  default     = {}
}

variable "secrets" {
  description = "Secret environment variables (reference to Secret Manager)"
  type = map(object({
    secret_id = string
    version   = string
  }))
  default = {}
}

variable "min_instances" {
  description = "Minimum number of instances"
  type        = number
  default     = 0
}

variable "max_instances" {
  description = "Maximum number of instances"
  type        = number
  default     = 10
}

variable "cpu" {
  description = "CPU allocation (e.g., '1' or '2')"
  type        = string
  default     = "1"
}

variable "memory" {
  description = "Memory allocation (e.g., '512Mi' or '1Gi')"
  type        = string
  default     = "512Mi"
}

variable "concurrency" {
  description = "Maximum concurrent requests per instance"
  type        = number
  default     = 80
}
