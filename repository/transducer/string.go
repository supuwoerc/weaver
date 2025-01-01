package transducer

import "database/sql"

func NullString(val sql.NullString) *string {
	if val.Valid {
		str := new(string)
		*str = val.String
		return str
	}
	return nil
}

func SqlNullString(val *string) sql.NullString {
	if val == nil {
		return sql.NullString{
			Valid: false,
		}
	}
	return sql.NullString{
		Valid:  true,
		String: *val,
	}
}
