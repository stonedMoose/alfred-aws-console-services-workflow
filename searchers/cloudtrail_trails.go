package searchers

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudtrail"
	"github.com/aws/aws-sdk-go-v2/service/cloudtrail/types"
	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/awsworkflow"
	"github.com/rkoval/alfred-aws-console-services-workflow/caching"
	"github.com/rkoval/alfred-aws-console-services-workflow/searchers/searchutil"
	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

type CloudTrailTrailSearcher struct{}

func (s CloudTrailTrailSearcher) Search(wf *aw.Workflow, searchArgs searchutil.SearchArgs) error {
	cacheName := util.GetCurrentFilename()
	entities := caching.LoadEntityArrayFromCache(wf, searchArgs, cacheName, s.fetch)
	for _, entity := range entities {
		s.addToWorkflow(wf, searchArgs, entity)
	}
	return nil
}

func (s CloudTrailTrailSearcher) fetch(cfg aws.Config) ([]types.TrailInfo, error) {
	svc := cloudtrail.NewFromConfig(cfg)

	var entities []types.TrailInfo
	nextToken := ""
	for {
		params := &cloudtrail.ListTrailsInput{}
		if nextToken != "" {
			params.NextToken = aws.String(nextToken)
		}
		resp, err := svc.ListTrails(context.TODO(), params)
		if err != nil {
			return nil, err
		}

		entities = append(entities, resp.Trails...)

		if resp.NextToken != nil && *resp.NextToken != "" {
			nextToken = *resp.NextToken
		} else {
			break
		}
	}

	return entities, nil
}

func (s CloudTrailTrailSearcher) addToWorkflow(wf *aw.Workflow, searchArgs searchutil.SearchArgs, entity types.TrailInfo) {
	title := ""
	if entity.Name != nil {
		title = *entity.Name
	}

	subtitleArray := []string{}
	subtitleArray = util.AppendString(subtitleArray, entity.HomeRegion)
	subtitle := strings.Join(subtitleArray, " – ")

	arn := ""
	if entity.TrailARN != nil {
		arn = *entity.TrailARN
	}
	if title == "" {
		title = util.GetEndOfArn(arn)
	}

	path := fmt.Sprintf("/cloudtrail/home#/trails/%s", arn)
	item := util.NewURLItem(wf, title).
		Subtitle(subtitle).
		Arg(util.ConstructAWSConsoleUrl(path, searchArgs.GetRegion())).
		Icon(awsworkflow.GetImageIcon("cloudtrail")).
		Valid(true)

	searchArgs.AddMatch(item, "arn:", arn, title)
}
