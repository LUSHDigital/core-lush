package nullable

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
)

// RawJSON aliases json.RawMessage
type RawJSON json.RawMessage

// MarshalJSON for String
func (n RawJSON) MarshalJSON() ([]byte, error) {
	if len(n) == 0 {
		return []byte("null"), nil
	}
	a := json.RawMessage(n)
	return a.MarshalJSON()
}

// UnmarshalJSON for String
func (n *RawJSON) UnmarshalJSON(b []byte) error {
	var a json.RawMessage
	if err := json.Unmarshal(b, &a); err != nil {
		return err
	}
	c := RawJSON(a)
	*n = c
	return nil
}

// Scan for String
func (n *RawJSON) Scan(src interface{}) error {
	var a sql.NullString
	if err := a.Scan(src); err != nil {
		return err
	}
	jsn := RawJSON([]byte(a.String))
	*n = jsn
	return nil
}

// Value returns the database/sql driver value for RawJson
func (n RawJSON) Value() (driver.Value, error) {
	return string(n), nil
}

