package caching

import (
	"errors"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/smithy-go"
	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/awsconfig"
	"github.com/rkoval/alfred-aws-console-services-workflow/searchers/searchutil"
	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

// The background fetch cannot talk to Alfred, so it leaves its error in a file
// for the next run to report.
const lastFetchErrorFile = "last-fetch-err.txt"

const (
	noCredentialsError = "NoCredentialProviders"
	missingRegionError = "MissingRegion"
	accessDeniedError  = "You do not have access to fetch these. Check your IAM permissions"
)

func lastFetchErrorPath(wf *aw.Workflow) string {
	return wf.CacheDir() + "/" + lastFetchErrorFile
}

func recordFetchError(wf *aw.Workflow, err error) {
	path := lastFetchErrorPath(wf)
	log.Printf("fetch error occurred. writing to %s ...", path)
	_ = os.WriteFile(path, []byte(describeFetchError(err)), 0600)
}

func clearFetchError(wf *aw.Workflow) {
	os.Remove(lastFetchErrorPath(wf))
}

// describeFetchError turns a fetch error into the line shown to the user.
func describeFetchError(err error) string {
	description := describeKnownFetchError(err)
	if description == "" {
		description = err.Error()
	}
	// without a credentials file, aws-sdk-go-v2 falls back to the instance
	// metadata service, which a workflow running on a Mac never reaches; the
	// resulting error is replaced by one that says what to configure
	if strings.Contains(description, "failed to retrieve credentials") {
		return noCredentialsError
	}
	return description
}

func describeKnownFetchError(err error) string {
	var missingRegion *aws.MissingRegionError
	if errors.As(err, &missingRegion) {
		return missingRegionError
	}
	var apiError smithy.APIError
	if !errors.As(err, &apiError) {
		return ""
	}
	if apiError.ErrorCode() == "AccessDeniedException" {
		return accessDeniedError
	}
	if message := apiError.ErrorMessage(); message != "" {
		return apiError.ErrorCode() + ": " + message
	}
	return apiError.ErrorCode()
}

// reportLastFetchError shows the error the last background fetch left behind,
// if any, and returns it.
// TODO fix the "no results" display when the cause is really a fetch error
func reportLastFetchError(wf *aw.Workflow, searchArgs searchutil.SearchArgs) error {
	description, err := os.ReadFile(lastFetchErrorPath(wf))
	if err != nil {
		// the file usually does not exist, so its absence is not worth a log line
		if !os.IsNotExist(err) {
			log.Println(err)
		}
		return nil
	}
	wf.Configure(aw.SuppressUIDs(true))
	addFetchErrorItem(wf, string(description), profileDescription(searchArgs.Profile))
	return errors.New(string(description))
}

func addFetchErrorItem(wf *aw.Workflow, description, forProfile string) {
	switch {
	case strings.HasPrefix(description, noCredentialsError):
		credentialsFile := homeRelative(awsconfig.GetAwsCredentialsFilePath())
		addConfigurationErrorItem(wf, "AWS credentials not set in "+credentialsFile+" "+forProfile, awsconfig.CredentialsFileDocsURL)
	case strings.HasPrefix(description, missingRegionError):
		configFile := homeRelative(awsconfig.GetAwsProfileFilePath())
		addConfigurationErrorItem(wf, "AWS region not set in "+configFile+" "+forProfile, awsconfig.ConfigFileDocsURL)
	default:
		wf.NewItem(description).Icon(aw.IconError)
	}
}

func addConfigurationErrorItem(wf *aw.Workflow, title, docsURL string) {
	util.NewURLItem(wf, title).
		Subtitle("Press enter to open AWS docs on how to configure").
		Arg(docsURL).
		Icon(aw.IconError)
}

func profileDescription(profile string) string {
	if profile == "" {
		return "for default profile"
	}
	return "for profile \"" + profile + "\""
}

func homeRelative(path string) string {
	home := os.Getenv("HOME")
	if home == "" {
		return path
	}
	return strings.Replace(path, home, "~", 1)
}
