package composerunner

import (
	"fmt"
	"regexp"
	"sync"
	"time"

	"github.com/Originate/exosphere/src/docker/compose"
	"github.com/pkg/errors"
)

func runImageGroup(options RunOptions, imageGroup ImageGroup) error {
	if len(imageGroup.Names) == 0 {
		return nil
	}
	cmdPlus, err := compose.RunImages(compose.ImagesOptions{
		DockerComposeDir: options.DockerComposeDir,
		ImageNames:       imageGroup.Names,
		Logger:           options.Logger,
		Env:              []string{fmt.Sprintf("COMPOSE_PROJECT_NAME=%s", options.DockerComposeProjectName)},
	})
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Failed to run %s\nOutput: %s", imageGroup.ID, cmdPlus.GetOutput()))
	}
	var wg sync.WaitGroup
	var onlineTextRegex *regexp.Regexp
	for role, onlineText := range imageGroup.OnlineTexts {
		wg.Add(1)
		onlineTextRegex, err = regexp.Compile(fmt.Sprintf("%s.*%s", role, onlineText))
		if err != nil {
			return err
		}
		go func(role string, onlineTextRegex *regexp.Regexp) {
			err := cmdPlus.WaitForRegexp(onlineTextRegex, time.Hour) // No user will actually wait this long
			if err != nil {
				options.Logger.Logf("'%s' failed to come online", role)
			} else if role != "" {
				options.Logger.Logf("'%s' is running", role)
			}
			wg.Done()
		}(role, onlineTextRegex)
	}
	wg.Wait()
	options.Logger.Logf("all %s online", imageGroup.ID)
	return nil
}
