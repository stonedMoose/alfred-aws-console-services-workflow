package searchers

import (
	"context"
	"net/url"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge/types"
	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/awspaging"
	"github.com/rkoval/alfred-aws-console-services-workflow/searchers/searchutil"
)

type EventBridgeRuleSearcher struct{}

func (s EventBridgeRuleSearcher) Search(wf *aw.Workflow, searchArgs searchutil.SearchArgs) error {
	return searchEntities(wf, searchArgs, "eventbridge_rules", s.fetch, s.addToWorkflow)
}

func (EventBridgeRuleSearcher) fetch(cfg aws.Config) ([]types.Rule, error) {
	client := eventbridge.NewFromConfig(cfg)
	return awspaging.FetchAllPages(func(pageToken string) ([]types.Rule, *string, error) {
		resp, err := client.ListRules(context.TODO(), &eventbridge.ListRulesInput{
			Limit:     aws.Int32(100), // max allowed by this API
			NextToken: awspaging.TokenOrNil(pageToken),
		})
		if err != nil {
			return nil, nil, err
		}
		return resp.Rules, resp.NextToken, nil
	})
}

func (EventBridgeRuleSearcher) addToWorkflow(wf *aw.Workflow, searchArgs searchutil.SearchArgs, rule types.Rule) {
	name := aws.ToString(rule.Name)
	searchArgs.AddConsoleResource(wf, searchutil.ConsoleResource{
		Title:       name,
		ConsolePath: "/events/home#/rules/" + url.PathEscape(name),
		ServiceID:   "eventbridge",
		ID:          aws.ToString(rule.Arn),
		IDPrefix:    searchutil.ARNPrefix,
	}).Subtitle(subtitleFrom(
		strings.ToLower(string(rule.State)),
		aws.ToString(rule.ScheduleExpression),
		aws.ToString(rule.Description),
	))
}
