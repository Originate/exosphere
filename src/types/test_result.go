package types

// TestResult represents the result of a test
type TestResult struct {
	Passed      bool
	Interrupted bool
	Error       error
}
