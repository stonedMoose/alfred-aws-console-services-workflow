package searchers

import (
	"context"
	"net/url"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge/types"
	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/awspaging"
	"github.com/rkoval/alfred-aws-console-services-workflow/searchers/searchutil"
)

type EventBridgeEventBusSearcher struct{}

func (s EventBridgeEventBusSearcher) Search(wf *aw.Workflow, searchArgs searchutil.SearchArgs) error {
	return searchEntities(wf, searchArgs, "eventbridge_event_buses", s.fetch, s.addToWorkflow)
}

func (EventBridgeEventBusSearcher) fetch(cfg aws.Config) ([]types.EventBus, error) {
	client := eventbridge.NewFromConfig(cfg)
	return awspaging.FetchAllPages(func(pageToken string) ([]types.EventBus, *string, error) {
		resp, err := client.ListEventBuses(context.TODO(), &eventbridge.ListEventBusesInput{
			Limit:     aws.Int32(100), // max allowed by this API
			NextToken: awspaging.TokenOrNil(pageToken),
		})
		if err != nil {
			return nil, nil, err
		}
		return resp.EventBuses, resp.NextToken, nil
	})
}

func (EventBridgeEventBusSearcher) addToWorkflow(wf *aw.Workflow, searchArgs searchutil.SearchArgs, eventBus types.EventBus) {
	name := aws.ToString(eventBus.Name)
	searchArgs.AddConsoleResource(wf, searchutil.ConsoleResource{
		Title:       name,
		ConsolePath: "/events/home#/eventbus/" + url.PathEscape(name),
		ServiceID:   "eventbridge",
		ID:          aws.ToString(eventBus.Arn),
		IDPrefix:    searchutil.ARNPrefix,
	}).Subtitle(aws.ToString(eventBus.Description))
}
