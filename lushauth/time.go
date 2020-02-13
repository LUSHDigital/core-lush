package lushauth

import (
	"time"
)

var (
	// TimeFunc is a variable with a function to determine the current time.
	// Can be overridden in a test environment to set the current time to whatever you want it to be.
	TimeFunc = time.Now

	// DefaultValidPeriod is the period a set of claims are valid.
	DefaultValidPeriod = 60 * time.Minute
)
