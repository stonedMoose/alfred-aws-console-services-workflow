package searchers

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/apprunner"
	"github.com/aws/aws-sdk-go-v2/service/apprunner/types"
	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/awspaging"
	"github.com/rkoval/alfred-aws-console-services-workflow/searchers/searchutil"
)

type AppRunnerServiceSearcher struct{}

func (s AppRunnerServiceSearcher) Search(wf *aw.Workflow, searchArgs searchutil.SearchArgs) error {
	return searchEntities(wf, searchArgs, "apprunner_services", s.fetch, s.addToWorkflow)
}

func (AppRunnerServiceSearcher) fetch(cfg aws.Config) ([]types.ServiceSummary, error) {
	client := apprunner.NewFromConfig(cfg)
	return awspaging.FetchAllPages(func(pageToken string) ([]types.ServiceSummary, *string, error) {
		resp, err := client.ListServices(context.TODO(), &apprunner.ListServicesInput{
			MaxResults: aws.Int32(20), // max allowed by this API
			NextToken:  awspaging.TokenOrNil(pageToken),
		})
		if err != nil {
			return nil, nil, err
		}
		return resp.ServiceSummaryList, resp.NextToken, nil
	})
}

func (AppRunnerServiceSearcher) addToWorkflow(wf *aw.Workflow, searchArgs searchutil.SearchArgs, service types.ServiceSummary) {
	arn := aws.ToString(service.ServiceArn)
	searchArgs.AddConsoleResource(wf, searchutil.ConsoleResource{
		Title:       aws.ToString(service.ServiceName),
		ConsolePath: "/apprunner/home#/services/" + arn,
		ServiceID:   "apprunner",
		ID:          arn,
		IDPrefix:    searchutil.ARNPrefix,
	}).Subtitle(subtitleFrom(
		strings.ToLower(string(service.Status)),
		aws.ToString(service.ServiceUrl),
		dateDetail("Updated", service.UpdatedAt),
	))
}
