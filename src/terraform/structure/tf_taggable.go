package structure

var TfTaggableResourceTypes = []string{
	"aws_accessanalyzer_analyzer",
	"aws_acm_certificate",
	"aws_acmpca_certificate_authority",
	"aws_alb",
	"aws_alb_listener",
	"aws_alb_listener_rule",
	"aws_alb_target_group",
	"aws_ami",
	"aws_ami_copy",
	"aws_ami_from_instance",
	"aws_amplify_app",
	"aws_amplify_branch",
	"aws_api_gateway_api_key",
	"aws_api_gateway_client_certificate",
	"aws_api_gateway_domain_name",
	"aws_api_gateway_rest_api",
	"aws_api_gateway_stage",
	"aws_api_gateway_usage_plan",
	"aws_api_gateway_vpc_link",
	"aws_apigatewayv2_api",
	"aws_apigatewayv2_domain_name",
	"aws_apigatewayv2_stage",
	"aws_apigatewayv2_vpc_link",
	"aws_appconfig_application",
	"aws_appconfig_configuration_profile",
	"aws_appconfig_deployment",
	"aws_appconfig_deployment_strategy",
	"aws_appconfig_environment",
	"aws_appmesh_gateway_route",
	"aws_appmesh_mesh",
	"aws_appmesh_route",
	"aws_appmesh_virtual_gateway",
	"aws_appmesh_virtual_node",
	"aws_appmesh_virtual_router",
	"aws_appmesh_virtual_service",
	"aws_apprunner_auto_scaling_configuration_version",
	"aws_apprunner_connection",
	"aws_apprunner_service",
	"aws_appsync_graphql_api",
	"aws_athena_workgroup",
	"aws_autoscaling_group",
	"aws_backup_plan",
	"aws_backup_vault",
	"aws_batch_compute_environment",
	"aws_batch_job_definition",
	"aws_batch_job_queue",
	"aws_cloud9_environment_ec2",
	"aws_cloudformation_stack",
	"aws_cloudformation_stack_set",
	"aws_cloudfront_distribution",
	"aws_cloudhsm_v2_cluster",
	"aws_cloudtrail",
	"aws_cloudwatch_composite_alarm",
	"aws_cloudwatch_event_bus",
	"aws_cloudwatch_event_rule",
	"aws_cloudwatch_log_group",
	"aws_cloudwatch_metric_alarm",
	"aws_cloudwatch_metric_stream",
	"aws_codeartifact_domain",
	"aws_codeartifact_repository",
	"aws_codebuild_project",
	"aws_codebuild_report_group",
	"aws_codecommit_repository",
	"aws_codedeploy_app",
	"aws_codedeploy_deployment_group",
	"aws_codepipeline",
	"aws_codepipeline_webhook",
	"aws_codestarconnections_connection",
	"aws_codestarnotifications_notification_rule",
	"aws_cognito_identity_pool",
	"aws_cognito_user_pool",
	"aws_config_aggregate_authorization",
	"aws_config_config_rule",
	"aws_config_configuration_aggregator",
	"aws_customer_gateway",
	"aws_datapipeline_pipeline",
	"aws_datasync_agent",
	"aws_datasync_location_efs",
	"aws_datasync_location_fsx_windows_file_system",
	"aws_datasync_location_nfs",
	"aws_datasync_location_s3",
	"aws_datasync_location_smb",
	"aws_datasync_task",
	"aws_dax_cluster",
	"aws_db_cluster_snapshot",
	"aws_db_event_subscription",
	"aws_db_instance",
	"aws_db_option_group",
	"aws_db_parameter_group",
	"aws_db_proxy",
	"aws_db_proxy_endpoint",
	"aws_db_security_group",
	"aws_db_snapshot",
	"aws_db_subnet_group",
	"aws_default_network_acl",
	"aws_default_route_table",
	"aws_default_security_group",
	"aws_default_subnet",
	"aws_default_vpc",
	"aws_default_vpc_dhcp_options",
	"aws_devicefarm_project",
	"aws_directory_service_directory",
	"aws_dlm_lifecycle_policy",
	"aws_dms_certificate",
	"aws_dms_endpoint",
	"aws_dms_event_subscription",
	"aws_dms_replication_instance",
	"aws_dms_replication_subnet_group",
	"aws_dms_replication_task",
	"aws_docdb_cluster",
	"aws_docdb_cluster_instance",
	"aws_docdb_cluster_parameter_group",
	"aws_docdb_subnet_group",
	"aws_dx_connection",
	"aws_dx_hosted_private_virtual_interface_accepter",
	"aws_dx_hosted_public_virtual_interface_accepter",
	"aws_dx_hosted_transit_virtual_interface_accepter",
	"aws_dx_lag",
	"aws_dx_private_virtual_interface",
	"aws_dx_public_virtual_interface",
	"aws_dx_transit_virtual_interface",
	"aws_dynamodb_table",
	"aws_ebs_snapshot",
	"aws_ebs_snapshot_copy",
	"aws_ebs_snapshot_import",
	"aws_ebs_volume",
	"aws_ec2_capacity_reservation",
	"aws_ec2_carrier_gateway",
	"aws_ec2_client_vpn_endpoint",
	"aws_ec2_fleet",
	"aws_ec2_local_gateway_route_table_vpc_association",
	"aws_ec2_managed_prefix_list",
	"aws_ec2_traffic_mirror_filter",
	"aws_ec2_traffic_mirror_session",
	"aws_ec2_traffic_mirror_target",
	"aws_ec2_transit_gateway",
	"aws_ec2_transit_gateway_peering_attachment",
	"aws_ec2_transit_gateway_peering_attachment_accepter",
	"aws_ec2_transit_gateway_route_table",
	"aws_ec2_transit_gateway_vpc_attachment",
	"aws_ec2_transit_gateway_vpc_attachment_accepter",
	"aws_ecr_repository",
	"aws_ecs_capacity_provider",
	"aws_ecs_cluster",
	"aws_ecs_service",
	"aws_ecs_task_definition",
	"aws_efs_access_point",
	"aws_efs_file_system",
	"aws_egress_only_internet_gateway",
	"aws_eip",
	"aws_eks_addon",
	"aws_eks_cluster",
	"aws_eks_fargate_profile",
	"aws_eks_identity_provider_config",
	"aws_eks_node_group",
	"aws_elastic_beanstalk_application",
	"aws_elastic_beanstalk_application_version",
	"aws_elastic_beanstalk_environment",
	"aws_elasticache_cluster",
	"aws_elasticache_parameter_group",
	"aws_elasticache_replication_group",
	"aws_elasticache_subnet_group",
	"aws_elasticache_user",
	"aws_elasticache_user_group",
	"aws_elasticsearch_domain",
	"aws_elb",
	"aws_emr_cluster",
	"aws_flow_log",
	"aws_fsx_backup",
	"aws_fsx_lustre_file_system",
	"aws_fsx_windows_file_system",
	"aws_gamelift_alias",
	"aws_gamelift_build",
	"aws_gamelift_fleet",
	"aws_gamelift_game_session_queue",
	"aws_glacier_vault",
	"aws_globalaccelerator_accelerator",
	"aws_glue_crawler",
	"aws_glue_dev_endpoint",
	"aws_glue_job",
	"aws_glue_ml_transform",
	"aws_glue_registry",
	"aws_glue_schema",
	"aws_glue_trigger",
	"aws_glue_workflow",
	"aws_guardduty_detector",
	"aws_guardduty_filter",
	"aws_guardduty_ipset",
	"aws_guardduty_threatintelset",
	"aws_iam_instance_profile",
	"aws_iam_openid_connect_provider",
	"aws_iam_policy",
	"aws_iam_role",
	"aws_iam_saml_provider",
	"aws_iam_server_certificate",
	"aws_iam_user",
	"aws_imagebuilder_component",
	"aws_imagebuilder_distribution_configuration",
	"aws_imagebuilder_image",
	"aws_imagebuilder_image_pipeline",
	"aws_imagebuilder_image_recipe",
	"aws_imagebuilder_infrastructure_configuration",
	"aws_inspector_assessment_template",
	"aws_inspector_resource_group",
	"aws_instance",
	"aws_internet_gateway",
	"aws_iot_topic_rule",
	"aws_key_pair",
	"aws_kinesis_analytics_application",
	"aws_kinesis_firehose_delivery_stream",
	"aws_kinesis_stream",
	"aws_kinesis_video_stream",
	"aws_kinesisanalyticsv2_application",
	"aws_kms_external_key",
	"aws_kms_key",
	"aws_lambda_function",
	"aws_launch_template",
	"aws_lb",
	"aws_lb_listener",
	"aws_lb_listener_rule",
	"aws_lb_target_group",
	"aws_licensemanager_license_configuration",
	"aws_lightsail_instance",
	"aws_macie2_classification_job",
	"aws_macie2_custom_data_identifier",
	"aws_macie2_findings_filter",
	"aws_macie2_member",
	"aws_media_convert_queue",
	"aws_media_package_channel",
	"aws_media_store_container",
	"aws_mq_broker",
	"aws_mq_configuration",
	"aws_msk_cluster",
	"aws_mwaa_environment",
	"aws_nat_gateway",
	"aws_neptune_cluster",
	"aws_neptune_cluster_endpoint",
	"aws_neptune_cluster_instance",
	"aws_neptune_cluster_parameter_group",
	"aws_neptune_event_subscription",
	"aws_neptune_parameter_group",
	"aws_neptune_subnet_group",
	"aws_network_acl",
	"aws_network_interface",
	"aws_networkfirewall_firewall",
	"aws_networkfirewall_firewall_policy",
	"aws_networkfirewall_rule_group",
	"aws_opsworks_custom_layer",
	"aws_opsworks_ganglia_layer",
	"aws_opsworks_haproxy_layer",
	"aws_opsworks_java_app_layer",
	"aws_opsworks_memcached_layer",
	"aws_opsworks_mysql_layer",
	"aws_opsworks_nodejs_app_layer",
	"aws_opsworks_php_app_layer",
	"aws_opsworks_rails_app_layer",
	"aws_opsworks_stack",
	"aws_opsworks_static_web_layer",
	"aws_organizations_account",
	"aws_organizations_organizational_unit",
	"aws_organizations_policy",
	"aws_pinpoint_app",
	"aws_placement_group",
	"aws_qldb_ledger",
	"aws_ram_resource_share",
	"aws_rds_cluster",
	"aws_rds_cluster_endpoint",
	"aws_rds_cluster_instance",
	"aws_rds_cluster_parameter_group",
	"aws_redshift_cluster",
	"aws_redshift_event_subscription",
	"aws_redshift_parameter_group",
	"aws_redshift_snapshot_copy_grant",
	"aws_redshift_snapshot_schedule",
	"aws_redshift_subnet_group",
	"aws_resourcegroups_group",
	"aws_route53_health_check",
	"aws_route53_resolver_endpoint",
	"aws_route53_resolver_firewall_domain_list",
	"aws_route53_resolver_firewall_rule_group",
	"aws_route53_resolver_firewall_rule_group_association",
	"aws_route53_resolver_query_log_config",
	"aws_route53_resolver_rule",
	"aws_route53_zone",
	"aws_route_table",
	"aws_s3_bucket",
	"aws_s3_bucket_object",
	"aws_s3_object_copy",
	"aws_s3control_bucket",
	"aws_sagemaker_app",
	"aws_sagemaker_device_fleet",
	"aws_sagemaker_domain",
	"aws_sagemaker_endpoint",
	"aws_sagemaker_endpoint_configuration",
	"aws_sagemaker_feature_group",
	"aws_sagemaker_human_task_ui",
	"aws_sagemaker_image",
	"aws_sagemaker_model",
	"aws_sagemaker_model_package_group",
	"aws_sagemaker_notebook_instance",
	"aws_sagemaker_user_profile",
	"aws_sagemaker_workteam",
	"aws_schemas_discoverer",
	"aws_schemas_registry",
	"aws_schemas_schema",
	"aws_secretsmanager_secret",
        "aws_secretsmanager_secret_rotation",
	"aws_security_group",
	"aws_serverlessapplicationrepository_cloudformation_stack",
	"aws_service_discovery_http_namespace",
	"aws_service_discovery_private_dns_namespace",
	"aws_service_discovery_public_dns_namespace",
	"aws_service_discovery_service",
	"aws_servicecatalog_portfolio",
	"aws_servicecatalog_product",
	"aws_servicecatalog_provisioned_product",
	"aws_sfn_activity",
	"aws_sfn_state_machine",
	"aws_shield_protection",
	"aws_shield_protection_group",
	"aws_signer_signing_profile",
	"aws_sns_topic",
	"aws_spot_fleet_request",
	"aws_spot_instance_request",
	"aws_sqs_queue",
	"aws_ssm_activation",
	"aws_ssm_document",
	"aws_ssm_maintenance_window",
	"aws_ssm_parameter",
	"aws_ssm_patch_baseline",
	"aws_ssoadmin_permission_set",
	"aws_storagegateway_cached_iscsi_volume",
	"aws_storagegateway_file_system_association",
	"aws_storagegateway_gateway",
	"aws_storagegateway_nfs_file_share",
	"aws_storagegateway_smb_file_share",
	"aws_storagegateway_stored_iscsi_volume",
	"aws_storagegateway_tape_pool",
	"aws_subnet",
	"aws_swf_domain",
	"aws_synthetics_canary",
	"aws_timestreamwrite_database",
	"aws_timestreamwrite_table",
	"aws_transfer_server",
	"aws_transfer_user",
	"aws_vpc",
	"aws_vpc_dhcp_options",
	"aws_vpc_endpoint",
	"aws_vpc_endpoint_service",
	"aws_vpc_peering_connection",
	"aws_vpc_peering_connection_accepter",
	"aws_vpn_connection",
	"aws_vpn_gateway",
	"aws_waf_rate_based_rule",
	"aws_waf_rule",
	"aws_waf_rule_group",
	"aws_waf_web_acl",
	"aws_wafregional_rate_based_rule",
	"aws_wafregional_rule",
	"aws_wafregional_rule_group",
	"aws_wafregional_web_acl",
	"aws_wafv2_ip_set",
	"aws_wafv2_regex_pattern_set",
	"aws_wafv2_rule_group",
	"aws_wafv2_web_acl",
	"aws_workspaces_directory",
	"aws_workspaces_ip_group",
	"aws_workspaces_workspace",
	"aws_xray_group",
	"aws_xray_sampling_rule",
	"azurerm_active_directory_domain_service",
	"azurerm_analysis_services_server",
	"azurerm_api_management",
	"azurerm_api_management_named_value",
	"azurerm_api_management_property",
	"azurerm_app_configuration",
	"azurerm_app_service",
	"azurerm_app_service_certificate",
	"azurerm_app_service_certificate_order",
	"azurerm_app_service_environment",
	"azurerm_app_service_environment_v3",
	"azurerm_app_service_managed_certificate",
	"azurerm_app_service_plan",
	"azurerm_app_service_slot",
	"azurerm_application_gateway",
	"azurerm_application_insights",
	"azurerm_application_insights_web_test",
	"azurerm_application_security_group",
	"azurerm_attestation_provider",
	"azurerm_automation_account",
	"azurerm_automation_dsc_configuration",
	"azurerm_automation_runbook",
	"azurerm_availability_set",
	"azurerm_backup_policy_file_share",
	"azurerm_backup_policy_vm",
	"azurerm_backup_protected_vm",
	"azurerm_bastion_host",
	"azurerm_batch_account",
	"azurerm_bot_channels_registration",
	"azurerm_bot_connection",
	"azurerm_bot_web_app",
	"azurerm_cdn_endpoint",
	"azurerm_cdn_profile",
	"azurerm_cognitive_account",
	"azurerm_communication_service",
	"azurerm_container_group",
	"azurerm_container_registry",
	"azurerm_container_registry_webhook",
	"azurerm_cosmosdb_account",
	"azurerm_custom_provider",
	"azurerm_dashboard",
	"azurerm_data_factory",
	"azurerm_data_lake_analytics_account",
	"azurerm_data_lake_store",
	"azurerm_data_protection_backup_vault",
	"azurerm_data_share_account",
	"azurerm_database_migration_project",
	"azurerm_database_migration_service",
	"azurerm_databox_edge_device",
	"azurerm_databricks_workspace",
	"azurerm_dedicated_hardware_security_module",
	"azurerm_dedicated_host",
	"azurerm_dedicated_host_group",
	"azurerm_dev_test_global_vm_shutdown_schedule",
	"azurerm_dev_test_lab",
	"azurerm_dev_test_linux_virtual_machine",
	"azurerm_dev_test_policy",
	"azurerm_dev_test_schedule",
	"azurerm_dev_test_virtual_network",
	"azurerm_dev_test_windows_virtual_machine",
	"azurerm_devspace_controller",
	"azurerm_digital_twins_instance",
	"azurerm_disk_access",
	"azurerm_disk_encryption_set",
	"azurerm_dns_a_record",
	"azurerm_dns_aaaa_record",
	"azurerm_dns_caa_record",
	"azurerm_dns_cname_record",
	"azurerm_dns_mx_record",
	"azurerm_dns_ns_record",
	"azurerm_dns_ptr_record",
	"azurerm_dns_srv_record",
	"azurerm_dns_txt_record",
	"azurerm_dns_zone",
	"azurerm_eventgrid_domain",
	"azurerm_eventgrid_system_topic",
	"azurerm_eventgrid_topic",
	"azurerm_eventhub_cluster",
	"azurerm_eventhub_namespace",
	"azurerm_express_route_circuit",
	"azurerm_express_route_gateway",
	"azurerm_express_route_port",
	"azurerm_firewall",
	"azurerm_firewall_policy",
	"azurerm_frontdoor",
	"azurerm_frontdoor_firewall_policy",
	"azurerm_function_app",
	"azurerm_function_app_slot",
	"azurerm_hdinsight_hadoop_cluster",
	"azurerm_hdinsight_hbase_cluster",
	"azurerm_hdinsight_interactive_query_cluster",
	"azurerm_hdinsight_kafka_cluster",
	"azurerm_hdinsight_ml_services_cluster",
	"azurerm_hdinsight_rserver_cluster",
	"azurerm_hdinsight_spark_cluster",
	"azurerm_hdinsight_storm_cluster",
	"azurerm_healthbot",
	"azurerm_healthcare_service",
	"azurerm_hpc_cache",
	"azurerm_image",
	"azurerm_integration_service_environment",
	"azurerm_iot_security_solution",
	"azurerm_iot_time_series_insights_event_source_iothub",
	"azurerm_iot_time_series_insights_gen2_environment",
	"azurerm_iot_time_series_insights_reference_data_set",
	"azurerm_iot_time_series_insights_standard_environment",
	"azurerm_iotcentral_application",
	"azurerm_iothub",
	"azurerm_iothub_dps",
	"azurerm_ip_group",
	"azurerm_key_vault",
	"azurerm_key_vault_certificate",
	"azurerm_key_vault_key",
	"azurerm_key_vault_managed_hardware_security_module",
	"azurerm_key_vault_secret",
	"azurerm_kubernetes_cluster",
	"azurerm_kubernetes_cluster_node_pool",
	"azurerm_kusto_cluster",
	"azurerm_lb",
	"azurerm_linux_virtual_machine",
	"azurerm_linux_virtual_machine_scale_set",
	"azurerm_local_network_gateway",
	"azurerm_log_analytics_cluster",
	"azurerm_log_analytics_linked_service",
	"azurerm_log_analytics_saved_search",
	"azurerm_log_analytics_solution",
	"azurerm_log_analytics_storage_insights",
	"azurerm_log_analytics_workspace",
	"azurerm_logic_app_integration_account",
	"azurerm_logic_app_workflow",
	"azurerm_machine_learning_compute_cluster",
	"azurerm_machine_learning_compute_instance",
	"azurerm_machine_learning_inference_cluster",
	"azurerm_machine_learning_synapse_spark",
	"azurerm_machine_learning_workspace",
	"azurerm_maintenance_configuration",
	"azurerm_managed_application",
	"azurerm_managed_application_definition",
	"azurerm_managed_disk",
	"azurerm_management_group_template_deployment",
	"azurerm_maps_account",
	"azurerm_mariadb_server",
	"azurerm_media_live_event",
	"azurerm_media_services_account",
	"azurerm_media_streaming_endpoint",
	"azurerm_monitor_action_group",
	"azurerm_monitor_action_rule_action_group",
	"azurerm_monitor_action_rule_suppression",
	"azurerm_monitor_activity_log_alert",
	"azurerm_monitor_autoscale_setting",
	"azurerm_monitor_metric_alert",
	"azurerm_monitor_scheduled_query_rules_alert",
	"azurerm_monitor_scheduled_query_rules_log",
	"azurerm_monitor_smart_detector_alert_rule",
	"azurerm_mssql_database",
	"azurerm_mssql_elasticpool",
	"azurerm_mssql_job_agent",
	"azurerm_mssql_server",
	"azurerm_mssql_virtual_machine",
	"azurerm_mysql_server",
	"azurerm_nat_gateway",
	"azurerm_netapp_account",
	"azurerm_netapp_pool",
	"azurerm_netapp_snapshot",
	"azurerm_netapp_volume",
	"azurerm_network_connection_monitor",
	"azurerm_network_ddos_protection_plan",
	"azurerm_network_interface",
	"azurerm_network_profile",
	"azurerm_network_security_group",
	"azurerm_network_watcher",
	"azurerm_network_watcher_flow_log",
	"azurerm_notification_hub",
	"azurerm_notification_hub_namespace",
	"azurerm_orchestrated_virtual_machine_scale_set",
	"azurerm_point_to_site_vpn_gateway",
	"azurerm_postgresql_flexible_server",
	"azurerm_postgresql_server",
	"azurerm_powerbi_embedded",
	"azurerm_private_dns_a_record",
	"azurerm_private_dns_aaaa_record",
	"azurerm_private_dns_cname_record",
	"azurerm_private_dns_mx_record",
	"azurerm_private_dns_ptr_record",
	"azurerm_private_dns_srv_record",
	"azurerm_private_dns_txt_record",
	"azurerm_private_dns_zone",
	"azurerm_private_dns_zone_virtual_network_link",
	"azurerm_private_endpoint",
	"azurerm_private_link_service",
	"azurerm_proximity_placement_group",
	"azurerm_public_ip",
	"azurerm_public_ip_prefix",
	"azurerm_purview_account",
	"azurerm_recovery_services_vault",
	"azurerm_redis_cache",
	"azurerm_redis_enterprise_cluster",
	"azurerm_relay_namespace",
	"azurerm_resource_group",
	"azurerm_resource_group_template_deployment",
	"azurerm_route_filter",
	"azurerm_route_table",
	"azurerm_search_service",
	"azurerm_security_center_automation",
	"azurerm_service_fabric_cluster",
	"azurerm_service_fabric_mesh_application",
	"azurerm_service_fabric_mesh_local_network",
	"azurerm_service_fabric_mesh_secret",
	"azurerm_service_fabric_mesh_secret_value",
	"azurerm_servicebus_namespace",
	"azurerm_shared_image",
	"azurerm_shared_image_gallery",
	"azurerm_shared_image_version",
	"azurerm_signalr_service",
	"azurerm_snapshot",
	"azurerm_spatial_anchors_account",
	"azurerm_spring_cloud_service",
	"azurerm_sql_database",
	"azurerm_sql_elasticpool",
	"azurerm_sql_failover_group",
	"azurerm_sql_server",
	"azurerm_ssh_public_key",
	"azurerm_stack_hci_cluster",
	"azurerm_static_site",
	"azurerm_storage_account",
	"azurerm_storage_sync",
	"azurerm_stream_analytics_job",
	"azurerm_subnet_service_endpoint_storage_policy",
	"azurerm_subscription",
	"azurerm_subscription_template_deployment",
	"azurerm_synapse_private_link_hub",
	"azurerm_synapse_spark_pool",
	"azurerm_synapse_sql_pool",
	"azurerm_synapse_workspace",
	"azurerm_tenant_template_deployment",
	"azurerm_traffic_manager_profile",
	"azurerm_user_assigned_identity",
	"azurerm_video_analyzer",
	"azurerm_virtual_desktop_application_group",
	"azurerm_virtual_desktop_host_pool",
	"azurerm_virtual_desktop_workspace",
	"azurerm_virtual_hub",
	"azurerm_virtual_hub_security_partner_provider",
	"azurerm_virtual_machine",
	"azurerm_virtual_machine_extension",
	"azurerm_virtual_machine_scale_set",
	"azurerm_virtual_network",
	"azurerm_virtual_network_gateway",
	"azurerm_virtual_network_gateway_connection",
	"azurerm_virtual_wan",
	"azurerm_vmware_private_cloud",
	"azurerm_vpn_gateway",
	"azurerm_vpn_server_configuration",
	"azurerm_vpn_site",
	"azurerm_web_application_firewall_policy",
	"azurerm_windows_virtual_machine",
	"azurerm_windows_virtual_machine_scale_set",
	"google_active_directory_domain",
	"google_assured_workloads_workload",
	"google_bigquery_dataset",
	"google_bigquery_job",
	"google_bigquery_table",
	"google_bigtable_instance",
	"google_cloud_identity_group",
	"google_cloudfunctions_function",
	"google_composer_environment",
	"google_compute_disk",
	"google_compute_image",
	"google_compute_instance",
	"google_compute_instance_from_template",
	"google_compute_instance_template",
	"google_compute_region_disk",
	"google_compute_snapshot",
	"google_dataflow_job",
	"google_dataproc_cluster",
	"google_dataproc_job",
	"google_dataproc_workflow_template",
	"google_dialogflow_cx_intent",
	"google_dns_managed_zone",
	"google_eventarc_trigger",
	"google_filestore_instance",
	"google_game_services_game_server_cluster",
	"google_game_services_game_server_config",
	"google_game_services_game_server_deployment",
	"google_game_services_realm",
	"google_gke_hub_membership",
	"google_healthcare_consent_store",
	"google_healthcare_dicom_store",
	"google_healthcare_fhir_store",
	"google_healthcare_hl7_v2_store",
	"google_kms_crypto_key",
	"google_memcache_instance",
	"google_ml_engine_model",
	"google_monitoring_notification_channel",
	"google_network_management_connectivity_test",
	"google_network_services_edge_cache_keyset",
	"google_network_services_edge_cache_origin",
	"google_network_services_edge_cache_service",
	"google_notebooks_instance",
	"google_privateca_ca_pool",
	"google_privateca_certificate",
	"google_privateca_certificate_authority",
	"google_project",
	"google_pubsub_subscription",
	"google_pubsub_topic",
	"google_redis_instance",
	"google_secret_manager_secret",
	"google_spanner_instance",
	"google_storage_bucket",
	"google_tpu_node",
	"google_vertex_ai_dataset",
	"google_workflows_workflow",
}
