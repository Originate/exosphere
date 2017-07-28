package application

import (
	"fmt"

	"github.com/Originate/exosphere/exo-go/src/docker_compose"
	"github.com/fsnotify/fsnotify"
	"github.com/pkg/errors"
)

// ServiceRestarter watches the given local service for changes and restarts it
type serviceRestarter struct {
	ServiceName      string
	ServiceDir       string
	DockerComposeDir string
	Env              []string
	LogChannel       chan string
	watcher          *fsnotify.Watcher
}

// Watch watches the service directory for changes
func (s *serviceRestarter) Watch(watcherErrChannel chan<- error) {
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
				if isCreateEvent(event) {
					s.LogChannel <- fmt.Sprintf("Restarting service '%s' because %s was created", s.ServiceName, event.Name)
				} else if isRemoveEvent(event) {
					s.LogChannel <- fmt.Sprintf("Restarting service '%s' because %s was deleted", s.ServiceName, event.Name)
				} else {
					s.LogChannel <- fmt.Sprintf("Restarting service '%s' because %s was changed", s.ServiceName, event.Name)
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

func (s *serviceRestarter) restart(watcherErrChannel chan<- error) {
	if err := s.watcher.Close(); err != nil {
		watcherErrChannel <- err
		return
	}
	if err := dockerCompose.KillContainer(s.ServiceName, s.DockerComposeDir, s.LogChannel); err != nil {
		watcherErrChannel <- errors.Wrap(err, fmt.Sprintf("Docker failed to kill container %s", s.ServiceName))
		return
	}
	s.LogChannel <- "Docker container stopped"
	if err := dockerCompose.CreateNewContainer(s.ServiceName, s.Env, s.DockerComposeDir, s.LogChannel); err != nil {
		watcherErrChannel <- errors.Wrap(err, fmt.Sprintf("Docker image failed to rebuild %s", s.ServiceName))
		return
	}
	s.LogChannel <- "Docker image rebuilt"
	if err := dockerCompose.RestartContainer(s.ServiceName, s.Env, s.DockerComposeDir, s.LogChannel); err != nil {
		watcherErrChannel <- errors.Wrap(err, fmt.Sprintf("Docker container failed to restart %s", s.ServiceName))
		return
	}
	s.LogChannel <- fmt.Sprintf("'%s' restarted successfully", s.ServiceName)
	s.Watch(watcherErrChannel)
}

// Helpers

func isCreateEvent(event fsnotify.Event) bool {
	return event.Op&fsnotify.Create == fsnotify.Create
}

func isRemoveEvent(event fsnotify.Event) bool {
	return event.Op&fsnotify.Remove == fsnotify.Remove
}
