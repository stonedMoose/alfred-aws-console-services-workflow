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

type EC2SnapshotSearcher struct{}

func (s EC2SnapshotSearcher) Search(wf *aw.Workflow, searchArgs searchutil.SearchArgs) error {
	cacheName := util.GetCurrentFilename()
	entities := caching.LoadEntityArrayFromCache(wf, searchArgs, cacheName, s.fetch)
	for _, entity := range entities {
		s.addToWorkflow(wf, searchArgs, entity)
	}
	return nil
}

// restricted to snapshots this account owns; every public snapshot in the
// region would otherwise come back
func (s EC2SnapshotSearcher) fetch(cfg aws.Config) ([]types.Snapshot, error) {
	svc := ec2.NewFromConfig(cfg)

	var entities []types.Snapshot
	nextToken := ""
	for {
		params := &ec2.DescribeSnapshotsInput{
			MaxResults: aws.Int32(1000), // get as many as we can
			OwnerIds:   []string{"self"},
		}
		if nextToken != "" {
			params.NextToken = aws.String(nextToken)
		}
		resp, err := svc.DescribeSnapshots(context.TODO(), params)
		if err != nil {
			return nil, err
		}

		entities = append(entities, resp.Snapshots...)

		if resp.NextToken != nil && *resp.NextToken != "" {
			nextToken = *resp.NextToken
		} else {
			break
		}
	}

	return entities, nil
}

func (s EC2SnapshotSearcher) addToWorkflow(wf *aw.Workflow, searchArgs searchutil.SearchArgs, entity types.Snapshot) {
	id := *entity.SnapshotId
	title := id

	subtitleArray := []string{}
	if name := util.GetEC2TagValue(entity.Tags, "Name"); name != "" {
		title = name
		subtitleArray = append(subtitleArray, id)
	}
	subtitleArray = util.AppendString(subtitleArray, entity.Description)
	if entity.VolumeSize != nil {
		subtitleArray = append(subtitleArray, strconv.Itoa(int(*entity.VolumeSize))+" GiB")
	}
	if entity.State != "" {
		subtitleArray = append(subtitleArray, string(entity.State))
	}
	subtitle := strings.Join(subtitleArray, " – ")

	path := fmt.Sprintf("/ec2/home#SnapshotDetails:snapshotId=%s", id)
	item := util.NewURLItem(wf, title).
		Subtitle(subtitle).
		Arg(util.ConstructAWSConsoleUrl(path, searchArgs.GetRegion())).
		Icon(awsworkflow.GetImageIcon("ec2")).
		Valid(true)

	searchArgs.AddMatch(item, "", "", title)
}
