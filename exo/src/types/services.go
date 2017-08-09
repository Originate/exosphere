package types

// Services represents the mapping of protection level to services
type Services struct {
	Public  map[string]ServiceData
	Private map[string]ServiceData
}
