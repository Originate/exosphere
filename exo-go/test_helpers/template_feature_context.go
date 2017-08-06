package testHelpers

import (
	"fmt"
	"path"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/Originate/exosphere/exo-go/src/osplus"
	"github.com/Originate/exosphere/exo-go/src/runplus"
	"github.com/pkg/errors"
)

// FetchTemplatesFeatureContext defines the festure context for features/fetch_templates.feature
// nolint gocyclo
func TemplateFeatureContext(s *godog.Suite) {
	s.Step(`^my application is a Git repository$`, func() error {
		if _, err := runplus.Run(appDir, "git", "init"); err != nil {
			return errors.Wrap(err, fmt.Sprintf("Failed to creates .git for %s:%s\n", appDir, err))
		}
		return nil
	})

	s.Step(`^my application has the templates:$`, func(table *gherkin.DataTable) error {
		for _, row := range table.Rows[1:] {
			templateName, gitURL := row.Cells[0].Value, row.Cells[1].Value
			if _, err := runplus.Run(appDir, "exo", "template", "add", templateName, gitURL); err != nil {
				return errors.Wrap(err, fmt.Sprintf("Failed to creates the template %s:%s\n", appDir, err))
			}
		}
		return nil
	})

	s.Step(`^my application contains the directory "([^"]*)"`, func(directory string) error {
		dirPath := path.Join(appDir, directory)
		if !osplus.DoesDirectoryExist(dirPath) {
			return fmt.Errorf("%s does not exist", directory)
		}
		if osplus.IsEmptyDirectory(directory) {
			return fmt.Errorf("%s is empty", directory)
		}
		return nil
	})

	s.Step(`^my git repository has a submodule "([^"]*)" with remote "([^"]*)"$`, func(submodulePath, gitURL string) error {
		if childOutput, err := runplus.Run(appDir, "git", "config", "--file", ".gitmodules", "--get-regexp", "path"); err != nil {
			return errors.Wrap(err, fmt.Sprintf("Failed to read git config file: %s", childOutput))
		} else {
			return validateTextContains(childOutput, submodulePath)
		}
		if childOutput, err := runplus.Run(appDir, "git", "config", "--file", ".gitmodules", "--get-regexp", "url"); err != nil {
			return errors.Wrap(err, "Failed to read git config file")
		} else {
			return validateTextContains(childOutput, gitURL)
		}
	})

	s.Step(`^my git repository does not have any submodules$`, func() error {
		if !osplus.IsEmptyDirectory(path.Join(appDir, ".exosphere")) || !osplus.IsEmptyFile(path.Join(appDir, ".gitmodules")) || !osplus.IsEmptyDirectory(path.Join(appDir, ".git", "modules")) {
			return fmt.Errorf("Expected the git reposity to not have any submodules")
		}
		return nil
	})

}
