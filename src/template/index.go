package template

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strings"
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

// AddService (used by exo template test) runs exo-add to add template
// at templateDir to the app at appDir and returns an error if any
func AddService(appDir, templateDir string, outputWriter io.Writer) error {
	cmd := execplus.NewCmdPlus("exo", "add")
	cmd.SetDir(appDir)
	go sendCmdOutputToWriter(cmd, outputWriter)
	if err := cmd.Start(); err != nil {
		return err
	}
	numFields, err := getNumFields(templateDir)
	if err != nil {
		return err
	}
	if err := selectFirstOption(cmd, "template"); err != nil {
		return err
	}
	if err := enterEmptyInputs(cmd, numFields); err != nil {
		return err
	}
	if err := selectFirstOption(cmd, "Protection Level"); err != nil {
		return err
	}
	if err := cmd.WaitForText("done", time.Second*5); err != nil {
		return err
	}
	return cmd.Wait()
}

// CreateEmptyApp (used by exo template test) runs exo-create to create an
// empty app with the default name "my-app" at dirPath and returns an error if any
func CreateEmptyApp(dirPath string) error {
	cmd := execplus.NewCmdPlus("exo", "create")
	cmd.SetDir(dirPath)
	if err := cmd.Start(); err != nil {
		return err
	}
	fields := []string{"AppName", "AppDescription", "AppVersion", "ExocomVersion"}
	if err := enterEmptyInputs(cmd, len(fields)); err != nil {
		return err
	}
	if err := cmd.WaitForText("done", time.Second*5); err != nil {
		return err
	}
	return cmd.Wait()
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
		isValid, _, err := IsValidTemplateDir(path.Join(templatesDir, directory))
		if err != nil {
			return result, err
		}
		if isValid {
			result = append(result, directory)
		}
	}
	return result, nil
}

// HasTemplatesDir returns whether or not there is an ".exosphere" folder
func HasTemplatesDir(appDir string) (bool, error) {
	return util.DoesDirectoryExist(path.Join(appDir, templatesDir))
}

// IsValidTemplateDir returns whether or not the template at templateDir
// is a valid exosphere template, a reason if invalid, and an
func IsValidTemplateDir(templateDir string) (bool, string, error) {
	doesExist, err := util.DoesFileExist(path.Join(templateDir, "project.json"))
	if err != nil {
		return false, "", err
	}
	if !doesExist {
		return false, "Missing file: 'project.json'", nil
	}
	doesExist, err = util.DoesDirectoryExist(path.Join(templateDir, "template"))
	if err != nil {
		return false, "", err
	}
	if !doesExist {
		return false, "Missing directory: 'template'", nil
	}
	files, err := ioutil.ReadDir(path.Join(templateDir, "template"))
	if err != nil {
		return false, "", err
	}
	if len(files) != 1 {
		return false, "'template' directory must contain a single subdirectory", nil
	}
	serviceDirName := files[0].Name()
	serviceDirPath := path.Join(templateDir, "template", serviceDirName)
	requiredFiles := []string{"service.yml", "Dockerfile.prod", "Dockerfile.dev"}
	for _, file := range requiredFiles {
		doesExist, err = util.DoesFileExist(path.Join(serviceDirPath, file))
		if err != nil {
			return false, "", err
		}
		if !doesExist {
			return false, fmt.Sprintf("Missing file: '%s'", path.Join("template", serviceDirName, file)), nil
		}
	}
	return true, "", nil
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

// RunTests (used by exo template test) runs exo-test in appDir and
// returns whether or not the tests pass and the output of running the tests
func RunTests(appDir string, outputWriter io.Writer) error {
	cmd := execplus.NewCmdPlus("exo", "test", "--no-mount")
	cmd.SetDir(appDir)
	go sendCmdOutputToWriter(cmd, outputWriter)
	return cmd.Run()
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

func enterEmptyInputs(cmd *execplus.CmdPlus, numFields int) error {
	for i := 1; i <= numFields; i++ {
		promptRegex, err := regexp.Compile(strings.Repeat(`\[\?\].*\:.*(.*\n)*`, i))
		if err != nil {
			return err
		}
		if err := cmd.WaitForRegexp(promptRegex, time.Second*5); err != nil {
			return err
		}
		if _, err := cmd.StdinPipe.Write([]byte("\n" + "\n")); err != nil {
			return err
		}
	}
	return nil
}

func getNumFields(templateDir string) (int, error) {
	projectJSON, err := ioutil.ReadFile(path.Join(templateDir, "project.json"))
	if err != nil {
		return 0, err
	}
	var defaults map[string]string
	if err := json.Unmarshal(projectJSON, &defaults); err != nil {
		return 0, err
	}
	fields := []string{}
	for field := range defaults {
		fields = append(fields, field)
	}
	return len(fields), nil
}

func selectFirstOption(cmd *execplus.CmdPlus, field string) error {
	if err := cmd.WaitForText(field+":", time.Second*5); err != nil {
		return err
	}
	_, err := cmd.StdinPipe.Write([]byte("1" + "\n"))
	return err
}

func sendCmdOutputToWriter(cmd *execplus.CmdPlus, outputWriter io.Writer) {
	outputChannel, _ := cmd.GetOutputChannel()
	for {
		outputChunk := <-outputChannel
		_, err := fmt.Fprintln(outputWriter, outputChunk.Chunk)
		if err != nil {
			fmt.Printf("Error sending cmd output to writer: %s\n", err)
		}
	}
}
