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

type WAFIPSetSearcher struct{}

func (s WAFIPSetSearcher) Search(wf *aw.Workflow, searchArgs searchutil.SearchArgs) error {
	return searchEntities(wf, searchArgs, "waf_ip_sets", s.fetch, s.addToWorkflow)
}

func (WAFIPSetSearcher) fetch(cfg aws.Config) ([]types.IPSetSummary, error) {
	client := wafv2.NewFromConfig(cfg)
	return awspaging.FetchAllPages(func(marker string) ([]types.IPSetSummary, *string, error) {
		resp, err := client.ListIPSets(context.TODO(), &wafv2.ListIPSetsInput{
			Limit:      aws.Int32(100), // get as many as we can
			Scope:      wafScope,
			NextMarker: awspaging.TokenOrNil(marker),
		})
		if err != nil {
			return nil, nil, err
		}
		return resp.IPSets, resp.NextMarker, nil
	})
}

func (WAFIPSetSearcher) addToWorkflow(wf *aw.Workflow, searchArgs searchutil.SearchArgs, ipSet types.IPSetSummary) {
	name := aws.ToString(ipSet.Name)
	searchArgs.AddConsoleResource(wf, searchutil.ConsoleResource{
		Title:       name,
		ConsolePath: wafConsolePath("/wafv2/homev2/ip-set/%s/%s?region=%s", name, aws.ToString(ipSet.Id), searchArgs.Cfg.Region),
		ServiceID:   "waf",
		ID:          aws.ToString(ipSet.ARN),
		IDPrefix:    searchutil.ARNPrefix,
	}).Subtitle(aws.ToString(ipSet.Description))
}
