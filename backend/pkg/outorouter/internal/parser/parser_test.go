package parser

import (
	"strings"
	"testing"
)

func TestParse_TableDriven(t *testing.T) {
	tests := []struct {
		name    string
		json    string
		wantErr bool
		assert  func(t *testing.T, meta *Metadata)
	}{
		{
			name: "valid single endpoint",
			json: `{
			  "user": {
			    "1": [
			      {
			        "kind": "JSON",
			        "domain": "user",
			        "version": 1,
			        "method_name": "CreateUser",
			        "http_method": "POST",
			        "request_type": "CreateUserRequest",
			        "response_type": "CreateUserResponse",
			        "summary": "ユーザー作成",
			        "description": "",
			        "tags": ["User"]
			      }
			    ]
			  }
			}`,
			assert: func(t *testing.T, meta *Metadata) {
				if got := len(meta.All); got != 1 {
					t.Fatalf("expected 1 endpoint, got %d", got)
				}
				if meta.All[0].Kind != KindUnaryJSON {
					t.Fatalf("unexpected kind: %s", meta.All[0].Kind)
				}
				if meta.All[0].Path() != "user/v1/CreateUser" {
					t.Fatalf("unexpected path: %s", meta.All[0].Path())
				}
			},
		},
		{
			name:    "invalid version",
			json:    `{"user":{"x":[]}}`,
			wantErr: true,
		},
		{
			name:    "unknown kind",
			json:    `{"user":{"1":[{"kind":"UNKNOWN","method_name":"X","http_method":"POST"}]}}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			meta, err := Parse(strings.NewReader(tt.json))
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if tt.assert != nil {
				tt.assert(t, meta)
			}
		})
	}
}
