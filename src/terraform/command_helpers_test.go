package terraform_test

import (
	"encoding/json"
	"strings"

	"github.com/Originate/exosphere/src/terraform"
	"github.com/Originate/exosphere/src/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("CompileVarFlags", func() {
	var _ = Describe("no dependencies", func() {
		service1EnvVars := types.EnvVars{
			Default: map[string]string{
				"env1": "val1",
			},
			Secrets: []string{"secret1"},
		}
		service1Config := types.ServiceConfig{
			Environment: service1EnvVars,
		}
		secrets := map[string]string{
			"secret1": "secret_value1",
		}
		deployConfig := types.DeployConfig{
			ServiceConfigs: map[string]types.ServiceConfig{
				"service1": service1Config,
			},
		}

		It("should compile the proper var flags", func() {
			vars, err := terraform.CompileVarFlags(deployConfig, secrets)
			Expect(err).NotTo(HaveOccurred())
			Expect(vars[0]).To(Equal("-var"))
			Expect(vars[1]).To(Equal("secret1=secret_value1"))
			Expect(vars[2]).To(Equal("-var"))

			varName := strings.Split(vars[3], "=")[0]
			varVal := strings.Split(vars[3], "=")[1]
			var escapedVal string
			actualVal := []map[string]string{}
			expectedVal := []map[string]string{
				{
					"name":  "ROLE",
					"value": "service1",
				},
				{
					"name":  "secret1",
					"value": "secret_value1",
				},
				{
					"name":  "env1",
					"value": "val1",
				},
			}
			err = json.Unmarshal([]byte(varVal), &escapedVal)
			Expect(err).NotTo(HaveOccurred())
			err = json.Unmarshal([]byte(escapedVal), &actualVal)
			Expect(err).NotTo(HaveOccurred())
			Expect(varName).To(Equal("service1_env_vars"))
			Expect(expectedVal).To(ConsistOf(actualVal))
		})
	})

	var _ = Describe("with exocom dependency", func() {
		deployConfig := types.DeployConfig{
			AppConfig: types.AppConfig{
				Dependencies: []types.DependencyConfig{
					{Name: "exocom"},
				},
				Name: "my-app",
			},
			ServiceConfigs: map[string]types.ServiceConfig{
				"service1": {},
			},
		}

		It("should add the dependency service env vars to each service", func() {
			vars, err := terraform.CompileVarFlags(deployConfig, map[string]string{})
			Expect(err).NotTo(HaveOccurred())
			Expect(vars[0]).To(Equal("-var"))
			varName := strings.Split(vars[1], "=")[0]
			varVal := strings.Split(vars[1], "=")[1]
			var escapedVal string
			actualVal := []map[string]string{}
			expectedVal := []map[string]string{
				{
					"name":  "ROLE",
					"value": "service1",
				},
				{
					"name":  "EXOCOM_HOST",
					"value": "exocom.my-app.local",
				},
			}
			err = json.Unmarshal([]byte(varVal), &escapedVal)
			Expect(err).NotTo(HaveOccurred())
			err = json.Unmarshal([]byte(escapedVal), &actualVal)
			Expect(err).NotTo(HaveOccurred())
			Expect(varName).To(Equal("service1_env_vars"))
			Expect(expectedVal).To(ConsistOf(actualVal))
		})
	})
})
