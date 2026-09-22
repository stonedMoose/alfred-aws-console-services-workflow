package searchers

import (
	"context"
	"net/url"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/awspaging"
	"github.com/rkoval/alfred-aws-console-services-workflow/searchers/searchutil"
)

type CloudWatchLogInsightsQuerySearcher struct{}

func (s CloudWatchLogInsightsQuerySearcher) Search(wf *aw.Workflow, searchArgs searchutil.SearchArgs) error {
	return searchEntities(wf, searchArgs, "cloud_watch_log_insights_queries", s.fetch, s.addToWorkflow)
}

func (CloudWatchLogInsightsQuerySearcher) fetch(cfg aws.Config) ([]types.QueryDefinition, error) {
	client := cloudwatchlogs.NewFromConfig(cfg)
	return awspaging.FetchAllPages(func(pageToken string) ([]types.QueryDefinition, *string, error) {
		resp, err := client.DescribeQueryDefinitions(context.TODO(), &cloudwatchlogs.DescribeQueryDefinitionsInput{
			MaxResults: aws.Int32(1000),
			NextToken:  awspaging.TokenOrNil(pageToken),
		})
		if err != nil {
			return nil, nil, err
		}
		return resp.QueryDefinitions, resp.NextToken, nil
	})
}

func (CloudWatchLogInsightsQuerySearcher) addToWorkflow(wf *aw.Workflow, searchArgs searchutil.SearchArgs, query types.QueryDefinition) {
	searchArgs.AddConsoleResource(wf, searchutil.ConsoleResource{
		Title:       aws.ToString(query.Name),
		ConsolePath: "/cloudwatch/home#logsV2:logs-insights$3FqueryDetail$3D" + insightsQueryDetail(query),
		ServiceID:   "cloudwatch",
	}).Subtitle(subtitleFrom(singleLine(aws.ToString(query.QueryString))))
}

func singleLine(text string) string {
	return strings.ReplaceAll(text, "\n", " ")
}

// insightsQueryDetail encodes a saved query the way the logs-insights console
// URL wants it: a URL escape wrapped around another, then "%" rewritten as "$".
// The console references queries by this scheme rather than by id. Adapted from
// https://stackoverflow.com/questions/60796991/is-there-a-way-to-generate-the-aws-console-urls-for-cloudwatch-log-group-filters
func insightsQueryDetail(query types.QueryDefinition) string {
	var detail strings.Builder
	detail.WriteString(outerEscape("~(end~0~start~-3600~timeType~'RELATIVE~unit~'seconds~editorString~'"))
	detail.WriteString(innerEscape(aws.ToString(query.QueryString)))
	detail.WriteString(outerEscape("~isLiveTail~false~queryId~'"))
	detail.WriteString(innerEscape(aws.ToString(query.QueryDefinitionId)))
	detail.WriteString(outerEscape("~source~("))
	for _, logGroupName := range query.LogGroupNames {
		detail.WriteString(outerEscape("~'"))
		detail.WriteString(innerEscape(logGroupName))
	}
	detail.WriteString(outerEscape("))"))
	return strings.ReplaceAll(detail.String(), "%", "$")
}

// outerEscape escapes tildes by hand: QueryEscape leaves them alone, browsers
// do not (https://github.com/golang/go/issues/47379).
// TODO do the outer escaping manually for performance?
func outerEscape(s string) string {
	return url.QueryEscape(strings.ReplaceAll(url.QueryEscape(s), "~", "%7E"))
}

func innerEscape(s string) string {
	return strings.ReplaceAll(url.QueryEscape(s), "%", "*")
}
