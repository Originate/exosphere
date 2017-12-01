package types

import (
	"fmt"

	"github.com/pkg/errors"
)

// RemoteDependency represents a production dependency
type RemoteDependency struct {
	Config  RemoteDependencyConfig `yaml:",omitempty"`
	Name    string
	Version string
}

//DbDependencies is a map from db engines to their underlying dependency
var DbDependencies = map[string]string{
	"postgres": "rds",
	"mysql":    "rds",
}

// GetDbDependency returns a map of db engines to the underlying dependency
func (p *RemoteDependency) GetDbDependency() string {
	return DbDependencies[p.Name]
}

// ValidateFields validates that a production config contains all required fields
func (p *RemoteDependency) ValidateFields() error {
	if p.GetDbDependency() != "" {
		err := p.Config.Rds.ValidateFields()
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("production dependency %s:%s has issues", p.Name, p.Version))
		}
	}
	return nil
}
