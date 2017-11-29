package terraform_test

import (
	"io/ioutil"

	"github.com/Originate/exosphere/src/config"
	"github.com/Originate/exosphere/src/terraform"
	"github.com/Originate/exosphere/src/types"
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
			Production: types.AppProductionConfig{
				URL: "example-app.com",
			},
		}
		serviceConfigs := map[string]types.ServiceConfig{}

		deployConfig := types.DeployConfig{
			AppContext: &types.AppContext{
				Config: appConfig,
			},
			ServiceConfigs: serviceConfigs,
			AwsConfig: types.AwsConfig{
				TerraformStateBucket: "example-app-terraform",
				TerraformLockTable:   "TerraformLocks",
				Region:               "us-west-2",
				AccountID:            "12345",
			},
			TerraformModulesRef: "TERRAFORM_MODULES_REF",
		}

		It("should generate an AWS module only", func() {
			result, err := terraform.Generate(deployConfig)
			Expect(err).To(BeNil())
			hclFile, err := hcl.GetHCLFileFromTerraform(result)
			Expect(err).To(BeNil())
			Expect(hclFile.GetModuleNames()).To(Equal([]string{"aws"}))
			Expect(hclFile.Module["aws"]).To(Equal(hcl.Module{
				"source":            "git@github.com:Originate/exosphere.git//terraform//aws?ref=TERRAFORM_MODULES_REF",
				"key_name":          "${var.key_name}",
				"name":              "example-app",
				"env":               "production",
				"external_dns_name": "example-app.com",
			}))
		})
	})

	var _ = Describe("Given an application with public and worker services", func() {
		var hclFile *hcl.File
		appConfig := types.AppConfig{
			Name: "example-app",
			Services: map[string]types.ServiceData{
				"public-service": types.ServiceData{},
				"worker-service": types.ServiceData{},
			},
		}
		serviceConfigs := map[string]types.ServiceConfig{
			"public-service": {
				Type: "public",
				Production: types.ServiceProductionConfig{
					Port:        "3000",
					CPU:         "128",
					URL:         "originate.com",
					HealthCheck: "/health-check",
					Memory:      "128",
				},
			},
			"worker-service": {
				Type: "worker",
				Production: types.ServiceProductionConfig{
					CPU:    "128",
					Memory: "128",
				},
			},
		}

		deployConfig := types.DeployConfig{
			AppContext: &types.AppContext{
				Config: appConfig,
			},
			ServiceConfigs: serviceConfigs,
			AwsConfig: types.AwsConfig{
				SslCertificateArn: "sslcert123",
			},
			TerraformModulesRef: "TERRAFORM_MODULES_REF",
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
			Expect(hclFile.Module["public-service"]).To(Equal(hcl.Module{
				"source":             "git@github.com:Originate/exosphere.git//terraform//aws//public-service?ref=TERRAFORM_MODULES_REF",
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
				"external_dns_name":     "originate.com",
				"external_zone_id":      "${module.aws.external_zone_id}",
				"health_check_endpoint": "/health-check",
				"internal_dns_name":     "public-service",
				"internal_zone_id":      "${module.aws.internal_zone_id}",
				"log_bucket":            "${module.aws.log_bucket_id}",
				"memory_reservation":    "128",
				"name":                  "public-service",
				"region":                "${module.aws.region}",
				"ssl_certificate_arn":   "sslcert123",
				"vpc_id":                "${module.aws.vpc_id}",
			}))
		})

		It("should generate a worker service module", func() {
			Expect(hclFile).To(matchers.HaveHCLVariable("worker-service_env_vars"))
			Expect(hclFile).To(matchers.HaveHCLVariable("worker-service_docker_image"))
			Expect(hclFile.Module["worker-service"]).To(Equal(hcl.Module{
				"source":        "git@github.com:Originate/exosphere.git//terraform//aws//worker-service?ref=TERRAFORM_MODULES_REF",
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
			appContext, err := types.GetAppContext(appDir)
			Expect(err).NotTo(HaveOccurred())
			serviceConfigs, err := config.GetServiceConfigs(appContext.Location, appContext.Config)
			Expect(err).NotTo(HaveOccurred())

			deployConfig := types.DeployConfig{
				AppContext:          appContext,
				ServiceConfigs:      serviceConfigs,
				TerraformModulesRef: "TERRAFORM_MODULES_REF",
			}
			result, err := terraform.Generate(deployConfig)
			Expect(err).To(BeNil())
			hclFile, err := hcl.GetHCLFileFromTerraform(result)
			Expect(err).To(BeNil())
			Expect(hclFile).To(matchers.HaveHCLVariable("exocom_env_vars"))
			Expect(hclFile.Module["exocom_cluster"]).To(Equal(hcl.Module{
				"source":                      "git@github.com:Originate/exosphere.git//terraform//aws//dependencies//exocom//exocom-cluster?ref=TERRAFORM_MODULES_REF",
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
				"source":       "git@github.com:Originate/exosphere.git//terraform//aws//dependencies//exocom//exocom-service?ref=TERRAFORM_MODULES_REF",
				"cluster_id":   "${module.exocom_cluster.cluster_id}",
				"cpu_units":    "128",
				"docker_image": "${var.exocom_docker_image}",
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
			appContext, err := types.GetAppContext(appDir)
			Expect(err).NotTo(HaveOccurred())
			serviceConfigs, err := config.GetServiceConfigs(appContext.Location, appContext.Config)
			Expect(err).NotTo(HaveOccurred())

			deployConfig := types.DeployConfig{
				AppContext:          appContext,
				ServiceConfigs:      serviceConfigs,
				TerraformModulesRef: "TERRAFORM_MODULES_REF",
			}
			result, err := terraform.Generate(deployConfig)
			Expect(err).To(BeNil())
			hclFile, err := hcl.GetHCLFileFromTerraform(result)
			Expect(err).To(BeNil())
			By("generating rds modules for application dependencies", func() {
				Expect(hclFile.Module["my-db_rds_instance"]).To(Equal(hcl.Module{
					"source":                  "git@github.com:Originate/exosphere.git//terraform//aws//dependencies//rds?ref=TERRAFORM_MODULES_REF",
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
					"source":                  "git@github.com:Originate/exosphere.git//terraform//aws//dependencies//rds?ref=TERRAFORM_MODULES_REF",
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
