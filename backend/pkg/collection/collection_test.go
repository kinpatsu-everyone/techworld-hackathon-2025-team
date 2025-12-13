package collection

import (
	"slices"
	"testing"
)

type person struct {
	ID   int
	Team string
	Age  int
}

func TestMapFilterAndFlatMap(t *testing.T) {
	src := []int{1, 2, 3, 4}

	doubled := Map(src, func(v int) int { return v * 2 })
	if !slices.Equal(doubled, []int{2, 4, 6, 8}) {
		t.Fatalf("Map: got %v", doubled)
	}

	even := Filter(doubled, func(v int) bool { return v%4 == 0 })
	if !slices.Equal(even, []int{4, 8}) {
		t.Fatalf("Filter: got %v", even)
	}

	flat := FlatMap([]string{"go", "lang"}, func(s string) []rune { return []rune(s) })
	expected := []rune{'g', 'o', 'l', 'a', 'n', 'g'}
	if !slices.Equal(flat, expected) {
		t.Fatalf("FlatMap: got %v", flat)
	}

	if len(Map[int, int](nil, func(v int) int { return v })) != 0 {
		t.Fatalf("Map should return empty slice for nil input")
	}
}

func TestReduceFoldAndSum(t *testing.T) {
	sum, ok := Reduce([]int{1, 2, 3}, func(acc, v int) int { return acc + v })
	if !ok || sum != 6 {
		t.Fatalf("Reduce: got (%d, %v)", sum, ok)
	}

	if _, ok := Reduce([]int{}, func(acc, v int) int { return acc + v }); ok {
		t.Fatalf("Reduce should report false on empty slice")
	}

	folded := Fold([]string{"go", "lang"}, "", func(acc, v string) string { return acc + v })
	if folded != "golang" {
		t.Fatalf("Fold: got %s", folded)
	}

	if total := Sum([]int{3, 4, 5}); total != 12 {
		t.Fatalf("Sum: got %d", total)
	}

	persons := []person{{1, "core", 20}, {2, "core", 25}, {3, "infra", 30}}
	ageSum := SumBy(persons, func(p person) int { return p.Age })
	if ageSum != 75 {
		t.Fatalf("SumBy: got %d", ageSum)
	}
}

func TestPredicates(t *testing.T) {
	nums := []int{2, 4, 6}
	if !All(nums, func(v int) bool { return v%2 == 0 }) {
		t.Fatalf("All should return true")
	}
	if Any(nums, func(v int) bool { return v == 3 }) {
		t.Fatalf("Any should return false")
	}
	if !None(nums, func(v int) bool { return v > 10 }) {
		t.Fatalf("None should return true")
	}
	if Count(nums, func(v int) bool { return v >= 4 }) != 2 {
		t.Fatalf("Count mismatch")
	}

	v, ok := Find(nums, func(v int) bool { return v == 4 })
	if !ok || v != 4 {
		t.Fatalf("Find mismatch: %v, %v", v, ok)
	}
}

func TestGroupByAndAssociate(t *testing.T) {
	data := []person{
		{ID: 1, Team: "core", Age: 20},
		{ID: 2, Team: "core", Age: 21},
		{ID: 3, Team: "infra", Age: 30},
	}

	grouped := GroupBy(data, func(p person) string { return p.Team })
	expected := map[string][]person{
		"core":  {{ID: 1, Team: "core", Age: 20}, {ID: 2, Team: "core", Age: 21}},
		"infra": {{ID: 3, Team: "infra", Age: 30}},
	}
	if len(grouped) != len(expected) {
		t.Fatalf("GroupBy length mismatch: %#v", grouped)
	}
	for key, exp := range expected {
		got, ok := grouped[key]
		if !ok {
			t.Fatalf("GroupBy missing key %s", key)
		}
		if !slices.Equal(got, exp) {
			t.Fatalf("GroupBy mismatch for key %s: %v", key, got)
		}
	}

	groupedAges := GroupByFn(data, func(p person) string { return p.Team }, func(p person) int { return p.Age })
	expectedAges := map[string][]int{"core": {20, 21}, "infra": {30}}
	if len(groupedAges) != len(expectedAges) {
		t.Fatalf("GroupByFn length mismatch: %#v", groupedAges)
	}
	for key, exp := range expectedAges {
		got, ok := groupedAges[key]
		if !ok {
			t.Fatalf("GroupByFn missing key %s", key)
		}
		if !slices.Equal(got, exp) {
			t.Fatalf("GroupByFn mismatch for key %s: %v", key, got)
		}
	}

	asMap := Associate(data, func(p person) (int, string) { return p.ID, p.Team })
	if len(asMap) != len(data) || asMap[1] != "core" {
		t.Fatalf("Associate mismatch: %#v", asMap)
	}

	teamMap := ToMap(data, func(p person) int { return p.ID }, func(p person) string { return p.Team })
	if !equalIntStringMap(teamMap, map[int]string{1: "core", 2: "core", 3: "infra"}) {
		t.Fatalf("ToMap mismatch: %#v", teamMap)
	}
}

func TestSortAndDistinct(t *testing.T) {
	words := []string{"go", "scala", "kotlin", "go"}
	sorted := SortBy(words, func(s string) int { return len(s) })
	if !slices.Equal(sorted, []string{"go", "go", "scala", "kotlin"}) {
		t.Fatalf("SortBy mismatch: %v", sorted)
	}

	sortedLex := SortWith(words, func(a, b string) bool { return a < b })
	if !slices.Equal(sortedLex, []string{"go", "go", "kotlin", "scala"}) {
		t.Fatalf("SortWith mismatch: %v", sortedLex)
	}

	distinct := Distinct(words)
	if !slices.Equal(distinct, []string{"go", "scala", "kotlin"}) {
		t.Fatalf("Distinct mismatch: %v", distinct)
	}

	distinctLen := DistinctBy(words, func(s string) int { return len(s) })
	if len(distinctLen) != 3 {
		t.Fatalf("DistinctBy should keep first occurrences per key")
	}
}

