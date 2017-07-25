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
		return
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
				s.restart(watcherErrChannel)
				return
			case err := <-watcher.Errors:
				watcherErrChannel <- err
				return
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

func (s *ServiceRestarter) restart(watcherErrChannel chan<- error) {
	if err := s.watcher.Close(); err != nil {
		watcherErrChannel <- err
		return
	}
	if err := dockerCompose.KillContainer(s.ServiceName, s.DockerComposeDir, s.Log); err != nil {
		watcherErrChannel <- errors.Wrap(err, fmt.Sprintf("Docker failed to kill container %s", s.ServiceName))
		return
	}
	s.Log("Docker container stopped")
	if err := dockerCompose.CreateNewContainer(s.ServiceName, s.Env, s.DockerComposeDir, s.Log); err != nil {
		watcherErrChannel <- errors.Wrap(err, fmt.Sprintf("Docker image failed to rebuild %s", s.ServiceName))
		return
	}
	s.Log("Docker image rebuilt")
	if err := dockerCompose.RestartContainer(s.ServiceName, s.Env, s.DockerComposeDir, s.Log); err != nil {
		watcherErrChannel <- errors.Wrap(err, fmt.Sprintf("Docker container failed to restart %s", s.ServiceName))
		return
	}
	s.Log(fmt.Sprintf("'%s' restarted successfully", s.ServiceName))
	s.Watch(watcherErrChannel)
}
