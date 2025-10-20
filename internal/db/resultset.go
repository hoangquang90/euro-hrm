package db

import (
	"database/sql/driver"
	"time"
)

type ResultSet struct {
	colNums map[string]int
	vals    []driver.Value
	rows    driver.Rows
}

func NewResultSet(rows driver.Rows) *ResultSet {
	cols := rows.(driver.RowsColumnTypeScanType).Columns()
	colNums := make(map[string]int, len(cols))
	for i := 0; i < len(cols); i++ {
		colNums[cols[i]] = i
	}
	vals := make([]driver.Value, len(cols))

	return &ResultSet{
		colNums: colNums,
		vals:    vals,
		rows:    rows,
	}
}

func (rs *ResultSet) Close() {
	rs.rows.Close()
}

func (rs *ResultSet) Next() error {
	return rs.rows.Next(rs.vals)
}

func (rs *ResultSet) GetString(fieldName string) string {
	return rs.vals[rs.colNums[fieldName]].(string)
}

func (rs *ResultSet) GetTime(fieldName string) *time.Time {
	v := rs.vals[rs.colNums[fieldName]]
	if v != nil {
		t := v.(time.Time)
		return &t
	} else {
		return nil
	}
}
