package searchers

import (
	"context"
	"net/url"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/awspaging"
	"github.com/rkoval/alfred-aws-console-services-workflow/searchers/searchutil"
)

type SQSQueueSearcher struct{}

func (s SQSQueueSearcher) Search(wf *aw.Workflow, searchArgs searchutil.SearchArgs) error {
	return searchEntities(wf, searchArgs, "sqs_queues", s.fetch, s.addToWorkflow)
}

// fetch returns queue URLs, which is all ListQueues gives and all the console
// link needs.
func (SQSQueueSearcher) fetch(cfg aws.Config) ([]string, error) {
	client := sqs.NewFromConfig(cfg)
	return awspaging.FetchAllPages(func(pageToken string) ([]string, *string, error) {
		resp, err := client.ListQueues(context.TODO(), &sqs.ListQueuesInput{
			MaxResults: aws.Int32(1000), // max allowed by this API
			NextToken:  awspaging.TokenOrNil(pageToken),
		})
		if err != nil {
			return nil, nil, err
		}
		return resp.QueueUrls, resp.NextToken, nil
	})
}

func (SQSQueueSearcher) addToWorkflow(wf *aw.Workflow, searchArgs searchutil.SearchArgs, queueURL string) {
	searchArgs.AddConsoleResource(wf, searchutil.ConsoleResource{
		Title:       queueNameOf(queueURL),
		ConsolePath: "/sqs/v2/home#/queues/" + url.QueryEscape(queueURL),
		ServiceID:   "sqs",
	}).Subtitle(queueURL)
}

// queueNameOf is the last segment of a queue URL, which looks like
// https://sqs.us-east-1.amazonaws.com/123456789012/my-queue.
func queueNameOf(queueURL string) string {
	return queueURL[strings.LastIndex(queueURL, "/")+1:]
}
