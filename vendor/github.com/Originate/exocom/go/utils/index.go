package utils

import (
	"errors"
	"time"
)

// WaitFor waits a maximum of 10 seconds the given condition to become true
// returning an error with the given message if and only if it does not
func WaitFor(condition func() bool, errorMessage string) error {
	success := make(chan bool, 1)

	go func() {
		for {
			time.Sleep(time.Millisecond)
			if condition() {
				success <- true
				break
			}
		}
	}()

	select {
	case <-success:
		return nil
	case <-time.After(time.Second * 10):
		return errors.New(errorMessage)
	}
}
