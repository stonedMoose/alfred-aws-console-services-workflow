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

type VPCSearcher struct{}

func (s VPCSearcher) Search(wf *aw.Workflow, searchArgs searchutil.SearchArgs) error {
	return searchEntities(wf, searchArgs, "vpc_vpcs", s.fetch, s.addToWorkflow)
}

func (VPCSearcher) fetch(cfg aws.Config) ([]types.Vpc, error) {
	client := ec2.NewFromConfig(cfg)
	return awspaging.FetchAllPages(func(pageToken string) ([]types.Vpc, *string, error) {
		resp, err := client.DescribeVpcs(context.TODO(), &ec2.DescribeVpcsInput{
			MaxResults: aws.Int32(1000), // get as many as we can
			NextToken:  awspaging.TokenOrNil(pageToken),
		})
		if err != nil {
			return nil, nil, err
		}
		return resp.Vpcs, resp.NextToken, nil
	})
}

func (VPCSearcher) addToWorkflow(wf *aw.Workflow, searchArgs searchutil.SearchArgs, vpc types.Vpc) {
	id := aws.ToString(vpc.VpcId)
	title, idDetail := namedOrID(ec2NameTag(vpc.Tags), id)
	searchArgs.AddConsoleResource(wf, searchutil.ConsoleResource{
		Title:       title,
		ConsolePath: "/vpc/home#VpcDetails:VpcId=" + id,
		ServiceID:   "vpc",
	}).Subtitle(subtitleFrom(
		idDetail,
		aws.ToString(vpc.CidrBlock),
		defaultVPCDetail(vpc.IsDefault),
	))
}

func defaultVPCDetail(isDefault *bool) string {
	if aws.ToBool(isDefault) {
		return "default"
	}
	return ""
}
