package searchers

import (
	"context"
	"fmt"
	"net/url"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kinesis"
	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/awsworkflow"
	"github.com/rkoval/alfred-aws-console-services-workflow/caching"
	"github.com/rkoval/alfred-aws-console-services-workflow/searchers/searchutil"
	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

type KinesisStreamSearcher struct{}

func (s KinesisStreamSearcher) Search(wf *aw.Workflow, searchArgs searchutil.SearchArgs) error {
	cacheName := util.GetCurrentFilename()
	entities := caching.LoadEntityArrayFromCache(wf, searchArgs, cacheName, s.fetch)
	for _, entity := range entities {
		s.addToWorkflow(wf, searchArgs, entity)
	}
	return nil
}

func (s KinesisStreamSearcher) fetch(cfg aws.Config) ([]string, error) {
	svc := kinesis.NewFromConfig(cfg)

	var entities []string
	nextToken := ""
	for {
		params := &kinesis.ListStreamsInput{
			Limit: aws.Int32(100), // max allowed by this API
		}
		if nextToken != "" {
			params.NextToken = aws.String(nextToken)
		}
		resp, err := svc.ListStreams(context.TODO(), params)
		if err != nil {
			return nil, err
		}

		entities = append(entities, resp.StreamNames...)

		if resp.NextToken != nil && *resp.NextToken != "" {
			nextToken = *resp.NextToken
		} else {
			break
		}
	}

	return entities, nil
}

func (s KinesisStreamSearcher) addToWorkflow(wf *aw.Workflow, searchArgs searchutil.SearchArgs, entity string) {
	title := entity

	path := fmt.Sprintf("/kinesis/home#/streams/details/%s", url.PathEscape(title))
	item := util.NewURLItem(wf, title).
		Arg(util.ConstructAWSConsoleUrl(path, searchArgs.GetRegion())).
		Icon(awsworkflow.GetImageIcon("kinesis")).
		Valid(true)

	searchArgs.AddMatch(item, "", "", title)
}
