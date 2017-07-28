package testHelpers

import (
	"fmt"
	"path"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/Originate/exosphere/exo-go/src/os_helpers"
	"github.com/Originate/exosphere/exo-go/src/util"
	"github.com/pkg/errors"
)

// FetchTemplatesFeatureContext defines the festure context for features/fetch_templates.feature
// nolint gocyclo
func TemplateFeatureContext(s *godog.Suite) {
	s.Step(`^my application is a Git repository$`, func() error {
		if _, err := util.Run(appDir, "git", "init"); err != nil {
			return errors.Wrap(err, fmt.Sprintf("Failed to creates .git for %s:%s\n", appDir, err))
		}
		return nil
	})

	s.Step(`^my application has the templates:$`, func(table *gherkin.DataTable) error {
		for _, row := range table.Rows[1:] {
			templateName, gitURL := row.Cells[0].Value, row.Cells[1].Value
			if _, err := util.Run(appDir, "exo", "template", "add", templateName, gitURL); err != nil {
				return errors.Wrap(err, fmt.Sprintf("Failed to creates the template %s:%s\n", appDir, err))
			}
		}
		return nil
	})

	s.Step(`^my application contains the directory "([^"]*)"`, func(directory string) error {
		dirPath := path.Join(appDir, directory)
		if !osHelpers.DirectoryExists(dirPath) {
			return fmt.Errorf("%s does not exist", directory)
		}
		if osHelpers.IsEmpty(directory) {
			return fmt.Errorf("%s is empty", directory)
		}
		return nil
	})

	s.Step(`^my git repository has a submodule "([^"]*)" with remote "([^"]*)"$`, func(submodulePath, gitURL string) error {
		if childOutput, err := util.Run(appDir, "git", "config", "--file", ".gitmodules", "--get-regexp", "path"); err != nil {
			return errors.Wrap(err, fmt.Sprintf("Failed to read git config file: %s", childOutput))
		} else {
			return validateTextContains(childOutput, submodulePath)
		}
		if childOutput, err := util.Run(appDir, "git", "config", "--file", ".gitmodules", "--get-regexp", "url"); err != nil {
			return errors.Wrap(err, "Failed to read git config file")
		} else {
			return validateTextContains(childOutput, gitURL)
		}
	})

	s.Step(`^my git repository does not have any submodules$`, func() error {
		if !osHelpers.IsEmpty(path.Join(appDir, ".exosphere")) || !osHelpers.IsEmptyFile(path.Join(appDir, ".gitmodules")) || !osHelpers.IsEmpty(path.Join(appDir, ".git", "modules")) {
			return fmt.Errorf("Expected the git reposity to not have any submodules")
		}
		return nil
	})

}
