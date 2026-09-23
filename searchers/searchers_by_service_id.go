package searchers

//go:generate go run ../generators/searchers_by_service_id_sorter/main.go

var cloudFormationStackSearcher = &CloudFormationStackSearcher{}
var cloudWatchLogGroupSearcher = &CloudWatchLogGroupSearcher{}
var cloudwatchLogInsightsQuerySearcher = &CloudWatchLogInsightsQuerySearcher{}
var codePipelinePipelineSearcher = &CodePipelinePipelinesSearcher{}
var ec2InstanceSearcher = &EC2InstanceSearcher{}
var ec2LoadBalancerSearcher = &EC2LoadBalancerSearcher{}
var ec2SecurityGroupSearcher = &EC2SecurityGroupSearcher{}
var elasticBeanstalkEnvironmentSearcher = &ElasticBeanstalkEnvironmentSearcher{}
var elasticacheMemcachedClusterSearcher = &ElasticacheMemcachedClusterSearcher{}
var elasticacheRedisClusterSearcher = &ElasticacheRedisClusterSearcher{}
var elasticbeanstalkApplicationSearcher = &ElasticBeanstalkApplicationSearcher{}
var lambdaFunctionSearcher = &LambdaFunctionSearcher{}
var rdsDatabaseSearcher = &RDSDatabaseSearcher{}
var route53HostedZoneSearcher = &Route53HostedZoneSearcher{}
var s3BucketSearcher = &S3BucketSearcher{}
var snsSubscriptionSearcher = &SNSSubscriptionSearcher{}
var snsTopicSearcher = &SNSTopicSearcher{}
var wafIPSetSearcher = &WAFIPSetSearcher{}
var wafWebACLSearcher = &WAFWebACLSearcher{}

var ecsClusterSearcher = &ECSClusterSearcher{}
var dynamoDBTableSearcher = &DynamoDBTableSearcher{}
var ecrRepositorySearcher = &ECRRepositorySearcher{}
var secretsManagerSecretSearcher = &SecretsManagerSecretSearcher{}
var sqsQueueSearcher = &SQSQueueSearcher{}
var ssmParameterSearcher = &SSMParameterSearcher{}
var eksClusterSearcher = &EKSClusterSearcher{}
var stepFunctionsStateMachineSearcher = &StepFunctionsStateMachineSearcher{}
var ec2TargetGroupSearcher = &EC2TargetGroupSearcher{}
var ec2AutoScalingGroupSearcher = &EC2AutoScalingGroupSearcher{}
var ecsTaskDefinitionSearcher = &ECSTaskDefinitionSearcher{}
var batchJobQueueSearcher = &BatchJobQueueSearcher{}
var appRunnerServiceSearcher = &AppRunnerServiceSearcher{}
var vpcSearcher = &VPCSearcher{}
var vpcSubnetSearcher = &VPCSubnetSearcher{}
var ec2VolumeSearcher = &EC2VolumeSearcher{}
var ec2SnapshotSearcher = &EC2SnapshotSearcher{}
var ec2ImageSearcher = &EC2ImageSearcher{}
var ec2KeyPairSearcher = &EC2KeyPairSearcher{}
var efsFileSystemSearcher = &EFSFileSystemSearcher{}
var iamRoleSearcher = &IAMRoleSearcher{}
var iamUserSearcher = &IAMUserSearcher{}
var iamPolicySearcher = &IAMPolicySearcher{}
var iamGroupSearcher = &IAMGroupSearcher{}
var kmsKeySearcher = &KMSKeySearcher{}
var acmCertificateSearcher = &ACMCertificateSearcher{}
var cloudFrontDistributionSearcher = &CloudFrontDistributionSearcher{}
var apiGatewayAPISearcher = &APIGatewayAPISearcher{}
var eventBridgeRuleSearcher = &EventBridgeRuleSearcher{}
var eventBridgeEventBusSearcher = &EventBridgeEventBusSearcher{}
var kinesisStreamSearcher = &KinesisStreamSearcher{}
var codeBuildProjectSearcher = &CodeBuildProjectSearcher{}
var glueJobSearcher = &GlueJobSearcher{}
var redshiftClusterSearcher = &RedshiftClusterSearcher{}
var cloudTrailTrailSearcher = &CloudTrailTrailSearcher{}
var cognitoUserPoolSearcher = &CognitoUserPoolSearcher{}
var cognitoIdentityPoolSearcher = &CognitoIdentityPoolSearcher{}

