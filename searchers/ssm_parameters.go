package searchers

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/awsworkflow"
	"github.com/rkoval/alfred-aws-console-services-workflow/caching"
	"github.com/rkoval/alfred-aws-console-services-workflow/searchers/searchutil"
	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

type SSMParameterSearcher struct{}

func (s SSMParameterSearcher) Search(wf *aw.Workflow, searchArgs searchutil.SearchArgs) error {
	cacheName := util.GetCurrentFilename()
	entities := caching.LoadEntityArrayFromCache(wf, searchArgs, cacheName, s.fetch)
	for _, entity := range entities {
		s.addToWorkflow(wf, searchArgs, entity)
	}
	return nil
}

func (s SSMParameterSearcher) fetch(cfg aws.Config) ([]types.ParameterMetadata, error) {
	svc := ssm.NewFromConfig(cfg)

	var entities []types.ParameterMetadata
	nextToken := ""
	for {
		params := &ssm.DescribeParametersInput{
			MaxResults: aws.Int32(50), // max allowed by this API
		}
		if nextToken != "" {
			params.NextToken = aws.String(nextToken)
		}
		resp, err := svc.DescribeParameters(context.TODO(), params)
		if err != nil {
			return nil, err
		}

		entities = append(entities, resp.Parameters...)

		if resp.NextToken != nil && *resp.NextToken != "" {
			nextToken = *resp.NextToken
		} else {
			break
		}
	}

	return entities, nil
}

func (s SSMParameterSearcher) addToWorkflow(wf *aw.Workflow, searchArgs searchutil.SearchArgs, entity types.ParameterMetadata) {
	title := *entity.Name

	subtitleArray := []string{}
	if entity.Type != "" {
		subtitleArray = append(subtitleArray, string(entity.Type))
	}
	subtitleArray = util.AppendString(subtitleArray, entity.Description)
	if entity.LastModifiedDate != nil {
		subtitleArray = append(subtitleArray, "Modified "+entity.LastModifiedDate.Format(time.UnixDate))
	}
	subtitle := strings.Join(subtitleArray, " – ")

	// the console expects the parameter name double-encoded, so `/db/password`
	// must be written as `%252Fdb%252Fpassword`
	path := fmt.Sprintf("/systems-manager/parameters/%s/description", url.PathEscape(url.PathEscape(title)))
	item := util.NewURLItem(wf, title).
		Subtitle(subtitle).
		Arg(util.ConstructAWSConsoleUrl(path, searchArgs.GetRegion())).
		Icon(awsworkflow.GetImageIcon("systemsmanager")).
		Valid(true)

	arn := ""
	if entity.ARN != nil {
		arn = *entity.ARN
	}
	searchArgs.AddMatch(item, "arn:", arn, title)
}
