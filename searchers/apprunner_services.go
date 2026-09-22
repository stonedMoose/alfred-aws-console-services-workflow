package searchers

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/apprunner"
	"github.com/aws/aws-sdk-go-v2/service/apprunner/types"
	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/awsworkflow"
	"github.com/rkoval/alfred-aws-console-services-workflow/caching"
	"github.com/rkoval/alfred-aws-console-services-workflow/searchers/searchutil"
	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

type AppRunnerServiceSearcher struct{}

func (s AppRunnerServiceSearcher) Search(wf *aw.Workflow, searchArgs searchutil.SearchArgs) error {
	cacheName := util.GetCurrentFilename()
	entities := caching.LoadEntityArrayFromCache(wf, searchArgs, cacheName, s.fetch)
	for _, entity := range entities {
		s.addToWorkflow(wf, searchArgs, entity)
	}
	return nil
}

func (s AppRunnerServiceSearcher) fetch(cfg aws.Config) ([]types.ServiceSummary, error) {
	svc := apprunner.NewFromConfig(cfg)

	var entities []types.ServiceSummary
	nextToken := ""
	for {
		params := &apprunner.ListServicesInput{
			MaxResults: aws.Int32(20), // max allowed by this API
		}
		if nextToken != "" {
			params.NextToken = aws.String(nextToken)
		}
		resp, err := svc.ListServices(context.TODO(), params)
		if err != nil {
			return nil, err
		}

		entities = append(entities, resp.ServiceSummaryList...)

		if resp.NextToken != nil && *resp.NextToken != "" {
			nextToken = *resp.NextToken
		} else {
			break
		}
	}

	return entities, nil
}

func (s AppRunnerServiceSearcher) addToWorkflow(wf *aw.Workflow, searchArgs searchutil.SearchArgs, entity types.ServiceSummary) {
	title := *entity.ServiceName

	subtitleArray := []string{}
	if entity.Status != "" {
		subtitleArray = append(subtitleArray, strings.ToLower(string(entity.Status)))
	}
	subtitleArray = util.AppendString(subtitleArray, entity.ServiceUrl)
	if entity.UpdatedAt != nil {
		subtitleArray = append(subtitleArray, "Updated "+entity.UpdatedAt.Format(time.UnixDate))
	}
	subtitle := strings.Join(subtitleArray, " – ")

	path := fmt.Sprintf("/apprunner/home#/services/%s", *entity.ServiceArn)
	item := util.NewURLItem(wf, title).
		Subtitle(subtitle).
		Arg(util.ConstructAWSConsoleUrl(path, searchArgs.GetRegion())).
		Icon(awsworkflow.GetImageIcon("apprunner")).
		Valid(true)

	searchArgs.AddMatch(item, "arn:", *entity.ServiceArn, title)
}
