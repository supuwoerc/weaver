package transducer

import "database/sql"

func NullValue[T any](val sql.Null[T]) *T {
	if val.Valid {
		result := new(T)
		*result = val.V
		return result
	}
	return nil
}

func SqlNullValue[T any](val *T) sql.Null[T] {
	if val == nil {
		return sql.Null[T]{
			Valid: false,
		}
	}
	return sql.Null[T]{
		Valid: true,
		V:     *val,
	}
}
