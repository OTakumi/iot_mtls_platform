// Package datetime provides a value object for handling time.
package datetime

import "time"

// Datetime is a value object for unifying timestamp formats.
type Datetime struct {
	Value time.Time
}

// Now returns the current time as a Datetime object.
func Now() Datetime {
	return Datetime{Value: time.Now()}
}
