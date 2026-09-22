package searchers

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/batch"
	"github.com/aws/aws-sdk-go-v2/service/batch/types"
	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/awspaging"
	"github.com/rkoval/alfred-aws-console-services-workflow/searchers/searchutil"
)

type BatchJobQueueSearcher struct{}

func (s BatchJobQueueSearcher) Search(wf *aw.Workflow, searchArgs searchutil.SearchArgs) error {
	return searchEntities(wf, searchArgs, "batch_job_queues", s.fetch, s.addToWorkflow)
}

func (BatchJobQueueSearcher) fetch(cfg aws.Config) ([]types.JobQueueDetail, error) {
	client := batch.NewFromConfig(cfg)
	return awspaging.FetchAllPages(func(pageToken string) ([]types.JobQueueDetail, *string, error) {
		resp, err := client.DescribeJobQueues(context.TODO(), &batch.DescribeJobQueuesInput{
			MaxResults: aws.Int32(100), // max allowed by this API
			NextToken:  awspaging.TokenOrNil(pageToken),
		})
		if err != nil {
			return nil, nil, err
		}
		return resp.JobQueues, resp.NextToken, nil
	})
}

func (BatchJobQueueSearcher) addToWorkflow(wf *aw.Workflow, searchArgs searchutil.SearchArgs, queue types.JobQueueDetail) {
	arn := aws.ToString(queue.JobQueueArn)
	searchArgs.AddConsoleResource(wf, searchutil.ConsoleResource{
		Title:       aws.ToString(queue.JobQueueName),
		ConsolePath: "/batch/home#queues/detail/" + arn,
		ServiceID:   "batch",
		ID:          arn,
		IDPrefix:    searchutil.ARNPrefix,
	}).Subtitle(subtitleFrom(
		strings.ToLower(string(queue.State)),
		strings.ToLower(string(queue.Status)),
		formatIfSet("priority %d", queue.Priority),
	))
}
