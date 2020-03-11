package nullable

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"reflect"
)

// Bool defines a nullable bool
type Bool struct {
	Bool  bool
	Valid bool
}

// MarshalJSON for Bool
func (n Bool) MarshalJSON() ([]byte, error) {
	var a *bool
	if n.Valid {
		a = &n.Bool
	}
	return json.Marshal(a)
}

// UnmarshalJSON for Bool
func (n *Bool) UnmarshalJSON(b []byte) error {
	var field *bool
	err := json.Unmarshal(b, &field)
	if field != nil {
		n.Valid = true
		n.Bool = *field
	}
	return err
}

// Scan implements the Scanner interface from database/sql
func (n *Bool) Scan(src interface{}) error {
	// Set initial state for subsequent scans.
	n.Valid = false

	var a sql.NullBool
	if err := a.Scan(src); err != nil {
		return err
	}
	n.Bool = a.Bool
	if reflect.TypeOf(src) != nil {
		n.Valid = true
	}
	return nil
}

// Value returns the database/sql driver value for Bool
func (n Bool) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Bool, nil
}

