package collection

import (
	"math/rand"
	"testing"
	"time"
)

// ベンチマーク用の入力データをグローバルに保持し、計測ループ内での生成コストを回避。
var benchInts []int
var benchStructs []benchItem

type benchItem struct {
	ID    int
	Group int
	Value int
}

func init() {
	n := 10_000
	benchInts = make([]int, n)
	for i := 0; i < n; i++ {
		benchInts[i] = rand.Intn(1_000_000)
	}
	benchStructs = make([]benchItem, n)
	for i := 0; i < n; i++ {
		benchStructs[i] = benchItem{ID: i, Group: i % 128, Value: rand.Intn(10_000)}
	}
}

// Baseline: 手書き for ループで平方数へ変換
func BenchmarkBaselineManualMap(b *testing.B) {
	src := benchInts
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		out := make([]int, len(src))
		for j, v := range src {
			out[j] = v * v
		}
		_ = out
	}
}

func BenchmarkMap(b *testing.B) {
	src := benchInts
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Map(src, func(v int) int { return v * v })
	}
}

func BenchmarkFilter(b *testing.B) {
	src := benchInts
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Filter(src, func(v int) bool { return v%2 == 0 && v%3 == 0 })
	}
}

func BenchmarkDistinct(b *testing.B) {
	src := benchInts
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Distinct(src)
	}
}

func BenchmarkGroupBy(b *testing.B) {
	src := benchStructs
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = GroupBy(src, func(it benchItem) int { return it.Group })
	}
}

// 複合的操作: Map → Filter → GroupBy
func BenchmarkPipeline(b *testing.B) {
	src := benchStructs
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mapped := Map(src, func(it benchItem) benchItem { it.Value = it.Value * 2; return it })
		filtered := Filter(mapped, func(it benchItem) bool { return it.Value%5 == 0 })
		_ = GroupBy(filtered, func(it benchItem) int { return it.Group })
	}
}

// 生成コストを含むケース（参考）
func BenchmarkMapWithAllocation(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tmp := make([]int, 5_000)
		for j := range tmp {
			tmp[j] = j
		}
		_ = Map(tmp, func(v int) int { return v + 1 })
	}
}

// 時間計測の一例（模擬遅延を追加）
func BenchmarkGroupByWithSleep(b *testing.B) {
	src := benchStructs[:5_000]
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		g := GroupBy(src, func(it benchItem) int { return it.Group })
		if len(g) == 0 {
			b.Fatalf("unexpected empty grouping")
		}
		time.Sleep(time.Microsecond)
	}
}
