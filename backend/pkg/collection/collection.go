package collection

import (
	"cmp"
	"slices"
)

// Number は Sum 系 API で扱う整数・浮動小数の列挙型制約です。
type Number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64
}

// Map はスライスの各要素に変換関数 f を適用した新しいスライスを返します。
// len(src) == 0 または src == nil の場合は長さ 0 のスライスを返します。
func Map[T any, R any](src []T, f func(T) R) []R {
	if len(src) == 0 {
		return []R{}
	}

	out := make([]R, len(src))
	for i, v := range src {
		out[i] = f(v)
	}
	return out
}

// FlatMap は要素ごとにスライスへ展開し、それらを 1 本のスライスへ連結します。
func FlatMap[T any, R any](src []T, f func(T) []R) []R {
	if len(src) == 0 {
		return []R{}
	}

	out := make([]R, 0, len(src))
	for _, v := range src {
		out = append(out, f(v)...)
	}
	if len(out) == 0 {
		return []R{}
	}
	return out
}

// Filter は述語 pred を満たす要素のみで構成される新しいスライスを返します。
func Filter[T any](src []T, pred func(T) bool) []T {
	if len(src) == 0 {
		return []T{}
	}

	out := make([]T, 0, len(src))
	for _, v := range src {
		if pred(v) {
			out = append(out, v)
		}
	}
	if len(out) == 0 {
		return []T{}
	}
	return out
}

// FilterNot は述語を満たさない要素のみを返します。
func FilterNot[T any](src []T, pred func(T) bool) []T {
	return Filter(src, func(v T) bool { return !pred(v) })
}

// Partition は述語 pred を基準にスライスを 2 つに分割します。
func Partition[T any](src []T, pred func(T) bool) (matched []T, unmatched []T) {
	if len(src) == 0 {
		return []T{}, []T{}
	}

	matched = make([]T, 0, len(src)/2+1)
	unmatched = make([]T, 0, len(src)/2+1)
	for _, v := range src {
		if pred(v) {
			matched = append(matched, v)
		} else {
			unmatched = append(unmatched, v)
		}
	}
	if len(matched) == 0 {
		matched = []T{}
	}
	if len(unmatched) == 0 {
		unmatched = []T{}
	}
	return matched, unmatched
}

// Reduce はスライスの左畳み込みを行います。空の場合は false を返します。
func Reduce[T any](src []T, f func(T, T) T) (T, bool) {
	var zero T
	if len(src) == 0 {
		return zero, false
	}

	acc := src[0]
	for i := 1; i < len(src); i++ {
		acc = f(acc, src[i])
	}
	return acc, true
}

// Fold は初期値 init から畳み込みを行います。
func Fold[T any, R any](src []T, init R, f func(R, T) R) R {
	acc := init
	for _, v := range src {
		acc = f(acc, v)
	}
	return acc
}

// Find は述語を満たす最初の要素を返します。
func Find[T any](src []T, pred func(T) bool) (T, bool) {
	var zero T
	for _, v := range src {
		if pred(v) {
			return v, true
		}
	}
	return zero, false
}

// Any は少なくとも 1 要素が述語を満たす場合に true を返します。
func Any[T any](src []T, pred func(T) bool) bool {
	for _, v := range src {
		if pred(v) {
			return true
		}
	}
	return false
}

// All はすべての要素が述語を満たす場合に true を返します。
func All[T any](src []T, pred func(T) bool) bool {
	for _, v := range src {
		if !pred(v) {
			return false
		}
	}
	return true
}

// None はいずれの要素も述語を満たさない場合に true を返します。
func None[T any](src []T, pred func(T) bool) bool {
	return !Any(src, pred)
}

// Count は述語を満たす要素数を返します。
func Count[T any](src []T, pred func(T) bool) int {
	c := 0
	for _, v := range src {
		if pred(v) {
			c++
		}
	}
	return c
}

// Sum は数値スライスの総和を返します。
func Sum[T Number](src []T) T {
	var total T
	for _, v := range src {
		total += v
	}
	return total
}

