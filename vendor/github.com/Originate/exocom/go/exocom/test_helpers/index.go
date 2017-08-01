package testHelpers

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/Originate/exocom/go/utils"
)

// WaitForText reads from the stream waiting for the given text
func WaitForText(stream io.Reader, text string) error {
	var output string
	var err error
	return utils.WaitForf(func() bool {
		buffer := make([]byte, 1000)
		var count int
		count, err = stream.Read(buffer)
		if count == 0 {
			if err == io.EOF {
				return false
			}
			if err != nil {
				return false
			}
		}
		output = output + string(buffer)
		return strings.Contains(output, text)
	}, func() error {
		if err != nil && err != io.EOF {
			return err
		}
		return fmt.Errorf("Expected '%s' to include '%s'", output, text)
	})
}

// StartMockServer starts a mock server on the given port
func StartMockServer(port int) *http.Server {
	server := http.Server{Addr: fmt.Sprintf(":%d", port)}
	go func() {
		err := server.ListenAndServe()
		if err != nil && err.Error() != "Server Closed" {
			fmt.Println("Mock server error", err)
		}
	}()
	return &server
}
