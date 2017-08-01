package testHelpers

import (
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"

	"github.com/Originate/exocom/go/structs"
	"github.com/Originate/exocom/go/utils"
)

// MockService is a mock of a real service that would connect to exocom
type MockService struct {
	exocomPort int
	clientName string
	socket     *websocket.Conn
}

// NewMockService returns a new MockService
func NewMockService(exocomPort int, clientName string) *MockService {
	return &MockService{
		clientName: clientName,
		exocomPort: exocomPort,
	}
}

// Connect connects the mock service to exocom
func (m *MockService) Connect() error {
	url := fmt.Sprintf("ws://%s:%d/services", "localhost", m.exocomPort)
	var err error
	m.socket, err = utils.ConnectWithRetry(url, 100)
	if err != nil {
		return err
	}
	serializedBytes, err := json.Marshal(structs.Message{
		Name:    "exocom.register-service",
		Payload: map[string]string{"clientName": m.clientName},
	})
	if err != nil {
		return err
	}
	return m.socket.WriteMessage(websocket.TextMessage, serializedBytes)
}

// Close disconnects the mock service from exocom
func (m *MockService) Close() error {
	if m.socket != nil {
		err := m.socket.Close()
		m.socket = nil
		return err
	}
	return nil
}
