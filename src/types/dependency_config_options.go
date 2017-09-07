package types

// DependencyConfigOptions represents the configuration of a dependency
type DependencyConfigOptions struct {
	Ports                 []string          `yaml:",omitempty"`
	Volumes               []string          `yaml:",omitempty"`
	OnlineText            string            `yaml:"online-text,omitempty"`
	DependencyEnvironment map[string]string `yaml:"dependency-environment,omitempty"`
	ServiceEnvironment    map[string]string `yaml:"service-environment,omitempty"`
	ExternalInProduction  bool              `yaml:"external-in-production,omitempty"`
}

// IsEmpty returns true if the given dependencyConfig object is empty
func (d *DependencyConfigOptions) IsEmpty() bool {
	return len(d.Ports) == 0 &&
		len(d.Volumes) == 0 &&
		d.OnlineText == "" &&
		len(d.DependencyEnvironment) == 0 &&
		len(d.ServiceEnvironment) == 0
}
