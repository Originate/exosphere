package utils

import "time"

// Retry will execute the given function until it returns true with the given
// delay between calls
func Retry(delayInMilliseconds int, attemptFn func() bool) {
	for {
		if attemptFn() {
			return
		}
		time.Sleep(time.Duration(delayInMilliseconds) * time.Millisecond)
	}
}
