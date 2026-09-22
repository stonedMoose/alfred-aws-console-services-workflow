package searchers

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront/types"
	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/awspaging"
	"github.com/rkoval/alfred-aws-console-services-workflow/searchers/searchutil"
)

type CloudFrontDistributionSearcher struct{}

func (s CloudFrontDistributionSearcher) Search(wf *aw.Workflow, searchArgs searchutil.SearchArgs) error {
	return searchEntities(wf, searchArgs, "cloudfront_distributions", s.fetch, s.addToWorkflow)
}

func (CloudFrontDistributionSearcher) fetch(cfg aws.Config) ([]types.DistributionSummary, error) {
	client := cloudfront.NewFromConfig(cfg)
	return awspaging.FetchAllPages(func(marker string) ([]types.DistributionSummary, *string, error) {
		resp, err := client.ListDistributions(context.TODO(), &cloudfront.ListDistributionsInput{
			MaxItems: aws.Int32(100), // max allowed by this API
			Marker:   awspaging.TokenOrNil(marker),
		})
		if err != nil {
			return nil, nil, err
		}
		list := resp.DistributionList
		if list == nil {
			return nil, nil, nil
		}
		return list.Items, awspaging.NextPageTokenIf(aws.ToBool(list.IsTruncated), list.NextMarker), nil
	})
}

func (CloudFrontDistributionSearcher) addToWorkflow(wf *aw.Workflow, searchArgs searchutil.SearchArgs, distribution types.DistributionSummary) {
	id := aws.ToString(distribution.Id)
	title, otherNames := distributionNames(id, distributionAliases(distribution))
	searchArgs.AddConsoleResource(wf, searchutil.ConsoleResource{
		Title:       title,
		ConsolePath: "/cloudfront/v4/home#/distributions/" + id,
		ServiceID:   "cloudfront",
		ID:          aws.ToString(distribution.ARN),
		IDPrefix:    searchutil.ARNPrefix,
	}).Subtitle(subtitleFrom(append(otherNames,
		aws.ToString(distribution.DomainName),
		disabledDetail(distribution.Enabled),
	)...))
}

// distributionNames titles a distribution by its first alias when it has one,
// the id and the remaining aliases becoming details; otherwise by its id.
func distributionNames(id string, aliases []string) (title string, otherNames []string) {
	if len(aliases) == 0 {
		return id, nil
	}
	return aliases[0], append([]string{id}, aliases[1:]...)
}

func distributionAliases(distribution types.DistributionSummary) []string {
	if distribution.Aliases == nil {
		return nil
	}
	return distribution.Aliases.Items
}

func disabledDetail(enabled *bool) string {
	if enabled != nil && !*enabled {
		return "disabled"
	}
	return ""
}
