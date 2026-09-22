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

type EC2SecurityGroupSearcher struct{}

func (s EC2SecurityGroupSearcher) Search(wf *aw.Workflow, searchArgs searchutil.SearchArgs) error {
	return searchEntities(wf, searchArgs, "ec2_security_groups", s.fetch, s.addToWorkflow)
}

func (EC2SecurityGroupSearcher) fetch(cfg aws.Config) ([]types.SecurityGroup, error) {
	client := ec2.NewFromConfig(cfg)
	return awspaging.FetchAllPages(func(pageToken string) ([]types.SecurityGroup, *string, error) {
		resp, err := client.DescribeSecurityGroups(context.TODO(), &ec2.DescribeSecurityGroupsInput{
			MaxResults: aws.Int32(1000), // get as many as we can
			NextToken:  aws.String(pageToken),
		})
		if err != nil {
			return nil, nil, err
		}
		return resp.SecurityGroups, resp.NextToken, nil
	})
}

func (EC2SecurityGroupSearcher) addToWorkflow(wf *aw.Workflow, searchArgs searchutil.SearchArgs, group types.SecurityGroup) {
	id := aws.ToString(group.GroupId)
	title, idDetail := namedOrID(ec2NameTag(group.Tags), id)
	searchArgs.AddConsoleResource(wf, searchutil.ConsoleResource{
		Title:       title,
		ConsolePath: "/ec2/v2/home#SecurityGroups:group-id=" + id,
		ServiceID:   "ec2",
		ID:          id,
		IDPrefix:    "sg-",
	}).Subtitle(joinKnown(" ", idDetail, aws.ToString(group.Description)))
}
