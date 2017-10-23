package composebuilder

import (
	"github.com/Originate/exosphere/src/types"
)

// ApplicationOptions are the options to
type ApplicationOptions struct {
	AppConfig types.AppConfig
	AppDir    string
	BuildMode BuildMode
	HomeDir   string
}
