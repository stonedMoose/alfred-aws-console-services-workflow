module github.com/rkoval/alfred-aws-console-services-workflow

go 1.24

require (
	github.com/aws/aws-sdk-go-v2 v1.47.0
	github.com/aws/aws-sdk-go-v2/config v1.29.14
	github.com/aws/aws-sdk-go-v2/service/cloudformation v1.59.2
	github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs v1.49.0
	github.com/aws/aws-sdk-go-v2/service/codepipeline v1.41.0
	github.com/aws/aws-sdk-go-v2/service/ec2 v1.218.0
	github.com/aws/aws-sdk-go-v2/service/ecs v1.57.0
	github.com/aws/aws-sdk-go-v2/service/elasticache v1.46.0
	github.com/aws/aws-sdk-go-v2/service/elasticbeanstalk v1.29.2
	github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2 v1.45.2
	github.com/aws/aws-sdk-go-v2/service/lambda v1.71.2
	github.com/aws/aws-sdk-go-v2/service/rds v1.95.0
	github.com/aws/aws-sdk-go-v2/service/route53 v1.51.1
	github.com/aws/aws-sdk-go-v2/service/s3 v1.79.3
	github.com/aws/aws-sdk-go-v2/service/sns v1.34.4
	github.com/aws/aws-sdk-go-v2/service/wafv2 v1.60.1
	github.com/aws/smithy-go v1.28.1
	github.com/bradleyjkemp/cupaloy v2.3.0+incompatible
	github.com/deanishe/awgo v0.29.1
	github.com/iancoleman/strcase v0.3.0
	github.com/stretchr/testify v1.7.0
	gopkg.in/dnaeon/go-vcr.v4 v4.0.2
	gopkg.in/ini.v1 v1.67.0
	gopkg.in/yaml.v2 v2.4.0
)

require (
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.7.20 // indirect
	github.com/aws/aws-sdk-go-v2/credentials v1.17.67 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.16.30 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.5.3 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.8.3 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.8.3 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.5.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/acm v1.50.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/apigateway v1.48.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/apigatewayv2 v1.43.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/apprunner v1.48.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/autoscaling v1.78.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/batch v1.77.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/cloudfront v1.73.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/cloudtrail v1.65.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/codebuild v1.78.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/cognitoidentity v1.42.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider v1.74.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/dynamodb v1.69.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/ecr v1.66.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/efs v1.49.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/eks v1.100.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/eventbridge v1.54.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/glue v1.161.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/iam v1.64.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.13.19 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.7.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/endpoint-discovery v1.13.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.12.15 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.18.15 // indirect
	github.com/aws/aws-sdk-go-v2/service/kinesis v1.55.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/kms v1.61.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/redshift v1.71.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/secretsmanager v1.50.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/sfn v1.51.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/sqs v1.52.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssm v1.78.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.25.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.30.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.33.19 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/magefile/mage v1.15.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	go.deanishe.net/env v0.5.1 // indirect
	go.deanishe.net/fuzzy v1.0.0 // indirect
	golang.org/x/text v0.25.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
