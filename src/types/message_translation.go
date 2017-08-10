package types

// MessageTranslation is the mapping from public to internal message names
type MessageTranslation struct {
	Public   string `yaml:"public"`
	Internal string `yaml:"internal"`
}
