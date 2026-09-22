package searchers

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/awsworkflow"
	"github.com/rkoval/alfred-aws-console-services-workflow/caching"
	"github.com/rkoval/alfred-aws-console-services-workflow/searchers/searchutil"
	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

type VPCSearcher struct{}

func (s VPCSearcher) Search(wf *aw.Workflow, searchArgs searchutil.SearchArgs) error {
	cacheName := util.GetCurrentFilename()
	entities := caching.LoadEntityArrayFromCache(wf, searchArgs, cacheName, s.fetch)
	for _, entity := range entities {
		s.addToWorkflow(wf, searchArgs, entity)
	}
	return nil
}

func (s VPCSearcher) fetch(cfg aws.Config) ([]types.Vpc, error) {
	svc := ec2.NewFromConfig(cfg)

	var entities []types.Vpc
	nextToken := ""
	for {
		params := &ec2.DescribeVpcsInput{
			MaxResults: aws.Int32(1000), // get as many as we can
		}
		if nextToken != "" {
			params.NextToken = aws.String(nextToken)
		}
		resp, err := svc.DescribeVpcs(context.TODO(), params)
		if err != nil {
			return nil, err
		}

		entities = append(entities, resp.Vpcs...)

		if resp.NextToken != nil && *resp.NextToken != "" {
			nextToken = *resp.NextToken
		} else {
			break
		}
	}

	return entities, nil
}

func (s VPCSearcher) addToWorkflow(wf *aw.Workflow, searchArgs searchutil.SearchArgs, entity types.Vpc) {
	id := *entity.VpcId
	title := id

	subtitleArray := []string{}
	if name := util.GetEC2TagValue(entity.Tags, "Name"); name != "" {
		title = name
		subtitleArray = append(subtitleArray, id)
	}
	subtitleArray = util.AppendString(subtitleArray, entity.CidrBlock)
	if entity.IsDefault != nil && *entity.IsDefault {
		subtitleArray = append(subtitleArray, "default")
	}
	subtitle := strings.Join(subtitleArray, " – ")

	path := fmt.Sprintf("/vpc/home#VpcDetails:VpcId=%s", id)
	item := util.NewURLItem(wf, title).
		Subtitle(subtitle).
		Arg(util.ConstructAWSConsoleUrl(path, searchArgs.GetRegion())).
		Icon(awsworkflow.GetImageIcon("vpc")).
		Valid(true)

	searchArgs.AddMatch(item, "", "", title)
}
