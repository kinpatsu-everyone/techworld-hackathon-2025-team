package linter

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/kinpatsu-everyone/backend-template/pkg/outorouter/internal/parser"
)

// LintLevel は指摘の重要度を表す。
type LintLevel string

const (
	LevelError   LintLevel = "error"
	LevelWarning LintLevel = "warning"
)

// Issue は1件の指摘を表す。
type Issue struct {
	Level   LintLevel
	Path    string
	Message string
}

// Lint は中間表現の品質を検査し、問題点を返す。
func Lint(meta *parser.Metadata) []Issue {
	issues := make([]Issue, 0)

	if meta == nil {
		return []Issue{{Level: LevelError, Path: "<root>", Message: "metadata is nil"}}
	}

	seen := make(map[string]struct{})
	methodRe := regexp.MustCompile(`^[A-Z][A-Za-z0-9]*$`)

	for _, ep := range meta.All {
		path := ep.Path()

		if ep.RequestType == "" {
			issues = append(issues, Issue{Level: LevelError, Path: path, Message: "request_type が空です"})
		}
		if ep.ResponseType == "" {
			issues = append(issues, Issue{Level: LevelError, Path: path, Message: "response_type が空です"})
		}
		if !methodRe.MatchString(ep.MethodName) {
			issues = append(issues, Issue{Level: LevelWarning, Path: path, Message: "MethodName はパスカルケースが推奨です"})
		}
		if ep.HTTPMethod == "" {
			issues = append(issues, Issue{Level: LevelError, Path: path, Message: "HTTPメソッドが未指定です"})
		} else if strings.ToUpper(ep.HTTPMethod) != ep.HTTPMethod {
			issues = append(issues, Issue{Level: LevelWarning, Path: path, Message: "HTTPメソッドは大文字で記述してください"})
		}

		key := fmt.Sprintf("%s#%s", ep.HTTPMethod, path)
		if _, ok := seen[key]; ok {
			issues = append(issues, Issue{Level: LevelError, Path: path, Message: "同一メソッド・パスのエンドポイントが重複しています"})
		} else {
			seen[key] = struct{}{}
		}

		if len(ep.Tags) == 0 {
			issues = append(issues, Issue{Level: LevelWarning, Path: path, Message: "Tags が設定されていません"})
		}
		if strings.TrimSpace(ep.Summary) == "" {
			issues = append(issues, Issue{Level: LevelWarning, Path: path, Message: "Summary が空です"})
		}
	}

	return issues
}
