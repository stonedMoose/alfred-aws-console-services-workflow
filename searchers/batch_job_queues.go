package searchers

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/batch"
	"github.com/aws/aws-sdk-go-v2/service/batch/types"
	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/awsworkflow"
	"github.com/rkoval/alfred-aws-console-services-workflow/caching"
	"github.com/rkoval/alfred-aws-console-services-workflow/searchers/searchutil"
	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

type BatchJobQueueSearcher struct{}

func (s BatchJobQueueSearcher) Search(wf *aw.Workflow, searchArgs searchutil.SearchArgs) error {
	cacheName := util.GetCurrentFilename()
	entities := caching.LoadEntityArrayFromCache(wf, searchArgs, cacheName, s.fetch)
	for _, entity := range entities {
		s.addToWorkflow(wf, searchArgs, entity)
	}
	return nil
}

func (s BatchJobQueueSearcher) fetch(cfg aws.Config) ([]types.JobQueueDetail, error) {
	svc := batch.NewFromConfig(cfg)

	var entities []types.JobQueueDetail
	nextToken := ""
	for {
		params := &batch.DescribeJobQueuesInput{
			MaxResults: aws.Int32(100), // max allowed by this API
		}
		if nextToken != "" {
			params.NextToken = aws.String(nextToken)
		}
		resp, err := svc.DescribeJobQueues(context.TODO(), params)
		if err != nil {
			return nil, err
		}

		entities = append(entities, resp.JobQueues...)

		if resp.NextToken != nil && *resp.NextToken != "" {
			nextToken = *resp.NextToken
		} else {
			break
		}
	}

	return entities, nil
}

func (s BatchJobQueueSearcher) addToWorkflow(wf *aw.Workflow, searchArgs searchutil.SearchArgs, entity types.JobQueueDetail) {
	title := *entity.JobQueueName

	subtitleArray := []string{}
	if entity.State != "" {
		subtitleArray = append(subtitleArray, strings.ToLower(string(entity.State)))
	}
	if entity.Status != "" {
		subtitleArray = append(subtitleArray, strings.ToLower(string(entity.Status)))
	}
	if entity.Priority != nil {
		subtitleArray = append(subtitleArray, "priority "+strconv.Itoa(int(*entity.Priority)))
	}
	subtitle := strings.Join(subtitleArray, " – ")

	path := fmt.Sprintf("/batch/home#queues/detail/%s", *entity.JobQueueArn)
	item := util.NewURLItem(wf, title).
		Subtitle(subtitle).
		Arg(util.ConstructAWSConsoleUrl(path, searchArgs.GetRegion())).
		Icon(awsworkflow.GetImageIcon("batch")).
		Valid(true)

	searchArgs.AddMatch(item, "arn:", *entity.JobQueueArn, title)
}
