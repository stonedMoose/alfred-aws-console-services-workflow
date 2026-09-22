package searchers

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront/types"
	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/awsworkflow"
	"github.com/rkoval/alfred-aws-console-services-workflow/caching"
	"github.com/rkoval/alfred-aws-console-services-workflow/searchers/searchutil"
	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

type CloudFrontDistributionSearcher struct{}

func (s CloudFrontDistributionSearcher) Search(wf *aw.Workflow, searchArgs searchutil.SearchArgs) error {
	cacheName := util.GetCurrentFilename()
	entities := caching.LoadEntityArrayFromCache(wf, searchArgs, cacheName, s.fetch)
	for _, entity := range entities {
		s.addToWorkflow(wf, searchArgs, entity)
	}
	return nil
}

func (s CloudFrontDistributionSearcher) fetch(cfg aws.Config) ([]types.DistributionSummary, error) {
	svc := cloudfront.NewFromConfig(cfg)

	var entities []types.DistributionSummary
	marker := ""
	for {
		params := &cloudfront.ListDistributionsInput{
			MaxItems: aws.Int32(100), // max allowed by this API
		}
		if marker != "" {
			params.Marker = aws.String(marker)
		}
		resp, err := svc.ListDistributions(context.TODO(), params)
		if err != nil {
			return nil, err
		}

		if resp.DistributionList == nil {
			break
		}
		entities = append(entities, resp.DistributionList.Items...)

		if resp.DistributionList.IsTruncated != nil && *resp.DistributionList.IsTruncated && resp.DistributionList.NextMarker != nil {
			marker = *resp.DistributionList.NextMarker
		} else {
			break
		}
	}

	return entities, nil
}

func (s CloudFrontDistributionSearcher) addToWorkflow(wf *aw.Workflow, searchArgs searchutil.SearchArgs, entity types.DistributionSummary) {
	id := *entity.Id
	title := id

	subtitleArray := []string{}
	for _, alias := range entity.Aliases.Items {
		subtitleArray = append(subtitleArray, alias)
	}
	if len(subtitleArray) > 0 {
		title = subtitleArray[0]
		subtitleArray = append([]string{id}, subtitleArray[1:]...)
	}
	subtitleArray = util.AppendString(subtitleArray, entity.DomainName)
	if entity.Enabled != nil && !*entity.Enabled {
		subtitleArray = append(subtitleArray, "disabled")
	}
	subtitle := strings.Join(subtitleArray, " – ")

	path := fmt.Sprintf("/cloudfront/v4/home#/distributions/%s", id)
	item := util.NewURLItem(wf, title).
		Subtitle(subtitle).
		Arg(util.ConstructAWSConsoleUrl(path, searchArgs.GetRegion())).
		Icon(awsworkflow.GetImageIcon("cloudfront")).
		Valid(true)

	arn := ""
	if entity.ARN != nil {
		arn = *entity.ARN
	}
	searchArgs.AddMatch(item, "arn:", arn, title)
}
