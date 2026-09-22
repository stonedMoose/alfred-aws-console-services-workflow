// Package aliases holds the prefixes that give a query word a special
// meaning, each overridable through the environment.
package aliases

import "os"

var (
	// OverrideAwsRegion introduces a region to search in, e.g. "$us-east-1".
	OverrideAwsRegion string
	// OverrideAwsProfile introduces a profile to search with, e.g. "@work".
	OverrideAwsProfile string
	// Search introduces a resource search within a service, e.g. "ec2 ,web".
	Search string
)

func init() {
	OverrideAwsRegion = envOrDefault("ALFRED_AWS_CONSOLE_SERVICES_OVERRIDE_AWS_REGION_ALIAS", "$")
	OverrideAwsProfile = envOrDefault("ALFRED_AWS_CONSOLE_SERVICES_OVERRIDE_AWS_PROFILE_ALIAS", "@")
	Search = envOrDefault("ALFRED_AWS_CONSOLE_SERVICES_WORKFLOW_SEARCH_ALIAS", ",")
}

func envOrDefault(name, defaultValue string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}
	return defaultValue
}
