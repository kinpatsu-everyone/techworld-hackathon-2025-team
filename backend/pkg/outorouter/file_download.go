package outorouter

import (
	"context"
	"net/http"
)

type FileDownloadResponseObject struct {
	Filename    string
	ContentType string
	Content     []byte
}

type FileDownloadHandlerFunc func(ctx context.Context, r *http.Request) (*FileDownloadResponseObject, error)

type FileDownloadEndpoint struct {
	Domain     string
	Version    uint8
	MethodName string

	Summary     string
	Description string
	Tags        []Tag

	Handler FileDownloadHandlerFunc
}
