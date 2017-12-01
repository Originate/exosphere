package types

// ServiceSource represents the service info as provided in application.yml
type ServiceSource struct {
	Location            string               `yaml:",omitempty"`
	DockerImage         string               `yaml:"docker-image,omitempty"`
	MessageTranslations []MessageTranslation `yaml:"message-translation,omitempty"`
}
