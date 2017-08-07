package config

// GetServiceBuiltDependencies returns the AppDependency objects for the
// dependencies defined in the given serviceConfig
func GetServiceBuiltDependencies(serviceConfig ServiceConfig, appConfig AppConfig, appDir, homeDir string) map[string]AppDependency {
	appBuiltDependencies := GetAppBuiltDependencies(appConfig, appDir, homeDir)
	result := map[string]AppDependency{}
	for _, dependency := range serviceConfig.Dependencies {
		if !dependency.Config.IsEmpty() {
			builtDependency := NewAppDependency(dependency, appConfig, appDir, homeDir)
			result[dependency.Name] = builtDependency
		} else {
			result[dependency.Name] = appBuiltDependencies[dependency.Name]
		}
	}
	return result
}
