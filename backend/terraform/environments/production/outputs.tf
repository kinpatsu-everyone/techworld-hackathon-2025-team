output "cloud_run_url" {
  description = "Cloud Run service URL"
  value       = google_cloud_run_v2_service.main.uri
}

output "load_balancer_ip" {
  description = "Load Balancer IP address (set this in Cloudflare DNS)"
  value       = google_compute_global_address.main.address
}

output "service_name" {
  description = "Cloud Run service name"
  value       = google_cloud_run_v2_service.main.name
}

output "region" {
  description = "Deployment region"
  value       = var.region
}
