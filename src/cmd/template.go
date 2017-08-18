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
		if util.DoesDirectoryExist(templateDir) {
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
		if !util.DoesDirectoryExist(templateDir) {
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
	Short: "Tests the service template",
	Long:  "Tests the service template",
	Run: func(cmd *cobra.Command, args []string) {
		if printHelpIfNecessary(cmd, args) {
			return
		}
		templateDir, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		templateName := filepath.Base(templateDir)
		fmt.Printf("We are about to test the template \"%s\"\n\n", templateName)
		valid, err := template.IsValidTemplateDir(templateDir)
		if err != nil {
			panic(err)
		}
		if !valid {
			fmt.Println("Template fails")
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
		fmt.Println("appDir", appDir)
		if err := osutil.CopyRecursively(templateDir, path.Join(appDir, ".exosphere", templateName)); err != nil {
			panic(err)
		}
		if err := template.AddService(appDir, templateDir); err != nil {
			panic(err)
		}
		testPassed := template.RunTests(appDir)
		if !testPassed {
			fmt.Println("Template fails")
			os.Exit(1)
		}
		fmt.Println("Template passes")
		if err := os.RemoveAll(tempDir); err != nil {
			panic(err)
		}
	},
}

func init() {
	templateCmd.AddCommand(addTemplateCmd)
	templateCmd.AddCommand(removeTemplateCmd)
	templateCmd.AddCommand(fetchTemplatesCmd)
	templateCmd.AddCommand(testTemplateCmd)
	RootCmd.AddCommand(templateCmd)
}
