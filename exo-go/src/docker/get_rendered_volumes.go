package docker

import (
	"os"
	"path"
	"strings"

	"github.com/pkg/errors"
)

// GetRenderedVolumes returns the rendered paths to the given volumes
func GetRenderedVolumes(volumes []string, appName string, role string, homeDir string) ([]string, error) {
	dataPath := path.Join(homeDir, ".exosphere", appName, role, "data")
	renderedVolumes := []string{}
	if err := os.MkdirAll(dataPath, 0777); err != nil { //nolint gas
		return renderedVolumes, errors.Wrap(err, "Failed to create the necessary directories for the volumes")
	}
	for _, volume := range volumes {
		renderedVolumes = append(renderedVolumes, strings.Replace(volume, "{{EXO_DATA_PATH}}", dataPath, -1))
	}
	return renderedVolumes, nil
}
