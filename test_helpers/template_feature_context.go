package testHelpers

import (
	"fmt"
	"io/ioutil"
	"path"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/Originate/exosphere/src/util"
	"github.com/pkg/errors"
)

// TemplateFeatureContext defines the festure context for features/template/**/*.feature
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
		doesExist, err := util.DoesDirectoryExist(dirPath)
		if err != nil {
			return err
		}
		if !doesExist {
			return fmt.Errorf("%s does not exist", directory)
		}
		fileInfos, err := ioutil.ReadDir(dirPath)
		if err != nil {
			return err
		}
		if len(fileInfos) == 0 {
			return fmt.Errorf("%s is empty", directory)
		}
		return nil
	})

	s.Step(`^my git repository has a submodule "([^"]*)" with remote "([^"]*)"$`, func(submodulePath, gitURL string) error {
		pathOutput, err := util.Run(appDir, "git", "config", "--file", ".gitmodules", "--get-regexp", "path")
		if err != nil {
			return err
		}
		if err := validateTextContains(pathOutput, submodulePath); err != nil {
			return err
		}
		urlOutput, err := util.Run(appDir, "git", "config", "--file", ".gitmodules", "--get-regexp", "url")
		if err != nil {
			return err
		}
		return validateTextContains(urlOutput, gitURL)
	})

	s.Step(`^my git repository has a submodule "([^"]*)" at commit "([^"]*)"$`, func(submodulePath, commitSha string) error {
		fullSubmodulePath := path.Join(appDir, submodulePath)
		output, err := util.Run(fullSubmodulePath, "git", "rev-parse", "HEAD")
		if err != nil {
			return err
		}
		return validateTextContains(output, commitSha)
	})

	s.Step(`^my git repository does not have any submodules$`, func() error {
		hasNoGitModules, err := util.IsEmptyFile(path.Join(appDir, ".gitmodules"))
		if err != nil {
			return err
		}
		if !hasNoGitModules {
			return fmt.Errorf("Expected the git reposity to not have any submodules")
		}
		return nil
	})

}
