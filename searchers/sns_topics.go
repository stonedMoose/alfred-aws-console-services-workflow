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

type SNSTopicSearcher struct{}

func (s SNSTopicSearcher) Search(wf *aw.Workflow, searchArgs searchutil.SearchArgs) error {
	return searchEntities(wf, searchArgs, "sns_topics", s.fetch, s.addToWorkflow)
}

func (SNSTopicSearcher) fetch(cfg aws.Config) ([]types.Topic, error) {
	client := sns.NewFromConfig(cfg)
	return awspaging.FetchAllPages(func(pageToken string) ([]types.Topic, *string, error) {
		resp, err := client.ListTopics(context.TODO(), &sns.ListTopicsInput{
			NextToken: awspaging.TokenOrNil(pageToken),
		})
		if err != nil {
			return nil, nil, err
		}
		return resp.Topics, resp.NextToken, nil
	})
}

func (SNSTopicSearcher) addToWorkflow(wf *aw.Workflow, searchArgs searchutil.SearchArgs, topic types.Topic) {
	arn := aws.ToString(topic.TopicArn)
	searchArgs.AddConsoleResource(wf, searchutil.ConsoleResource{
		Title:       arnResourceName(arn),
		ConsolePath: "/sns/v3/home#/topic/" + arn,
		ServiceID:   "sns",
		ID:          arn,
		IDPrefix:    searchutil.ARNPrefix,
	}).Subtitle(arn)
}
