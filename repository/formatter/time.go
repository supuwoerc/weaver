package formatter

import "database/sql"

func SqlNullTimeFormat(t sql.NullTime, layout string) string {
	if !t.Valid {
		return ""
	}
	return t.Time.Format(layout)
}
