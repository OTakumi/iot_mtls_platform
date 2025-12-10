package datetime

import (
	"testing"
	"time"
)

func TestNow(t *testing.T) {
	// Call the Now() function
	dt := Now()

	// Get the current time for comparison
	now := time.Now()

	// Check if the returned datetime is close to the current time
	// Allow for a small difference due to execution time
	if dt.datetime.Before(now.Add(-time.Second)) || dt.datetime.After(now.Add(time.Second)) {
		t.Errorf("Now() returned %v, expected a time close to %v", dt.datetime, now)
	}

	// Ensure the time is not zero
	if dt.datetime.IsZero() {
		t.Errorf("Now() returned a zero time")
	}
}
