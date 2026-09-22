package searchers

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/awsworkflow"
	"github.com/rkoval/alfred-aws-console-services-workflow/caching"
	"github.com/rkoval/alfred-aws-console-services-workflow/searchers/searchutil"
	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

type SQSQueueSearcher struct{}

func (s SQSQueueSearcher) Search(wf *aw.Workflow, searchArgs searchutil.SearchArgs) error {
	cacheName := util.GetCurrentFilename()
	entities := caching.LoadEntityArrayFromCache(wf, searchArgs, cacheName, s.fetch)
	for _, entity := range entities {
		s.addToWorkflow(wf, searchArgs, entity)
	}
	return nil
}

// ListQueues returns queue URLs rather than queue objects, so that is what gets cached
func (s SQSQueueSearcher) fetch(cfg aws.Config) ([]string, error) {
	svc := sqs.NewFromConfig(cfg)

	var entities []string
	nextToken := ""
	for {
		params := &sqs.ListQueuesInput{
			MaxResults: aws.Int32(1000), // max allowed by this API
		}
		if nextToken != "" {
			params.NextToken = aws.String(nextToken)
		}
		resp, err := svc.ListQueues(context.TODO(), params)
		if err != nil {
			return nil, err
		}

		entities = append(entities, resp.QueueUrls...)

		if resp.NextToken != nil && *resp.NextToken != "" {
			nextToken = *resp.NextToken
		} else {
			break
		}
	}

	return entities, nil
}

func (s SQSQueueSearcher) addToWorkflow(wf *aw.Workflow, searchArgs searchutil.SearchArgs, entity string) {
	// a queue URL looks like https://sqs.us-east-1.amazonaws.com/123456789012/my-queue
	title := entity[strings.LastIndex(entity, "/")+1:]

	path := fmt.Sprintf("/sqs/v2/home#/queues/%s", url.QueryEscape(entity))
	item := util.NewURLItem(wf, title).
		Subtitle(entity).
		Arg(util.ConstructAWSConsoleUrl(path, searchArgs.GetRegion())).
		Icon(awsworkflow.GetImageIcon("sqs")).
		Valid(true)

	searchArgs.AddMatch(item, "", "", title)
}