// SumBy は任意型から数値を抽出して総和を計算します。
func SumBy[T any, N Number](src []T, f func(T) N) N {
	var total N
	for _, v := range src {
		total += f(v)
	}
	return total
}

// GroupBy は keySelector が返すキーごとに要素をグルーピングします。
func GroupBy[T any, K comparable](src []T, keySelector func(T) K) map[K][]T {
	if len(src) == 0 {
		return map[K][]T{}
	}

	out := make(map[K][]T, len(src))
	for _, v := range src {
		k := keySelector(v)
		out[k] = append(out[k], v)
	}
	return out
}

// GroupByFn はキーと値の両方を変換してグルーピングします。
func GroupByFn[T any, K comparable, V any](src []T, keySelector func(T) K, valueSelector func(T) V) map[K][]V {
	if len(src) == 0 {
		return map[K][]V{}
	}

	out := make(map[K][]V, len(src))
	for _, v := range src {
		k := keySelector(v)
		out[k] = append(out[k], valueSelector(v))
	}
	return out
}

// ToMap はスライスを map へ変換します。重複キーは後勝ちです。
func ToMap[T any, K comparable, V any](src []T, keySelector func(T) K, valueSelector func(T) V) map[K]V {
	out := make(map[K]V, len(src))
	for _, v := range src {
		out[keySelector(v)] = valueSelector(v)
	}
	return out
}

// Associate は要素から (key, value) を生成して map を構築します。
func Associate[T any, K comparable, V any](src []T, transform func(T) (K, V)) map[K]V {
	out := make(map[K]V, len(src))
	for _, v := range src {
		k, value := transform(v)
		out[k] = value
	}
	return out
}

// SortBy はキー抽出関数に従って昇順ソートした新しいスライスを返します。
func SortBy[T any, K cmp.Ordered](src []T, keySelector func(T) K) []T {
	if len(src) == 0 {
		return []T{}
	}

	out := slices.Clone(src)
	slices.SortFunc(out, func(a, b T) int {
		return cmp.Compare(keySelector(a), keySelector(b))
	})
	return out
}

// SortWith は less 関数を用いたソート結果を返します。
func SortWith[T any](src []T, less func(a, b T) bool) []T {
	if len(src) == 0 {
		return []T{}
	}

	out := slices.Clone(src)
	slices.SortFunc(out, func(a, b T) int {
		switch {
		case less(a, b):
			return -1
		case less(b, a):
			return 1
		default:
			return 0
		}
	})
	return out
}

// Distinct は comparable な要素の重複を除去し、先頭出現順を保ちます。
func Distinct[T comparable](src []T) []T {
	if len(src) == 0 {
		return []T{}
	}

	seen := make(map[T]struct{}, len(src))
	out := make([]T, 0, len(src))
	for _, v := range src {
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		out = append(out, v)
	}
	if len(out) == 0 {
		return []T{}
	}
	return out
}

// DistinctBy はキー抽出結果を元に重複を除去します。
func DistinctBy[T any, K comparable](src []T, keySelector func(T) K) []T {
	if len(src) == 0 {
		return []T{}
	}

	seen := make(map[K]struct{}, len(src))
	out := make([]T, 0, len(src))
	for _, v := range src {
		key := keySelector(v)
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		out = append(out, v)
	}
	if len(out) == 0 {
		return []T{}
	}
	return out
}

// Concat は任意個のスライスを連結します。
func Concat[T any](slicesList ...[]T) []T {
	total := 0
	for _, s := range slicesList {
		total += len(s)
	}

	if total == 0 {
		return []T{}
	}

	out := make([]T, 0, total)
	for _, s := range slicesList {
		out = append(out, s...)
	}
	return out
}

// Flatten は二次元スライスを一次元へ変換します。
func Flatten[T any](src [][]T) []T {
	return Concat(src...)
}

// Chunked は size 要素ごとに区切ったスライスを返します。
func Chunked[T any](src []T, size int) [][]T {
	if size <= 0 {
		panic("collection: chunk size must be positive")
	}
	if len(src) == 0 {
		return [][]T{}
	}

	out := make([][]T, 0, (len(src)+size-1)/size)
	for i := 0; i < len(src); i += size {
		end := i + size
		if end > len(src) {
			end = len(src)
		}
		chunk := make([]T, end-i)
		copy(chunk, src[i:end])
		out = append(out, chunk)
	}
	return out
}

