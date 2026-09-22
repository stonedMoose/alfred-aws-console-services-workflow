package searchers

import (
	"context"
	"net/url"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/aws/aws-sdk-go-v2/service/ecr/types"
	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/awspaging"
	"github.com/rkoval/alfred-aws-console-services-workflow/searchers/searchutil"
)

type ECRRepositorySearcher struct{}

func (s ECRRepositorySearcher) Search(wf *aw.Workflow, searchArgs searchutil.SearchArgs) error {
	return searchEntities(wf, searchArgs, "ecr_repositories", s.fetch, s.addToWorkflow)
}

func (ECRRepositorySearcher) fetch(cfg aws.Config) ([]types.Repository, error) {
	client := ecr.NewFromConfig(cfg)
	return awspaging.FetchAllPages(func(pageToken string) ([]types.Repository, *string, error) {
		resp, err := client.DescribeRepositories(context.TODO(), &ecr.DescribeRepositoriesInput{
			MaxResults: aws.Int32(1000), // max allowed by this API
			NextToken:  awspaging.TokenOrNil(pageToken),
		})
		if err != nil {
			return nil, nil, err
		}
		return resp.Repositories, resp.NextToken, nil
	})
}

func (ECRRepositorySearcher) addToWorkflow(wf *aw.Workflow, searchArgs searchutil.SearchArgs, repository types.Repository) {
	name := aws.ToString(repository.RepositoryName)
	searchArgs.AddConsoleResource(wf, searchutil.ConsoleResource{
		Title:       name,
		ConsolePath: "/ecr/repositories/private/" + aws.ToString(repository.RegistryId) + "/" + url.PathEscape(name),
		ServiceID:   "ecr",
		ID:          aws.ToString(repository.RepositoryArn),
		IDPrefix:    searchutil.ARNPrefix,
	}).Subtitle(subtitleFrom(
		aws.ToString(repository.RepositoryUri),
		strings.ToLower(string(repository.ImageTagMutability)),
		dateDetail("Created", repository.CreatedAt),
	))
}
