package types

// ShutdownConfig represents the configuration to use when
// shutting down an application
type ShutdownConfig struct {
	CloseMessage string
	ErrorMessage string
}
