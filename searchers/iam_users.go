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

type IAMUserSearcher struct{}

func (s IAMUserSearcher) Search(wf *aw.Workflow, searchArgs searchutil.SearchArgs) error {
	return searchEntities(wf, searchArgs, "iam_users", s.fetch, s.addToWorkflow)
}

func (IAMUserSearcher) fetch(cfg aws.Config) ([]types.User, error) {
	client := iam.NewFromConfig(cfg)
	return awspaging.FetchAllPages(func(marker string) ([]types.User, *string, error) {
		resp, err := client.ListUsers(context.TODO(), &iam.ListUsersInput{
			MaxItems: aws.Int32(1000), // get as many as we can
			Marker:   awspaging.TokenOrNil(marker),
		})
		if err != nil {
			return nil, nil, err
		}
		return resp.Users, awspaging.NextPageTokenIf(resp.IsTruncated, resp.Marker), nil
	})
}

func (IAMUserSearcher) addToWorkflow(wf *aw.Workflow, searchArgs searchutil.SearchArgs, user types.User) {
	name := aws.ToString(user.UserName)
	searchArgs.AddConsoleResource(wf, searchutil.ConsoleResource{
		Title:       name,
		ConsolePath: "/iamv2/home#/users/details/" + url.PathEscape(name),
		ServiceID:   "iam",
		ID:          aws.ToString(user.Arn),
		IDPrefix:    searchutil.ARNPrefix,
	}).Subtitle(subtitleFrom(aws.ToString(user.Path), userActivityDetail(user)))
}

// userActivityDetail prefers the last console login over the creation date.
func userActivityDetail(user types.User) string {
	if user.PasswordLastUsed != nil {
		return dateDetail("Last used", user.PasswordLastUsed)
	}
	return dateDetail("Created", user.CreateDate)
}
