package appSetup

import (
	"os"
	"path"

	"github.com/Originate/exosphere/exo-go/src/types"
	"github.com/Originate/exosphere/exo-go/test_helpers"
)

func joinDockerConfigMaps(maps ...map[string]types.DockerConfig) map[string]types.DockerConfig {
	result := map[string]types.DockerConfig{}
	for _, currentMap := range maps {
		for key, val := range currentMap {
			result[key] = val
		}
	}
	return result
}

// CheckoutApp copies the example app appName to cwd
func CheckoutApp(cwd, appName string) error {
	src := path.Join("..", "..", "..", "exosphere-shared", "example-apps", appName)
	dest := path.Join(cwd, appName)
	err := os.RemoveAll(dest)
	if err != nil {
		return err
	}
	return testHelpers.CopyDir(src, dest)
}
