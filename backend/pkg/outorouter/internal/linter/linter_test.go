package linter

import (
	"testing"

	"github.com/kinpatsu-everyone/backend-template/pkg/outorouter/internal/parser"
)

func TestLint_TableDriven(t *testing.T) {
	tests := []struct {
		name   string
		meta   *parser.Metadata
		assert func(t *testing.T, issues []Issue)
	}{
		{
			name: "multiple issues detected",
			meta: &parser.Metadata{All: []parser.Endpoint{
				{
					Kind:         parser.KindUnaryJSON,
					Domain:       "user",
					Version:      1,
					MethodName:   "createUser",
					HTTPMethod:   "post",
					RequestType:  "",
					ResponseType: "",
					Summary:      "",
				},
			}},
			assert: func(t *testing.T, issues []Issue) {
				if len(issues) < 4 {
					t.Fatalf("expected >=4 issues, got %d", len(issues))
				}
			},
		},
		{
			name: "no issues",
			meta: &parser.Metadata{All: []parser.Endpoint{
				{
					Kind:         parser.KindUnaryJSON,
					Domain:       "user",
					Version:      1,
					MethodName:   "CreateUser",
					HTTPMethod:   "POST",
					RequestType:  "Req",
					ResponseType: "Res",
					Summary:      "ok",
					Tags:         []parser.Tag{"User"},
				},
			}},
			assert: func(t *testing.T, issues []Issue) {
				if len(issues) != 0 {
					t.Fatalf("expected no issues, got %d", len(issues))
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			issues := Lint(tt.meta)
			if tt.assert != nil {
				tt.assert(t, issues)
			}
		})
	}
}