var SearchersByServiceId map[string]Searcher = map[string]Searcher{
	"apprunner":                     appRunnerServiceSearcher,
	"apprunner_services":            appRunnerServiceSearcher,
	"batch":                         batchJobQueueSearcher,
	"batch_jobqueues":               batchJobQueueSearcher,
	"cloudformation":                cloudFormationStackSearcher,
	"cloudformation_stacks":         cloudFormationStackSearcher,
	"cloudwatch":                    cloudWatchLogGroupSearcher,
	"cloudwatch_loggroups":          cloudWatchLogGroupSearcher,
	"cloudwatch_loginsights":        cloudwatchLogInsightsQuerySearcher,
	"codepipeline":                  codePipelinePipelineSearcher,
	"codepipeline_pipelines":        codePipelinePipelineSearcher,
	"dynamodb":                      dynamoDBTableSearcher,
	"dynamodb_tables":               dynamoDBTableSearcher,
	"ec2":                           ec2InstanceSearcher,
	"ec2_autoscalinggroups":         ec2AutoScalingGroupSearcher,
	"ec2_instances":                 ec2InstanceSearcher,
	"ec2_loadbalancers":             ec2LoadBalancerSearcher,
	"ec2_securitygroups":            ec2SecurityGroupSearcher,
	"ec2_targetgroups":              ec2TargetGroupSearcher,
	"ecr":                           ecrRepositorySearcher,
	"ecr_privaterepositories":       ecrRepositorySearcher,
	"ecs":                           ecsClusterSearcher,
	"ecs_clusters":                  ecsClusterSearcher,
	"ecs_taskdefinitions":           ecsTaskDefinitionSearcher,
	"eks":                           eksClusterSearcher,
	"eks_clusters":                  eksClusterSearcher,
	"elasticache":                   elasticacheRedisClusterSearcher,
	"elasticache_memcached":         elasticacheMemcachedClusterSearcher,
	"elasticache_redis":             elasticacheRedisClusterSearcher,
	"elasticbeanstalk":              elasticBeanstalkEnvironmentSearcher,
	"elasticbeanstalk_applications": elasticbeanstalkApplicationSearcher,
	"elasticbeanstalk_environments": elasticBeanstalkEnvironmentSearcher,
	"lambda":                        lambdaFunctionSearcher,
	"lambda_functions":              lambdaFunctionSearcher,
	"rds":                           rdsDatabaseSearcher,
	"rds_databases":                 rdsDatabaseSearcher,
	"route53":                       route53HostedZoneSearcher,
	"route53_hostedzones":           route53HostedZoneSearcher,
	"s3":                            s3BucketSearcher,
	"s3_buckets":                    s3BucketSearcher,
	"secretsmanager":                secretsManagerSecretSearcher,
	"secretsmanager_secrets":        secretsManagerSecretSearcher,
	"sns":                           snsTopicSearcher,
	"sns_subscriptions":             snsSubscriptionSearcher,
	"sns_topics":                    snsTopicSearcher,
	"sqs":                           sqsQueueSearcher,
	"sqs_queues":                    sqsQueueSearcher,
	"stepfunctions":                 stepFunctionsStateMachineSearcher,
	"stepfunctions_statemachines":   stepFunctionsStateMachineSearcher,
	"systemsmanager":                ssmParameterSearcher,
	"systemsmanager_parameterstore": ssmParameterSearcher,
	"waf":                           wafWebACLSearcher,
	"waf_ipsets":                    wafIPSetSearcher,
	"waf_webacls":                   wafWebACLSearcher,

	"vpc":                      vpcSearcher,
	"vpc_vpcs":                 vpcSearcher,
	"vpc_subnets":              vpcSubnetSearcher,
	"ec2_volumes":              ec2VolumeSearcher,
	"ec2_snapshots":            ec2SnapshotSearcher,
	"ec2_amis":                 ec2ImageSearcher,
	"ec2_keypairs":             ec2KeyPairSearcher,
	"efs":                      efsFileSystemSearcher,
	"efs_filesystems":          efsFileSystemSearcher,
	"iam":                      iamRoleSearcher,
	"iam_roles":                iamRoleSearcher,
	"iam_users":                iamUserSearcher,
	"iam_policies":             iamPolicySearcher,
	"iam_usergroups":           iamGroupSearcher,
	"kms":                      kmsKeySearcher,
	"kms_customermanagedkeys":  kmsKeySearcher,
	"acm":                      acmCertificateSearcher,
	"acm_certificatemanager":   acmCertificateSearcher,
	"cloudfront":               cloudFrontDistributionSearcher,
	"cloudfront_distributions": cloudFrontDistributionSearcher,
	"apigateway":               apiGatewayAPISearcher,
	"apigateway_apis":          apiGatewayAPISearcher,
	"eventbridge":              eventBridgeRuleSearcher,
	"eventbridge_rules":        eventBridgeRuleSearcher,
	"eventbridge_eventbuses":   eventBridgeEventBusSearcher,
	"kinesis":                  kinesisStreamSearcher,
	"kinesis_datastreams":      kinesisStreamSearcher,
	"codebuild":                codeBuildProjectSearcher,
	"codebuild_projects":       codeBuildProjectSearcher,
	"glue_jobs":                glueJobSearcher,
	"redshift":                 redshiftClusterSearcher,
	"redshift_clusters":        redshiftClusterSearcher,
	"cloudtrail":               cloudTrailTrailSearcher,
	"cloudtrail_trails":        cloudTrailTrailSearcher,
	"cognito":                  cognitoUserPoolSearcher,
	"cognito_userpools":        cognitoUserPoolSearcher,
	"cognito_identitypools":    cognitoIdentityPoolSearcher,
}
