package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/Originate/exosphere/src/aws"
	"github.com/Originate/exosphere/src/types"
	"github.com/Originate/exosphere/src/types/context"
	"github.com/Originate/exosphere/src/types/deploy"
)

func getSecrets(deployConfig deploy.Config) types.Secrets {
	secrets, err := aws.ReadSecrets(deployConfig.GetAwsOptions())
	if err != nil {
		log.Fatalf("Cannot read secrets: %s", err)
	}
	fmt.Print("Existing secrets:\n\n")
	prettyPrintSecrets(secrets)
	return secrets
}

func getBaseDeployConfig(remoteEnvironmentID, awsProfile string) deploy.Config {
	userContext, err := GetUserContext()
	if err != nil {
		log.Fatal(err)
	}
	err = validateRemoteEnvironmentID(userContext, remoteEnvironmentID)
	if err != nil {
		log.Fatal(err)
	}
	return deploy.Config{
		AppContext:          userContext.AppContext,
		AwsProfile:          awsProfile,
		RemoteEnvironmentID: remoteEnvironmentID,
		Writer:              os.Stdout,
	}
}

func prettyPrintSecrets(secrets map[string]string) {
	secretsPretty, err := json.MarshalIndent(secrets, "", "  ")
	if err != nil {
		log.Fatalf("Could not marshal secrets map: %s", err)
	}
	fmt.Printf("%s\n\n", string(secretsPretty))
}

func validateArgCount(args []string, number int) error {
	if len(args) != number {
		return errors.New("Wrong number of arguments")
	}
	return nil
}

func validateRemoteEnvironmentID(userContext *context.UserContext, remoteEnvironmentID string) error {
	if _, ok := userContext.AppContext.Config.Remote.Environments[remoteEnvironmentID]; !ok {
		validIDs := []string{}
		for validID := range userContext.AppContext.Config.Remote.Environments {
			validIDs = append(validIDs, validID)
		}
		return fmt.Errorf("Invalid remote environment id: %s. Valid remote environment ids: %s", remoteEnvironmentID, strings.Join(validIDs, ","))
	}
	return nil
}
