package searchers

import (
	"context"
	"net/url"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/awspaging"
	"github.com/rkoval/alfred-aws-console-services-workflow/searchers/searchutil"
)

type SSMParameterSearcher struct{}

func (s SSMParameterSearcher) Search(wf *aw.Workflow, searchArgs searchutil.SearchArgs) error {
	return searchEntities(wf, searchArgs, "ssm_parameters", s.fetch, s.addToWorkflow)
}

func (SSMParameterSearcher) fetch(cfg aws.Config) ([]types.ParameterMetadata, error) {
	client := ssm.NewFromConfig(cfg)
	return awspaging.FetchAllPages(func(pageToken string) ([]types.ParameterMetadata, *string, error) {
		resp, err := client.DescribeParameters(context.TODO(), &ssm.DescribeParametersInput{
			MaxResults: aws.Int32(50), // max allowed by this API
			NextToken:  awspaging.TokenOrNil(pageToken),
		})
		if err != nil {
			return nil, nil, err
		}
		return resp.Parameters, resp.NextToken, nil
	})
}

func (SSMParameterSearcher) addToWorkflow(wf *aw.Workflow, searchArgs searchutil.SearchArgs, parameter types.ParameterMetadata) {
	name := aws.ToString(parameter.Name)
	searchArgs.AddConsoleResource(wf, searchutil.ConsoleResource{
		Title: name,
		// the console expects the parameter name double-encoded, so `/db/password`
		// must be written as `%252Fdb%252Fpassword`
		ConsolePath: "/systems-manager/parameters/" + url.PathEscape(url.PathEscape(name)) + "/description",
		ServiceID:   "systemsmanager",
		ID:          aws.ToString(parameter.ARN),
		IDPrefix:    searchutil.ARNPrefix,
	}).Subtitle(subtitleFrom(
		string(parameter.Type),
		aws.ToString(parameter.Description),
		dateDetail("Modified", parameter.LastModifiedDate),
	))
}
