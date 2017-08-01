package messageTranslator

import "strings"

// GetInternalMessageNameOptions are the options to the GetInternalMessageName func
type GetInternalMessageNameOptions struct {
	Namespace         string
	PublicMessageName string
}

// GetPublicMessageNameOptions are the options to the GetPublicMessageName func
type GetPublicMessageNameOptions struct {
	Namespace           string
	ClientName          string
	InternalMessageName string
}

// GetInternalMessageName returns the internal message name for the given options
func GetInternalMessageName(opts *GetInternalMessageNameOptions) string {
	if !strings.Contains(opts.PublicMessageName, ".") || opts.Namespace == "" {
		return opts.PublicMessageName
	}
	return opts.Namespace + "." + strings.Split(opts.PublicMessageName, ".")[1]
}

// GetPublicMessageName returns the public message name for the given options
func GetPublicMessageName(opts *GetPublicMessageNameOptions) string {
	if !strings.Contains(opts.InternalMessageName, ".") || opts.Namespace == "" {
		return opts.InternalMessageName
	}
	return opts.ClientName + "." + strings.Split(opts.InternalMessageName, ".")[1]
}
