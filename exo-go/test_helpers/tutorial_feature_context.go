package testHelpers

import (
	"io/ioutil"
  "net/http"
	"os"
	"path"
	"strings"
  "url"

	"github.com/DATA-DOG/godog"
)

// TutorialFeatureContext defines the festure context for the tutorial
// nolint gocyclo
func TutorialFeatureContext(s *godog.Suite) {

	s.Step(`^I see "([^"]*)" in the terminal$`, func(expectedText string) {
		process.WaitForTextWithTimeout(expectedText, 100
	})

	s.Step(`^I am in an empty folder$`, func() {
		appDir = path.Join(os.TempDir())
		return nil
	})

	s.Step(`^I cd into "([^"]*)"$`, func(dir string) {
		appName = dir
		appDir = dir
		return nil
	})

	s.Step(`^my application contains the template folder "([^"]*)"$`, func(templateDir string) {
		templateName = strings.Split(templateDir, "/")[1]
		checkoutServiceTemplate(appDir, templateName)
	})

	s.Step(`^waiting until I see "([^"]*)" in the terminal$`, func(expectedText string) error {
		return process.WaitForTextWithTimeout(expectedText, 5000)
	})

	s.Step(`^I stop all running processes$`, func() error {
    return process.Kill()
	})

	s.Step(`^the file "([^"]*)":$`, func(filePath, content string) error {
		return ioutil.WriteFile(path.Join(appDir, filePath), []byte(content), 0777)
	})

	s.Step(`^http:\/\/localhost:(\d+) displays:$`, func(content string) error {
    resp, err := http.Get("http://example.com/")
    responseData,err := ioutil.ReadAll(response.Body)
    return validateTextContains(string(responseData, content)
	})

	s.Step(`^adding a todo entry called "([^"]*)" via the web application$`, func(entry) error {
    url := "localhost:3000/todos"
    req, err := http.PostForm(url, url.Values{"text":entry})
	})

}
