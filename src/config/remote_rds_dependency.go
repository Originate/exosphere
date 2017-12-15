package config

import (
	"github.com/Originate/exosphere/src/types"
	"github.com/Originate/exosphere/src/types/context"
)

type remoteRdsDependency struct {
	name       string
	config     types.RemoteDependency
	appContext *context.AppContext
}

//GetDeploymentConfig returns configuration needed in deployment
func (r *remoteRdsDependency) GetDeploymentConfig() (map[string]string, error) {
	config := map[string]string{
		"engine":             r.config.Config.Rds.Engine,
		"engineVersion":      r.config.Config.Rds.EngineVersion,
		"allocatedStorage":   r.config.Config.Rds.AllocatedStorage,
		"instanceClass":      r.config.Config.Rds.InstanceClass,
		"name":               r.config.Config.Rds.DbName,
		"username":           r.config.Config.Rds.Username,
		"passwordSecretName": r.config.Config.Rds.PasswordSecretName,
		"storageType":        r.config.Config.Rds.StorageType,
	}
	return config, nil
}
