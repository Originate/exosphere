package types

import (
	"fmt"

	"github.com/pkg/errors"
)

// RemoteDependency represents a production dependency
type RemoteDependency struct {
	Config RemoteDependencyConfig `yaml:",omitempty"`
	Type   string
}

// ValidateFields validates that a remote config contains all required fields
func (p *RemoteDependency) ValidateFields() error {
	if p.Type == "rds" {
		err := p.Config.Rds.ValidateFields()
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("remote dependency %s has issues", p.Type))
		}
	}
	if p.Type == "exocom" {
		if p.Config.Version == "" {
			return errors.New("exocom dependency missing 'version' field")
		}
	}
	return nil
}
