package searchers

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/wafv2"
	"github.com/aws/aws-sdk-go-v2/service/wafv2/types"
	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/awspaging"
	"github.com/rkoval/alfred-aws-console-services-workflow/searchers/searchutil"
)

type WAFWebACLSearcher struct{}

func (s WAFWebACLSearcher) Search(wf *aw.Workflow, searchArgs searchutil.SearchArgs) error {
	return searchEntities(wf, searchArgs, "waf_web_acls", s.fetch, s.addToWorkflow)
}

func (WAFWebACLSearcher) fetch(cfg aws.Config) ([]types.WebACLSummary, error) {
	client := wafv2.NewFromConfig(cfg)
	return awspaging.FetchAllPages(func(marker string) ([]types.WebACLSummary, *string, error) {
		resp, err := client.ListWebACLs(context.TODO(), &wafv2.ListWebACLsInput{
			Limit:      aws.Int32(100), // get as many as we can
			Scope:      wafScope,
			NextMarker: awspaging.TokenOrNil(marker),
		})
		if err != nil {
			return nil, nil, err
		}
		return resp.WebACLs, resp.NextMarker, nil
	})
}

func (WAFWebACLSearcher) addToWorkflow(wf *aw.Workflow, searchArgs searchutil.SearchArgs, webACL types.WebACLSummary) {
	name := aws.ToString(webACL.Name)
	searchArgs.AddConsoleResource(wf, searchutil.ConsoleResource{
		Title:       name,
		ConsolePath: wafConsolePath("/wafv2/homev2/web-acl/%s/%s/overview?region=%s", name, aws.ToString(webACL.Id), searchArgs.Cfg.Region),
		ServiceID:   "waf",
		ID:          aws.ToString(webACL.ARN),
		IDPrefix:    searchutil.ARNPrefix,
	}).Subtitle(aws.ToString(webACL.Description))
}
