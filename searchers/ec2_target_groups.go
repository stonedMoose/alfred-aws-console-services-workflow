package searchers

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2/types"
	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/awsworkflow"
	"github.com/rkoval/alfred-aws-console-services-workflow/caching"
	"github.com/rkoval/alfred-aws-console-services-workflow/searchers/searchutil"
	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

type EC2TargetGroupSearcher struct{}

func (s EC2TargetGroupSearcher) Search(wf *aw.Workflow, searchArgs searchutil.SearchArgs) error {
	cacheName := util.GetCurrentFilename()
	entities := caching.LoadEntityArrayFromCache(wf, searchArgs, cacheName, s.fetch)
	for _, entity := range entities {
		s.addToWorkflow(wf, searchArgs, entity)
	}
	return nil
}

func (s EC2TargetGroupSearcher) fetch(cfg aws.Config) ([]types.TargetGroup, error) {
	client := elasticloadbalancingv2.NewFromConfig(cfg)

	entities := []types.TargetGroup{}
	pageToken := ""
	for {
		params := &elasticloadbalancingv2.DescribeTargetGroupsInput{
			PageSize: aws.Int32(400),
		}
		if pageToken != "" {
			params.Marker = &pageToken
		}
		resp, err := client.DescribeTargetGroups(context.TODO(), params)

		if err != nil {
			return nil, err
		}

		entities = append(entities, resp.TargetGroups...)

		if resp.NextMarker != nil {
			pageToken = *resp.NextMarker
		} else {
			break
		}
	}

	return entities, nil
}

func (s EC2TargetGroupSearcher) addToWorkflow(wf *aw.Workflow, searchArgs searchutil.SearchArgs, entity types.TargetGroup) {
	title := *entity.TargetGroupName

	subtitleArray := []string{}
	if entity.TargetType != "" {
		subtitleArray = append(subtitleArray, string(entity.TargetType))
	}
	if entity.Protocol != "" && entity.Port != nil {
		subtitleArray = append(subtitleArray, string(entity.Protocol)+":"+strconv.Itoa(int(*entity.Port)))
	}
	subtitleArray = util.AppendString(subtitleArray, entity.VpcId)
	subtitle := strings.Join(subtitleArray, " – ")

	path := fmt.Sprintf("/ec2/home#TargetGroup:targetGroupArn=%s", *entity.TargetGroupArn)
	item := util.NewURLItem(wf, title).
		Subtitle(subtitle).
		Arg(util.ConstructAWSConsoleUrl(path, searchArgs.GetRegion())).
		Icon(awsworkflow.GetImageIcon("ec2")).
		Valid(true)

	searchArgs.AddMatch(item, "arn:", *entity.TargetGroupArn, title)
}
