package searchers

import (
	"context"
	"net/url"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/awspaging"
	"github.com/rkoval/alfred-aws-console-services-workflow/searchers/searchutil"
)

type IAMGroupSearcher struct{}

func (s IAMGroupSearcher) Search(wf *aw.Workflow, searchArgs searchutil.SearchArgs) error {
	return searchEntities(wf, searchArgs, "iam_groups", s.fetch, s.addToWorkflow)
}

func (IAMGroupSearcher) fetch(cfg aws.Config) ([]types.Group, error) {
	client := iam.NewFromConfig(cfg)
	return awspaging.FetchAllPages(func(marker string) ([]types.Group, *string, error) {
		resp, err := client.ListGroups(context.TODO(), &iam.ListGroupsInput{
			MaxItems: aws.Int32(1000), // get as many as we can
			Marker:   awspaging.TokenOrNil(marker),
		})
		if err != nil {
			return nil, nil, err
		}
		return resp.Groups, awspaging.NextPageTokenIf(resp.IsTruncated, resp.Marker), nil
	})
}

func (IAMGroupSearcher) addToWorkflow(wf *aw.Workflow, searchArgs searchutil.SearchArgs, group types.Group) {
	name := aws.ToString(group.GroupName)
	searchArgs.AddConsoleResource(wf, searchutil.ConsoleResource{
		Title:       name,
		ConsolePath: "/iamv2/home#/groups/details/" + url.PathEscape(name),
		ServiceID:   "iam",
		ID:          aws.ToString(group.Arn),
		IDPrefix:    searchutil.ARNPrefix,
	}).Subtitle(subtitleFrom(
		aws.ToString(group.Path),
		dateDetail("Created", group.CreateDate),
	))
}
