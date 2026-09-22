package searchers

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling/types"
	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/awsworkflow"
	"github.com/rkoval/alfred-aws-console-services-workflow/caching"
	"github.com/rkoval/alfred-aws-console-services-workflow/searchers/searchutil"
	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

type EC2AutoScalingGroupSearcher struct{}

func (s EC2AutoScalingGroupSearcher) Search(wf *aw.Workflow, searchArgs searchutil.SearchArgs) error {
	cacheName := util.GetCurrentFilename()
	entities := caching.LoadEntityArrayFromCache(wf, searchArgs, cacheName, s.fetch)
	for _, entity := range entities {
		s.addToWorkflow(wf, searchArgs, entity)
	}
	return nil
}

func (s EC2AutoScalingGroupSearcher) fetch(cfg aws.Config) ([]types.AutoScalingGroup, error) {
	svc := autoscaling.NewFromConfig(cfg)

	var entities []types.AutoScalingGroup
	nextToken := ""
	for {
		params := &autoscaling.DescribeAutoScalingGroupsInput{
			MaxRecords: aws.Int32(100), // max allowed by this API
		}
		if nextToken != "" {
			params.NextToken = aws.String(nextToken)
		}
		resp, err := svc.DescribeAutoScalingGroups(context.TODO(), params)
		if err != nil {
			return nil, err
		}

		entities = append(entities, resp.AutoScalingGroups...)

		if resp.NextToken != nil && *resp.NextToken != "" {
			nextToken = *resp.NextToken
		} else {
			break
		}
	}

	return entities, nil
}

func (s EC2AutoScalingGroupSearcher) addToWorkflow(wf *aw.Workflow, searchArgs searchutil.SearchArgs, entity types.AutoScalingGroup) {
	title := *entity.AutoScalingGroupName

	subtitleArray := []string{}
	if entity.DesiredCapacity != nil {
		subtitleArray = append(subtitleArray, strconv.Itoa(len(entity.Instances))+"/"+strconv.Itoa(int(*entity.DesiredCapacity))+" instances")
	}
	if entity.MinSize != nil && entity.MaxSize != nil {
		subtitleArray = append(subtitleArray, "min "+strconv.Itoa(int(*entity.MinSize))+" – max "+strconv.Itoa(int(*entity.MaxSize)))
	}
	subtitleArray = util.AppendString(subtitleArray, entity.LaunchConfigurationName)
	subtitle := strings.Join(subtitleArray, " – ")

	path := fmt.Sprintf("/ec2autoscaling/home#/details/%s", url.PathEscape(title))
	item := util.NewURLItem(wf, title).
		Subtitle(subtitle).
		Arg(util.ConstructAWSConsoleUrl(path, searchArgs.GetRegion())).
		Icon(awsworkflow.GetImageIcon("ec2")).
		Valid(true)

	arn := ""
	if entity.AutoScalingGroupARN != nil {
		arn = *entity.AutoScalingGroupARN
	}
	searchArgs.AddMatch(item, "arn:", arn, title)
}
