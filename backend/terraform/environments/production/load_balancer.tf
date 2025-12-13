# =============================================================================
# Serverless NEG for Cloud Run
# =============================================================================

resource "google_compute_region_network_endpoint_group" "serverless_neg" {
  name                  = "${var.service_name}-neg"
  network_endpoint_type = "SERVERLESS"
  region                = var.region

  cloud_run {
    service = google_cloud_run_v2_service.main.name
  }
}

# =============================================================================
# SSL Certificate (Cloudflare Origin CA)
# =============================================================================

resource "google_compute_ssl_certificate" "cloudflare_origin" {
  name        = "${var.service_name}-cloudflare-origin-cert"
  private_key = var.cloudflare_origin_ca_key
  certificate = var.cloudflare_origin_ca_cert
}

# =============================================================================
# Backend Service
# =============================================================================

resource "google_compute_backend_service" "main" {
  name                  = "${var.service_name}-backend"
  protocol              = "HTTPS"
  port_name             = "http"
  timeout_sec           = 30
  load_balancing_scheme = "EXTERNAL_MANAGED"

  backend {
    group = google_compute_region_network_endpoint_group.serverless_neg.id
  }

  log_config {
    enable      = true
    sample_rate = 1.0
  }
}

# =============================================================================
# URL Map
# =============================================================================

resource "google_compute_url_map" "main" {
  name            = "${var.service_name}-url-map"
  default_service = google_compute_backend_service.main.id
}

# =============================================================================
# HTTPS Proxy
# =============================================================================

resource "google_compute_target_https_proxy" "main" {
  name             = "${var.service_name}-https-proxy"
  url_map          = google_compute_url_map.main.id
  ssl_certificates = [google_compute_ssl_certificate.cloudflare_origin.id]
}

# =============================================================================
# Global IP Address
# =============================================================================

resource "google_compute_global_address" "main" {
  name = "${var.service_name}-ip"
}

# =============================================================================
# Global Forwarding Rule (HTTPS)
# =============================================================================

resource "google_compute_global_forwarding_rule" "https" {
  name                  = "${var.service_name}-https-rule"
  ip_protocol           = "TCP"
  load_balancing_scheme = "EXTERNAL_MANAGED"
  port_range            = "443"
  target                = google_compute_target_https_proxy.main.id
  ip_address            = google_compute_global_address.main.id
}

# =============================================================================
# HTTP to HTTPS Redirect
# =============================================================================

resource "google_compute_url_map" "http_redirect" {
  name = "${var.service_name}-http-redirect"

  default_url_redirect {
    https_redirect         = true
    redirect_response_code = "MOVED_PERMANENTLY_DEFAULT"
    strip_query            = false
  }
}

resource "google_compute_target_http_proxy" "http_redirect" {
  name    = "${var.service_name}-http-proxy"
  url_map = google_compute_url_map.http_redirect.id
}

resource "google_compute_global_forwarding_rule" "http_redirect" {
  name                  = "${var.service_name}-http-rule"
  ip_protocol           = "TCP"
  load_balancing_scheme = "EXTERNAL_MANAGED"
  port_range            = "80"
  target                = google_compute_target_http_proxy.http_redirect.id
  ip_address            = google_compute_global_address.main.id
}
