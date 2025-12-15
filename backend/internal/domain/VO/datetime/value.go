package datetime

import "time"

// Object for unifying timestamp formats
type Datetime struct {
	Value time.Time
}

func Now() Datetime {
	return Datetime{Value: time.Now()}
}
