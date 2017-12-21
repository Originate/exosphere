package helpers

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"regexp"
	"time"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/jaytaylor/html2text"
)

// TutorialFeatureContext defines the festure context for the tutorial
// nolint gocyclo
func TutorialFeatureContext(s *godog.Suite) {

	s.Step(`^I am in an empty folder$`, func() error {
		return nil
	})

	s.Step(`^I cd into "([^"]*)"$`, func(dir string) error {
		appDir = path.Join(appDir, dir)
		return nil
	})

	s.Step(`^I add the "([^"]*)" template$`, func(templateName string) error {
		return checkoutTemplate(path.Join(appDir, ".exosphere", "service_templates", templateName), templateName)
	})

	s.Step(`^waiting until I see "([^"]*)" in the terminal$`, func(expectedText string) error {
		return childCmdPlus.WaitForText(expectedText, time.Minute)
	})

	s.Step(`^I stop all running processes$`, func() error {
		if childCmdPlus != nil {
			return childCmdPlus.Kill()
		}
		return nil
	})

	s.Step(`^the file "([^"]*)":$`, func(filePath string, expectedContent *gherkin.DocString) error {
		return ioutil.WriteFile(path.Join(appDir, filePath), []byte(expectedContent.Content), 0777)
	})

	s.Step(`^http:\/\/localhost:(\d+) displays:$`, func(port string, expectedContent *gherkin.DocString) error {
		resp, err := http.Get(fmt.Sprintf("http://localhost:%s", port))
		if err != nil {
			return err
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		formattedText, err := html2text.FromString(string(body))
		if err != nil {
			return err
		}
		text := regexp.MustCompile(`\*+\n|\-+\n`).ReplaceAllString(formattedText, "")
		return validateTextContains(text, expectedContent.Content)
	})

	s.Step(`^adding a todo entry called "([^"]*)" via the web application$`, func(entry string) error {
		_, err := http.PostForm("http://localhost:3000/todos", url.Values{"text": {entry}})
		return err
	})

}
