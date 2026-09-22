package util

import (
	"errors"
	"strings"
)

// AWSConsoleDomain is the host of the AWS console, which depends on the
// partition (commercial, China, GovCloud). It is set once at startup.
var AWSConsoleDomain string

// ConstructAWSConsoleUrl turns a console path into a full URL for the given
// region. Paths that are already absolute URLs are returned untouched. An empty
// region addresses the global console, as global services do.
func ConstructAWSConsoleUrl(path, region string) string {
	if strings.HasPrefix(path, "http") {
		return path
	}
	if AWSConsoleDomain == "" {
		panic(errors.New("AWSConsoleDomain was not initialized"))
	}
	return "https://" + regionalHost(region) + withRegionParameter(path, region)
}

func regionalHost(region string) string {
	if region == "" {
		return AWSConsoleDomain
	}
	return region + "." + AWSConsoleDomain
}

// withRegionParameter adds the region query parameter the console expects,
// before the fragment when there is one, unless the path already carries it.
func withRegionParameter(path, region string) string {
	if region == "" || strings.Contains(path, "region=") {
		return path
	}
	regionParameter := "?region=" + region
	if strings.Contains(path, "#") {
		return strings.Replace(path, "#", regionParameter+"#", 1)
	}
	return path + regionParameter
}
