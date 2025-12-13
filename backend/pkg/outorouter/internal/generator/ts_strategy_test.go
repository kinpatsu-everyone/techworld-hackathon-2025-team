package generator

import (
	"strings"
	"testing"

	"github.com/kinpatsu-everyone/backend-template/pkg/outorouter/internal/parser"
)

func TestTypeScriptStrategy_TableDriven(t *testing.T) {
	meta := &parser.Metadata{All: []parser.Endpoint{
		{
			Kind:         parser.KindUnaryJSON,
			Domain:       "user",
			Version:      1,
			MethodName:   "CreateUser",
			HTTPMethod:   "POST",
			RequestType:  "CreateUserRequest",
			ResponseType: "CreateUserResponse",
			Summary:      "summary",
			Tags:         []parser.Tag{"User"},
		},
	}}

	tests := []struct {
		name string
		gen  Strategy
	}{
		{name: "typescript", gen: TypeScriptStrategy{}},
		{name: "go", gen: GoTypeStrategy{PackageName: "p"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := New(tt.gen)
			code, err := g.Generate(meta)
			if err != nil {
				t.Fatalf("generate failed: %v", err)
			}
			if !strings.Contains(code, "user/v1/CreateUser") {
				t.Fatalf("generated code missing path: %s", code)
			}
			if !strings.Contains(code, "CreateUserRequest") {
				t.Fatalf("generated code missing request type")
			}
		})
	}
}

func TestTypeScriptClientStrategy(t *testing.T) {
	meta := &parser.Metadata{All: []parser.Endpoint{
		{
			Kind:         parser.KindUnaryJSON,
			Domain:       "user",
			Version:      1,
			MethodName:   "CreateUser",
			HTTPMethod:   "POST",
			RequestType:  "CreateUserRequest",
			ResponseType: "CreateUserResponse",
			Summary:      "Create a new user",
			Tags:         []parser.Tag{"User"},
		},
		{
			Kind:         parser.KindUnaryJSON,
			Domain:       "auth",
			Version:      1,
			MethodName:   "Login",
			HTTPMethod:   "POST",
			RequestType:  "LoginRequest",
			ResponseType: "LoginResponse",
			Summary:      "Login to the system",
			Tags:         []parser.Tag{"Auth"},
		},
	}}

	strategy := TypeScriptClientStrategy{BaseURL: "https://api.example.com"}
	gen := New(strategy)
	code, err := gen.Generate(meta)
	if err != nil {
		t.Fatalf("generate failed: %v", err)
	}

	// Verify the generated code contains expected elements
	expectedStrings := []string{
		// Configuration
		`const DEFAULT_BASE_URL = "https://api.example.com"`,
		`export interface ApiClientConfig`,
		`export function configureApiClient`,

		// Types
		`export interface ApiResponse<T>`,
		`export class ApiError extends Error`,

		// Request/Response types
		`export interface CreateUserRequest`,
		`export interface CreateUserResponse`,
		`export interface LoginRequest`,
		`export interface LoginResponse`,

		// Endpoints constant
		`export const Endpoints = {`,
		`CreateUser: "/user/v1/CreateUser"`,
		`Login: "/auth/v1/Login"`,

		// EndpointTypes mapping
		`export interface EndpointTypes`,
		`"/user/v1/CreateUser": {`,
		`request: CreateUserRequest`,
		`response: CreateUserResponse`,

		// API function
		`export async function api<P extends keyof EndpointTypes>`,

		// Pre-configured callers
		`export const apiCallers = {`,
		`CreateUser: createApiCaller(Endpoints.CreateUser)`,
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(code, expected) {
			t.Errorf("generated code missing expected string: %q\n\nGenerated code:\n%s", expected, code)
		}
	}
}

func TestTypeScriptClientStrategy_EmptyEndpoints(t *testing.T) {
	meta := &parser.Metadata{All: []parser.Endpoint{}}

	strategy := TypeScriptClientStrategy{}
	gen := New(strategy)
	code, err := gen.Generate(meta)
	if err != nil {
		t.Fatalf("generate failed: %v", err)
	}

	// Should still generate the base structure
	if !strings.Contains(code, "export interface ApiClientConfig") {
		t.Error("generated code missing ApiClientConfig interface")
	}
	if !strings.Contains(code, "export async function api") {
		t.Error("generated code missing api function")
	}
}
