package daemon

import (
	"errors"
	"strconv"
	"time"

	"golang.org/x/net/context"

	"github.com/Sirupsen/logrus"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/backend"
	containertypes "github.com/docker/docker/api/types/container"
	timetypes "github.com/docker/docker/api/types/time"
	"github.com/docker/docker/container"
	"github.com/docker/docker/daemon/logger"
)

// ContainerLogs copies the container's log channel to the channel provided in
// the config. If ContainerLogs returns an error, no messages have been copied.
// and the channel will be closed without data.
//
// if it returns nil, the config channel will be active and return log
// messages until it runs out or the context is canceled.
func (daemon *Daemon) ContainerLogs(ctx context.Context, containerName string, config *types.ContainerLogsOptions) (<-chan *backend.LogMessage, error) {
	lg := logrus.WithFields(logrus.Fields{
		"module":    "daemon",
		"method":    "(*Daemon).ContainerLogs",
		"container": containerName,
	})

	if !(config.ShowStdout || config.ShowStderr) {
		return nil, errors.New("You must choose at least one stream")
	}
	container, err := daemon.GetContainer(containerName)
	if err != nil {
		return nil, err
	}

	if container.RemovalInProgress || container.Dead {
		return nil, errors.New("can not get logs from container which is dead or marked for removal")
	}

	if container.HostConfig.LogConfig.Type == "none" {
		return nil, logger.ErrReadLogsNotSupported
	}

	cLog, err := daemon.getLogger(container)
	if err != nil {
		return nil, err
	}

	logReader, ok := cLog.(logger.LogReader)
	if !ok {
		return nil, logger.ErrReadLogsNotSupported
	}

	follow := config.Follow && container.IsRunning()
	tailLines, err := strconv.Atoi(config.Tail)
	if err != nil {
		tailLines = -1
	}

	var since time.Time
	if config.Since != "" {
		s, n, err := timetypes.ParseTimestamps(config.Since, 0)
		if err != nil {
			return nil, err
		}
		since = time.Unix(s, n)
	}

	readConfig := logger.ReadConfig{
		Since:  since,
		Tail:   tailLines,
		Follow: follow,
	}

	logs := logReader.ReadLogs(readConfig)

	// past this point, we can't possibly return any errors, so we can just
	// start a goroutine and return to tell the caller not to expect errors
	// (if the caller wants to give up on logs, they have to cancel the context)
	// this goroutine functions as a shim between the logger and the caller.
	messageChan := make(chan *backend.LogMessage, 1)
	go func() {
		// set up some defers
		defer func() {
			// ok so this function, originally, was placed right after that
			// logger.ReadLogs call above. I THINK that means it sets off the
			// chain of events that results in the logger needing to be closed.
			// i do not know if an error in time parsing above causing an early
			// return will result in leaking the logger. if that is the case,
			// it would also have been a bug in the original code
			logs.Close()
			if cLog != container.LogDriver {
				// Since the logger isn't cached in the container, which
				// occurs if it is running, it must get explicitly closed
				// here to avoid leaking it and any file handles it has.
				if err := cLog.Close(); err != nil {
					logrus.Errorf("Error closing logger: %v", err)
				}
			}
		}()
		// close the messages channel. closing is the only way to signal above
		// that we're doing with logs (other than context cancel i guess).
		defer close(messageChan)

		lg.Debug("begin logs")
		for {
			select {
			// i do not believe as the system is currently designed any error
			// is possible, but we should be prepared to handle it anyway. if
			// we do get an error, copy only the error field to a new object so
			// we don't end up with partial data in the other fields
			case err := <-logs.Err:
				lg.Errorf("Error streaming logs: %v", err)
				select {
				case <-ctx.Done():
				case messageChan <- &backend.LogMessage{Err: err}:
				}
				return
			case <-ctx.Done():
				lg.Debug("logs: end stream, ctx is done: %v", ctx.Err())
				return
			case msg, ok := <-logs.Msg:
				// there is some kind of pool or ring buffer in the logger that
				// produces these messages, and a possible future optimization
				// might be to use that pool and reuse message objects
				if !ok {
					lg.Debug("end logs")
					return
				}
				m := msg.AsLogMessage() // just a pointer conversion, does not copy data

				// there could be a case where the reader stops accepting
				// messages and the context is canceled. we need to check that
				// here, or otherwise we risk blocking forever on the message
				// send.
				select {
				case <-ctx.Done():
					return
				case messageChan <- m:
				}
			}
		}
	}()
	return messageChan, nil
}

func (daemon *Daemon) getLogger(container *container.Container) (logger.Logger, error) {
	if container.LogDriver != nil && container.IsRunning() {
		return container.LogDriver, nil
	}
	return container.StartLogger()
}

// mergeLogConfig merges the daemon log config to the container's log config if the container's log driver is not specified.
func (daemon *Daemon) mergeAndVerifyLogConfig(cfg *containertypes.LogConfig) error {
	if cfg.Type == "" {
		cfg.Type = daemon.defaultLogConfig.Type
	}

	if cfg.Config == nil {
		cfg.Config = make(map[string]string)
	}

	if cfg.Type == daemon.defaultLogConfig.Type {
		for k, v := range daemon.defaultLogConfig.Config {
			if _, ok := cfg.Config[k]; !ok {
				cfg.Config[k] = v
			}
		}
	}

	return logger.ValidateLogOpts(cfg.Type, cfg.Config)
}
