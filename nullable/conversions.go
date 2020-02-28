package nullable

import (
	"time"
)

// ToString returns a new String
func ToString(s *string) String {
	if s == nil {
		return String{Valid: false}
	}
	return String{String: *s, Valid: true}
}

// ToInt64 returns a new Int64
func ToInt64(i *int64) Int64 {
	if i == nil {
		return Int64{Valid: false}
	}
	return Int64{Int64: *i, Valid: true}
}

// ToFloat64 returns a new Float64
func ToFloat64(i *float64) Float64 {
	if i == nil {
		return Float64{Valid: false}
	}
	return Float64{Float64: *i, Valid: true}
}

// ToBool creates a new Bool
func ToBool(b *bool) Bool {
	if b == nil {
		return Bool{Valid: false}
	}
	return Bool{Bool: *b, Valid: true}
}

// ToTime creates a new NullTime
func ToTime(t time.Time) Time {
	if t == emptyTime {
		return Time{Valid: false}
	}
	return Time{Time: t, Valid: true}
}
