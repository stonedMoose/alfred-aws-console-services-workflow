package searchers

// searchersByServiceID maps the id of a service, or "<service>_<sub-service>",
// to the searcher listing its resources. A bare service id points at the
// searcher of the sub-service people look for most.
var searchersByServiceID = map[string]Searcher{
	"acm":                           ACMCertificateSearcher{},
	"acm_certificatemanager":        ACMCertificateSearcher{},
	"apigateway":                    APIGatewayAPISearcher{},
	"apigateway_apis":               APIGatewayAPISearcher{},
	"apprunner":                     AppRunnerServiceSearcher{},
	"apprunner_services":            AppRunnerServiceSearcher{},
	"batch":                         BatchJobQueueSearcher{},
	"batch_jobqueues":               BatchJobQueueSearcher{},
	"cloudformation":                CloudFormationStackSearcher{},
	"cloudformation_stacks":         CloudFormationStackSearcher{},
	"cloudfront":                    CloudFrontDistributionSearcher{},
	"cloudfront_distributions":      CloudFrontDistributionSearcher{},
	"cloudtrail":                    CloudTrailTrailSearcher{},
	"cloudtrail_trails":             CloudTrailTrailSearcher{},
	"cloudwatch":                    CloudWatchLogGroupSearcher{},
	"cloudwatch_loggroups":          CloudWatchLogGroupSearcher{},
	"cloudwatch_loginsights":        CloudWatchLogInsightsQuerySearcher{},
	"cognito":                       CognitoUserPoolSearcher{},
	"cognito_identitypools":         CognitoIdentityPoolSearcher{},
	"cognito_userpools":             CognitoUserPoolSearcher{},
	"codebuild":                     CodeBuildProjectSearcher{},
	"codebuild_projects":            CodeBuildProjectSearcher{},
	"codepipeline":                  CodePipelinePipelinesSearcher{},
	"codepipeline_pipelines":        CodePipelinePipelinesSearcher{},
	"dynamodb":                      DynamoDBTableSearcher{},
	"dynamodb_tables":               DynamoDBTableSearcher{},
	"ec2":                           EC2InstanceSearcher{},
	"ec2_amis":                      EC2ImageSearcher{},
	"ec2_autoscalinggroups":         EC2AutoScalingGroupSearcher{},
	"ec2_instances":                 EC2InstanceSearcher{},
	"ec2_keypairs":                  EC2KeyPairSearcher{},
	"ec2_loadbalancers":             EC2LoadBalancerSearcher{},
	"ec2_securitygroups":            EC2SecurityGroupSearcher{},
	"ec2_snapshots":                 EC2SnapshotSearcher{},
	"ec2_targetgroups":              EC2TargetGroupSearcher{},
	"ec2_volumes":                   EC2VolumeSearcher{},
	"ecr":                           ECRRepositorySearcher{},
	"ecr_privaterepositories":       ECRRepositorySearcher{},
	"ecs":                           ECSClusterSearcher{},
	"ecs_clusters":                  ECSClusterSearcher{},
	"ecs_taskdefinitions":           ECSTaskDefinitionSearcher{},
	"efs":                           EFSFileSystemSearcher{},
	"efs_filesystems":               EFSFileSystemSearcher{},
	"eks":                           EKSClusterSearcher{},
	"eks_clusters":                  EKSClusterSearcher{},
	"elasticache":                   ElasticacheRedisClusterSearcher{},
	"elasticache_memcached":         ElasticacheMemcachedClusterSearcher{},
	"elasticache_redis":             ElasticacheRedisClusterSearcher{},
	"elasticbeanstalk":              ElasticBeanstalkEnvironmentSearcher{},
	"elasticbeanstalk_applications": ElasticBeanstalkApplicationSearcher{},
	"elasticbeanstalk_environments": ElasticBeanstalkEnvironmentSearcher{},
	"eventbridge":                   EventBridgeRuleSearcher{},
	"eventbridge_eventbuses":        EventBridgeEventBusSearcher{},
	"eventbridge_rules":             EventBridgeRuleSearcher{},
	"glue_jobs":                     GlueJobSearcher{},
	"iam":                           IAMRoleSearcher{},
	"iam_policies":                  IAMPolicySearcher{},
	"iam_roles":                     IAMRoleSearcher{},
	"iam_usergroups":                IAMGroupSearcher{},
	"iam_users":                     IAMUserSearcher{},
	"kinesis":                       KinesisStreamSearcher{},
	"kinesis_datastreams":           KinesisStreamSearcher{},
	"kms":                           KMSKeySearcher{},
	"kms_customermanagedkeys":       KMSKeySearcher{},
	"lambda":                        LambdaFunctionSearcher{},
	"lambda_functions":              LambdaFunctionSearcher{},
	"rds":                           RDSDatabaseSearcher{},
	"rds_databases":                 RDSDatabaseSearcher{},
	"redshift":                      RedshiftClusterSearcher{},
	"redshift_clusters":             RedshiftClusterSearcher{},
	"route53":                       Route53HostedZoneSearcher{},
	"route53_hostedzones":           Route53HostedZoneSearcher{},
	"s3":                            S3BucketSearcher{},
	"s3_buckets":                    S3BucketSearcher{},
	"secretsmanager":                SecretsManagerSecretSearcher{},
	"secretsmanager_secrets":        SecretsManagerSecretSearcher{},
	"sns":                           SNSTopicSearcher{},
	"sns_subscriptions":             SNSSubscriptionSearcher{},
	"sns_topics":                    SNSTopicSearcher{},
	"sqs":                           SQSQueueSearcher{},
	"sqs_queues":                    SQSQueueSearcher{},
	"stepfunctions":                 StepFunctionsStateMachineSearcher{},
	"stepfunctions_statemachines":   StepFunctionsStateMachineSearcher{},
	"systemsmanager":                SSMParameterSearcher{},
	"systemsmanager_parameterstore": SSMParameterSearcher{},
	"vpc":                           VPCSearcher{},
	"vpc_subnets":                   VPCSubnetSearcher{},
	"vpc_vpcs":                      VPCSearcher{},
	"waf":                           WAFWebACLSearcher{},
	"waf_ipsets":                    WAFIPSetSearcher{},
	"waf_webacls":                   WAFWebACLSearcher{},
}

// SearcherFor returns the searcher of a sub-service, or of the service itself
// when subServiceID is empty. It returns nil when no searcher exists yet.
func SearcherFor(serviceID, subServiceID string) Searcher {
	return searchersByServiceID[searcherKey(serviceID, subServiceID)]
}

// IsDefaultSearcher tells whether the sub-service is the one a bare service id
// searches. Searchers carry no state, so comparing them compares their kinds.
func IsDefaultSearcher(serviceID, subServiceID string) bool {
	searcher := SearcherFor(serviceID, subServiceID)
	return searcher != nil && searcher == SearcherFor(serviceID, "")
}

func searcherKey(serviceID, subServiceID string) string {
	if subServiceID == "" {
		return serviceID
	}
	return serviceID + "_" + subServiceID
}
