package searchers

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/awspaging"
	"github.com/rkoval/alfred-aws-console-services-workflow/searchers/searchutil"
)

type VPCSubnetSearcher struct{}

func (s VPCSubnetSearcher) Search(wf *aw.Workflow, searchArgs searchutil.SearchArgs) error {
	return searchEntities(wf, searchArgs, "vpc_subnets", s.fetch, s.addToWorkflow)
}

func (VPCSubnetSearcher) fetch(cfg aws.Config) ([]types.Subnet, error) {
	client := ec2.NewFromConfig(cfg)
	return awspaging.FetchAllPages(func(pageToken string) ([]types.Subnet, *string, error) {
		resp, err := client.DescribeSubnets(context.TODO(), &ec2.DescribeSubnetsInput{
			MaxResults: aws.Int32(1000), // get as many as we can
			NextToken:  awspaging.TokenOrNil(pageToken),
		})
		if err != nil {
			return nil, nil, err
		}
		return resp.Subnets, resp.NextToken, nil
	})
}

func (VPCSubnetSearcher) addToWorkflow(wf *aw.Workflow, searchArgs searchutil.SearchArgs, subnet types.Subnet) {
	id := aws.ToString(subnet.SubnetId)
	title, idDetail := namedOrID(ec2NameTag(subnet.Tags), id)
	searchArgs.AddConsoleResource(wf, searchutil.ConsoleResource{
		Title:       title,
		ConsolePath: "/vpc/home#SubnetDetails:subnetId=" + id,
		ServiceID:   "vpc",
	}).Subtitle(subtitleFrom(
		idDetail,
		aws.ToString(subnet.CidrBlock),
		aws.ToString(subnet.AvailabilityZone),
		aws.ToString(subnet.VpcId),
	))
}
