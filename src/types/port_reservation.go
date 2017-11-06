package types

import "strconv"

// PortReservation reserves a free port for services
type PortReservation struct {
	port int
}

// NewPortReservation instantiates a port reservation struct
func NewPortReservation() *PortReservation {
	return &PortReservation{port: 3000}
}

// GetAvailablePort returns an available port as a string
func (p *PortReservation) GetAvailablePort() string {
	availablePort := p.port
	p.port = availablePort + 10
	return strconv.Itoa(availablePort)

}
