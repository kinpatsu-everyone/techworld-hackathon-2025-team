package handler

import "context"

type HealthzRequest struct{}

func (h HealthzRequest) Validate() error {
	return nil
}

type HealthzResponse struct {
	Version string `json:"version"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

func Healthz(ctx context.Context, _ *HealthzRequest) (*HealthzResponse, error) {
	return &HealthzResponse{
		Version: "1.0.0",
		Status:  "ok",
		Message: "Service is healthy",
	}, nil
}
