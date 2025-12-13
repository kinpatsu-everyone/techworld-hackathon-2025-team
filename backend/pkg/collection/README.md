# collection パッケージ

Scala/Kotlin 風のコレクション操作を Go ジェネリクスで提供するユーティリティ群です。DDD のユースケース・ドメインサービス内で大量の for ループ記述を減らし、意図を明確にします。入力スライスは常に不変扱いで、`nil` は空スライスとして処理されます。

## 特徴
- 関数型スタイル: メソッドチェーンではなく `collection.Map(slice, f)` のような呼び出し。
- 不変性: すべて新しいスライス・map を返し、入力を変更しない。
- 型安全: Go 1.18+ ジェネリクス活用。
- 一貫したエッジケース: 空 / nil 入力は常に空スライスや `(zero,false)`。

## 主な関数一覧

### 変換 / フィルタ
- `Map`, `FlatMap`, `Filter`, `FilterNot`, `Partition`
- `Concat`, `Flatten`

### 集約 / 探索
- `Reduce`, `Fold`, `Find`, `Any`, `All`, `None`, `Count`
- `Sum`, `SumBy`, `Avg`, `AvgBy`
- `Min`, `Max`, `MinBy`, `MaxBy`

### グルーピング / 辞書変換
- `GroupBy`, `GroupByFn`
- `Associate`, `ToMap`

### ソート / 重複除去
- `SortBy`, `SortWith`
- `Distinct`, `DistinctBy`

### スライス操作
- `Chunked`, `ChunkWhile`
- `Take`, `TakeLast`, `Drop`, `DropLast`, `Reverse`
- `Zip`, `Unzip`

### map ユーティリティ
- `Keys`, `Values`, `MapKeys`, `MapValues`, `FilterMap`

### 副作用
- `ForEach`

## 使用例

```go
package hoge
import "github.com/kinpatsu-everyone/backend-template/pkg/collection"

// User ドメインを仮定
type User struct {
    ID    int
    Name  string
    Age   int
    Group string
}

func Example() {
    users := []User{
        {ID:1, Name:"Alice", Age:25, Group:"core"},
        {ID:2, Name:"Bob", Age:30, Group:"core"},
        {ID:3, Name:"Carol", Age:28, Group:"infra"},
    }

    // 年齢リストへ変換
    ages := collection.Map(users, func(u User) int { return u.Age })
    // 平均年齢
    avg, _ := collection.Avg(ages)

    // 28歳以上をフィルタ
    seniors := collection.Filter(users, func(u User) bool { return u.Age >= 28 })

    // グルーピング (Group -> []User)
    grouped := collection.GroupBy(users, func(u User) string { return u.Group })

    // Group 毎の年齢合計 map[string]int
    ageSumByGroup := collection.ToMap(grouped["core"], func(u User) int { return u.ID }, func(u User) int { return u.Age })

    // 名前と年齢の Pair スライス
    pairs := collection.Zip(collection.Map(users, func(u User) string { return u.Name }), ages)
    _ = pairs
    _ = avg
    _ = seniors
    _ = ageSumByGroup
    _ = grouped
}
```

## ベンチマーク

`collection_bench_test.go` に代表的な操作 (Map / Filter / Distinct / GroupBy / Pipeline) のベンチマークを収録。

```bash
go test -bench=. -run=^$ ./backend/go/pkg/collection
```

## 設計上の注意
- ソートは `slices.SortFunc` を利用し安定ソートではない（必要なら Stable 版拡張を検討）。
- `Chunked` は size <= 0 で panic。呼び出し側で入力検証する前提。
- 並列化は将来拡張 (ParallelMap など) へ委譲。

## 今後の拡張候補
- 並列版 (`ParallelMap`, `ParallelFilter`)
- `StableSortBy` / `SortedDescending` などのソートバリエーション
- イテレータ / ストリーム形式 (lazy evaluation) の軽量導入
- エラー対応版 (`MapE`, `FoldE`) で副作用を扱うパターン

## ライセンス
内部ユーティリティのため特別なライセンス指定なし。外部公開時は OSS ライセンス明記を検討。

