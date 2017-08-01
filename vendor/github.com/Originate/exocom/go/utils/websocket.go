package utils

import (
	"fmt"

	"github.com/gorilla/websocket"
)

// ConnectWithRetry will attempt to connect a websocket to the given url
// retrying if needed with the given delay (milliseconds) between calls
func ConnectWithRetry(url string, delay int) (socket *websocket.Conn, err error) {
	Retry(delay, func() bool {
		socket, _, err = websocket.DefaultDialer.Dial(url, nil)
		if err != nil {
			fmt.Printf("Unable to connect (%s). Retrying...\n", err)
			return false
		}
		return true
	})
	return socket, err
}
