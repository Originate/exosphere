package terraform_test

import (
	"fmt"

	"github.com/Originate/exosphere/src/terraform"
	"github.com/Originate/exosphere/src/types"
	"github.com/Originate/exosphere/src/types/context"
	"github.com/Originate/exosphere/src/types/deploy"
	"github.com/Originate/exosphere/src/types/hcl"
	"github.com/Originate/exosphere/test/matchers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("GenerateServices", func() {
	var _ = Describe("no services", func() {
		appConfig := &types.AppConfig{
			Name: "example-app",
			Remote: types.AppRemoteConfig{
				Environments: map[string]types.AppRemoteEnvironment{
					"qa": {
						URL:       "example-app.com",
						AccountID: "12345",
						Region:    "us-west-2",
					},
				},
			},
		}

		deployConfig := deploy.Config{
			AppContext: &context.AppContext{
				Config: appConfig,
			},
			RemoteEnvironmentID: "qa",
		}

		It("should generate an AWS module only", func() {
			result, err := terraform.GenerateServices(deployConfig)
			Expect(err).To(BeNil())
			hclFile, err := hcl.GetHCLFileFromTerraform(result)
			Expect(err).To(BeNil())
			Expect(hclFile).To(matchers.HaveHCLVariable("aws_profile"))
			Expect(hclFile).To(matchers.HaveHCLVariable("aws_region"))
			Expect(hclFile).To(matchers.HaveHCLVariable("aws_account_id"))
			Expect(hclFile).To(matchers.HaveHCLVariable("aws_ssl_certificate_arn"))
			Expect(hclFile).To(matchers.HaveHCLVariable("application_url"))
			Expect(hclFile).To(matchers.HaveHCLVariable("env"))
			Expect(hclFile.GetDataNames()).To(Equal([]string{"terraform_remote_state"}))
			Expect(hclFile.Data["terraform_remote_state"]["main_infrastructure"]).To(Equal([]map[string]interface{}{
				{
					"backend": "s3",
					"config": []map[string]interface{}{
						{
							"key":            "${var.aws_account_id}-example-app-${var.env}-terraform/infrastructure.tfstate",
							"dynamodb_table": "TerraformLocks",
							"region":         "${var.aws_region}",
							"profile":        "${var.aws_profile}",
						},
					},
				},
			}))
		})
	})

	var _ = Describe("public and worker services", func() {
		var hclFile *hcl.File
		appConfig := &types.AppConfig{
			Name: "example-app",
			Services: map[string]types.ServiceSource{
				"public-service": types.ServiceSource{},
				"worker-service": types.ServiceSource{},
			},
			Remote: types.AppRemoteConfig{
				Environments: map[string]types.AppRemoteEnvironment{
					"qa": {
						SslCertificateArn: "sslcert123",
					},
				},
			},
		}
		serviceContexts := map[string]*context.ServiceContext{
			"public-service": {
				Config: types.ServiceConfig{
					Type: "public",
					Production: types.ServiceProductionConfig{
						Port:        "3000",
						HealthCheck: "/health-check",
					},
					Remote: types.ServiceRemoteConfig{
						CPU:    "128",
						Memory: "128",
						Environments: map[string]types.ServiceRemoteEnvironment{
							"qa": {
								URL: "originate.com",
							},
						},
					},
				},
			},
			"worker-service": {
				Config: types.ServiceConfig{
					Type: "worker",
					Remote: types.ServiceRemoteConfig{
						CPU:    "128",
						Memory: "128",
					},
				},
			},
		}

		deployConfig := deploy.Config{
			AppContext: &context.AppContext{
				Config:          appConfig,
				ServiceContexts: serviceContexts,
			},
			RemoteEnvironmentID: "qa",
		}

		BeforeEach(func() {
			var err error
			result, err := terraform.GenerateServices(deployConfig)
			Expect(err).To(BeNil())
			hclFile, err = hcl.GetHCLFileFromTerraform(result)
			Expect(err).To(BeNil())
		})

		It("should generate a public service module", func() {
			Expect(hclFile).To(matchers.HaveHCLVariable("public-service_env_vars"))
			Expect(hclFile).To(matchers.HaveHCLVariable("public-service_docker_image"))
			Expect(hclFile).To(matchers.HaveHCLVariable("public-service_url"))
			Expect(hclFile.Module["public-service"]).To(Equal(hcl.Module{
				"source":             fmt.Sprintf("github.com/Originate/exosphere.git//terraform//aws//public-service?ref=%s", terraform.TerraformModulesRef),
				"alb_security_group": "${data.terraform_remote_state.main_infrastructure.external_alb_security_group}",
				"alb_subnet_ids":     []interface{}{"${data.terraform_remote_state.main_infrastructure.public_subnet_ids}"},
				"cluster_id":         "${data.terraform_remote_state.main_infrastructure.ecs_cluster_id}",
				"container_port":     "3000",
				"cpu":                "128",
				"desired_count":      1,
				"docker_image":       "${var.public-service_docker_image}",
				"ecs_role_arn":       "${data.terraform_remote_state.main_infrastructure.ecs_service_iam_role_arn}",
				"env":                "${var.env}",
				"environment_variables": "${var.public-service_env_vars}",
				"external_dns_name":     "${var.public-service_url}",
				"external_zone_id":      "${data.terraform_remote_state.main_infrastructure.external_zone_id}",
				"health_check_endpoint": "/health-check",
				"internal_dns_name":     "public-service",
				"internal_zone_id":      "${data.terraform_remote_state.main_infrastructure.internal_zone_id}",
				"log_bucket":            "${data.terraform_remote_state.main_infrastructure.log_bucket_id}",
				"memory_reservation":    "128",
				"name":                  "public-service",
				"region":                "${data.terraform_remote_state.main_infrastructure.region}",
				"ssl_certificate_arn":   "${var.aws_ssl_certificate_arn}",
				"vpc_id":                "${data.terraform_remote_state.main_infrastructure.vpc_id}",
			}))
		})

		It("should generate a worker service module", func() {
			Expect(hclFile).To(matchers.HaveHCLVariable("worker-service_env_vars"))
			Expect(hclFile).To(matchers.HaveHCLVariable("worker-service_docker_image"))
			Expect(hclFile.Module["worker-service"]).To(Equal(hcl.Module{
				"source":        fmt.Sprintf("github.com/Originate/exosphere.git//terraform//aws//worker-service?ref=%s", terraform.TerraformModulesRef),
				"cluster_id":    "${data.terraform_remote_state.main_infrastructure.ecs_cluster_id}",
				"cpu":           "128",
				"desired_count": 1,
				"docker_image":  "${var.worker-service_docker_image}",
				"env":           "${var.env}",
				"environment_variables": "${var.worker-service_env_vars}",
				"memory_reservation":    "128",
				"name":                  "worker-service",
				"region":                "${data.terraform_remote_state.main_infrastructure.region}",
			}))
		})
	})
})
