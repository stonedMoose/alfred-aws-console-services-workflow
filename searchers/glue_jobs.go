package searchers

import (
	"context"
	"net/url"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/glue"
	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/awspaging"
	"github.com/rkoval/alfred-aws-console-services-workflow/searchers/searchutil"
)

type GlueJobSearcher struct{}

func (s GlueJobSearcher) Search(wf *aw.Workflow, searchArgs searchutil.SearchArgs) error {
	return searchEntities(wf, searchArgs, "glue_jobs", s.fetch, s.addToWorkflow)
}

func (GlueJobSearcher) fetch(cfg aws.Config) ([]string, error) {
	client := glue.NewFromConfig(cfg)
	return awspaging.FetchAllPages(func(pageToken string) ([]string, *string, error) {
		resp, err := client.ListJobs(context.TODO(), &glue.ListJobsInput{
			MaxResults: aws.Int32(1000), // get as many as we can
			NextToken:  awspaging.TokenOrNil(pageToken),
		})
		if err != nil {
			return nil, nil, err
		}
		return resp.JobNames, resp.NextToken, nil
	})
}

func (GlueJobSearcher) addToWorkflow(wf *aw.Workflow, searchArgs searchutil.SearchArgs, jobName string) {
	searchArgs.AddConsoleResource(wf, searchutil.ConsoleResource{
		Title:       jobName,
		ConsolePath: "/gluestudio/home#/editor/job/" + url.PathEscape(jobName) + "/script",
		ServiceID:   "glue",
	})
}
