package utils

import "github.com/gorilla/websocket"

// ConnectWithRetry will attempt to connect a websocket to the given url
// retrying if needed for the given tiemout
func ConnectWithRetry(url string, timeout int) (socket *websocket.Conn, err error) {
	Retry(timeout, func() bool {
		socket, _, err = websocket.DefaultDialer.Dial(url, nil)
		return err == nil
	})
	return socket, err
}
