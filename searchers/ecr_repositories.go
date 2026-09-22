package searchers

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/aws/aws-sdk-go-v2/service/ecr/types"
	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/awsworkflow"
	"github.com/rkoval/alfred-aws-console-services-workflow/caching"
	"github.com/rkoval/alfred-aws-console-services-workflow/searchers/searchutil"
	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

type ECRRepositorySearcher struct{}

func (s ECRRepositorySearcher) Search(wf *aw.Workflow, searchArgs searchutil.SearchArgs) error {
	cacheName := util.GetCurrentFilename()
	entities := caching.LoadEntityArrayFromCache(wf, searchArgs, cacheName, s.fetch)
	for _, entity := range entities {
		s.addToWorkflow(wf, searchArgs, entity)
	}
	return nil
}

func (s ECRRepositorySearcher) fetch(cfg aws.Config) ([]types.Repository, error) {
	svc := ecr.NewFromConfig(cfg)

	var entities []types.Repository
	nextToken := ""
	for {
		params := &ecr.DescribeRepositoriesInput{
			MaxResults: aws.Int32(1000), // max allowed by this API
		}
		if nextToken != "" {
			params.NextToken = aws.String(nextToken)
		}
		resp, err := svc.DescribeRepositories(context.TODO(), params)
		if err != nil {
			return nil, err
		}

		entities = append(entities, resp.Repositories...)

		if resp.NextToken != nil && *resp.NextToken != "" {
			nextToken = *resp.NextToken
		} else {
			break
		}
	}

	return entities, nil
}

func (s ECRRepositorySearcher) addToWorkflow(wf *aw.Workflow, searchArgs searchutil.SearchArgs, entity types.Repository) {
	title := *entity.RepositoryName

	subtitleArray := []string{}
	subtitleArray = util.AppendString(subtitleArray, entity.RepositoryUri)
	if entity.ImageTagMutability != "" {
		subtitleArray = append(subtitleArray, strings.ToLower(string(entity.ImageTagMutability)))
	}
	if entity.CreatedAt != nil {
		subtitleArray = append(subtitleArray, "Created "+entity.CreatedAt.Format(time.UnixDate))
	}
	subtitle := strings.Join(subtitleArray, " – ")

	registryId := ""
	if entity.RegistryId != nil {
		registryId = *entity.RegistryId
	}
	path := fmt.Sprintf("/ecr/repositories/private/%s/%s", registryId, url.PathEscape(title))
	item := util.NewURLItem(wf, title).
		Subtitle(subtitle).
		Arg(util.ConstructAWSConsoleUrl(path, searchArgs.GetRegion())).
		Icon(awsworkflow.GetImageIcon("ecr")).
		Valid(true)

	arn := ""
	if entity.RepositoryArn != nil {
		arn = *entity.RepositoryArn
	}
	searchArgs.AddMatch(item, "arn:", arn, title)
}
