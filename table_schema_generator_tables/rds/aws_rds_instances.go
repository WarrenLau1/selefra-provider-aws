package rds

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/aws/aws-sdk-go-v2/service/rds/types"
	"github.com/selefra/selefra-provider-aws/aws_client"
	"github.com/selefra/selefra-provider-sdk/table_schema_generator"
	"github.com/selefra/selefra-provider-sdk/provider/schema"
	"github.com/selefra/selefra-provider-sdk/provider/transformer/column_value_extractor"
)

type TableAwsRdsInstancesGenerator struct {
}

var _ table_schema_generator.TableSchemaGenerator = &TableAwsRdsInstancesGenerator{}

func (x *TableAwsRdsInstancesGenerator) GetTableName() string {
	return "aws_rds_instances"
}

func (x *TableAwsRdsInstancesGenerator) GetTableDescription() string {
	return ""
}

func (x *TableAwsRdsInstancesGenerator) GetVersion() uint64 {
	return 0
}

func (x *TableAwsRdsInstancesGenerator) GetOptions() *schema.TableOptions {
	return &schema.TableOptions{
		PrimaryKeys: []string{
			"arn",
		},
	}
}

func (x *TableAwsRdsInstancesGenerator) GetDataSource() *schema.DataSource {
	return &schema.DataSource{
		Pull: func(ctx context.Context, clientMeta *schema.ClientMeta, client any, task *schema.DataSourcePullTask, resultChannel chan<- any) *schema.Diagnostics {
			var config rds.DescribeDBInstancesInput
			c := client.(*aws_client.Client)
			svc := c.AwsServices().Rds
			for {
				response, err := svc.DescribeDBInstances(ctx, &config)
				if err != nil {
					return schema.NewDiagnosticsErrorPullTable(task.Table, err)

				}
				resultChannel <- response.DBInstances
				if aws.ToString(response.Marker) == "" {
					break
				}
				config.Marker = response.Marker
			}
			return nil
		},
	}
}

func (x *TableAwsRdsInstancesGenerator) GetExpandClientTask() func(ctx context.Context, clientMeta *schema.ClientMeta, client any, task *schema.DataSourcePullTask) []*schema.ClientTaskContext {
	return aws_client.ExpandByPartitionAndRegion("rds")
}

