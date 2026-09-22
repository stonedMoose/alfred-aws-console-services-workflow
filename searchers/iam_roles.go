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

type IAMRoleSearcher struct{}

func (s IAMRoleSearcher) Search(wf *aw.Workflow, searchArgs searchutil.SearchArgs) error {
	return searchEntities(wf, searchArgs, "iam_roles", s.fetch, s.addToWorkflow)
}

func (IAMRoleSearcher) fetch(cfg aws.Config) ([]types.Role, error) {
	client := iam.NewFromConfig(cfg)
	return awspaging.FetchAllPages(func(marker string) ([]types.Role, *string, error) {
		resp, err := client.ListRoles(context.TODO(), &iam.ListRolesInput{
			MaxItems: aws.Int32(1000), // get as many as we can
			Marker:   awspaging.TokenOrNil(marker),
		})
		if err != nil {
			return nil, nil, err
		}
		return resp.Roles, awspaging.NextPageTokenIf(resp.IsTruncated, resp.Marker), nil
	})
}

func (IAMRoleSearcher) addToWorkflow(wf *aw.Workflow, searchArgs searchutil.SearchArgs, role types.Role) {
	name := aws.ToString(role.RoleName)
	searchArgs.AddConsoleResource(wf, searchutil.ConsoleResource{
		Title:       name,
		ConsolePath: "/iamv2/home#/roles/details/" + url.PathEscape(name),
		ServiceID:   "iam",
		ID:          aws.ToString(role.Arn),
		IDPrefix:    searchutil.ARNPrefix,
	}).Subtitle(subtitleFrom(
		aws.ToString(role.Description),
		aws.ToString(role.Path),
		dateDetail("Created", role.CreateDate),
	))
}
