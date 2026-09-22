// Package awsworkflow loads the AWS configuration the workflow runs with and
// describes the console services it knows.
package awsworkflow

import (
	"context"
	"errors"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/rkoval/alfred-aws-console-services-workflow/awsconfig"
	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

// InitAWS loads the AWS configuration, for the given profile and region when
// they override the defaults, and points the console URLs at the matching
// partition. transport, when set, replaces the HTTP transport of the clients.
func InitAWS(transport http.RoundTripper, profile *awsconfig.Profile, region *awsconfig.Region) aws.Config {
	cfg, err := config.LoadDefaultConfig(context.TODO(), loadOptionsFor(profile, region)...)
	if err != nil {
		panic(err)
	}
	if transport != nil {
		cfg.HTTPClient = &http.Client{Transport: transport}
	}
	InitAWSConsoleDomain(cfg.Region)
	return cfg
}

// loadOptionsFor applies the profile, the profile's own region, then the
// region override, which wins.
func loadOptionsFor(profile *awsconfig.Profile, region *awsconfig.Region) []func(*config.LoadOptions) error {
	var options []func(*config.LoadOptions) error
	if profile != nil {
		options = append(options, config.WithSharedConfigProfile(profile.Name))
		if profile.Region != "" {
			options = append(options, config.WithRegion(profile.Region))
		}
	}
	if region != nil {
		options = append(options, config.WithRegion(region.Name))
	}
	return options
}

// ProfileName is the name of the shared profile the configuration was loaded
// from, or empty when it came from elsewhere.
func ProfileName(cfg aws.Config) string {
	for _, source := range cfg.ConfigSources {
		if sharedConfig, ok := source.(config.SharedConfig); ok {
			return sharedConfig.Profile
		}
	}
	return ""
}

const (
	consoleDomainEnvVar           = "ALFRED_AWS_CONSOLE_SERVICES_WORKFLOW_AWS_CONSOLE_DOMAIN"
	misspelledConsoleDomainEnvVar = "ALRED_AWS_CONSOLE_SERVICES_WORKFLOW_AWS_CONSOLE_DOMAIN"

	commercialConsoleDomain = "console.aws.amazon.com"
	chinaConsoleDomain      = "console.amazonaws.cn"
	usGovConsoleDomain      = "console.amazonaws-us-gov.com"
)

// InitAWSConsoleDomain points the console URLs at the partition the region
// belongs to, unless the environment names a domain.
func InitAWSConsoleDomain(region string) {
	rejectMisspelledConsoleDomainEnvVar()
	util.AWSConsoleDomain = configuredConsoleDomain(region)
}

func rejectMisspelledConsoleDomainEnvVar() {
	if os.Getenv(misspelledConsoleDomainEnvVar) != "" {
		panic(errors.New("`" + misspelledConsoleDomainEnvVar + "` env var was renamed to `" + consoleDomainEnvVar + "` due to the misspelling. Please update your config"))
	}
}

func configuredConsoleDomain(region string) string {
	if domain := os.Getenv(consoleDomainEnvVar); domain != "" {
		return domain
	}
	return partitionConsoleDomain(region)
}

func partitionConsoleDomain(region string) string {
	switch {
	case strings.HasPrefix(region, "cn-"):
		return chinaConsoleDomain
	case strings.HasPrefix(region, "us-gov-"):
		return usGovConsoleDomain
	default:
		return commercialConsoleDomain
	}
}
