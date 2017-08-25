package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/Originate/exosphere/src/template"
	"github.com/Originate/exosphere/src/util"
	"github.com/spf13/cobra"
	"github.com/tmrts/boilr/pkg/util/osutil"
)

var templateCmd = &cobra.Command{
	Use:   "template",
	Short: "Manages remote service templates",
	Long:  "Manages remote service templates",
}

var addTemplateCmd = &cobra.Command{
	Use:   "add <name> <url> [<commit-ish>]",
	Short: "Adds a remote service template to .exosphere",
	Long:  "Adds a remote service template to .exosphere",
	Run: func(cmd *cobra.Command, args []string) {
		if printHelpIfNecessary(cmd, args) {
			return
		}
		if len(args) != 2 && len(args) != 3 {
			fmt.Println("wrong number of arguments")
			os.Exit(1)
		}
		fmt.Print("We are about to add a new service template\n")
		templateName, gitURL := args[0], args[1]
		commitIsh := ""
		if len(args) == 3 {
			commitIsh = args[2]
		}
		templateDir := path.Join(".exosphere", templateName)
		hasTemplate, err := util.DoesDirectoryExist(templateDir)
		if err != nil {
			panic(err)
		}
		if hasTemplate {
			fmt.Printf(`The template "%s" already exists\n`, templateName)
			os.Exit(1)
		} else {
			if err := template.Add(gitURL, templateName, templateDir, commitIsh); err != nil {
				panic(err)
			}
		}
		fmt.Println("\ndone")
	},
}

var fetchTemplatesCmd = &cobra.Command{
	Use:   "fetch <name>",
	Short: "Fetches updates for an existing service template in .exosphere",
	Long:  "Fetches updates for an existing service template in .exosphere",
	Run: func(cmd *cobra.Command, args []string) {
		if printHelpIfNecessary(cmd, args) {
			return
		}
		if len(args) != 1 {
			fmt.Println("wrong number of arguments")
			os.Exit(1)
		}
		templateName := args[0]
		templateDir := path.Join(".exosphere", templateName)
		fmt.Printf("We are about to fetch updates for the template  \"%s\"\n\n", templateName)
		if err := template.Fetch(templateDir); err != nil {
			panic(err)
		}
		fmt.Println("\ndone")
	},
}

var removeTemplateCmd = &cobra.Command{
	Use:   "remove <name>",
	Short: "Removes an existing service template from .exosphere",
	Long:  "Removes an existing service template from .exosphere",
	Run: func(cmd *cobra.Command, args []string) {
		if printHelpIfNecessary(cmd, args) {
			return
		}
		if len(args) != 1 {
			fmt.Println("wrong number of arguments")
			os.Exit(1)
		}
		templateName := args[0]
		templateDir := path.Join(".exosphere", templateName)
		fmt.Printf("We are about to remove the template \"%s\"\n\n", templateName)
		hasTemplate, err := util.DoesDirectoryExist(templateDir)
		if err != nil {
			panic(err)
		}
		if !hasTemplate {
			fmt.Println("Error: template does not exist")
			os.Exit(1)
		} else {
			if err := template.Remove(templateName, templateDir); err != nil {
				panic(err)
			}
			fmt.Println("\ndone")
		}
	},
}

var testTemplateCmd = &cobra.Command{
	Use:   "test",
	Short: "Tests service templates",
	Long: `Tests service templates by adding to an exopshere application and running tests.
This command must be run in the directory of an exosphere template.`,
	Run: func(cmd *cobra.Command, args []string) {
		if printHelpIfNecessary(cmd, args) {
			return
		}
		templateDir, err := os.Getwd()
		if err != nil {
			panic(err)
		}
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
		fmt.Print("Adding a service with this template to an empty exosphere application...\n\n")
		addServiceErr := template.AddService(appDir, templateDir, os.Stdout)
		fmt.Print("\n\n")
		if addServiceErr != nil {
			fmt.Printf("Error adding service: %s", addServiceErr)
			os.Exit(1)
		}
		fmt.Print("Successfully added service\n\n")
		fmt.Println("Running tests...\n\n")
		testErr := template.RunTests(appDir, os.Stdout)
		fmt.Print("\n\n")
		if testErr != nil {
			fmt.Printf("Tests failed: %s", testErr)
			os.Exit(1)
		}
		fmt.Print("Tests passed")
	},
}

func init() {
	templateCmd.AddCommand(addTemplateCmd)
	templateCmd.AddCommand(removeTemplateCmd)
	templateCmd.AddCommand(fetchTemplatesCmd)
	templateCmd.AddCommand(testTemplateCmd)
	RootCmd.AddCommand(templateCmd)
}
