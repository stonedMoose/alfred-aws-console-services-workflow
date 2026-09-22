package searchers

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/efs"
	"github.com/aws/aws-sdk-go-v2/service/efs/types"
	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/awspaging"
	"github.com/rkoval/alfred-aws-console-services-workflow/searchers/searchutil"
	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

type EFSFileSystemSearcher struct{}

func (s EFSFileSystemSearcher) Search(wf *aw.Workflow, searchArgs searchutil.SearchArgs) error {
	return searchEntities(wf, searchArgs, "efs_file_systems", s.fetch, s.addToWorkflow)
}

func (EFSFileSystemSearcher) fetch(cfg aws.Config) ([]types.FileSystemDescription, error) {
	client := efs.NewFromConfig(cfg)
	return awspaging.FetchAllPages(func(marker string) ([]types.FileSystemDescription, *string, error) {
		resp, err := client.DescribeFileSystems(context.TODO(), &efs.DescribeFileSystemsInput{
			MaxItems: aws.Int32(100), // max allowed by this API
			Marker:   awspaging.TokenOrNil(marker),
		})
		if err != nil {
			return nil, nil, err
		}
		return resp.FileSystems, resp.NextMarker, nil
	})
}

func (EFSFileSystemSearcher) addToWorkflow(wf *aw.Workflow, searchArgs searchutil.SearchArgs, fileSystem types.FileSystemDescription) {
	id := aws.ToString(fileSystem.FileSystemId)
	title, idDetail := namedOrID(aws.ToString(fileSystem.Name), id)
	searchArgs.AddConsoleResource(wf, searchutil.ConsoleResource{
		Title:       title,
		ConsolePath: "/efs/home#/file-systems/" + id,
		ServiceID:   "efs",
		ID:          aws.ToString(fileSystem.FileSystemArn),
		IDPrefix:    searchutil.ARNPrefix,
	}).Subtitle(subtitleFrom(
		idDetail,
		strings.ToLower(string(fileSystem.LifeCycleState)),
		fileSystemSizeDetail(fileSystem.SizeInBytes),
	))
}

func fileSystemSizeDetail(size *types.FileSystemSize) string {
	if size == nil || size.Value == 0 {
		return ""
	}
	return util.ByteFormat(size.Value, 2)
}
