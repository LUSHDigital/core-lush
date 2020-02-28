package nullable

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"reflect"
)

// Float64 aliases sql.Float64
type Float64 struct {
	Float64 float64
	Valid   bool
}

// MarshalJSON for Float64
func (n Float64) MarshalJSON() ([]byte, error) {
	var a *float64
	if n.Valid {
		a = &n.Float64
	}
	return json.Marshal(a)
}

// Value for Float64
func (n Float64) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Float64, nil
}

// UnmarshalJSON for Float64
func (n *Float64) UnmarshalJSON(b []byte) error {
	if bytes.EqualFold(b, nullLiteral) {
		n.Valid = false
		return nil
	}
	err := json.Unmarshal(b, &n.Float64)
	n.Valid = err == nil
	return err
}

// Scan for Float64
func (n *Float64) Scan(src interface{}) error {
	// Set initial state for subsequent scans.
	n.Valid = false

	var a sql.NullFloat64
	if err := a.Scan(src); err != nil {
		return err
	}
	n.Float64 = a.Float64
	if reflect.TypeOf(src) != nil {
		n.Valid = true
	}
	return nil
}
