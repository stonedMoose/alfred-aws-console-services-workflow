package searchers

import (
	"context"
	"fmt"
	"strconv"
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

type EC2VolumeSearcher struct{}

func (s EC2VolumeSearcher) Search(wf *aw.Workflow, searchArgs searchutil.SearchArgs) error {
	cacheName := util.GetCurrentFilename()
	entities := caching.LoadEntityArrayFromCache(wf, searchArgs, cacheName, s.fetch)
	for _, entity := range entities {
		s.addToWorkflow(wf, searchArgs, entity)
	}
	return nil
}

func (s EC2VolumeSearcher) fetch(cfg aws.Config) ([]types.Volume, error) {
	svc := ec2.NewFromConfig(cfg)

	var entities []types.Volume
	nextToken := ""
	for {
		params := &ec2.DescribeVolumesInput{
			MaxResults: aws.Int32(500), // max allowed by this API
		}
		if nextToken != "" {
			params.NextToken = aws.String(nextToken)
		}
		resp, err := svc.DescribeVolumes(context.TODO(), params)
		if err != nil {
			return nil, err
		}

		entities = append(entities, resp.Volumes...)

		if resp.NextToken != nil && *resp.NextToken != "" {
			nextToken = *resp.NextToken
		} else {
			break
		}
	}

	return entities, nil
}

func (s EC2VolumeSearcher) addToWorkflow(wf *aw.Workflow, searchArgs searchutil.SearchArgs, entity types.Volume) {
	id := *entity.VolumeId
	title := id

	subtitleArray := []string{}
	if name := util.GetEC2TagValue(entity.Tags, "Name"); name != "" {
		title = name
		subtitleArray = append(subtitleArray, id)
	}
	if entity.Size != nil {
		subtitleArray = append(subtitleArray, strconv.Itoa(int(*entity.Size))+" GiB")
	}
	if entity.VolumeType != "" {
		subtitleArray = append(subtitleArray, string(entity.VolumeType))
	}
	if entity.State != "" {
		subtitleArray = append(subtitleArray, string(entity.State))
	}
	subtitle := strings.Join(subtitleArray, " – ")

	path := fmt.Sprintf("/ec2/home#VolumeDetails:volumeId=%s", id)
	item := util.NewURLItem(wf, title).
		Subtitle(subtitle).
		Arg(util.ConstructAWSConsoleUrl(path, searchArgs.GetRegion())).
		Icon(awsworkflow.GetImageIcon("ec2")).
		Valid(true)

	searchArgs.AddMatch(item, "", "", title)
}
