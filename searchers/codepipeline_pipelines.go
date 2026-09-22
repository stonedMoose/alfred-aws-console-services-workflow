package searchers

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/codepipeline"
	"github.com/aws/aws-sdk-go-v2/service/codepipeline/types"
	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/awspaging"
	"github.com/rkoval/alfred-aws-console-services-workflow/searchers/searchutil"
)

type CodePipelinePipelinesSearcher struct{}

func (s CodePipelinePipelinesSearcher) Search(wf *aw.Workflow, searchArgs searchutil.SearchArgs) error {
	return searchEntities(wf, searchArgs, "codepipeline_pipelines", s.fetch, s.addToWorkflow)
}

func (CodePipelinePipelinesSearcher) fetch(cfg aws.Config) ([]types.PipelineSummary, error) {
	client := codepipeline.NewFromConfig(cfg)
	return awspaging.FetchAllPages(func(pageToken string) ([]types.PipelineSummary, *string, error) {
		resp, err := client.ListPipelines(context.TODO(), &codepipeline.ListPipelinesInput{
			MaxResults: aws.Int32(100),
			NextToken:  awspaging.TokenOrNil(pageToken),
		})
		if err != nil {
			return nil, nil, err
		}
		return resp.Pipelines, resp.NextToken, nil
	})
}

func (CodePipelinePipelinesSearcher) addToWorkflow(wf *aw.Workflow, searchArgs searchutil.SearchArgs, pipeline types.PipelineSummary) {
	name := aws.ToString(pipeline.Name)
	searchArgs.AddConsoleResource(wf, searchutil.ConsoleResource{
		Title:       name,
		ConsolePath: "/codesuite/codepipeline/pipelines/" + name + "/view",
		ServiceID:   "codepipeline",
	}).Subtitle(subtitleFrom(
		formatIfSet("Version %d", pipeline.Version),
		dateDetail("Created", pipeline.Created),
	))
}
