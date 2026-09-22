package searchers

import (
	"context"
	"fmt"
	"net/url"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge/types"
	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/awsworkflow"
	"github.com/rkoval/alfred-aws-console-services-workflow/caching"
	"github.com/rkoval/alfred-aws-console-services-workflow/searchers/searchutil"
	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

type EventBridgeEventBusSearcher struct{}

func (s EventBridgeEventBusSearcher) Search(wf *aw.Workflow, searchArgs searchutil.SearchArgs) error {
	cacheName := util.GetCurrentFilename()
	entities := caching.LoadEntityArrayFromCache(wf, searchArgs, cacheName, s.fetch)
	for _, entity := range entities {
		s.addToWorkflow(wf, searchArgs, entity)
	}
	return nil
}

func (s EventBridgeEventBusSearcher) fetch(cfg aws.Config) ([]types.EventBus, error) {
	svc := eventbridge.NewFromConfig(cfg)

	var entities []types.EventBus
	nextToken := ""
	for {
		params := &eventbridge.ListEventBusesInput{
			Limit: aws.Int32(100), // max allowed by this API
		}
		if nextToken != "" {
			params.NextToken = aws.String(nextToken)
		}
		resp, err := svc.ListEventBuses(context.TODO(), params)
		if err != nil {
			return nil, err
		}

		entities = append(entities, resp.EventBuses...)

		if resp.NextToken != nil && *resp.NextToken != "" {
			nextToken = *resp.NextToken
		} else {
			break
		}
	}

	return entities, nil
}

func (s EventBridgeEventBusSearcher) addToWorkflow(wf *aw.Workflow, searchArgs searchutil.SearchArgs, entity types.EventBus) {
	title := *entity.Name

	subtitle := ""
	if entity.Description != nil {
		subtitle = *entity.Description
	}

	path := fmt.Sprintf("/events/home#/eventbus/%s", url.PathEscape(title))
	item := util.NewURLItem(wf, title).
		Subtitle(subtitle).
		Arg(util.ConstructAWSConsoleUrl(path, searchArgs.GetRegion())).
		Icon(awsworkflow.GetImageIcon("eventbridge")).
		Valid(true)

	arn := ""
	if entity.Arn != nil {
		arn = *entity.Arn
	}
	searchArgs.AddMatch(item, "arn:", arn, title)
}
