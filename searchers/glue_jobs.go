package searchers

import (
	"context"
	"fmt"
	"net/url"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/glue"
	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/awsworkflow"
	"github.com/rkoval/alfred-aws-console-services-workflow/caching"
	"github.com/rkoval/alfred-aws-console-services-workflow/searchers/searchutil"
	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

type GlueJobSearcher struct{}

func (s GlueJobSearcher) Search(wf *aw.Workflow, searchArgs searchutil.SearchArgs) error {
	cacheName := util.GetCurrentFilename()
	entities := caching.LoadEntityArrayFromCache(wf, searchArgs, cacheName, s.fetch)
	for _, entity := range entities {
		s.addToWorkflow(wf, searchArgs, entity)
	}
	return nil
}

func (s GlueJobSearcher) fetch(cfg aws.Config) ([]string, error) {
	svc := glue.NewFromConfig(cfg)

	var entities []string
	nextToken := ""
	for {
		params := &glue.ListJobsInput{
			MaxResults: aws.Int32(1000), // get as many as we can
		}
		if nextToken != "" {
			params.NextToken = aws.String(nextToken)
		}
		resp, err := svc.ListJobs(context.TODO(), params)
		if err != nil {
			return nil, err
		}

		entities = append(entities, resp.JobNames...)

		if resp.NextToken != nil && *resp.NextToken != "" {
			nextToken = *resp.NextToken
		} else {
			break
		}
	}

	return entities, nil
}

func (s GlueJobSearcher) addToWorkflow(wf *aw.Workflow, searchArgs searchutil.SearchArgs, entity string) {
	title := entity

	path := fmt.Sprintf("/gluestudio/home#/editor/job/%s/script", url.PathEscape(title))
	item := util.NewURLItem(wf, title).
		Arg(util.ConstructAWSConsoleUrl(path, searchArgs.GetRegion())).
		Icon(awsworkflow.GetImageIcon("glue")).
		Valid(true)

	searchArgs.AddMatch(item, "", "", title)
}
