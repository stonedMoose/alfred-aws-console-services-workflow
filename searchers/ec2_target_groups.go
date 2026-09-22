package searchers

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2/types"
	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/awspaging"
	"github.com/rkoval/alfred-aws-console-services-workflow/searchers/searchutil"
)

type EC2TargetGroupSearcher struct{}

func (s EC2TargetGroupSearcher) Search(wf *aw.Workflow, searchArgs searchutil.SearchArgs) error {
	return searchEntities(wf, searchArgs, "ec2_target_groups", s.fetch, s.addToWorkflow)
}

func (EC2TargetGroupSearcher) fetch(cfg aws.Config) ([]types.TargetGroup, error) {
	client := elasticloadbalancingv2.NewFromConfig(cfg)
	return awspaging.FetchAllPages(func(marker string) ([]types.TargetGroup, *string, error) {
		resp, err := client.DescribeTargetGroups(context.TODO(), &elasticloadbalancingv2.DescribeTargetGroupsInput{
			PageSize: aws.Int32(400),
			Marker:   awspaging.TokenOrNil(marker),
		})
		if err != nil {
			return nil, nil, err
		}
		return resp.TargetGroups, resp.NextMarker, nil
	})
}

func (EC2TargetGroupSearcher) addToWorkflow(wf *aw.Workflow, searchArgs searchutil.SearchArgs, targetGroup types.TargetGroup) {
	arn := aws.ToString(targetGroup.TargetGroupArn)
	searchArgs.AddConsoleResource(wf, searchutil.ConsoleResource{
		Title:       aws.ToString(targetGroup.TargetGroupName),
		ConsolePath: "/ec2/home#TargetGroup:targetGroupArn=" + arn,
		ServiceID:   "ec2",
		ID:          arn,
		IDPrefix:    searchutil.ARNPrefix,
	}).Subtitle(subtitleFrom(
		string(targetGroup.TargetType),
		protocolPortDetail(targetGroup),
		aws.ToString(targetGroup.VpcId),
	))
}

func protocolPortDetail(targetGroup types.TargetGroup) string {
	if targetGroup.Protocol == "" || targetGroup.Port == nil {
		return ""
	}
	return fmt.Sprintf("%s:%d", targetGroup.Protocol, *targetGroup.Port)
}
