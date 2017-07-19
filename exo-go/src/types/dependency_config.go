package types

// DependencyConfig represents the configuration of a dependency
type DependencyConfig struct {
	Ports                 []string          `yaml:",omitempty"`
	Volumes               []string          `yaml:",omitempty"`
	OnlineText            string            `yaml:"online-text,omitempty"`
	DependencyEnvironment map[string]string `yaml:"dependency-environment,omitempty"`
	ServiceEnvironment    map[string]string `yaml:"service-environment,omitempty"`
}

// IsEmpty returns true if the given dependencyConfig object is empty
func (dependencyConfig *DependencyConfig) IsEmpty() bool {
	return len(dependencyConfig.Ports) == 0 && len(dependencyConfig.Volumes) == 0 && len(dependencyConfig.OnlineText) == 0 && len(dependencyConfig.DependencyEnvironment) == 0 && len(dependencyConfig.ServiceEnvironment) == 0
}
