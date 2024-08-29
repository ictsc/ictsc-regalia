package optional

import "database/sql"

// Of オプショナル型
type Of[T any] struct {
	V     T
	Valid bool
}

// New オプショナル型を作成
func New[T any](v T, valid bool) Of[T] {
	return Of[T]{
		V:     v,
		Valid: valid,
	}
}

// NewValid Validなオプショナル型を作成
func NewValid[T any](v T) Of[T] {
	return Of[T]{
		V:     v,
		Valid: true,
	}
}

// ToPointer ポインタに変換する
func ToPointer[T any](of Of[T]) *T {
	if of.Valid {
		return &of.V
	}

	return nil
}

// ToSQLNull sql.Nullに変換する
func ToSQLNull[T any](of Of[T]) sql.Null[T] {
	return sql.Null[T]{
		V:     of.V,
		Valid: of.Valid,
	}
}
