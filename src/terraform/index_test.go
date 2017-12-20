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

var _ = Describe("Template builder", func() {
	var _ = Describe("Given an application with no services", func() {
		appConfig := types.AppConfig{
			Name: "example-app",
			Remote: types.AppRemoteConfig{
				URL: "example-app.com",
			},
		}

		deployConfig := deploy.Config{
			AppContext: &context.AppContext{
				Config: appConfig,
			},
			AwsConfig: types.AwsConfig{
				TerraformStateBucket: "example-app-terraform",
				TerraformLockTable:   "TerraformLocks",
				Region:               "us-west-2",
				AccountID:            "12345",
			},
		}

		It("should generate an AWS module only", func() {
			result, err := terraform.Generate(deployConfig)
			Expect(err).To(BeNil())
			hclFile, err := hcl.GetHCLFileFromTerraform(result)
			Expect(err).To(BeNil())
			Expect(hclFile).To(matchers.HaveHCLVariable("aws_profile"))
			Expect(hclFile).To(matchers.HaveHCLVariable("aws_region"))
			Expect(hclFile).To(matchers.HaveHCLVariable("aws_account_id"))
			Expect(hclFile).To(matchers.HaveHCLVariable("aws_ssl_certificate_arn"))
			Expect(hclFile).To(matchers.HaveHCLVariable("application_url"))
			Expect(hclFile).To(matchers.HaveHCLVariable("key_name"))
			Expect(hclFile.GetModuleNames()).To(Equal([]string{"aws"}))
			Expect(hclFile.Module["aws"]).To(Equal(hcl.Module{
				"source":            fmt.Sprintf("github.com/Originate/exosphere.git//terraform//aws?ref=%s", terraform.TerraformModulesRef),
				"key_name":          "${var.key_name}",
				"name":              "example-app",
				"env":               "production",
				"external_dns_name": "${var.application_url}",
			}))
		})
	})

	var _ = Describe("Given an application with public and worker services", func() {
		var hclFile *hcl.File
		appConfig := types.AppConfig{
			Name: "example-app",
			Services: map[string]types.ServiceSource{
				"public-service": types.ServiceSource{},
				"worker-service": types.ServiceSource{},
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
						URL:    "originate.com",
						Memory: "128",
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
			AwsConfig: types.AwsConfig{
				SslCertificateArn: "sslcert123",
			},
		}

		BeforeEach(func() {
			var err error
			result, err := terraform.Generate(deployConfig)
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
				"alb_security_group": "${module.aws.external_alb_security_group}",
				"alb_subnet_ids":     []interface{}{"${module.aws.public_subnet_ids}"},
				"cluster_id":         "${module.aws.ecs_cluster_id}",
				"container_port":     "3000",
				"cpu":                "128",
				"desired_count":      1,
				"docker_image":       "${var.public-service_docker_image}",
				"ecs_role_arn":       "${module.aws.ecs_service_iam_role_arn}",
				"env":                "production",
				"environment_variables": "${var.public-service_env_vars}",
				"external_dns_name":     "${var.public-service_url}",
				"external_zone_id":      "${module.aws.external_zone_id}",
				"health_check_endpoint": "/health-check",
				"internal_dns_name":     "public-service",
				"internal_zone_id":      "${module.aws.internal_zone_id}",
				"log_bucket":            "${module.aws.log_bucket_id}",
				"memory_reservation":    "128",
				"name":                  "public-service",
				"region":                "${module.aws.region}",
				"ssl_certificate_arn":   "${var.aws_ssl_certificate_arn}",
				"vpc_id":                "${module.aws.vpc_id}",
			}))
		})

		It("should generate a worker service module", func() {
			Expect(hclFile).To(matchers.HaveHCLVariable("worker-service_env_vars"))
			Expect(hclFile).To(matchers.HaveHCLVariable("worker-service_docker_image"))
			Expect(hclFile.Module["worker-service"]).To(Equal(hcl.Module{
				"source":        fmt.Sprintf("github.com/Originate/exosphere.git//terraform//aws//worker-service?ref=%s", terraform.TerraformModulesRef),
				"cluster_id":    "${module.aws.ecs_cluster_id}",
				"cpu":           "128",
				"desired_count": 1,
				"docker_image":  "${var.worker-service_docker_image}",
				"env":           "production",
				"environment_variables": "${var.worker-service_env_vars}",
				"memory_reservation":    "128",
				"name":                  "worker-service",
				"region":                "${module.aws.region}",
			}))
		})
	})

	var _ = Describe("Given an application with dependencies", func() {
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
			result, err := terraform.Generate(deployConfig)
			Expect(err).To(BeNil())
			hclFile, err := hcl.GetHCLFileFromTerraform(result)
			Expect(err).To(BeNil())
			Expect(hclFile).To(matchers.HaveHCLVariable("exocom_env_vars"))
			Expect(hclFile.Module["exocom_cluster"]).To(Equal(hcl.Module{
				"source":                      fmt.Sprintf("github.com/Originate/exosphere.git//remote-dependency-templates//exocom//modules//exocom-cluster?ref=%s", terraform.TerraformModulesRef),
				"availability_zones":          "${module.aws.availability_zones}",
				"bastion_security_group":      []interface{}{"${module.aws.bastion_security_group}"},
				"ecs_cluster_security_groups": []interface{}{"${module.aws.ecs_cluster_security_group}", "${module.aws.external_alb_security_group}"},
				"env":                     "production",
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
				"env":          "production",
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
			result, err := terraform.Generate(deployConfig)
			Expect(err).To(BeNil())
			hclFile, err := hcl.GetHCLFileFromTerraform(result)
			Expect(err).To(BeNil())
			By("generating rds modules for application dependencies", func() {
				Expect(hclFile.Module["my-db_rds_instance"]).To(Equal(hcl.Module{
					"source":                  fmt.Sprintf("github.com/Originate/exosphere.git//remote-dependency-templates//rds//module?ref=%s", terraform.TerraformModulesRef),
					"allocated_storage":       "10",
					"bastion_security_group":  "${module.aws.bastion_security_group}",
					"ecs_security_group":      "${module.aws.ecs_cluster_security_group}",
					"engine":                  "postgres",
					"engine_version":          "9.6.4",
					"env":                     "production",
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
					"source":                  fmt.Sprintf("github.com/Originate/exosphere.git//remote-dependency-templates//rds//module?ref=%s", terraform.TerraformModulesRef),
					"allocated_storage":       "10",
					"bastion_security_group":  "${module.aws.bastion_security_group}",
					"ecs_security_group":      "${module.aws.ecs_cluster_security_group}",
					"engine":                  "mysql",
					"engine_version":          "5.6.17",
					"env":                     "production",
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
