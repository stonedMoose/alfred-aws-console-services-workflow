package searchers

import (
	"context"
	"net/url"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/awspaging"
	"github.com/rkoval/alfred-aws-console-services-workflow/searchers/searchutil"
	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

type CloudWatchLogGroupSearcher struct{}

func (s CloudWatchLogGroupSearcher) Search(wf *aw.Workflow, searchArgs searchutil.SearchArgs) error {
	return searchEntities(wf, searchArgs, "cloudwatch_log_groups", s.fetch, s.addToWorkflow)
}

func (CloudWatchLogGroupSearcher) fetch(cfg aws.Config) ([]types.LogGroup, error) {
	client := cloudwatchlogs.NewFromConfig(cfg)
	return awspaging.FetchAllPages(func(pageToken string) ([]types.LogGroup, *string, error) {
		resp, err := client.DescribeLogGroups(context.TODO(), &cloudwatchlogs.DescribeLogGroupsInput{
			Limit:     aws.Int32(50), // get as many as we can
			NextToken: awspaging.TokenOrNil(pageToken),
		})
		if err != nil {
			return nil, nil, err
		}
		return resp.LogGroups, resp.NextToken, nil
	})
}

func (CloudWatchLogGroupSearcher) addToWorkflow(wf *aw.Workflow, searchArgs searchutil.SearchArgs, logGroup types.LogGroup) {
	name := aws.ToString(logGroup.LogGroupName)
	searchArgs.AddConsoleResource(wf, searchutil.ConsoleResource{
		Title:       name,
		ConsolePath: "/cloudwatch/home#logsV2:log-groups/log-group/" + url.PathEscape(name) + "/log-events",
		ServiceID:   "cloudwatch",
		ID:          aws.ToString(logGroup.Arn),
		IDPrefix:    searchutil.ARNPrefix,
	}).Subtitle(subtitleFrom(
		storedBytesDetail(logGroup.StoredBytes),
		formatIfSet("%d day retention", logGroup.RetentionInDays),
	))
}

func storedBytesDetail(storedBytes *int64) string {
	if storedBytes == nil {
		return ""
	}
	return util.ByteFormat(*storedBytes, 2) + " stored"
}
