package types

// ServiceData represents the service info as provided in application.yml
type ServiceData struct {
	Location            string               `yaml:",omitempty"`
	DockerImage         string               `yaml:"docker-image,omitempty"`
	MessageTranslations []MessageTranslation `yaml:"message-translation,omitempty"`
}