func TestChunkingAndChunkWhile(t *testing.T) {
	src := []int{1, 2, 3, 4, 5}
	chunks := Chunked(src, 2)
	expected := [][]int{{1, 2}, {3, 4}, {5}}
	if len(chunks) != len(expected) {
		t.Fatalf("Chunked length mismatch: %v", chunks)
	}
	for i := range chunks {
		if !slices.Equal(chunks[i], expected[i]) {
			t.Fatalf("Chunked mismatch at %d: %v", i, chunks[i])
		}
	}

	chunkWhile := ChunkWhile([]int{1, 2, 4, 5, 10}, func(prev, curr int) bool { return curr-prev == 1 })
	target := [][]int{{1, 2}, {4, 5}, {10}}
	if len(chunkWhile) != len(target) {
		t.Fatalf("ChunkWhile length mismatch: %v", chunkWhile)
	}
	for i := range chunkWhile {
		if !slices.Equal(chunkWhile[i], target[i]) {
			t.Fatalf("ChunkWhile mismatch at %d: %v", i, chunkWhile[i])
		}
	}
}

func TestTakeDropReverseAndIndex(t *testing.T) {
	src := []int{1, 2, 3, 4, 5}
	if !slices.Equal(Take(src, 3), []int{1, 2, 3}) {
		t.Fatalf("Take mismatch")
	}
	if !slices.Equal(TakeLast(src, 2), []int{4, 5}) {
		t.Fatalf("TakeLast mismatch")
	}
	if !slices.Equal(Drop(src, 2), []int{3, 4, 5}) {
		t.Fatalf("Drop mismatch")
	}
	if !slices.Equal(DropLast(src, 3), []int{1, 2}) {
		t.Fatalf("DropLast mismatch")
	}
	if !slices.Equal(Reverse(src), []int{5, 4, 3, 2, 1}) {
		t.Fatalf("Reverse mismatch")
	}
	if idx := IndexOf(src, 3); idx != 2 {
		t.Fatalf("IndexOf mismatch: %d", idx)
	}
	if !Contains(src, 4) || Contains(src, 10) {
		t.Fatalf("Contains mismatch")
	}
}

func TestMapHelpers(t *testing.T) {
	m := map[string]int{"go": 1, "scala": 2}
	keys := Keys(m)
	if len(keys) != len(m) {
		t.Fatalf("Keys mismatch: %v", keys)
	}
	if len(Values(m)) != len(m) {
		t.Fatalf("Values mismatch")
	}

	mapped := MapKeys(m, func(k string) string { return k + "!" })
	if mapped["go!"] != 1 {
		t.Fatalf("MapKeys mismatch: %#v", mapped)
	}

	mappedValues := MapValues(m, func(v int) string { return string(rune('a' + v)) })
	if mappedValues["go"] != "b" {
		t.Fatalf("MapValues mismatch: %#v", mappedValues)
	}

	filtered := FilterMap(m, func(k string, v int) bool { return k == "go" })
	if len(filtered) != 1 || filtered["go"] != 1 {
		t.Fatalf("FilterMap mismatch: %#v", filtered)
	}
}

func TestChunkedPanicsOnInvalidSize(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("Chunked should panic on non-positive size")
		}
	}()
	_ = Chunked([]int{1}, 0)
}

func TestZipUnzipAndMinMax(t *testing.T) {
	a := []int{1, 2, 3}
	b := []string{"a", "b"}
	pairs := Zip(a, b)
	if len(pairs) != 2 || pairs[0].First != 1 || pairs[1].Second != "b" {
		t.Fatalf("Zip mismatch: %#v", pairs)
	}
	aa, bb := Unzip(pairs)
	if !slices.Equal(aa, []int{1, 2}) || !slices.Equal(bb, []string{"a", "b"}) {
		t.Fatalf("Unzip mismatch: %v %v", aa, bb)
	}

	if v, ok := Min([]int{5, 2, 9}); !ok || v != 2 {
		t.Fatalf("Min mismatch: %v %v", v, ok)
	}
	if v, ok := Max([]int{5, 2, 9}); !ok || v != 9 {
		t.Fatalf("Max mismatch: %v %v", v, ok)
	}

	people := []person{{ID: 1, Team: "x", Age: 30}, {ID: 2, Team: "x", Age: 20}, {ID: 3, Team: "y", Age: 25}}
	if p, ok := MinBy(people, func(p person) int { return p.Age }); !ok || p.ID != 2 {
		t.Fatalf("MinBy mismatch: %#v", p)
	}
	if p, ok := MaxBy(people, func(p person) int { return p.Age }); !ok || p.ID != 1 {
		t.Fatalf("MaxBy mismatch: %#v", p)
	}

	if avg, ok := Avg([]int{2, 4, 6}); !ok || avg != 4 {
		t.Fatalf("Avg mismatch: %v %v", avg, ok)
	}
	if avg, ok := AvgBy(people, func(p person) int { return p.Age }); !ok || avg < 24.9 || avg > 25.1 {
		t.Fatalf("AvgBy mismatch: %v", avg)
	}
}

func TestForEach(t *testing.T) {
	acc := 0
	ForEach([]int{1, 2, 3}, func(v int) { acc += v })
	if acc != 6 {
		t.Fatalf("ForEach side effect mismatch: %d", acc)
	}
}

func equalIntStringMap(a, b map[int]string) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if b[k] != v {
			return false
		}
	}
	return true
}
