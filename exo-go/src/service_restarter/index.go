package serviceRestarter

import (
	"fmt"

	"github.com/Originate/exosphere/exo-go/src/docker_compose"
	"github.com/fsnotify/fsnotify"
	"github.com/pkg/errors"
)

// ServiceRestarter watches the given local service for changes and restarts it
type ServiceRestarter struct {
	ServiceName      string
	ServiceDir       string
	DockerComposeDir string
	Env              []string
	Log              func(string)
	watcher          *fsnotify.Watcher
}

// Watch watches the service directory for changes
func (s *ServiceRestarter) Watch(watcherErrChannel chan<- error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		watcherErrChannel <- err
	}
	s.watcher = watcher
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				if isCreate(event) {
					s.Log(fmt.Sprintf("Restarting service '%s' because %s was created", s.ServiceName, event.Name))
				} else if isRemove(event) {
					s.Log(fmt.Sprintf("Restarting service '%s' because %s was deleted", s.ServiceName, event.Name))
				} else if isChange(event) {
					s.Log(fmt.Sprintf("Restarting service '%s' because %s was changed", s.ServiceName, event.Name))
				}
				if err := s.restart(watcherErrChannel); err != nil {
					watcherErrChannel <- err
				}
			case err := <-watcher.Errors:
				watcherErrChannel <- err
			}
		}
	}()
	if err := watcher.Add(s.ServiceDir); err != nil {
		watcherErrChannel <- err
	}
}

func isChange(event fsnotify.Event) bool {
	return (event.Op&fsnotify.Write == fsnotify.Write) || (event.Op&fsnotify.Rename == fsnotify.Rename) || (event.Op&fsnotify.Chmod == fsnotify.Chmod)
}

func isCreate(event fsnotify.Event) bool {
	return event.Op&fsnotify.Create == fsnotify.Create
}

func isRemove(event fsnotify.Event) bool {
	return event.Op&fsnotify.Remove == fsnotify.Remove
}

func (s *ServiceRestarter) restart(watcherErrChannel chan<- error) error {
	if err := s.watcher.Close(); err != nil {
		return err
	}
	if err := dockerCompose.KillContainer(s.ServiceName, s.DockerComposeDir, s.Log); err != nil {
		return errors.Wrap(err, fmt.Sprintf("Docker failed to kill container %s", s.ServiceName))
	}
	s.Log("Docker container stopped")
	if err := dockerCompose.CreateNewContainer(s.ServiceName, s.Env, s.DockerComposeDir, s.Log); err != nil {
		return errors.Wrap(err, fmt.Sprintf("Docker image failed to rebuild %s", s.ServiceName))
	}
	s.Log("Docker image rebuilt")
	if err := dockerCompose.RestartContainer(s.ServiceName, s.Env, s.DockerComposeDir, s.Log); err != nil {
		return errors.Wrap(err, fmt.Sprintf("Docker container failed to restart %s", s.ServiceName))
	}
	s.Log(fmt.Sprintf("'%s' restarted successfully", s.ServiceName))
	s.Watch(watcherErrChannel)
	return nil
}
