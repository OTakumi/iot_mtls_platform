package datetime

import "time"

// Object for unifying timestamp formats
type Datetime struct {
	datetime time.Time
}

func Now() Datetime {
	now := time.Now()

	return Datetime{datetime: now}
}
