package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"template"
	"log"

	"github.com/spf13/cobra"
	"github.com/tcnksm/go-input"
)

var rootCmd = &cobra.Command{
	Use:   "exo create",
	Short: "Create a new Exosphere application\n- options: \"[<app-name>] [<app-version>] [<exocom-version>] [<app-description>]\"",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 1 && args[0] == "help" {
			err := cmd.Help()
			if err != nil {
				panic(err)
			}
			return
		}
		fmt.Print("We are about to create a new Exosphere application\n\n")
		appConfig := getAppConfig(args)
		templatePath := path.Join("..", "exosphere-shared", "templates", "create-app", "application.yml")
		t, err := template.ParseFiles(templatePath)
		if err != nil {
		    panic(err)
		}
		f, err := os.Create(path.Join(appConfig.AppName, "application.yml"))
		if err != nil {
		    panic(err)
		}
		err = t.Execute(f, appConfig)
		if err != nil {
		    panic(err)
		}
		f.Close()
		f, err := os.Mkdir(path.Join(appConfig.AppName, ".exosphere"), os.FileMode(0522))
		if err != nil {
    	panic(err)
		}
		f.Close()
	},
}

type AppConfig struct {
    AppName, AppVersion, ExocomVersion, AppDescription  string
}

func ask(query string, defaultVal string, optional bool) string {
  ui := &input.UI{
    Writer: os.Stdout,
    Reader: os.Stdin,
  }
  answer, err := ui.Ask(query, &input.Options{
    Default: defaultVal,
    Required: true,
    Loop:     true,
  })
  if err != nil {
    panic(err)
  }
  answer = strings.TrimSpace(answer)
  if len(answer) == 0 && !optional {
  	panic("expect a non-empty string")
  }
  return answer
}

func getAppConfig(args []string) AppConfig {
  if len(args) > 1 {
    appName = args[1]
  } else {
    appName = ask("Name of the application to create:", "", false)
  }
  if len(args) > 2 {
    appVersion = args[2]
  } else {
    appVersion = ask("Initial version:", "0.0.1", false)
  }
  if len(args) > 3 {
    exocomVersion = args[3]
  } else {
    exocomVersion = ask("ExoCom version:", "0.22.1", false)
  }
  if len(args) > 4 {
    appDescription = strings.Join(args[4:], " ")
  } else {
    appDescription = ask("Description:", "", true)
  }
  return AppConfig{appName, appVersion, exocomVersion, appDescription}
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
