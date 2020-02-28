package nullable

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"reflect"
)

// Int64 defines a nullable int64
type Int64 struct {
	Int64 int64
	Valid bool // Valid is true if Time is not NULL
}

// MarshalJSON for Int64
func (n Int64) MarshalJSON() ([]byte, error) {
	var a *int64
	if n.Valid {
		a = &n.Int64
	}
	return json.Marshal(a)
}

// Value for Int64
func (n Int64) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Int64, nil
}

// UnmarshalJSON for Int64
func (n *Int64) UnmarshalJSON(b []byte) error {
	if bytes.EqualFold(b, nullLiteral) {
		n.Valid = false
		return nil
	}
	err := json.Unmarshal(b, &n.Int64)
	n.Valid = err == nil
	return err
}

// Scan for Int64
func (n *Int64) Scan(src interface{}) error {
	// Set initial state for subsequent scans.
	n.Valid = false

	var a sql.NullInt64
	if err := a.Scan(src); err != nil {
		return err
	}
	n.Int64 = a.Int64
	if reflect.TypeOf(src) != nil {
		n.Valid = true
	}
	return nil
}
