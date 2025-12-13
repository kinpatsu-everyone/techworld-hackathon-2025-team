#!/bin/bash
set -euo pipefail

# =============================================================================
# Deploy script for GCP Artifact Registry + Cloud Run
# =============================================================================

# Configuration (override via environment variables)
GCP_PROJECT="${GCP_PROJECT:-fanlav}"
GCP_REGION="${GCP_REGION:-asia-northeast1}"
ARTIFACT_REPO="${ARTIFACT_REPO:-kinpatsu}"
IMAGE_NAME="${IMAGE_NAME:-backend-api}"
SERVICE_NAME="${SERVICE_NAME:-backend-api}"
TAG="${1:-latest}"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Get git info for build args
get_version() {
    if git describe --tags --exact-match 2>/dev/null; then
        return
    fi
    echo "${TAG}"
}

get_commit_sha() {
    git rev-parse --short HEAD 2>/dev/null || echo "unknown"
}

# Validate required variables
validate_config() {
    if [[ -z "${GCP_PROJECT}" ]]; then
        log_error "GCP_PROJECT is required. Set it via environment variable."
        exit 1
    fi
}

# Main
main() {
    validate_config

    local registry="${GCP_REGION}-docker.pkg.dev"
    local full_image="${registry}/${GCP_PROJECT}/${ARTIFACT_REPO}/${IMAGE_NAME}:${TAG}"
    local version
    version=$(get_version)
    local commit_sha
    commit_sha=$(get_commit_sha)

    log_info "Building and deploying to GCP..."
    log_info "  Project:  ${GCP_PROJECT}"
    log_info "  Region:   ${GCP_REGION}"
    log_info "  Image:    ${full_image}"
    log_info "  Version:  ${version}"
    log_info "  Commit:   ${commit_sha}"
    echo ""

    # Configure Docker for Artifact Registry
    log_info "Configuring Docker authentication..."
    gcloud auth configure-docker "${registry}" --quiet

    # Build
    log_info "Building Docker image..."
    docker build \
        -f Dockerfile.production \
        --build-arg VERSION="${version}" \
        --build-arg COMMIT_SHA="${commit_sha}" \
        --platform linux/amd64 \
        -t "${full_image}" \
        .

    # Also tag as latest if not already
    if [[ "${TAG}" != "latest" ]]; then
        docker tag "${full_image}" "${registry}/${GCP_PROJECT}/${ARTIFACT_REPO}/${IMAGE_NAME}:latest"
    fi

    log_info "Build complete: ${full_image}"
    echo ""

    # Push
    log_info "Pushing to Artifact Registry..."
    docker push "${full_image}"

    if [[ "${TAG}" != "latest" ]]; then
        docker push "${registry}/${GCP_PROJECT}/${ARTIFACT_REPO}/${IMAGE_NAME}:latest"
    fi

    echo ""
    log_info "Push complete!"
    echo ""

    # # Deploy to Cloud Run
    # log_info "Deploying to Cloud Run..."
    # gcloud run services update "${SERVICE_NAME}" \
    #     --project="${GCP_PROJECT}" \
    #     --region="${GCP_REGION}" \
    #     --image="${full_image}" \
    #     --quiet

    # echo ""
    # log_info "Deploy complete!"
    # log_info "Image: ${full_image}"
    # log_info "Service: ${SERVICE_NAME}"
}

# Help
if [[ "${1:-}" == "-h" || "${1:-}" == "--help" ]]; then
    cat << EOF
Usage: $(basename "$0") [TAG]

Build and push Docker image to GCP Artifact Registry, then deploy to Cloud Run.

Arguments:
  TAG    Image tag (default: latest)

Required environment variables:
  GCP_PROJECT     GCP Project ID

Optional environment variables:
  GCP_REGION      GCP Region (default: asia-northeast1)
  ARTIFACT_REPO   Artifact Registry repository name (default: backend)
  IMAGE_NAME      Image name (default: api)
  SERVICE_NAME    Cloud Run service name (default: backend-api)

Examples:
  GCP_PROJECT=my-project $(basename "$0")           # Deploy with 'latest' tag
  GCP_PROJECT=my-project $(basename "$0") v1.0.0    # Deploy with 'v1.0.0' tag

Before running:
  # Authenticate with GCP
  gcloud auth login
  gcloud config set project YOUR_PROJECT_ID

  # Create Artifact Registry repository (first time only)
  gcloud artifacts repositories create backend \\
    --repository-format=docker \\
    --location=asia-northeast1
EOF
    exit 0
fi

main "$@"
