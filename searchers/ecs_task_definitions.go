package searchers

import (
	"context"
	"net/url"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/awspaging"
	"github.com/rkoval/alfred-aws-console-services-workflow/searchers/searchutil"
)

type ECSTaskDefinitionSearcher struct{}

func (s ECSTaskDefinitionSearcher) Search(wf *aw.Workflow, searchArgs searchutil.SearchArgs) error {
	return searchEntities(wf, searchArgs, "ecs_task_definitions", s.fetch, s.addToWorkflow)
}

// fetch lists families rather than individual revisions, since the console
// groups revisions under the family and there can be thousands of them.
func (ECSTaskDefinitionSearcher) fetch(cfg aws.Config) ([]string, error) {
	client := ecs.NewFromConfig(cfg)
	return awspaging.FetchAllPages(func(pageToken string) ([]string, *string, error) {
		resp, err := client.ListTaskDefinitionFamilies(context.TODO(), &ecs.ListTaskDefinitionFamiliesInput{
			MaxResults: aws.Int32(100), // max allowed by this API
			Status:     types.TaskDefinitionFamilyStatusActive,
			NextToken:  awspaging.TokenOrNil(pageToken),
		})
		if err != nil {
			return nil, nil, err
		}
		return resp.Families, resp.NextToken, nil
	})
}

func (ECSTaskDefinitionSearcher) addToWorkflow(wf *aw.Workflow, searchArgs searchutil.SearchArgs, family string) {
	searchArgs.AddConsoleResource(wf, searchutil.ConsoleResource{
		Title:       family,
		ConsolePath: "/ecs/v2/task-definitions/" + url.PathEscape(family),
		ServiceID:   "ecs",
	})
}
