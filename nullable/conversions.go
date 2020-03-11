package nullable

import (
	"time"
)

// MakeString returns a new String
func MakeString(s *string) String {
	if s == nil {
		return String{Valid: false}
	}
	return String{String: *s, Valid: true}
}

// MakeInt64 returns a new Int64
func MakeInt64(i *int64) Int64 {
	if i == nil {
		return Int64{Valid: false}
	}
	return Int64{Int64: *i, Valid: true}
}

// MakeFloat64 returns a new Float64
func MakeFloat64(i *float64) Float64 {
	if i == nil {
		return Float64{Valid: false}
	}
	return Float64{Float64: *i, Valid: true}
}

// MakeBool creates a new Bool
func MakeBool(b *bool) Bool {
	if b == nil {
		return Bool{Valid: false}
	}
	return Bool{Bool: *b, Valid: true}
}

// MakeTime creates a new NullTime
func MakeTime(t time.Time) Time {
	if t == emptyTime {
		return Time{Valid: false}
	}
	return Time{Time: t, Valid: true}
}