// ChunkWhile は隣接要素の関係 predicate が false になるたびに新しいチャンクを切ります。
func ChunkWhile[T any](src []T, predicate func(prev, curr T) bool) [][]T {
	if len(src) == 0 {
		return [][]T{}
	}

	out := make([][]T, 0)
	current := []T{src[0]}
	for i := 1; i < len(src); i++ {
		prev := src[i-1]
		curr := src[i]
		if predicate(prev, curr) {
			current = append(current, curr)
			continue
		}
		out = append(out, current)
		current = []T{curr}
	}
	out = append(out, current)
	return out
}

// Take は先頭から最大 n 要素を返します。
func Take[T any](src []T, n int) []T {
	if n <= 0 || len(src) == 0 {
		return []T{}
	}
	if n >= len(src) {
		out := make([]T, len(src))
		copy(out, src)
		return out
	}
	out := make([]T, n)
	copy(out, src[:n])
	return out
}

// TakeLast は末尾から最大 n 要素を返します。
func TakeLast[T any](src []T, n int) []T {
	if n <= 0 || len(src) == 0 {
		return []T{}
	}
	if n >= len(src) {
		out := make([]T, len(src))
		copy(out, src)
		return out
	}
	out := make([]T, n)
	copy(out, src[len(src)-n:])
	return out
}

// Drop は先頭から n 要素を捨てた残りを返します。
func Drop[T any](src []T, n int) []T {
	if n <= 0 {
		out := make([]T, len(src))
		copy(out, src)
		return out
	}
	if n >= len(src) {
		return []T{}
	}
	out := make([]T, len(src)-n)
	copy(out, src[n:])
	return out
}

// DropLast は末尾から n 要素を捨てた残りを返します。
func DropLast[T any](src []T, n int) []T {
	if n <= 0 {
		out := make([]T, len(src))
		copy(out, src)
		return out
	}
	if n >= len(src) {
		return []T{}
	}
	out := make([]T, len(src)-n)
	copy(out, src[:len(src)-n])
	return out
}

// Reverse はスライスを逆順にした新しいスライスを返します。
func Reverse[T any](src []T) []T {
	if len(src) == 0 {
		return []T{}
	}
	out := make([]T, len(src))
	for i := range src {
		out[len(src)-1-i] = src[i]
	}
	return out
}

// IndexOf は指定要素の最初の位置を返し、存在しなければ -1 を返します。
func IndexOf[T comparable](src []T, target T) int {
	for i, v := range src {
		if v == target {
			return i
		}
	}
	return -1
}

// Contains は要素の存在を返します。
func Contains[T comparable](src []T, target T) bool {
	return IndexOf(src, target) >= 0
}

