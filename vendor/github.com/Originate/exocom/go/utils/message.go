package utils

import (
	"encoding/json"
	"net"

	"github.com/Originate/exocom/go/structs"
	"github.com/gorilla/websocket"
)

// ListenForMessages continuously reads from the given websocket calling the
// given function on each received message
func ListenForMessages(socket *websocket.Conn, onMessage func(structs.Message)) error {
	for {
		_, bytes, err := socket.ReadMessage()
		if err != nil {
			if opErr, ok := err.(*net.OpError); ok {
				if opErr.Err.Error() == "use of closed network connection" {
					return nil
				}
			}
			if websocket.IsCloseError(err, websocket.CloseAbnormalClosure) {
				return nil
			}
			return err
		}
		message, err := messageFromJSON(bytes)
		if err != nil {
			return err
		}
		onMessage(message)
	}
}

// Helpers

func messageFromJSON(bytes []byte) (structs.Message, error) {
	var unmarshaled structs.Message
	err := json.Unmarshal(bytes, &unmarshaled)
	if err != nil {
		return structs.Message{}, err
	}

	return unmarshaled, nil
}
