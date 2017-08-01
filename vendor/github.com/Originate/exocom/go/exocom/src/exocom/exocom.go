package exocom

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/Originate/exocom/go/exocom/src/client_registry"
	"github.com/Originate/exocom/go/exocom/src/logger"
	"github.com/Originate/exocom/go/structs"
	"github.com/Originate/exocom/go/utils"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

// ExoCom is the top level message broadcaster
type ExoCom struct {
	server         http.Server
	clientRegistry *clientRegistry.ClientRegistry
	logger         *logger.Logger
}

var upgrader = websocket.Upgrader{}

// New creates a new ExoCom instance
func New(serviceRoutes string) (*ExoCom, error) {
	result := new(ExoCom)
	var err error
	result.clientRegistry, err = clientRegistry.NewClientRegistry(serviceRoutes)
	result.logger = logger.NewLogger(os.Stdout)
	if err != nil {
		return result, err
	}
	var handler http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/services" {
			conn, err := upgrader.Upgrade(w, r, nil)
			if err != nil {
				fmt.Println("Error upgrading request to websocket:", err)
				return
			}
			go result.websocketHandler(conn)
		} else if r.URL.Path == "/config.json" {
			err := json.NewEncoder(w).Encode(map[string]interface{}{
				"clients": result.clientRegistry.Clients,
				"routes":  result.clientRegistry.Routing,
			})
			if err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		} else {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		}
	}
	result.server = http.Server{Handler: handler}
	return result, nil
}

// Close closes the server
func (e *ExoCom) Close() error {
	return e.server.Close()
}

// Listen opens a server on the given port
func (e *ExoCom) Listen(port int) error {
	e.server.Addr = fmt.Sprintf(":%d", port)
	fmt.Printf("ExoCom online at port %d\n", port)
	return e.server.ListenAndServe()
}

// Helpers

func (e *ExoCom) websocketHandler(socket *websocket.Conn) {
	var clientName string
	utils.ListenForMessages(socket, func(message structs.Message) error {
		if message.Name == "exocom.register-service" {
			var err error
			clientName, err = parseRegisterMessagePayload(message)
			if err == nil {
				e.clientRegistry.RegisterClient(clientName)
				printLogError(e.logger.Log(fmt.Sprintf("'%s' registered", clientName)))
			} else {
				printLogError(e.logger.Error(err.Error()))
			}
		}
		return nil
	}, func(err error) {
		fmt.Println(errors.Wrap(err, "Exocom listening for messages"))
	})
	if clientName != "" {
		e.clientRegistry.DeregisterClient(clientName)
	}
}

func printLogError(err error) {
	if err != nil {
		fmt.Println("Error logging to stdout", err)
	}
}

func parseRegisterMessagePayload(message structs.Message) (string, error) {
	if objectPayload, ok := message.Payload.(map[string]interface{}); ok {
		if clientName, ok := objectPayload["clientName"].(string); ok {
			return clientName, nil
		}
	}
	return "", fmt.Errorf("Invalid register message payload: %v", message.Payload)
}
