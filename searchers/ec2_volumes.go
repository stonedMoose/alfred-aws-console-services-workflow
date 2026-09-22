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

type EC2VolumeSearcher struct{}

func (s EC2VolumeSearcher) Search(wf *aw.Workflow, searchArgs searchutil.SearchArgs) error {
	return searchEntities(wf, searchArgs, "ec2_volumes", s.fetch, s.addToWorkflow)
}

func (EC2VolumeSearcher) fetch(cfg aws.Config) ([]types.Volume, error) {
	client := ec2.NewFromConfig(cfg)
	return awspaging.FetchAllPages(func(pageToken string) ([]types.Volume, *string, error) {
		resp, err := client.DescribeVolumes(context.TODO(), &ec2.DescribeVolumesInput{
			MaxResults: aws.Int32(500), // max allowed by this API
			NextToken:  awspaging.TokenOrNil(pageToken),
		})
		if err != nil {
			return nil, nil, err
		}
		return resp.Volumes, resp.NextToken, nil
	})
}

func (EC2VolumeSearcher) addToWorkflow(wf *aw.Workflow, searchArgs searchutil.SearchArgs, volume types.Volume) {
	id := aws.ToString(volume.VolumeId)
	title, idDetail := namedOrID(ec2NameTag(volume.Tags), id)
	searchArgs.AddConsoleResource(wf, searchutil.ConsoleResource{
		Title:       title,
		ConsolePath: "/ec2/home#VolumeDetails:volumeId=" + id,
		ServiceID:   "ec2",
	}).Subtitle(subtitleFrom(
		idDetail,
		formatIfSet("%d GiB", volume.Size),
		string(volume.VolumeType),
		string(volume.State),
	))
}
