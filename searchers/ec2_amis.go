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

type EC2ImageSearcher struct{}

func (s EC2ImageSearcher) Search(wf *aw.Workflow, searchArgs searchutil.SearchArgs) error {
	cacheName := util.GetCurrentFilename()
	entities := caching.LoadEntityArrayFromCache(wf, searchArgs, cacheName, s.fetch)
	for _, entity := range entities {
		s.addToWorkflow(wf, searchArgs, entity)
	}
	return nil
}

// restricted to images this account owns; the full public AMI catalog is
// hundreds of thousands of entries
func (s EC2ImageSearcher) fetch(cfg aws.Config) ([]types.Image, error) {
	svc := ec2.NewFromConfig(cfg)

	var entities []types.Image
	nextToken := ""
	for {
		params := &ec2.DescribeImagesInput{
			MaxResults: aws.Int32(1000), // get as many as we can
			Owners:     []string{"self"},
		}
		if nextToken != "" {
			params.NextToken = aws.String(nextToken)
		}
		resp, err := svc.DescribeImages(context.TODO(), params)
		if err != nil {
			return nil, err
		}

		entities = append(entities, resp.Images...)

		if resp.NextToken != nil && *resp.NextToken != "" {
			nextToken = *resp.NextToken
		} else {
			break
		}
	}

	return entities, nil
}

func (s EC2ImageSearcher) addToWorkflow(wf *aw.Workflow, searchArgs searchutil.SearchArgs, entity types.Image) {
	id := *entity.ImageId
	title := id

	subtitleArray := []string{}
	if entity.Name != nil && *entity.Name != "" {
		title = *entity.Name
		subtitleArray = append(subtitleArray, id)
	}
	subtitleArray = util.AppendString(subtitleArray, entity.Description)
	if entity.Architecture != "" {
		subtitleArray = append(subtitleArray, string(entity.Architecture))
	}
	subtitleArray = util.AppendString(subtitleArray, entity.CreationDate)
	subtitle := strings.Join(subtitleArray, " – ")

	path := fmt.Sprintf("/ec2/home#ImageDetails:imageId=%s", id)
	item := util.NewURLItem(wf, title).
		Subtitle(subtitle).
		Arg(util.ConstructAWSConsoleUrl(path, searchArgs.GetRegion())).
		Icon(awsworkflow.GetImageIcon("ec2")).
		Valid(true)

	searchArgs.AddMatch(item, "", "", title)
}
