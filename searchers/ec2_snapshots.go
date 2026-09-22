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

type EC2SnapshotSearcher struct{}

func (s EC2SnapshotSearcher) Search(wf *aw.Workflow, searchArgs searchutil.SearchArgs) error {
	return searchEntities(wf, searchArgs, "ec2_snapshots", s.fetch, s.addToWorkflow)
}

// fetch is restricted to the snapshots this account owns; every public
// snapshot in the region would otherwise come back.
func (EC2SnapshotSearcher) fetch(cfg aws.Config) ([]types.Snapshot, error) {
	client := ec2.NewFromConfig(cfg)
	return awspaging.FetchAllPages(func(pageToken string) ([]types.Snapshot, *string, error) {
		resp, err := client.DescribeSnapshots(context.TODO(), &ec2.DescribeSnapshotsInput{
			MaxResults: aws.Int32(1000), // get as many as we can
			OwnerIds:   []string{"self"},
			NextToken:  awspaging.TokenOrNil(pageToken),
		})
		if err != nil {
			return nil, nil, err
		}
		return resp.Snapshots, resp.NextToken, nil
	})
}

func (EC2SnapshotSearcher) addToWorkflow(wf *aw.Workflow, searchArgs searchutil.SearchArgs, snapshot types.Snapshot) {
	id := aws.ToString(snapshot.SnapshotId)
	title, idDetail := namedOrID(ec2NameTag(snapshot.Tags), id)
	searchArgs.AddConsoleResource(wf, searchutil.ConsoleResource{
		Title:       title,
		ConsolePath: "/ec2/home#SnapshotDetails:snapshotId=" + id,
		ServiceID:   "ec2",
	}).Subtitle(subtitleFrom(
		idDetail,
		aws.ToString(snapshot.Description),
		formatIfSet("%d GiB", snapshot.VolumeSize),
		string(snapshot.State),
	))
}
