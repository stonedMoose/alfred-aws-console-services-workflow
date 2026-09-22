package searchers

import (
	"context"
	"fmt"
	"net/url"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/codebuild"
	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/awsworkflow"
	"github.com/rkoval/alfred-aws-console-services-workflow/caching"
	"github.com/rkoval/alfred-aws-console-services-workflow/searchers/searchutil"
	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

type CodeBuildProjectSearcher struct{}

func (s CodeBuildProjectSearcher) Search(wf *aw.Workflow, searchArgs searchutil.SearchArgs) error {
	cacheName := util.GetCurrentFilename()
	entities := caching.LoadEntityArrayFromCache(wf, searchArgs, cacheName, s.fetch)
	for _, entity := range entities {
		s.addToWorkflow(wf, searchArgs, entity)
	}
	return nil
}

// ListProjects only returns names; batch describing them to enrich the
// subtitle would cost an extra round trip per 100 projects
func (s CodeBuildProjectSearcher) fetch(cfg aws.Config) ([]string, error) {
	svc := codebuild.NewFromConfig(cfg)

	var entities []string
	nextToken := ""
	for {
		params := &codebuild.ListProjectsInput{}
		if nextToken != "" {
			params.NextToken = aws.String(nextToken)
		}
		resp, err := svc.ListProjects(context.TODO(), params)
		if err != nil {
			return nil, err
		}

		entities = append(entities, resp.Projects...)

		if resp.NextToken != nil && *resp.NextToken != "" {
			nextToken = *resp.NextToken
		} else {
			break
		}
	}

	return entities, nil
}

func (s CodeBuildProjectSearcher) addToWorkflow(wf *aw.Workflow, searchArgs searchutil.SearchArgs, entity string) {
	title := entity

	path := fmt.Sprintf("/codesuite/codebuild/projects/%s/history", url.PathEscape(title))
	item := util.NewURLItem(wf, title).
		Arg(util.ConstructAWSConsoleUrl(path, searchArgs.GetRegion())).
		Icon(awsworkflow.GetImageIcon("codebuild")).
		Valid(true)

	searchArgs.AddMatch(item, "", "", title)
}
