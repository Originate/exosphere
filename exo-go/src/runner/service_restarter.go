package runner

import (
	"fmt"

	"github.com/Originate/exosphere/exo-go/src/dockercompose"
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
				s.LogChannel <- fmt.Sprintf("Restarting service '%s' because %s was %s", s.ServiceName, event.Name, getEventName(event))
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
	if err := dockercompose.KillContainer(s.ServiceName, s.DockerComposeDir, s.LogChannel); err != nil {
		watcherErrChannel <- errors.Wrap(err, fmt.Sprintf("Docker failed to kill container %s", s.ServiceName))
		return
	}
	s.LogChannel <- "Docker container stopped"
	if err := dockercompose.CreateNewContainer(s.ServiceName, s.Env, s.DockerComposeDir, s.LogChannel); err != nil {
		watcherErrChannel <- errors.Wrap(err, fmt.Sprintf("Docker image failed to rebuild %s", s.ServiceName))
		return
	}
	s.LogChannel <- "Docker image rebuilt"
	if err := dockercompose.RestartContainer(s.ServiceName, s.Env, s.DockerComposeDir, s.LogChannel); err != nil {
		watcherErrChannel <- errors.Wrap(err, fmt.Sprintf("Docker container failed to restart %s", s.ServiceName))
		return
	}
	s.LogChannel <- fmt.Sprintf("'%s' restarted successfully", s.ServiceName)
	s.Watch(watcherErrChannel)
}

// Helpers

func getEventName(event fsnotify.Event) string {
	if isCreateEvent(event) {
		return "created"
	}
	if isRemoveEvent(event) {
		return "deleted"
	}
	return "changed"
}

func isCreateEvent(event fsnotify.Event) bool {
	return event.Op&fsnotify.Create == fsnotify.Create
}

func isRemoveEvent(event fsnotify.Event) bool {
	return event.Op&fsnotify.Remove == fsnotify.Remove
}
