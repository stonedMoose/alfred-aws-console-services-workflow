package searchers

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudtrail"
	"github.com/aws/aws-sdk-go-v2/service/cloudtrail/types"
	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/awspaging"
	"github.com/rkoval/alfred-aws-console-services-workflow/searchers/searchutil"
)

type CloudTrailTrailSearcher struct{}

func (s CloudTrailTrailSearcher) Search(wf *aw.Workflow, searchArgs searchutil.SearchArgs) error {
	return searchEntities(wf, searchArgs, "cloudtrail_trails", s.fetch, s.addToWorkflow)
}

func (CloudTrailTrailSearcher) fetch(cfg aws.Config) ([]types.TrailInfo, error) {
	client := cloudtrail.NewFromConfig(cfg)
	return awspaging.FetchAllPages(func(pageToken string) ([]types.TrailInfo, *string, error) {
		resp, err := client.ListTrails(context.TODO(), &cloudtrail.ListTrailsInput{
			NextToken: awspaging.TokenOrNil(pageToken),
		})
		if err != nil {
			return nil, nil, err
		}
		return resp.Trails, resp.NextToken, nil
	})
}

func (CloudTrailTrailSearcher) addToWorkflow(wf *aw.Workflow, searchArgs searchutil.SearchArgs, trail types.TrailInfo) {
	arn := aws.ToString(trail.TrailARN)
	title := aws.ToString(trail.Name)
	if title == "" {
		title = arnResourceName(arn)
	}
	searchArgs.AddConsoleResource(wf, searchutil.ConsoleResource{
		Title:       title,
		ConsolePath: "/cloudtrail/home#/trails/" + arn,
		ServiceID:   "cloudtrail",
		ID:          arn,
		IDPrefix:    searchutil.ARNPrefix,
	}).Subtitle(subtitleFrom(aws.ToString(trail.HomeRegion)))
}
