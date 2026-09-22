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

type IAMPolicySearcher struct{}

func (s IAMPolicySearcher) Search(wf *aw.Workflow, searchArgs searchutil.SearchArgs) error {
	return searchEntities(wf, searchArgs, "iam_policies", s.fetch, s.addToWorkflow)
}

// fetch lists customer managed policies only; the AWS managed ones number in
// the thousands and are rarely what someone is looking for.
func (IAMPolicySearcher) fetch(cfg aws.Config) ([]types.Policy, error) {
	client := iam.NewFromConfig(cfg)
	return awspaging.FetchAllPages(func(marker string) ([]types.Policy, *string, error) {
		resp, err := client.ListPolicies(context.TODO(), &iam.ListPoliciesInput{
			MaxItems: aws.Int32(1000), // get as many as we can
			Scope:    types.PolicyScopeTypeLocal,
			Marker:   awspaging.TokenOrNil(marker),
		})
		if err != nil {
			return nil, nil, err
		}
		return resp.Policies, awspaging.NextPageTokenIf(resp.IsTruncated, resp.Marker), nil
	})
}

func (IAMPolicySearcher) addToWorkflow(wf *aw.Workflow, searchArgs searchutil.SearchArgs, policy types.Policy) {
	arn := aws.ToString(policy.Arn)
	searchArgs.AddConsoleResource(wf, searchutil.ConsoleResource{
		Title:       aws.ToString(policy.PolicyName),
		ConsolePath: "/iamv2/home#/policies/details/" + url.QueryEscape(arn),
		ServiceID:   "iam",
		ID:          arn,
		IDPrefix:    searchutil.ARNPrefix,
	}).Subtitle(subtitleFrom(
		aws.ToString(policy.Description),
		formatIfSet("attached to %d", policy.AttachmentCount),
		aws.ToString(policy.Path),
	))
}
