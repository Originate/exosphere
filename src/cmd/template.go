package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/Originate/exosphere/src/template"
	"github.com/spf13/cobra"
	"github.com/tmrts/boilr/pkg/util/osutil"
)

var templateCmd = &cobra.Command{
	Use:   "template",
	Short: "Manage service templates",
	Long:  "Manage service templates",
}

var testTemplateCmd = &cobra.Command{
	Use:   "test",
	Short: "Test a service templates",
	Long: `Test a service templates by adding to an exopshere application and running tests.
This command must be run in the directory of an exosphere template.`,
	Run: func(cmd *cobra.Command, args []string) {
		if printHelpIfNecessary(cmd, args) {
			return
		}
		templateDir, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		writer := os.Stdout
		templateName := filepath.Base(templateDir)
		valid, msg, err := template.IsValidTemplateDir(templateDir)
		if err != nil {
			panic(err)
		}
		if !valid {
			fmt.Println(msg)
			os.Exit(1)
		}
		tempDir, err := ioutil.TempDir("", "template-test")
		if err != nil {
			panic(err)
		}
		if err := template.CreateEmptyApp(tempDir); err != nil {
			panic(err)
		}
		appDir := path.Join(tempDir, "my-app")
		if err := osutil.CopyRecursively(templateDir, path.Join(appDir, ".exosphere", templateName)); err != nil {
			panic(err)
		}
		fmt.Fprintln(writer, "Adding a service with this template to an empty exosphere application...")
		addServiceErr := template.AddService(appDir, templateDir, os.Stdout)
		if addServiceErr != nil {
			fmt.Fprintln(writer, "Error adding service")
			os.Exit(1)
		}
		fmt.Fprintln(writer, "Running tests...")
		testErr := template.RunTests(appDir, os.Stdout)
		if testErr != nil {
			fmt.Fprintln(writer, "Tests failed")
			os.Exit(1)
		}
	},
}

func init() {
	templateCmd.AddCommand(testTemplateCmd)
	RootCmd.AddCommand(templateCmd)
}
