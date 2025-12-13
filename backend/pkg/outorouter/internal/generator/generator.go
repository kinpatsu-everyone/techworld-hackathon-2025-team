package generator

import "github.com/kinpatsu-everyone/backend-template/pkg/outorouter/internal/parser"

// Strategy は中間表現を任意フォーマットへ変換する戦略を表す。
type Strategy interface {
	Name() string
	Generate(meta *parser.Metadata) (string, error)
}

// Generator は Strategy パターンでコード生成を実行する。
type Generator struct {
	strategy Strategy
}

func New(strategy Strategy) *Generator {
	return &Generator{strategy: strategy}
}

// Generate は設定された Strategy で生成を行う。
func (g *Generator) Generate(meta *parser.Metadata) (string, error) {
	return g.strategy.Generate(meta)
}

func (g *Generator) StrategyName() string {
	return g.strategy.Name()
}
