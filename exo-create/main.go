package main

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "exo create",
	Short: "Usage: \"exo create application\"\nCreate a new Exosphere application\n  - options: \"[<app-name>] [<app-version>] [<exocom-version>] [<app-description>]\"",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Error: missing argument for 'create' command")
			return
		} else if args[0] == "help" {
			if err := cmd.Help(); err != nil {
				panic(err)
			}
			return
		} else if args[0] != "application" {
			fmt.Println("Error: cannot create '" + args[0] + "'")
			return
		}
		fmt.Print("We are about to create a new Exosphere application\n\n")
		appConfig := getAppConfig(args)
		if err := createApplicationYAML(appConfig); err != nil {
			panic(err)
		}
		if err := os.Mkdir(path.Join(appConfig.AppName, ".exosphere"), os.FileMode(0777)); err != nil {
			panic(err)
		}
		fmt.Println("\ndone")
	},
}

// AppConfig contains the configuration for the application
type AppConfig struct {
	AppName, AppVersion, ExocomVersion, AppDescription string
}

func ask(reader *bufio.Reader, query string, defaultVal string, required bool) string {
	if len(defaultVal) > 0 {
		query = query + "(" + defaultVal + ") "
	}
	fmt.Print(query)
	answer, err := reader.ReadString('\n')
	answer = strings.TrimSpace(answer)
	if err != nil {
		panic(err)
	}
	if len(answer) == 0 {
		if len(defaultVal) > 0 {
			answer = defaultVal
		} else if required {
			fmt.Println("expect a non-empty string")
			return ask(reader, query, defaultVal, required)
		}
	}
	return answer
}

func getAppConfig(args []string) AppConfig {
	var appName, appVersion, exocomVersion, appDescription string
	reader := bufio.NewReader(os.Stdin)
	if len(args) > 1 {
		appName = args[1]
	} else {
		appName = ask(reader, "Name of the application to create: ", "", true)
	}
	if len(args) > 2 {
		appVersion = args[2]
	} else {
		appVersion = ask(reader, "Initial version: ", "0.0.1", true)
	}
	if len(args) > 3 {
		exocomVersion = args[3]
	} else {
		exocomVersion = ask(reader, "ExoCom version: ", "0.22.1", true)
	}
	if len(args) > 4 {
		appDescription = strings.Join(args[4:], " ")
	} else {
		appDescription = ask(reader, "Description: ", "", false)
	}
	return AppConfig{appName, appVersion, exocomVersion, appDescription}
}

func createApplicationYAML(appConfig AppConfig) error {
	dir, err := os.Executable()
	if err != nil {
		return err
	}
	templatePath := path.Join(dir, "..", "..", "..", "exosphere-shared", "templates", "create-app", "application-go.yml")
	t, err := template.ParseFiles(templatePath)
	if err != nil {
		return err
	}
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	if err = os.Mkdir(path.Join(cwd, appConfig.AppName), os.FileMode(0777)); err != nil {
		return err
	}
	f, err := os.Create(path.Join(cwd, appConfig.AppName, "application.yml"))
	if err != nil {
		return err
	}
	if err = t.Execute(f, appConfig); err != nil {
		return err
	}
	if err = f.Close(); err != nil {
		return err
	}
	return nil
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
