package searchers

import (
	"context"
	"net/url"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kinesis"
	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/awspaging"
	"github.com/rkoval/alfred-aws-console-services-workflow/searchers/searchutil"
)

type KinesisStreamSearcher struct{}

func (s KinesisStreamSearcher) Search(wf *aw.Workflow, searchArgs searchutil.SearchArgs) error {
	return searchEntities(wf, searchArgs, "kinesis_streams", s.fetch, s.addToWorkflow)
}

func (KinesisStreamSearcher) fetch(cfg aws.Config) ([]string, error) {
	client := kinesis.NewFromConfig(cfg)
	return awspaging.FetchAllPages(func(pageToken string) ([]string, *string, error) {
		resp, err := client.ListStreams(context.TODO(), &kinesis.ListStreamsInput{
			Limit:     aws.Int32(100), // max allowed by this API
			NextToken: awspaging.TokenOrNil(pageToken),
		})
		if err != nil {
			return nil, nil, err
		}
		return resp.StreamNames, resp.NextToken, nil
	})
}

func (KinesisStreamSearcher) addToWorkflow(wf *aw.Workflow, searchArgs searchutil.SearchArgs, streamName string) {
	searchArgs.AddConsoleResource(wf, searchutil.ConsoleResource{
		Title:       streamName,
		ConsolePath: "/kinesis/home#/streams/details/" + url.PathEscape(streamName),
		ServiceID:   "kinesis",
	})
}
