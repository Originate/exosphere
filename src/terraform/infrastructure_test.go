package terraform_test

import (
	"fmt"
	"io/ioutil"

	"github.com/Originate/exosphere/src/terraform"
	"github.com/Originate/exosphere/src/types"
	"github.com/Originate/exosphere/src/types/context"
	"github.com/Originate/exosphere/src/types/deploy"
	"github.com/Originate/exosphere/src/types/hcl"
	"github.com/Originate/exosphere/test/helpers"
	"github.com/Originate/exosphere/test/matchers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("GenerateInfrastructure", func() {
	var _ = Describe("with no dependencies", func() {
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
			result, err := terraform.GenerateInfrastructure(deployConfig)
			Expect(err).To(BeNil())
			hclFile, err := hcl.GetHCLFileFromTerraform(result)
			Expect(err).To(BeNil())
			Expect(hclFile).To(matchers.HaveHCLVariable("aws_profile"))
			Expect(hclFile).To(matchers.HaveHCLVariable("aws_region"))
			Expect(hclFile).To(matchers.HaveHCLVariable("aws_account_id"))
			Expect(hclFile).To(matchers.HaveHCLVariable("aws_ssl_certificate_arn"))
			Expect(hclFile).To(matchers.HaveHCLVariable("application_url"))
			Expect(hclFile).To(matchers.HaveHCLVariable("env"))
			Expect(hclFile).To(matchers.HaveHCLVariable("key_name"))
			Expect(hclFile.GetModuleNames()).To(Equal([]string{"aws"}))
			Expect(hclFile.Module["aws"]).To(Equal(hcl.Module{
				"source":            fmt.Sprintf("github.com/Originate/exosphere.git//terraform//aws?ref=%s", terraform.TerraformModulesRef),
				"key_name":          "${var.key_name}",
				"name":              "example-app",
				"env":               "${var.env}",
				"external_dns_name": "${var.application_url}",
				"log_bucket_prefix": "${var.aws_account_id}-example-app-${var.env}",
			}))
		})
	})

	var _ = Describe("with dependencies", func() {
		It("should generate dependency modules for exocom", func() {
			appDir, err := ioutil.TempDir("", "")
			Expect(err).NotTo(HaveOccurred())
			err = helpers.CheckoutApp(appDir, "simple")
			Expect(err).NotTo(HaveOccurred())
			appContext, err := context.GetAppContext(appDir)
			Expect(err).NotTo(HaveOccurred())

			deployConfig := deploy.Config{
				AppContext: appContext,
			}
			result, err := terraform.GenerateInfrastructure(deployConfig)
			Expect(err).To(BeNil())
			hclFile, err := hcl.GetHCLFileFromTerraform(result)
			Expect(err).To(BeNil())
			Expect(hclFile).To(matchers.HaveHCLVariable("exocom_env_vars"))
			Expect(hclFile.Module["exocom_cluster"]).To(Equal(hcl.Module{
				"source":                      fmt.Sprintf("github.com/Originate/exosphere.git//remote-dependency-templates//exocom//modules//exocom-cluster?ref=%s", terraform.TerraformModulesRef),
				"availability_zones":          "${module.aws.availability_zones}",
				"bastion_security_group":      []interface{}{"${module.aws.bastion_security_group}"},
				"ecs_cluster_security_groups": []interface{}{"${module.aws.ecs_cluster_security_group}", "${module.aws.external_alb_security_group}"},
				"env":                     "${var.env}",
				"instance_type":           "t2.micro",
				"internal_hosted_zone_id": "${module.aws.internal_zone_id}",
				"key_name":                "${var.key_name}",
				"name":                    "exocom",
				"region":                  "${module.aws.region}",
				"subnet_ids":              "${module.aws.private_subnet_ids}",
				"vpc_id":                  "${module.aws.vpc_id}",
			}))
			Expect(hclFile.Module["exocom_service"]).To(Equal(hcl.Module{
				"source":       fmt.Sprintf("github.com/Originate/exosphere.git//remote-dependency-templates//exocom//modules//exocom-service?ref=%s", terraform.TerraformModulesRef),
				"cluster_id":   "${module.exocom_cluster.cluster_id}",
				"cpu_units":    "128",
				"docker_image": "originate/exocom:0.27.0",
				"env":          "${var.env}",
				"environment_variables": "${var.exocom_env_vars}",
				"memory_reservation":    "128",
				"name":                  "exocom",
				"region":                "${module.aws.region}",
			}))
		})

		It("should generate rds modules for dependencies", func() {
			appDir, err := ioutil.TempDir("", "")
			Expect(err).NotTo(HaveOccurred())
			err = helpers.CheckoutApp(appDir, "rds")
			Expect(err).NotTo(HaveOccurred())
			appContext, err := context.GetAppContext(appDir)
			Expect(err).NotTo(HaveOccurred())

			deployConfig := deploy.Config{
				AppContext: appContext,
			}
			result, err := terraform.GenerateInfrastructure(deployConfig)
			Expect(err).To(BeNil())
			hclFile, err := hcl.GetHCLFileFromTerraform(result)
			Expect(err).To(BeNil())
			By("generating rds modules for application dependencies", func() {
				Expect(hclFile.Module["my-db_rds_instance"]).To(Equal(hcl.Module{
					"source":                  fmt.Sprintf("github.com/Originate/exosphere.git//remote-dependency-templates//rds//modules//rds?ref=%s", terraform.TerraformModulesRef),
					"allocated_storage":       "10",
					"bastion_security_group":  "${module.aws.bastion_security_group}",
					"ecs_security_group":      "${module.aws.ecs_cluster_security_group}",
					"engine":                  "postgres",
					"engine_version":          "9.6.4",
					"env":                     "${var.env}",
					"instance_class":          "db.t2.micro",
					"internal_hosted_zone_id": "${module.aws.internal_zone_id}",
					"name":         "my-db",
					"password":     "${var.POSTGRES_PASSWORD}",
					"storage_type": "gp2",
					"subnet_ids":   []interface{}{"${module.aws.private_subnet_ids}"},
					"username":     "originate-user",
					"vpc_id":       "${module.aws.vpc_id}",
				}))
			})

			By("should generate rds modules for service dependencies", func() {
				Expect(hclFile.Module["my-sql-db_rds_instance"]).To(Equal(hcl.Module{
					"source":                  fmt.Sprintf("github.com/Originate/exosphere.git//remote-dependency-templates//rds//modules//rds?ref=%s", terraform.TerraformModulesRef),
					"allocated_storage":       "10",
					"bastion_security_group":  "${module.aws.bastion_security_group}",
					"ecs_security_group":      "${module.aws.ecs_cluster_security_group}",
					"engine":                  "mysql",
					"engine_version":          "5.6.17",
					"env":                     "${var.env}",
					"instance_class":          "db.t1.micro",
					"internal_hosted_zone_id": "${module.aws.internal_zone_id}",
					"name":         "my-sql-db",
					"password":     "${var.MYSQL_PASSWORD}",
					"storage_type": "gp2",
					"subnet_ids":   []interface{}{"${module.aws.private_subnet_ids}"},
					"username":     "originate-user",
					"vpc_id":       "${module.aws.vpc_id}",
				}))
			})
		})
	})
})
