package searchers

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sns/types"
	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/awspaging"
	"github.com/rkoval/alfred-aws-console-services-workflow/searchers/searchutil"
)

// pendingSubscriptionArn is what SNS reports instead of an ARN while the
// endpoint has not confirmed the subscription.
const pendingSubscriptionArn = "PendingConfirmation"

type SNSSubscriptionSearcher struct{}

func (s SNSSubscriptionSearcher) Search(wf *aw.Workflow, searchArgs searchutil.SearchArgs) error {
	return searchEntities(wf, searchArgs, "sns_subscriptions", s.fetch, s.addToWorkflow)
}

func (SNSSubscriptionSearcher) fetch(cfg aws.Config) ([]types.Subscription, error) {
	client := sns.NewFromConfig(cfg)
	return awspaging.FetchAllPages(func(pageToken string) ([]types.Subscription, *string, error) {
		resp, err := client.ListSubscriptions(context.TODO(), &sns.ListSubscriptionsInput{
			NextToken: awspaging.TokenOrNil(pageToken),
		})
		if err != nil {
			return nil, nil, err
		}
		return resp.Subscriptions, resp.NextToken, nil
	})
}

// addToWorkflow sends a confirmed subscription to its own page and a pending
// one to the subscription list, where it can be confirmed.
func (SNSSubscriptionSearcher) addToWorkflow(wf *aw.Workflow, searchArgs searchutil.SearchArgs, subscription types.Subscription) {
	protocol := aws.ToString(subscription.Protocol)
	endpoint := aws.ToString(subscription.Endpoint)
	if isPendingConfirmation(subscription) {
		addSubscriptionItem(wf, searchArgs, subscription, "/sns/v3/home#/subscriptions").
			Subtitle("🕘 " + subtitleFrom(protocol, endpoint))
		return
	}
	arn := aws.ToString(subscription.SubscriptionArn)
	addSubscriptionItem(wf, searchArgs, subscription, "/sns/v3/home#/subscription/"+arn).
		Subtitle("✅ " + subtitleFrom(protocol, endpoint, arnResourceName(arn)))
}

func isPendingConfirmation(subscription types.Subscription) bool {
	return subscription.SubscriptionArn == nil || *subscription.SubscriptionArn == pendingSubscriptionArn
}

func addSubscriptionItem(wf *aw.Workflow, searchArgs searchutil.SearchArgs, subscription types.Subscription, consolePath string) *aw.Item {
	topicArn := aws.ToString(subscription.TopicArn)
	return searchArgs.AddConsoleResource(wf, searchutil.ConsoleResource{
		Title:       arnResourceName(topicArn),
		ConsolePath: consolePath,
		ServiceID:   "sns",
		ID:          topicArn,
		IDPrefix:    searchutil.ARNPrefix,
	})
}
