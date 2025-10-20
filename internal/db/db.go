package db

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"reflect"
	"time"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var caser = cases.Title(language.Vietnamese)

// https://gist.github.com/keidrun/d1b2791f840753e25070771b857af7ba

// NullString is an alias for sql.NullString data type
type NullString struct {
	sql.NullString
}

// MarshalJSON for NullString
func (ns *NullString) MarshalJSON() []byte {
	if !ns.Valid {
		return []byte(`""`)
	}
	b, err := json.Marshal(ns.String)
	if err != nil {
		return []byte(`""`)
	}
	return b
}

func (s *NullString) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	s.String = str
	s.Valid = true
	return nil
}

func (s *NullString) Capitalize() {
	if s.Valid {
		s.String = caser.String(s.String)
	}
}

// NullTime is an alias for mysql.NullTime data type
type NullTime struct {
	sql.NullTime
}

// MarshalJSON for NullTime
func (nt *NullTime) MarshalJSON() []byte {
	if !nt.Valid {
		return []byte("null")
	}
	b, err := json.Marshal(nt.Time.Format(time.RFC3339))
	if err != nil {
		return []byte("null")
	}
	return b
}

// UnmarshalJSON for NullTime
func (nt *NullTime) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		nt.Valid = false
		return err
	}
	if s == "" {
		nt.Valid = false
		return nil
	}
	x, err := time.Parse(time.RFC3339, s)
	if err != nil {
		nt.Valid = false
		return err
	}
	nt.Time = x
	nt.Valid = true
	return nil
}

// NullInt64 is an alias for sql.NullInt64 data type
type NullInt64 struct {
	sql.NullInt64
}

// MarshalJSON for NullInt64
func (ni *NullInt64) MarshalJSON() []byte {
	if !ni.Valid {
		return []byte("0")
	}
	b, err := json.Marshal(ni.Int64)
	if err != nil {
		return []byte("0")
	}
	return b
}

// UnmarshalJSON for NullInt64
func (ni *NullInt64) UnmarshalJSON(b []byte) error {
	var i *int64
	if err := json.Unmarshal(b, &i); err != nil {
		return err
	}
	if i != nil {
		ni.Valid = true
		ni.Int64 = *i
	} else {
		ni.Valid = false
	}
	return nil
}

// Convert sql.Nullxxx => original go data type
// Reason: Validators library in gin doesn't support nullable type in sql
func NullTypeFunc(field reflect.Value) interface{} {

	if valuer, ok := field.Interface().(driver.Valuer); ok {

		val, err := valuer.Value()
		if err == nil {
			if val != nil {
				return val
			}
			switch valuer.(type) {
			case NullString:
				return ""
			case NullInt64:
				return 0
			default:
				return nil
			}
		}
		// handle the error how you want
	}

	return nil
}
