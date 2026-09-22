package searchers

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2/types"
	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/awspaging"
	"github.com/rkoval/alfred-aws-console-services-workflow/searchers/searchutil"
)

type EC2LoadBalancerSearcher struct{}

func (s EC2LoadBalancerSearcher) Search(wf *aw.Workflow, searchArgs searchutil.SearchArgs) error {
	return searchEntities(wf, searchArgs, "ec2_load_balancers", s.fetch, s.addToWorkflow)
}

func (EC2LoadBalancerSearcher) fetch(cfg aws.Config) ([]types.LoadBalancer, error) {
	client := elasticloadbalancingv2.NewFromConfig(cfg)
	return awspaging.FetchAllPages(func(marker string) ([]types.LoadBalancer, *string, error) {
		resp, err := client.DescribeLoadBalancers(context.TODO(), &elasticloadbalancingv2.DescribeLoadBalancersInput{
			PageSize: aws.Int32(400),
			Marker:   awspaging.TokenOrNil(marker),
		})
		if err != nil {
			return nil, nil, err
		}
		return resp.LoadBalancers, resp.NextMarker, nil
	})
}

func (EC2LoadBalancerSearcher) addToWorkflow(wf *aw.Workflow, searchArgs searchutil.SearchArgs, loadBalancer types.LoadBalancer) {
	arn := aws.ToString(loadBalancer.LoadBalancerArn)
	title := arn
	if loadBalancer.LoadBalancerName != nil {
		title = *loadBalancer.LoadBalancerName
	}
	searchArgs.AddConsoleResource(wf, searchutil.ConsoleResource{
		Title:       title,
		ConsolePath: "/ec2/home#LoadBalancers:search=" + arn + ";sort=loadBalancerName",
		ServiceID:   "ec2",
		ID:          arn,
		IDPrefix:    searchutil.ARNPrefix,
	}).Subtitle(subtitleFrom(string(loadBalancer.Type), aws.ToString(loadBalancer.DNSName)))
}
