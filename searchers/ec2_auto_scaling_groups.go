package searchers

import (
	"context"
	"fmt"
	"net/url"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling/types"
	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/awspaging"
	"github.com/rkoval/alfred-aws-console-services-workflow/searchers/searchutil"
)

type EC2AutoScalingGroupSearcher struct{}

func (s EC2AutoScalingGroupSearcher) Search(wf *aw.Workflow, searchArgs searchutil.SearchArgs) error {
	return searchEntities(wf, searchArgs, "ec2_auto_scaling_groups", s.fetch, s.addToWorkflow)
}

func (EC2AutoScalingGroupSearcher) fetch(cfg aws.Config) ([]types.AutoScalingGroup, error) {
	client := autoscaling.NewFromConfig(cfg)
	return awspaging.FetchAllPages(func(pageToken string) ([]types.AutoScalingGroup, *string, error) {
		resp, err := client.DescribeAutoScalingGroups(context.TODO(), &autoscaling.DescribeAutoScalingGroupsInput{
			MaxRecords: aws.Int32(100), // max allowed by this API
			NextToken:  awspaging.TokenOrNil(pageToken),
		})
		if err != nil {
			return nil, nil, err
		}
		return resp.AutoScalingGroups, resp.NextToken, nil
	})
}

func (EC2AutoScalingGroupSearcher) addToWorkflow(wf *aw.Workflow, searchArgs searchutil.SearchArgs, group types.AutoScalingGroup) {
	name := aws.ToString(group.AutoScalingGroupName)
	searchArgs.AddConsoleResource(wf, searchutil.ConsoleResource{
		Title:       name,
		ConsolePath: "/ec2autoscaling/home#/details/" + url.PathEscape(name),
		ServiceID:   "ec2",
		ID:          aws.ToString(group.AutoScalingGroupARN),
		IDPrefix:    searchutil.ARNPrefix,
	}).Subtitle(subtitleFrom(
		capacityDetail(group),
		sizeRangeDetail(group),
		aws.ToString(group.LaunchConfigurationName),
	))
}

func capacityDetail(group types.AutoScalingGroup) string {
	if group.DesiredCapacity == nil {
		return ""
	}
	return fmt.Sprintf("%d/%d instances", len(group.Instances), *group.DesiredCapacity)
}

func sizeRangeDetail(group types.AutoScalingGroup) string {
	if group.MinSize == nil || group.MaxSize == nil {
		return ""
	}
	return fmt.Sprintf("min %d – max %d", *group.MinSize, *group.MaxSize)
}