func (x *TableAwsRdsInstancesGenerator) GetColumns() []*schema.Column {
	return []*schema.Column{
		table_schema_generator.NewColumnBuilder().ColumnName("preferred_backup_window").ColumnType(schema.ColumnTypeString).
			Extractor(column_value_extractor.StructSelector("PreferredBackupWindow")).Build(),
		table_schema_generator.NewColumnBuilder().ColumnName("activity_stream_policy_status").ColumnType(schema.ColumnTypeString).
			Extractor(column_value_extractor.StructSelector("ActivityStreamPolicyStatus")).Build(),
		table_schema_generator.NewColumnBuilder().ColumnName("character_set_name").ColumnType(schema.ColumnTypeString).
			Extractor(column_value_extractor.StructSelector("CharacterSetName")).Build(),
		table_schema_generator.NewColumnBuilder().ColumnName("master_username").ColumnType(schema.ColumnTypeString).
			Extractor(column_value_extractor.StructSelector("MasterUsername")).Build(),
		table_schema_generator.NewColumnBuilder().ColumnName("publicly_accessible").ColumnType(schema.ColumnTypeBool).
			Extractor(column_value_extractor.StructSelector("PubliclyAccessible")).Build(),
		table_schema_generator.NewColumnBuilder().ColumnName("activity_stream_kinesis_stream_name").ColumnType(schema.ColumnTypeString).
			Extractor(column_value_extractor.StructSelector("ActivityStreamKinesisStreamName")).Build(),
		table_schema_generator.NewColumnBuilder().ColumnName("endpoint").ColumnType(schema.ColumnTypeJSON).
			Extractor(column_value_extractor.StructSelector("Endpoint")).Build(),
		table_schema_generator.NewColumnBuilder().ColumnName("processor_features").ColumnType(schema.ColumnTypeJSON).
			Extractor(column_value_extractor.WrapperExtractFunction(func(ctx context.Context, clientMeta *schema.ClientMeta, client any,
				task *schema.DataSourcePullTask, row *schema.Row, column *schema.Column, result any) (any, *schema.Diagnostics) {

				extractor := func() (any, error) {
					r := result.(types.DBInstance)
					processorFeatures := map[string]*string{}
					for _, t := range r.ProcessorFeatures {
						processorFeatures[*t.Name] = t.Value
					}
					return processorFeatures, nil
				}
				extractResultValue, err := extractor()
				if err != nil {
					return nil, schema.NewDiagnostics().AddErrorColumnValueExtractor(task.Table, column, err)
				} else {
					return extractResultValue, nil
				}
			})).Build(),
		table_schema_generator.NewColumnBuilder().ColumnName("read_replica_source_db_instance_identifier").ColumnType(schema.ColumnTypeString).
			Extractor(column_value_extractor.StructSelector("ReadReplicaSourceDBInstanceIdentifier")).Build(),
		table_schema_generator.NewColumnBuilder().ColumnName("selefra_id").ColumnType(schema.ColumnTypeString).SetUnique().Description("primary keys value md5").
			Extractor(column_value_extractor.PrimaryKeysID()).Build(),
		table_schema_generator.NewColumnBuilder().ColumnName("auto_minor_version_upgrade").ColumnType(schema.ColumnTypeBool).
			Extractor(column_value_extractor.StructSelector("AutoMinorVersionUpgrade")).Build(),
		table_schema_generator.NewColumnBuilder().ColumnName("deletion_protection").ColumnType(schema.ColumnTypeBool).
			Extractor(column_value_extractor.StructSelector("DeletionProtection")).Build(),
		table_schema_generator.NewColumnBuilder().ColumnName("aws_backup_recovery_point_arn").ColumnType(schema.ColumnTypeString).
			Extractor(column_value_extractor.StructSelector("AwsBackupRecoveryPointArn")).Build(),
		table_schema_generator.NewColumnBuilder().ColumnName("db_instance_class").ColumnType(schema.ColumnTypeString).
			Extractor(column_value_extractor.StructSelector("DBInstanceClass")).Build(),
		table_schema_generator.NewColumnBuilder().ColumnName("iam_database_authentication_enabled").ColumnType(schema.ColumnTypeBool).
			Extractor(column_value_extractor.StructSelector("IAMDatabaseAuthenticationEnabled")).Build(),
		table_schema_generator.NewColumnBuilder().ColumnName("storage_type").ColumnType(schema.ColumnTypeString).
			Extractor(column_value_extractor.StructSelector("StorageType")).Build(),
		table_schema_generator.NewColumnBuilder().ColumnName("activity_stream_kms_key_id").ColumnType(schema.ColumnTypeString).
			Extractor(column_value_extractor.StructSelector("ActivityStreamKmsKeyId")).Build(),
		table_schema_generator.NewColumnBuilder().ColumnName("automatic_restart_time").ColumnType(schema.ColumnTypeTimestamp).
			Extractor(column_value_extractor.StructSelector("AutomaticRestartTime")).Build(),
		table_schema_generator.NewColumnBuilder().ColumnName("db_instance_identifier").ColumnType(schema.ColumnTypeString).
			Extractor(column_value_extractor.StructSelector("DBInstanceIdentifier")).Build(),
		table_schema_generator.NewColumnBuilder().ColumnName("db_security_groups").ColumnType(schema.ColumnTypeJSON).
			Extractor(column_value_extractor.StructSelector("DBSecurityGroups")).Build(),
		table_schema_generator.NewColumnBuilder().ColumnName("db_system_id").ColumnType(schema.ColumnTypeString).
			Extractor(column_value_extractor.StructSelector("DBSystemId")).Build(),
		table_schema_generator.NewColumnBuilder().ColumnName("performance_insights_enabled").ColumnType(schema.ColumnTypeBool).
			Extractor(column_value_extractor.StructSelector("PerformanceInsightsEnabled")).Build(),
		table_schema_generator.NewColumnBuilder().ColumnName("backup_retention_period").ColumnType(schema.ColumnTypeBigInt).
			Extractor(column_value_extractor.StructSelector("BackupRetentionPeriod")).Build(),
		table_schema_generator.NewColumnBuilder().ColumnName("custom_iam_instance_profile").ColumnType(schema.ColumnTypeString).
			Extractor(column_value_extractor.StructSelector("CustomIamInstanceProfile")).Build(),
		table_schema_generator.NewColumnBuilder().ColumnName("enabled_cloudwatch_logs_exports").ColumnType(schema.ColumnTypeStringArray).
			Extractor(column_value_extractor.StructSelector("EnabledCloudwatchLogsExports")).Build(),
		table_schema_generator.NewColumnBuilder().ColumnName("instance_create_time").ColumnType(schema.ColumnTypeTimestamp).
			Extractor(column_value_extractor.StructSelector("InstanceCreateTime")).Build(),
		table_schema_generator.NewColumnBuilder().ColumnName("pending_modified_values").ColumnType(schema.ColumnTypeJSON).
			Extractor(column_value_extractor.StructSelector("PendingModifiedValues")).Build(),
		table_schema_generator.NewColumnBuilder().ColumnName("performance_insights_kms_key_id").ColumnType(schema.ColumnTypeString).
			Extractor(column_value_extractor.StructSelector("PerformanceInsightsKMSKeyId")).Build(),
		table_schema_generator.NewColumnBuilder().ColumnName("read_replica_db_cluster_identifiers").ColumnType(schema.ColumnTypeStringArray).
			Extractor(column_value_extractor.StructSelector("ReadReplicaDBClusterIdentifiers")).Build(),
		table_schema_generator.NewColumnBuilder().ColumnName("region").ColumnType(schema.ColumnTypeString).
			Extractor(aws_client.AwsRegionIDExtractor()).Build(),
		table_schema_generator.NewColumnBuilder().ColumnName("associated_roles").ColumnType(schema.ColumnTypeJSON).
			Extractor(column_value_extractor.StructSelector("AssociatedRoles")).Build(),
		table_schema_generator.NewColumnBuilder().ColumnName("db_instance_automated_backups_replications").ColumnType(schema.ColumnTypeJSON).
			Extractor(column_value_extractor.StructSelector("DBInstanceAutomatedBackupsReplications")).Build(),
		table_schema_generator.NewColumnBuilder().ColumnName("promotion_tier").ColumnType(schema.ColumnTypeBigInt).
			Extractor(column_value_extractor.StructSelector("PromotionTier")).Build(),
		table_schema_generator.NewColumnBuilder().ColumnName("tde_credential_arn").ColumnType(schema.ColumnTypeString).
			Extractor(column_value_extractor.StructSelector("TdeCredentialArn")).Build(),
		table_schema_generator.NewColumnBuilder().ColumnName("activity_stream_engine_native_audit_fields_included").ColumnType(schema.ColumnTypeBool).
			Extractor(column_value_extractor.StructSelector("ActivityStreamEngineNativeAuditFieldsIncluded")).Build(),
		table_schema_generator.NewColumnBuilder().ColumnName("option_group_memberships").ColumnType(schema.ColumnTypeJSON).
			Extractor(column_value_extractor.StructSelector("OptionGroupMemberships")).Build(),
		table_schema_generator.NewColumnBuilder().ColumnName("replica_mode").ColumnType(schema.ColumnTypeString).
			Extractor(column_value_extractor.StructSelector("ReplicaMode")).Build(),
		table_schema_generator.NewColumnBuilder().ColumnName("vpc_security_groups").ColumnType(schema.ColumnTypeJSON).
			Extractor(column_value_extractor.StructSelector("VpcSecurityGroups")).Build(),
		table_schema_generator.NewColumnBuilder().ColumnName("nchar_character_set_name").ColumnType(schema.ColumnTypeString).
			Extractor(column_value_extractor.StructSelector("NcharCharacterSetName")).Build(),
		table_schema_generator.NewColumnBuilder().ColumnName("preferred_maintenance_window").ColumnType(schema.ColumnTypeString).
			Extractor(column_value_extractor.StructSelector("PreferredMaintenanceWindow")).Build(),
		table_schema_generator.NewColumnBuilder().ColumnName("db_subnet_group").ColumnType(schema.ColumnTypeJSON).
			Extractor(column_value_extractor.StructSelector("DBSubnetGroup")).Build(),
		table_schema_generator.NewColumnBuilder().ColumnName("domain_memberships").ColumnType(schema.ColumnTypeJSON).
			Extractor(column_value_extractor.StructSelector("DomainMemberships")).Build(),
		table_schema_generator.NewColumnBuilder().ColumnName("arn").ColumnType(schema.ColumnTypeString).
			Extractor(column_value_extractor.StructSelector("DBInstanceArn")).Build(),
		table_schema_generator.NewColumnBuilder().ColumnName("activity_stream_status").ColumnType(schema.ColumnTypeString).
			Extractor(column_value_extractor.StructSelector("ActivityStreamStatus")).Build(),
		table_schema_generator.NewColumnBuilder().ColumnName("db_instance_arn").ColumnType(schema.ColumnTypeString).
			Extractor(column_value_extractor.StructSelector("DBInstanceArn")).Build(),
		table_schema_generator.NewColumnBuilder().ColumnName("db_parameter_groups").ColumnType(schema.ColumnTypeJSON).
			Extractor(column_value_extractor.StructSelector("DBParameterGroups")).Build(),
		table_schema_generator.NewColumnBuilder().ColumnName("engine").ColumnType(schema.ColumnTypeString).
			Extractor(column_value_extractor.StructSelector("Engine")).Build(),
		table_schema_generator.NewColumnBuilder().ColumnName("license_model").ColumnType(schema.ColumnTypeString).
			Extractor(column_value_extractor.StructSelector("LicenseModel")).Build(),
		table_schema_generator.NewColumnBuilder().ColumnName("multi_az").ColumnType(schema.ColumnTypeBool).
			Extractor(column_value_extractor.StructSelector("MultiAZ")).Build(),
		table_schema_generator.NewColumnBuilder().ColumnName("network_type").ColumnType(schema.ColumnTypeString).
			Extractor(column_value_extractor.StructSelector("NetworkType")).Build(),
		table_schema_generator.NewColumnBuilder().ColumnName("timezone").ColumnType(schema.ColumnTypeString).
			Extractor(column_value_extractor.StructSelector("Timezone")).Build(),
		table_schema_generator.NewColumnBuilder().ColumnName("copy_tags_to_snapshot").ColumnType(schema.ColumnTypeBool).
			Extractor(column_value_extractor.StructSelector("CopyTagsToSnapshot")).Build(),
		table_schema_generator.NewColumnBuilder().ColumnName("customer_owned_ip_enabled").ColumnType(schema.ColumnTypeBool).
			Extractor(column_value_extractor.StructSelector("CustomerOwnedIpEnabled")).Build(),
		table_schema_generator.NewColumnBuilder().ColumnName("account_id").ColumnType(schema.ColumnTypeString).
			Extractor(aws_client.AwsAccountIDExtractor()).Build(),
		table_schema_generator.NewColumnBuilder().ColumnName("dbi_resource_id").ColumnType(schema.ColumnTypeString).
			Extractor(column_value_extractor.StructSelector("DbiResourceId")).Build(),
		table_schema_generator.NewColumnBuilder().ColumnName("iops").ColumnType(schema.ColumnTypeBigInt).
			Extractor(column_value_extractor.StructSelector("Iops")).Build(),
		table_schema_generator.NewColumnBuilder().ColumnName("kms_key_id").ColumnType(schema.ColumnTypeString).
			Extractor(column_value_extractor.StructSelector("KmsKeyId")).Build(),
		table_schema_generator.NewColumnBuilder().ColumnName("read_replica_db_instance_identifiers").ColumnType(schema.ColumnTypeStringArray).
			Extractor(column_value_extractor.StructSelector("ReadReplicaDBInstanceIdentifiers")).Build(),
		table_schema_generator.NewColumnBuilder().ColumnName("resume_full_automation_mode_time").ColumnType(schema.ColumnTypeTimestamp).
			Extractor(column_value_extractor.StructSelector("ResumeFullAutomationModeTime")).Build(),
		table_schema_generator.NewColumnBuilder().ColumnName("storage_encrypted").ColumnType(schema.ColumnTypeBool).
			Extractor(column_value_extractor.StructSelector("StorageEncrypted")).Build(),
		table_schema_generator.NewColumnBuilder().ColumnName("automation_mode").ColumnType(schema.ColumnTypeString).
			Extractor(column_value_extractor.StructSelector("AutomationMode")).Build(),
		table_schema_generator.NewColumnBuilder().ColumnName("ca_certificate_identifier").ColumnType(schema.ColumnTypeString).
			Extractor(column_value_extractor.StructSelector("CACertificateIdentifier")).Build(),
		table_schema_generator.NewColumnBuilder().ColumnName("tag_list").ColumnType(schema.ColumnTypeJSON).
			Extractor(column_value_extractor.StructSelector("TagList")).Build(),
		table_schema_generator.NewColumnBuilder().ColumnName("db_cluster_identifier").ColumnType(schema.ColumnTypeString).
			Extractor(column_value_extractor.StructSelector("DBClusterIdentifier")).Build(),
		table_schema_generator.NewColumnBuilder().ColumnName("max_allocated_storage").ColumnType(schema.ColumnTypeBigInt).
			Extractor(column_value_extractor.StructSelector("MaxAllocatedStorage")).Build(),
		table_schema_generator.NewColumnBuilder().ColumnName("performance_insights_retention_period").ColumnType(schema.ColumnTypeBigInt).
			Extractor(column_value_extractor.StructSelector("PerformanceInsightsRetentionPeriod")).Build(),
		table_schema_generator.NewColumnBuilder().ColumnName("tags").ColumnType(schema.ColumnTypeJSON).Build(),
		table_schema_generator.NewColumnBuilder().ColumnName("availability_zone").ColumnType(schema.ColumnTypeString).
			Extractor(column_value_extractor.StructSelector("AvailabilityZone")).Build(),
		table_schema_generator.NewColumnBuilder().ColumnName("backup_target").ColumnType(schema.ColumnTypeString).
			Extractor(column_value_extractor.StructSelector("BackupTarget")).Build(),
		table_schema_generator.NewColumnBuilder().ColumnName("monitoring_role_arn").ColumnType(schema.ColumnTypeString).
			Extractor(column_value_extractor.StructSelector("MonitoringRoleArn")).Build(),
		table_schema_generator.NewColumnBuilder().ColumnName("allocated_storage").ColumnType(schema.ColumnTypeBigInt).
			Extractor(column_value_extractor.StructSelector("AllocatedStorage")).Build(),
		table_schema_generator.NewColumnBuilder().ColumnName("db_instance_port").ColumnType(schema.ColumnTypeBigInt).
			Extractor(column_value_extractor.StructSelector("DbInstancePort")).Build(),
		table_schema_generator.NewColumnBuilder().ColumnName("latest_restorable_time").ColumnType(schema.ColumnTypeTimestamp).
			Extractor(column_value_extractor.StructSelector("LatestRestorableTime")).Build(),
		table_schema_generator.NewColumnBuilder().ColumnName("certificate_details").ColumnType(schema.ColumnTypeJSON).
			Extractor(column_value_extractor.StructSelector("CertificateDetails")).Build(),
		table_schema_generator.NewColumnBuilder().ColumnName("engine_version").ColumnType(schema.ColumnTypeString).
			Extractor(column_value_extractor.StructSelector("EngineVersion")).Build(),
		table_schema_generator.NewColumnBuilder().ColumnName("listener_endpoint").ColumnType(schema.ColumnTypeJSON).
			Extractor(column_value_extractor.StructSelector("ListenerEndpoint")).Build(),
		table_schema_generator.NewColumnBuilder().ColumnName("monitoring_interval").ColumnType(schema.ColumnTypeBigInt).
			Extractor(column_value_extractor.StructSelector("MonitoringInterval")).Build(),
		table_schema_generator.NewColumnBuilder().ColumnName("db_instance_status").ColumnType(schema.ColumnTypeString).
			Extractor(column_value_extractor.StructSelector("DBInstanceStatus")).Build(),
		table_schema_generator.NewColumnBuilder().ColumnName("db_name").ColumnType(schema.ColumnTypeString).
			Extractor(column_value_extractor.StructSelector("DBName")).Build(),
		table_schema_generator.NewColumnBuilder().ColumnName("master_user_secret").ColumnType(schema.ColumnTypeJSON).
			Extractor(column_value_extractor.StructSelector("MasterUserSecret")).Build(),
		table_schema_generator.NewColumnBuilder().ColumnName("secondary_availability_zone").ColumnType(schema.ColumnTypeString).
			Extractor(column_value_extractor.StructSelector("SecondaryAvailabilityZone")).Build(),
		table_schema_generator.NewColumnBuilder().ColumnName("status_infos").ColumnType(schema.ColumnTypeJSON).
			Extractor(column_value_extractor.StructSelector("StatusInfos")).Build(),
		table_schema_generator.NewColumnBuilder().ColumnName("storage_throughput").ColumnType(schema.ColumnTypeBigInt).
			Extractor(column_value_extractor.StructSelector("StorageThroughput")).Build(),
		table_schema_generator.NewColumnBuilder().ColumnName("activity_stream_mode").ColumnType(schema.ColumnTypeString).
			Extractor(column_value_extractor.StructSelector("ActivityStreamMode")).Build(),
		table_schema_generator.NewColumnBuilder().ColumnName("enhanced_monitoring_resource_arn").ColumnType(schema.ColumnTypeString).
			Extractor(column_value_extractor.StructSelector("EnhancedMonitoringResourceArn")).Build(),
	}
}

func (x *TableAwsRdsInstancesGenerator) GetSubTables() []*schema.Table {
	return nil
}
