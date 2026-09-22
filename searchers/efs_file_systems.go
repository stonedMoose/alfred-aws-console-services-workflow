package searchers

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/efs"
	"github.com/aws/aws-sdk-go-v2/service/efs/types"
	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/awsworkflow"
	"github.com/rkoval/alfred-aws-console-services-workflow/caching"
	"github.com/rkoval/alfred-aws-console-services-workflow/searchers/searchutil"
	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

type EFSFileSystemSearcher struct{}

func (s EFSFileSystemSearcher) Search(wf *aw.Workflow, searchArgs searchutil.SearchArgs) error {
	cacheName := util.GetCurrentFilename()
	entities := caching.LoadEntityArrayFromCache(wf, searchArgs, cacheName, s.fetch)
	for _, entity := range entities {
		s.addToWorkflow(wf, searchArgs, entity)
	}
	return nil
}

func (s EFSFileSystemSearcher) fetch(cfg aws.Config) ([]types.FileSystemDescription, error) {
	svc := efs.NewFromConfig(cfg)

	var entities []types.FileSystemDescription
	marker := ""
	for {
		params := &efs.DescribeFileSystemsInput{
			MaxItems: aws.Int32(100), // max allowed by this API
		}
		if marker != "" {
			params.Marker = aws.String(marker)
		}
		resp, err := svc.DescribeFileSystems(context.TODO(), params)
		if err != nil {
			return nil, err
		}

		entities = append(entities, resp.FileSystems...)

		if resp.NextMarker != nil && *resp.NextMarker != "" {
			marker = *resp.NextMarker
		} else {
			break
		}
	}

	return entities, nil
}

func (s EFSFileSystemSearcher) addToWorkflow(wf *aw.Workflow, searchArgs searchutil.SearchArgs, entity types.FileSystemDescription) {
	id := *entity.FileSystemId
	title := id

	subtitleArray := []string{}
	if entity.Name != nil && *entity.Name != "" {
		title = *entity.Name
		subtitleArray = append(subtitleArray, id)
	}
	if entity.LifeCycleState != "" {
		subtitleArray = append(subtitleArray, strings.ToLower(string(entity.LifeCycleState)))
	}
	if entity.SizeInBytes != nil && entity.SizeInBytes.Value != 0 {
		subtitleArray = append(subtitleArray, util.ByteFormat(entity.SizeInBytes.Value, 2))
	}
	subtitle := strings.Join(subtitleArray, " – ")

	path := fmt.Sprintf("/efs/home#/file-systems/%s", id)
	item := util.NewURLItem(wf, title).
		Subtitle(subtitle).
		Arg(util.ConstructAWSConsoleUrl(path, searchArgs.GetRegion())).
		Icon(awsworkflow.GetImageIcon("efs")).
		Valid(true)

	arn := ""
	if entity.FileSystemArn != nil {
		arn = *entity.FileSystemArn
	}
	searchArgs.AddMatch(item, "arn:", arn, title)
}
