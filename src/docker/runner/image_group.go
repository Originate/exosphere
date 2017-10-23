package runner

// ImageGroup is a group of images to run at once
type ImageGroup struct {
	ID          string
	Names       []string
	OnlineTexts map[string]string
}
