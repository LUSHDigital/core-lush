package nullable

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"reflect"
)

// String defines a nullable string
type String struct {
	String string
	Valid  bool // Valid is true if String is not NULL
}

// MarshalJSON for String
func (n String) MarshalJSON() ([]byte, error) {
	var a *string
	if n.Valid {
		a = &n.String
	}
	return json.Marshal(a)
}

// UnmarshalJSON for String
func (n *String) UnmarshalJSON(b []byte) error {
	if bytes.EqualFold(b, nullLiteral) {
		n.Valid = false
		return nil
	}
	err := json.Unmarshal(b, &n.String)
	n.Valid = err == nil
	return err
}

// Scan implements the Scanner interface from database/sql
func (n *String) Scan(src interface{}) error {
	// Set initial state for subsequent scans.
	n.Valid = false

	var a sql.NullString
	if err := a.Scan(src); err != nil {
		return err
	}
	n.String = a.String
	if reflect.TypeOf(src) != nil {
		n.Valid = true
	}
	return nil
}

// Value returns the database/sql driver value for String
func (n String) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.String, nil
}
