package transducer

import (
	"database/sql"
	"time"
)

func NullTime(t sql.NullTime) time.Time {
	if !t.Valid {
		return time.Time{}
	}
	return t.Time
}

func SqlNullTime(t time.Time) sql.NullTime {
	if t.IsZero() {
		return sql.NullTime{Valid: false}
	}
	return sql.NullTime{Valid: true, Time: t}
}
