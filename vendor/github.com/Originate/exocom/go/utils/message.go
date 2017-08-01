package utils

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Originate/exocom/go/structs"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

// ListenForMessages continuously reads from the given websocket calling the
// given function on each received message
func ListenForMessages(socket *websocket.Conn, onMessage func(structs.Message) error, onError func(error)) {
	for {
		_, bytes, err := socket.ReadMessage()
		if err != nil {
			if isCloseError(err) {
				return
			}
			onError(errors.Wrap(err, "Error reading from socket"))
			continue
		}
		message, err := MessageFromJSON(bytes)
		if err != nil {
			onError(errors.Wrap(err, fmt.Sprintf("Error unmarshaling message: '%s'", bytes)))
		}
		err = onMessage(message)
		if err != nil {
			onError(errors.Wrap(err, "Error returned by on message handler"))
		}
	}
}

// MessageFromJSON unmarshals bytes into a message
func MessageFromJSON(bytes []byte) (structs.Message, error) {
	var unmarshaled structs.Message
	err := json.Unmarshal(bytes, &unmarshaled)
	if err != nil {
		return structs.Message{}, err
	}

	return unmarshaled, nil
}

func isCloseError(err error) bool {
	return websocket.IsCloseError(err, websocket.CloseAbnormalClosure) ||
		strings.Contains(err.Error(), "connection reset by peer") ||
		strings.Contains(err.Error(), "use of closed network connection")
}
