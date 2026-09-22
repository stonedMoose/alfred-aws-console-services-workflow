package searchers

import (
	"context"
	"net/url"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/codebuild"
	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/awspaging"
	"github.com/rkoval/alfred-aws-console-services-workflow/searchers/searchutil"
)

type CodeBuildProjectSearcher struct{}

func (s CodeBuildProjectSearcher) Search(wf *aw.Workflow, searchArgs searchutil.SearchArgs) error {
	return searchEntities(wf, searchArgs, "codebuild_projects", s.fetch, s.addToWorkflow)
}

// fetch lists project names only: describing them to enrich the subtitle would
// cost an extra round trip per 100 projects.
func (CodeBuildProjectSearcher) fetch(cfg aws.Config) ([]string, error) {
	client := codebuild.NewFromConfig(cfg)
	return awspaging.FetchAllPages(func(pageToken string) ([]string, *string, error) {
		resp, err := client.ListProjects(context.TODO(), &codebuild.ListProjectsInput{
			NextToken: awspaging.TokenOrNil(pageToken),
		})
		if err != nil {
			return nil, nil, err
		}
		return resp.Projects, resp.NextToken, nil
	})
}

func (CodeBuildProjectSearcher) addToWorkflow(wf *aw.Workflow, searchArgs searchutil.SearchArgs, projectName string) {
	searchArgs.AddConsoleResource(wf, searchutil.ConsoleResource{
		Title:       projectName,
		ConsolePath: "/codesuite/codebuild/projects/" + url.PathEscape(projectName) + "/history",
		ServiceID:   "codebuild",
	})
}
