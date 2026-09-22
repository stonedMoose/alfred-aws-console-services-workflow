package searchers

import (
	"context"
	"fmt"
	"net/url"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/awsworkflow"
	"github.com/rkoval/alfred-aws-console-services-workflow/caching"
	"github.com/rkoval/alfred-aws-console-services-workflow/searchers/searchutil"
	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

type ECSTaskDefinitionSearcher struct{}

func (s ECSTaskDefinitionSearcher) Search(wf *aw.Workflow, searchArgs searchutil.SearchArgs) error {
	cacheName := util.GetCurrentFilename()
	entities := caching.LoadEntityArrayFromCache(wf, searchArgs, cacheName, s.fetch)
	for _, entity := range entities {
		s.addToWorkflow(wf, searchArgs, entity)
	}
	return nil
}

// families rather than individual revisions, since the console groups
// revisions under the family and there can be thousands of them
func (s ECSTaskDefinitionSearcher) fetch(cfg aws.Config) ([]string, error) {
	svc := ecs.NewFromConfig(cfg)

	var entities []string
	nextToken := ""
	for {
		params := &ecs.ListTaskDefinitionFamiliesInput{
			MaxResults: aws.Int32(100), // max allowed by this API
			Status:     types.TaskDefinitionFamilyStatusActive,
		}
		if nextToken != "" {
			params.NextToken = aws.String(nextToken)
		}
		resp, err := svc.ListTaskDefinitionFamilies(context.TODO(), params)
		if err != nil {
			return nil, err
		}

		entities = append(entities, resp.Families...)

		if resp.NextToken != nil && *resp.NextToken != "" {
			nextToken = *resp.NextToken
		} else {
			break
		}
	}

	return entities, nil
}

func (s ECSTaskDefinitionSearcher) addToWorkflow(wf *aw.Workflow, searchArgs searchutil.SearchArgs, entity string) {
	title := entity

	path := fmt.Sprintf("/ecs/v2/task-definitions/%s", url.PathEscape(title))
	item := util.NewURLItem(wf, title).
		Arg(util.ConstructAWSConsoleUrl(path, searchArgs.GetRegion())).
		Icon(awsworkflow.GetImageIcon("ecs")).
		Valid(true)

	searchArgs.AddMatch(item, "", "", title)
}
