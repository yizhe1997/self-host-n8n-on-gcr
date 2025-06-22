package main

import (
	"github.com/pulumi/pulumi-gcp/sdk/v8/go/gcp/artifactregistry"
	"github.com/pulumi/pulumi-gcp/sdk/v8/go/gcp/cloudrunv2"
	"github.com/pulumi/pulumi-gcp/sdk/v8/go/gcp/projects"
	"github.com/pulumi/pulumi-gcp/sdk/v8/go/gcp/secretmanager"
	"github.com/pulumi/pulumi-gcp/sdk/v8/go/gcp/serviceaccount"
	"github.com/pulumi/pulumi-gcp/sdk/v8/go/gcp/sql"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// https://www.pulumi.com/registry/packages/gcp/api-docs/sql/databaseinstance/#inputs
		_, err := sql.NewDatabaseInstance(ctx, "dbInstance", &sql.DatabaseInstanceArgs{
			DatabaseVersion:    pulumi.String("POSTGRES_13"),
			InstanceType:       pulumi.String("CLOUD_SQL_INSTANCE"),
			MaintenanceVersion: pulumi.String("POSTGRES_13_21.R20250302.00_31"),
			Name:               pulumi.String("n8n-db"),
			Project:            pulumi.String("sunlit-alloy-458006-q2"),
			Region:             pulumi.String("asia-southeast1"),
			Settings: &sql.DatabaseInstanceSettingsArgs{
				BackupConfiguration: &sql.DatabaseInstanceSettingsBackupConfigurationArgs{
					BackupRetentionSettings: &sql.DatabaseInstanceSettingsBackupConfigurationBackupRetentionSettingsArgs{
						RetainedBackups: pulumi.Int(7),
					},
					StartTime:                   pulumi.String("00:00"),
					TransactionLogRetentionDays: pulumi.Int(7),
				},
				ConnectorEnforcement: pulumi.String("NOT_REQUIRED"),
				DiskSize:             pulumi.Int(10),
				DiskType:             pulumi.String("PD_HDD"),
				IpConfiguration: &sql.DatabaseInstanceSettingsIpConfigurationArgs{
					AuthorizedNetworks: sql.DatabaseInstanceSettingsIpConfigurationAuthorizedNetworkArray{
						&sql.DatabaseInstanceSettingsIpConfigurationAuthorizedNetworkArgs{
							Name:  pulumi.String("AllowAll"),
							Value: pulumi.String("0.0.0.0/0"),
						},
					},
					ServerCaMode: pulumi.String("GOOGLE_MANAGED_INTERNAL_CA"),
					SslMode:      pulumi.String("ALLOW_UNENCRYPTED_AND_ENCRYPTED"),
				},
				LocationPreference: &sql.DatabaseInstanceSettingsLocationPreferenceArgs{
					Zone: pulumi.String("asia-southeast1-c"),
				},
				Tier: pulumi.String("db-f1-micro"),
			},
		})
		if err != nil {
			return err
		}

		// https://www.pulumi.com/registry/packages/gcp/api-docs/sql/database/#inputs
		_, err = sql.NewDatabase(ctx, "db", &sql.DatabaseArgs{
			Charset:   pulumi.String("UTF8"),
			Collation: pulumi.String("en_US.UTF8"),
			Instance:  pulumi.String("n8n-db"),
			Name:      pulumi.String("n8n"),
			Project:   pulumi.String("sunlit-alloy-458006-q2"),
		})
		if err != nil {
			return err
		}

		// https://www.pulumi.com/registry/packages/gcp/api-docs/sql/user/
		_, err = sql.NewUser(ctx, "n8n-user", &sql.UserArgs{
			Instance: pulumi.String("n8n-db"),
			Name:     pulumi.String("n8n-user"),
			Project:  pulumi.String("sunlit-alloy-458006-q2"),
		})
		if err != nil {
			return err
		}

		//https://www.pulumi.com/registry/packages/gcp/api-docs/projects/service/
		_, err = projects.NewService(ctx, "api-artifactregistry", &projects.ServiceArgs{
			Project: pulumi.String("sunlit-alloy-458006-q2"),
			Service: pulumi.String("artifactregistry.googleapis.com"),
		})
		if err != nil {
			return err
		}
		_, err = projects.NewService(ctx, "api-googleapis", &projects.ServiceArgs{
			Project: pulumi.String("sunlit-alloy-458006-q2"),
			Service: pulumi.String("run.googleapis.com"),
		})
		if err != nil {
			return err
		}
		_, err = projects.NewService(ctx, "api-sqladmin", &projects.ServiceArgs{
			Project: pulumi.String("sunlit-alloy-458006-q2"),
			Service: pulumi.String("sqladmin.googleapis.com"),
		})
		if err != nil {
			return err
		}
		_, err = projects.NewService(ctx, "api-secretmanager", &projects.ServiceArgs{
			Project: pulumi.String("sunlit-alloy-458006-q2"),
			Service: pulumi.String("secretmanager.googleapis.com"),
		})
		if err != nil {
			return err
		}

		// https://www.pulumi.com/registry/packages/gcp/api-docs/artifactregistry/repository/
		_, err = artifactregistry.NewRepository(ctx, "artifactregistry", &artifactregistry.RepositoryArgs{
			Description:  pulumi.String("Repository for n8n workflow images"),
			Format:       pulumi.String("DOCKER"),
			Location:     pulumi.String("asia-southeast1"),
			Project:      pulumi.String("sunlit-alloy-458006-q2"),
			RepositoryId: pulumi.String("n8n-repo"),
		})
		if err != nil {
			return err
		}

		// https://www.pulumi.com/registry/packages/gcp/api-docs/secretmanager/secret/
		_, err = secretmanager.NewSecret(ctx, "secret-n8n-db-password", &secretmanager.SecretArgs{
			Project: pulumi.String("sunlit-alloy-458006-q2"),
			Replication: &secretmanager.SecretReplicationArgs{
				Auto: &secretmanager.SecretReplicationAutoArgs{},
			},
			SecretId: pulumi.String("n8n-db-password"),
		})
		if err != nil {
			return err
		}
		_, err = secretmanager.NewSecret(ctx, "secret-n8n-encryption-key", &secretmanager.SecretArgs{
			Project: pulumi.String("sunlit-alloy-458006-q2"),
			Replication: &secretmanager.SecretReplicationArgs{
				Auto: &secretmanager.SecretReplicationAutoArgs{},
			},
			SecretId: pulumi.String("n8n-encryption-key"),
		})
		if err != nil {
			return err
		}

		// https://www.pulumi.com/registry/packages/gcp/api-docs/serviceaccount/account/
		_, err = serviceaccount.NewAccount(ctx, "n8n-service-account", &serviceaccount.AccountArgs{
			AccountId:   pulumi.String("n8n-service-account"),
			DisplayName: pulumi.String("n8n Service Account"),
			Project:     pulumi.String("sunlit-alloy-458006-q2"),
		})
		if err != nil {
			return err
		}

		// https://www.pulumi.com/registry/packages/gcp/api-docs/secretmanager/secretiambinding/
		_, err = secretmanager.NewSecretIamBinding(ctx, "service-acc-iam-binding-n8n-db-password", &secretmanager.SecretIamBindingArgs{
			Members: pulumi.StringArray{
				pulumi.String("serviceAccount:n8n-service-account@sunlit-alloy-458006-q2.iam.gserviceaccount.com"),
			},
			Project:  pulumi.String("sunlit-alloy-458006-q2"),
			Role:     pulumi.String("roles/secretmanager.secretAccessor"),
			SecretId: pulumi.String("projects/sunlit-alloy-458006-q2/secrets/n8n-db-password"),
		})
		if err != nil {
			return err
		}
		_, err = secretmanager.NewSecretIamBinding(ctx, "service-acc-iam-binding-n8n-encryption-key", &secretmanager.SecretIamBindingArgs{
			Members: pulumi.StringArray{
				pulumi.String("serviceAccount:n8n-service-account@sunlit-alloy-458006-q2.iam.gserviceaccount.com"),
			},
			Project:  pulumi.String("sunlit-alloy-458006-q2"),
			Role:     pulumi.String("roles/secretmanager.secretAccessor"),
			SecretId: pulumi.String("projects/sunlit-alloy-458006-q2/secrets/n8n-encryption-key"),
		})
		if err != nil {
			return err
		}

		// https://www.pulumi.com/registry/packages/gcp/api-docs/serviceaccount/iambinding/
		_, err = serviceaccount.NewIAMBinding(ctx, "service-acc-iam-binding-cloudsql", &serviceaccount.IAMBindingArgs{
			Members:          pulumi.StringArray{},
			Role:             pulumi.String("roles/cloudsql.client"),
			ServiceAccountId: pulumi.String("projects/sunlit-alloy-458006-q2/serviceAccounts/n8n-service-account@sunlit-alloy-458006-q2.iam.gserviceaccount.com"),
		})
		if err != nil {
			return err
		}

		// https://www.pulumi.com/registry/packages/gcp/api-docs/cloudrunv2/service/
		_, err = cloudrunv2.NewService(ctx, "n8n", &cloudrunv2.ServiceArgs{
			Client:      pulumi.String("cloud-console"),
			Ingress:     pulumi.String("INGRESS_TRAFFIC_ALL"),
			LaunchStage: pulumi.String("GA"),
			Location:    pulumi.String("asia-southeast1"),
			Name:        pulumi.String("n8n"),
			Project:     pulumi.String("sunlit-alloy-458006-q2"),
			Scaling: &cloudrunv2.ServiceScalingArgs{
				MinInstanceCount: pulumi.Int(1),
			},
			Template: &cloudrunv2.ServiceTemplateArgs{
				Containers: cloudrunv2.ServiceTemplateContainerArray{
					&cloudrunv2.ServiceTemplateContainerArgs{
						Envs: cloudrunv2.ServiceTemplateContainerEnvArray{
							&cloudrunv2.ServiceTemplateContainerEnvArgs{
								Name:  pulumi.String("DB_POSTGRESDB_DATABASE"),
								Value: pulumi.String("n8n"),
							},
							&cloudrunv2.ServiceTemplateContainerEnvArgs{
								Name:  pulumi.String("DB_POSTGRESDB_HOST"),
								Value: pulumi.String("/cloudsql/sunlit-alloy-458006-q2:asia-southeast1:n8n-db"),
							},
							&cloudrunv2.ServiceTemplateContainerEnvArgs{
								Name: pulumi.String("DB_POSTGRESDB_PASSWORD"),
								ValueSource: &cloudrunv2.ServiceTemplateContainerEnvValueSourceArgs{
									SecretKeyRef: &cloudrunv2.ServiceTemplateContainerEnvValueSourceSecretKeyRefArgs{
										Secret:  pulumi.String("n8n-db-password"),
										Version: pulumi.String("latest"),
									},
								},
							},
							&cloudrunv2.ServiceTemplateContainerEnvArgs{
								Name:  pulumi.String("DB_POSTGRESDB_PORT"),
								Value: pulumi.String("5432"),
							},
							&cloudrunv2.ServiceTemplateContainerEnvArgs{
								Name:  pulumi.String("DB_POSTGRESDB_SCHEMA"),
								Value: pulumi.String("public"),
							},
							&cloudrunv2.ServiceTemplateContainerEnvArgs{
								Name:  pulumi.String("DB_POSTGRESDB_USER"),
								Value: pulumi.String("n8n-user"),
							},
							&cloudrunv2.ServiceTemplateContainerEnvArgs{
								Name:  pulumi.String("DB_TYPE"),
								Value: pulumi.String("postgresdb"),
							},
							&cloudrunv2.ServiceTemplateContainerEnvArgs{
								Name:  pulumi.String("GENERIC_TIMEZONE"),
								Value: pulumi.String("UTC"),
							},
							&cloudrunv2.ServiceTemplateContainerEnvArgs{
								Name:  pulumi.String("N8N_EDITOR_BASE_URL"),
								Value: pulumi.String("https://n8n-197184456803.asia-southeast1.run.app"),
							},
							&cloudrunv2.ServiceTemplateContainerEnvArgs{
								Name: pulumi.String("N8N_ENCRYPTION_KEY"),
								ValueSource: &cloudrunv2.ServiceTemplateContainerEnvValueSourceArgs{
									SecretKeyRef: &cloudrunv2.ServiceTemplateContainerEnvValueSourceSecretKeyRefArgs{
										Secret:  pulumi.String("n8n-encryption-key"),
										Version: pulumi.String("latest"),
									},
								},
							},
							&cloudrunv2.ServiceTemplateContainerEnvArgs{
								Name:  pulumi.String("N8N_HOST"),
								Value: pulumi.String("n8n-197184456803.asia-southeast1.run.app"),
							},
							&cloudrunv2.ServiceTemplateContainerEnvArgs{
								Name:  pulumi.String("N8N_PATH"),
								Value: pulumi.String("/"),
							},
							&cloudrunv2.ServiceTemplateContainerEnvArgs{
								Name:  pulumi.String("N8N_PORT"),
								Value: pulumi.String("443"),
							},
							&cloudrunv2.ServiceTemplateContainerEnvArgs{
								Name:  pulumi.String("N8N_PROTOCOL"),
								Value: pulumi.String("https"),
							},
							&cloudrunv2.ServiceTemplateContainerEnvArgs{
								Name:  pulumi.String("N8N_USER_FOLDER"),
								Value: pulumi.String("/home/node/.n8n"),
							},
							&cloudrunv2.ServiceTemplateContainerEnvArgs{
								Name:  pulumi.String("N8N_WEBHOOK_URL"),
								Value: pulumi.String("https://n8n-197184456803.asia-southeast1.run.app"),
							},
							&cloudrunv2.ServiceTemplateContainerEnvArgs{
								Name:  pulumi.String("QUEUE_HEALTH_CHECK_ACTIVE"),
								Value: pulumi.String("true"),
							},
							&cloudrunv2.ServiceTemplateContainerEnvArgs{
								Name:  pulumi.String("WEBHOOK_URL"),
								Value: pulumi.String("https://n8n-197184456803.asia-southeast1.run.app"),
							},
						},
						Image: pulumi.String("asia-southeast1-docker.pkg.dev/sunlit-alloy-458006-q2/n8n-repo/n8n:latest"),
						Name:  pulumi.String("n8n-1"),
						Ports: &cloudrunv2.ServiceTemplateContainerPortsArgs{
							ContainerPort: pulumi.Int(5678),
							Name:          pulumi.String("http1"),
						},
						Resources: &cloudrunv2.ServiceTemplateContainerResourcesArgs{
							CpuIdle: pulumi.Bool(true),
							Limits: pulumi.StringMap{
								"cpu":    pulumi.String("1"),
								"memory": pulumi.String("2Gi"),
							},
							StartupCpuBoost: pulumi.Bool(true),
						},
						StartupProbe: &cloudrunv2.ServiceTemplateContainerStartupProbeArgs{
							FailureThreshold: pulumi.Int(1),
							PeriodSeconds:    pulumi.Int(240),
							TcpSocket: &cloudrunv2.ServiceTemplateContainerStartupProbeTcpSocketArgs{
								Port: pulumi.Int(5678),
							},
							TimeoutSeconds: pulumi.Int(240),
						},
						VolumeMounts: cloudrunv2.ServiceTemplateContainerVolumeMountArray{
							&cloudrunv2.ServiceTemplateContainerVolumeMountArgs{
								MountPath: pulumi.String("/cloudsql"),
								Name:      pulumi.String("cloudsql"),
							},
						},
					},
				},
				MaxInstanceRequestConcurrency: pulumi.Int(80),
				Scaling: &cloudrunv2.ServiceTemplateScalingArgs{
					MaxInstanceCount: pulumi.Int(1),
				},
				ServiceAccount: pulumi.String("n8n-service-account@sunlit-alloy-458006-q2.iam.gserviceaccount.com"),
				Timeout:        pulumi.String("300s"),
				Volumes: cloudrunv2.ServiceTemplateVolumeArray{
					&cloudrunv2.ServiceTemplateVolumeArgs{
						CloudSqlInstance: &cloudrunv2.ServiceTemplateVolumeCloudSqlInstanceArgs{
							Instances: pulumi.StringArray{
								pulumi.String("sunlit-alloy-458006-q2:asia-southeast1:n8n-db"),
							},
						},
						Name: pulumi.String("cloudsql"),
					},
				},
			},
			Traffics: cloudrunv2.ServiceTrafficArray{
				&cloudrunv2.ServiceTrafficArgs{
					Percent: pulumi.Int(100),
					Type:    pulumi.String("TRAFFIC_TARGET_ALLOCATION_TYPE_LATEST"),
				},
			},
		})
		if err != nil {
			return err
		}

		return nil
	})
}
