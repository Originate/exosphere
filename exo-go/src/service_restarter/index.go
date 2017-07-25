package serviceRestarter

import (
	"fmt"

	"github.com/Originate/exosphere/exo-go/src/docker_compose"
	"github.com/fsnotify/fsnotify"
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
func (s *ServiceRestarter) Watch(watcherErr chan<- error) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	s.watcher = watcher
	go func() {
		select {
		case event := <-watcher.Events:
			if isCreate(event) {
				s.Log(fmt.Sprintf("Restarting service '%s' because %s was created", s.ServiceName, s.ServiceDir))
			} else if isRemove(event) {
				s.Log(fmt.Sprintf("Restarting service '%s' because %s was deleted", s.ServiceName, s.ServiceDir))
			} else if isChange(event) {
				s.Log(fmt.Sprintf("Restarting service '%s' because %s was changed", s.ServiceName, s.ServiceDir))
			}
			watcherErr <- s.restart()
		case err := <-watcher.Errors:
			watcherErr <- err
		}
	}()
	return watcher.Add(s.ServiceDir)
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

func (s *ServiceRestarter) restart() error {
	if err := s.watcher.Close(); err != nil {
		return err
	}
	if err := dockerCompose.KillContainer(s.ServiceName, s.DockerComposeDir, s.Log); err != nil {
		s.Log(fmt.Sprintf("Docker failed to kill container %s", s.ServiceName))
		return err
	}
	s.Log("Docker container stopped")
	if err := dockerCompose.CreateNewContainer(s.ServiceName, s.Env, s.DockerComposeDir, s.Log); err != nil {
		s.Log(fmt.Sprintf("Docker image failed to rebuild %s", s.ServiceName))
		return err
	}
	s.Log("Docker image rebuilt")
	if err := dockerCompose.StartContainer(s.ServiceName, s.Env, s.DockerComposeDir, s.Log); err != nil {
		s.Log(fmt.Sprintf("Docker container failed to restart %s", s.ServiceName))
		return err
	}
	s.Log(fmt.Sprintf("'%s' restarted successfully", s.ServiceName))
	return nil
}