// Keys は map のキーを昇順は保証せずに収集します。
func Keys[K comparable, V any](m map[K]V) []K {
	if len(m) == 0 {
		return []K{}
	}
	out := make([]K, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}

// Values は map の値を収集します。
func Values[K comparable, V any](m map[K]V) []V {
	if len(m) == 0 {
		return []V{}
	}
	out := make([]V, 0, len(m))
	for _, v := range m {
		out = append(out, v)
	}
	return out
}

// MapKeys はキーを変換した新しい map を返します。重複キーは後勝ちです。
func MapKeys[K1 comparable, K2 comparable, V any](m map[K1]V, keyFn func(K1) K2) map[K2]V {
	out := make(map[K2]V, len(m))
	for k, v := range m {
		out[keyFn(k)] = v
	}
	return out
}

// MapValues は値を変換した新しい map を返します。
func MapValues[K comparable, V1 any, V2 any](m map[K]V1, valueFn func(V1) V2) map[K]V2 {
	out := make(map[K]V2, len(m))
	for k, v := range m {
		out[k] = valueFn(v)
	}
	return out
}

// FilterMap は述語を満たすエントリのみで構成される map を返します。
func FilterMap[K comparable, V any](m map[K]V, pred func(K, V) bool) map[K]V {
	out := make(map[K]V)
	for k, v := range m {
		if pred(k, v) {
			out[k] = v
		}
	}
	return out
}

// Pair は 2 つの型要素をまとめる汎用的な組構造です。Zip/Unzip 等で利用します。
type Pair[A any, B any] struct {
	First  A
	Second B
}

// ForEach は副作用関数 f を全要素に適用します。戻り値はなく、入力スライスは不変です。
func ForEach[T any](src []T, f func(T)) {
	for _, v := range src {
		f(v)
	}
}

// Zip は 2 つのスライスを最短長に揃えて Pair のスライスへ結合します。
// 長さが異なる場合は短い方に合わせます。両方空または nil の場合は空スライスを返します。
func Zip[A any, B any](as []A, bs []B) []Pair[A, B] {
	n := len(as)
	if len(bs) < n {
		n = len(bs)
	}
	if n == 0 {
		return []Pair[A, B]{}
	}
	out := make([]Pair[A, B], n)
	for i := 0; i < n; i++ {
		out[i] = Pair[A, B]{First: as[i], Second: bs[i]}
	}
	return out
}

// Unzip は Pair のスライスを 2 つのスライスに分解します。
func Unzip[A any, B any](pairs []Pair[A, B]) ([]A, []B) {
	if len(pairs) == 0 {
		return []A{}, []B{}
	}
	as := make([]A, len(pairs))
	bs := make([]B, len(pairs))
	for i, p := range pairs {
		as[i] = p.First
		bs[i] = p.Second
	}
	return as, bs
}

// Min は cmp.Ordered な要素の最小値を返します。空の場合は (zero,false) を返します。
func Min[T cmp.Ordered](src []T) (T, bool) {
	var zero T
	if len(src) == 0 {
		return zero, false
	}
	m := src[0]
	for i := 1; i < len(src); i++ {
		if src[i] < m {
			m = src[i]
		}
	}
	return m, true
}

// Max は cmp.Ordered な要素の最大値を返します。空の場合は (zero,false) を返します。
func Max[T cmp.Ordered](src []T) (T, bool) {
	var zero T
	if len(src) == 0 {
		return zero, false
	}
	m := src[0]
	for i := 1; i < len(src); i++ {
		if src[i] > m {
			m = src[i]
		}
	}
	return m, true
}

// MinBy は keySelector の結果が最小となる要素を返します。空の場合は (zero,false)。
func MinBy[T any, K cmp.Ordered](src []T, keySelector func(T) K) (T, bool) {
	var zero T
	if len(src) == 0 {
		return zero, false
	}
	best := src[0]
	bestKey := keySelector(best)
	for i := 1; i < len(src); i++ {
		k := keySelector(src[i])
		if k < bestKey {
			best = src[i]
			bestKey = k
		}
	}
	return best, true
}

// MaxBy は keySelector の結果が最大となる要素を返します。空の場合は (zero,false)。
func MaxBy[T any, K cmp.Ordered](src []T, keySelector func(T) K) (T, bool) {
	var zero T
	if len(src) == 0 {
		return zero, false
	}
	best := src[0]
	bestKey := keySelector(best)
	for i := 1; i < len(src); i++ {
		k := keySelector(src[i])
		if k > bestKey {
			best = src[i]
			bestKey = k
		}
	}
	return best, true
}

// Avg は数値スライスの平均値を返します。空の場合は (0,false) を返します。
// 戻り値は float64 に正規化されます (整数型でも精度を保つため)。
func Avg[T Number](src []T) (float64, bool) {
	if len(src) == 0 {
		return 0, false
	}
	var sum float64
	for _, v := range src {
		sum += float64(v)
	}
	return sum / float64(len(src)), true
}

// AvgBy は任意型から抽出した数値の平均値を返します。空の場合は (0,false)。
func AvgBy[T any, N Number](src []T, f func(T) N) (float64, bool) {
	if len(src) == 0 {
		return 0, false
	}
	var sum float64
	for _, v := range src {
		sum += float64(f(v))
	}
	return sum / float64(len(src)), true
}
