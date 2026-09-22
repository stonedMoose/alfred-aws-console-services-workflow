package searchers

import (
	"context"
	"net/url"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/awspaging"
	"github.com/rkoval/alfred-aws-console-services-workflow/searchers/searchutil"
	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

type LambdaFunctionSearcher struct{}

func (s LambdaFunctionSearcher) Search(wf *aw.Workflow, searchArgs searchutil.SearchArgs) error {
	return searchEntities(wf, searchArgs, "lambda_functions", s.fetch, s.addToWorkflow)
}

func (LambdaFunctionSearcher) fetch(cfg aws.Config) ([]types.FunctionConfiguration, error) {
	client := lambda.NewFromConfig(cfg)
	return awspaging.FetchAllPages(func(marker string) ([]types.FunctionConfiguration, *string, error) {
		resp, err := client.ListFunctions(context.TODO(), &lambda.ListFunctionsInput{
			MaxItems: aws.Int32(200), // get as many as we can
			Marker:   awspaging.TokenOrNil(marker),
		})
		if err != nil {
			return nil, nil, err
		}
		return resp.Functions, resp.NextMarker, nil
	})
}

func (LambdaFunctionSearcher) addToWorkflow(wf *aw.Workflow, searchArgs searchutil.SearchArgs, function types.FunctionConfiguration) {
	name := aws.ToString(function.FunctionName)
	searchArgs.AddConsoleResource(wf, searchutil.ConsoleResource{
		Title:       name,
		ConsolePath: "/lambda/home#/functions/" + url.PathEscape(name) + "?tab=configuration",
		ServiceID:   "lambda",
		ID:          aws.ToString(function.FunctionArn),
		IDPrefix:    searchutil.ARNPrefix,
	}).Subtitle(subtitleFrom(
		aws.ToString(function.Description),
		string(function.Runtime),
		codeSizeDetail(function.CodeSize),
	))
}

func codeSizeDetail(codeSize int64) string {
	if codeSize == 0 {
		return ""
	}
	return util.ByteFormat(codeSize, 2)
}
