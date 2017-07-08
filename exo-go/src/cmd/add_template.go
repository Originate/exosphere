package cmd

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"

	yaml "gopkg.in/yaml.v2"

	"github.com/Originate/exosphere/exo-go/src/app_config_helpers"
	"github.com/Originate/exosphere/exo-go/src/os_helpers"
	"github.com/Originate/exosphere/exo-go/src/template_helpers"
	"github.com/Originate/exosphere/exo-go/src/user_input_helpers"
	"github.com/Originate/exosphere/exo-go/src/util"
	"github.com/spf13/cobra"
)

var addTemplateCmd = &cobra.Command{
	Use:   "add-template",
	Short: "Adds a remote service template to .exosphere",
	Long: `Adds a remote service template to .exosphere
This command must be called in the root directory of the application`,
	Run: func(cmd *cobra.Command, args []string) {
		if util.PrintHelpIfNecessary(cmd, args) {
			return
		}
		if len(args) != 2 {
			fmt.Println("not enough arguments")
      os.Exit(1)
		}
		fmt.Print("We are about to add a new service template to .exosphere folder\n\n")
		templateName, gitURL := args[0], args[1]
		templateDir := path.Join(".exosphere", templateName)
		if osHelpers.DirectoryExists(templateDir) {
			reader := bufio.NewReader(os.Stdin)
			query := fmt.Sprintf(`Template "%s" already exists, do you want to update its git URL to %s`, templateName, gitURL)
			if userInputHelpers.Confirm(reader, query) {
				if err := templateHelpers.UpdateTemplate(gitURL, templateName); err != nil {
					log.Fatalf(`Failed to update template "%s": %s`, templateName, err)
				}
			} else {
				fmt.Printf("\"%s\" not updated\n", templateName)
				return
			}
		} else {
			if err := templateHelpers.AddTemplate(gitURL, templateName, templateDir); err != nil {
				log.Fatalf(`Failed to add template "%s": %s`, templateName, err)
			}
		}
		appConfig := appConfigHelpers.GetAppConfig()
		if len(appConfig.Templates) == 0 {
			appConfig.Templates = make(map[string]string)
		}
		appConfig.Templates[templateName] = gitURL
		bytes, err := yaml.Marshal(appConfig)
		if err != nil {
			log.Fatalf("Failed to marshal application.yml: %s", err)
		}
		err = ioutil.WriteFile("application.yml", bytes, 0777)
		if err != nil {
			log.Fatalf("Failed to update application.yml: %s", err)
		}
		fmt.Println("\ndone")
	},
}

func init() {
	RootCmd.AddCommand(addTemplateCmd)
}
