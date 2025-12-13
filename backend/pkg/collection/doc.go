package collection

// Package collection は Scala や Kotlin の標準コレクション API を参考にした
// ジェネリック関数群を提供します。すべての関数はスライスを不変データとして扱い、
// 入力が nil の場合は長さ 0 のスライスとして処理します。
//
// # 設計方針
//   - チェーンスタイルではなく、collection.Map のような関数呼び出しスタイルを採用します。
//   - すべての変換は新しいスライス・マップを返し、入力データを破壊しません。
//   - predicate や変換関数内のエラー／panic は呼び出し側に委ねます。
//   - Kotlin/Scala 由来の API 名称を尊重しつつ、Go の命名規則を崩さない範囲で調整します。
//
// # API グループ
//   - 変換系: Map, FlatMap, Filter, FilterNot, Partition など
//   - 集約系: Reduce, Fold, Sum, SumBy, Count, Any, All, None など
//   - グルーピング・辞書系: GroupBy, GroupByFn, ToMap, Associate など
//   - 並び替え・重複除去: SortBy, SortWith, StableSort, Distinct, DistinctBy など
//   - スライス操作: Chunked, ChunkWhile, Window, Take, Drop, Reverse, Concat など
//   - マップユーティリティ: Keys, Values, MapKeys, MapValues, FilterMap など
//
// これらの関数は、DDD の値オブジェクトやエンティティ集合を扱うユースケースで
// 煩雑な for ループを減らし、意図を明確に表現するためのビルディングブロックです。
