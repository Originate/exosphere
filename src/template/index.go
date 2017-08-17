package template

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"time"

	"github.com/Originate/exosphere/src/util"
	execplus "github.com/Originate/go-execplus"
	"github.com/tmrts/boilr/pkg/template"
)

const templatesDir = ".exosphere"

// Add fetches a remote template from GitHub and stores it
// under templateDir, returns an error if any
func Add(gitURL, templateName, templateDir, commitIsh string) error {
	if _, err := util.Run("", "git", "submodule", "add", "--name", templateName, gitURL, templateDir); err != nil {
		return err
	}
	if commitIsh != "" {
		_, err := util.Run(templateDir, "git", "checkout", commitIsh)
		return err
	}
	return nil
}

// CreateTmpServiceDir makes bolir scaffold the template chosenTemplate
// and store the scaffoled service folder in a tmp folder, and finally
// returns the path to the tmp folder
func CreateTmpServiceDir(appDir, chosenTemplate string) (string, error) {
	templateDir := path.Join(appDir, templatesDir, chosenTemplate)
	template, err := template.Get(templateDir)
	if err != nil {
		return "", err
	}
	serviceTmpDir, err := ioutil.TempDir("", "service-tmp")
	if err != nil {
		return "", err
	}
	if err = template.Execute(serviceTmpDir); err != nil {
		return "", err
	}
	return serviceTmpDir, nil
}

// Fetch fetches updates for all existing remote templates
func Fetch(templateDir string) error {
	_, err := util.Run(templateDir, "git", "pull")
	return err
}

// GetTemplates returns a slice of all template names found in the ".exosphere"
// folder of the application
func GetTemplates(appDir string) (result []string, err error) {
	subdirectories, err := util.GetSubdirectories(path.Join(appDir, templatesDir))
	if err != nil {
		return result, err
	}
	for _, directory := range subdirectories {
		if isValidDir(path.Join(templatesDir, directory)) {
			result = append(result, directory)
		}
	}
	return result, nil
}

// HasTemplatesDir returns whether or not there is an ".exosphere" folder
func HasTemplatesDir(appDir string) bool {
	return util.DoesDirectoryExist(path.Join(appDir, templatesDir)) && !util.IsEmptyDirectory(templatesDir)
}

// IsValidTemplateDir returns whether or not the given template directory
// is a valid template
func IsValidTemplateDir(templateDir string) bool {
	return util.DoesFileExist(path.Join(templateDir, "project.json")) && util.DoesDirectoryExist(path.Join(templateDir, "template")) && util.DoesFileExist(path.Join(templateDir, "template", "{{serviceRole}}", "tests", "Dockerfle"))
}

// Run executes the boilr template from templateDir into resultDir
func Run(templateDir, resultDir string) error {
	t, err := template.Get(templateDir)
	if err != nil {
		return err
	}
	if err = t.Execute(resultDir); err != nil {
		return err
	}
	if err = os.RemoveAll(templateDir); err != nil {
		return err
	}
	return nil
}

// Remove removes the given template from the application
func Remove(templateName, templateDir string) error {
	denitCommand := []string{"git", "submodule", "deinit", "-f", templateDir}
	removeModulesCommand := []string{"rm", "-rf", fmt.Sprintf(".git/modules/%s", templateName)}
	gitRemoveCommand := []string{"git", "rm", "-f", templateDir}
	return util.RunSeries("", [][]string{denitCommand, removeModulesCommand, gitRemoveCommand})
}

// Helpers

func createProjectJSON(templateDir string, content string) error {
	return ioutil.WriteFile(path.Join(templateDir, "project.json"), []byte(content), 0777)
}

func isValidDir(templateDir string) bool {
	return util.DoesFileExist(path.Join(templateDir, "project.json")) && util.DoesDirectoryExist(path.Join(templateDir, "template"))
}

func enterEmptyInputs(cmd execplus.CmdPlus, fields string) error {
	for _, field := range fileds {
		if err := childCmdPlus.WaitForText(field, time.Second*5); err != nil {
			return err
		}
		if _, err := childCmdPlus.StdinPipe.Write([]byte("\n" + "\n")); err != nil {
			return err
		}
	}
	return nil
}

func CreateEmptyApp(appDir string) error {
	cmdPlus = execplus.NewCmdPlus("exo create")
	cmdPlus.SetDir(appDir)
	if err := cmdPlus.Start(); err != nil {
		return err
	}
	fields = []string{"AppName", "AppDescription", "AppVersion", "ExocomVersion"}
	return enterEmptyInputs(cmdPlus, fields)
}

func AddService(appDir string) error {
	cmdPlus = execplus.NewCmdPlus("exo create")
	cmdPlus.SetDir(appDir)
	if err := cmdPlus.Start(); err != nil {
		return err
	}
	projectJSON, err := ioutil.ReadFile(path.Join(appDir, "project.json"))
	if err != nil {
		return err
	}
	var defaults interface{}
	if err := json.Unmarshal(projectJSON, &defaults); err != nil {
		return err
	}
	fields = []string{} // get fields from project.json
	return enterEmptyInputs(cmdPlus, fields)
}
