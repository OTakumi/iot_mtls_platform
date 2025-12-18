package datetime_test

import (
	"testing"
	"time"

	"backend/internal/domain/VO/datetime"
)

func TestNow(t *testing.T) {
	t.Parallel()
	// Call the Now() function
	dt := datetime.Now()

	// Get the current time for comparison
	now := time.Now()

	// Check if the returned datetime is close to the current time
	// Allow for a small difference due to execution time
	if dt.Value.Before(now.Add(-time.Second)) || dt.Value.After(now.Add(time.Second)) {
		t.Errorf("Now() returned %v, expected a time close to %v", dt.Value, now)
	}

	// Ensure the time is not zero
	if dt.Value.IsZero() {
		t.Errorf("Now() returned a zero time")
	}
}
