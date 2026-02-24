package utils

import (
	"time"
)

// RunTicker executes the given function periodically based on the duration d.
// The function f is executed immediately, and then subsequently at every tick.
// Note: This blocks the current goroutine, so run it in a go routine.
func RunTicker(d time.Duration, f func()) {
	// Execute immediately first
	f()
	ticker := time.NewTicker(d)
	defer ticker.Stop()
	for range ticker.C {
		f()
	}
}
