package logger_test

import (
	"io"
	"io/ioutil"

	"github.com/Originate/exocom/go/exocom/src/client_registry"
	"github.com/Originate/exocom/go/exocom/src/logger"
	"github.com/Originate/exocom/go/structs"
	"github.com/fatih/color"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func GetOutput(f func(*logger.Logger)) string {
	pipeReader, pipeWriter := io.Pipe()
	testLogger := logger.NewLogger(pipeWriter)
	go func() {
		f(testLogger)
		pipeWriter.Close()
	}()
	result, err := ioutil.ReadAll(pipeReader)
	Expect(err).To(BeNil())
	return string(result)
}

var _ = Describe("Logger", func() {
	Describe("Log", func() {
		It("prints the text to the screen with a new line character appended", func() {
			result := GetOutput(func(testLogger *logger.Logger) {
				err := testLogger.Log("foo bar")
				Expect(err).To(BeNil())
			})
			Expect(result).To(Equal("foo bar\n"))
		})
	})
	Describe("Write", func() {
		It("prints the text to the screen without a new line character appended", func() {
			result := GetOutput(func(testLogger *logger.Logger) {
				err := testLogger.Write("foo bar")
				Expect(err).To(BeNil())
			})
			Expect(result).To(Equal("foo bar"))
		})
	})
	Describe("Error", func() {
		It("prints an error", func() {
			result := GetOutput(func(testLogger *logger.Logger) {
				err := testLogger.Error("error")
				Expect(err).To(BeNil())
			})
			Expect(result).To(Equal(color.RedString("error\n")))
		})
	})
	Describe("Header", func() {
		It("prints a header", func() {
			result := GetOutput(func(testLogger *logger.Logger) {
				err := testLogger.Header("header")
				Expect(err).To(BeNil())
			})
			expected := color.New(color.Faint).Sprint("header\n")
			Expect(result).To(Equal(expected))
		})
	})
	Describe("Warning", func() {
		It("prints a warning", func() {
			result := GetOutput(func(testLogger *logger.Logger) {
				err := testLogger.Warning("warning")
				Expect(err).To(BeNil())
			})
			Expect(result).To(Equal(color.YellowString("warning\n")))
		})
	})

	Describe("Routing Setup", func() {
		It("prints the routing setup received", func() {
			usersRoute := clientRegistry.Route{}
			usersRoute.Receives = []string{"users.create"}
			routes := clientRegistry.Routes{
				"users": usersRoute,
			}
			result := GetOutput(func(testLogger *logger.Logger) {
				err := testLogger.RoutingSetup(routes)
				Expect(err).To(BeNil())
			})
			Expect(result).To(Equal("Receiving Routing Setup\n  --[ users ]-> users.create\n"))
		})
	})
	Describe("Messages", func() {
		It("prints the message that is sent, the response time, and the payload", func() {
			message := structs.Message{
				Name:         "users.created",
				Sender:       "users",
				Payload:      map[string]interface{}{},
				ResponseTo:   "users.create",
				ResponseTime: 1000}
			receivers := []string{"users"}
			messages := []structs.Message{message}
			result := GetOutput(func(testLogger *logger.Logger) {
				err := testLogger.Messages(messages, receivers)
				Expect(err).To(BeNil())
			})
			Expect(result).To(Equal("users  --[ users.created ]->  users  ( 1µs )\n{}\n"))
		})
	})
})
