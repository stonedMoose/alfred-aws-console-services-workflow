package searchers

import (
	"context"
	"net/url"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/awspaging"
	"github.com/rkoval/alfred-aws-console-services-workflow/searchers/searchutil"
)

type EKSClusterSearcher struct{}

func (s EKSClusterSearcher) Search(wf *aw.Workflow, searchArgs searchutil.SearchArgs) error {
	return searchEntities(wf, searchArgs, "eks_clusters", s.fetch, s.addToWorkflow)
}

// fetch lists cluster names only: describing each one to enrich the subtitle
// would cost an API call per cluster.
func (EKSClusterSearcher) fetch(cfg aws.Config) ([]string, error) {
	client := eks.NewFromConfig(cfg)
	return awspaging.FetchAllPages(func(pageToken string) ([]string, *string, error) {
		resp, err := client.ListClusters(context.TODO(), &eks.ListClustersInput{
			MaxResults: aws.Int32(100), // max allowed by this API
			NextToken:  awspaging.TokenOrNil(pageToken),
		})
		if err != nil {
			return nil, nil, err
		}
		return resp.Clusters, resp.NextToken, nil
	})
}

func (EKSClusterSearcher) addToWorkflow(wf *aw.Workflow, searchArgs searchutil.SearchArgs, clusterName string) {
	searchArgs.AddConsoleResource(wf, searchutil.ConsoleResource{
		Title:       clusterName,
		ConsolePath: "/eks/home#/clusters/" + url.PathEscape(clusterName),
		ServiceID:   "eks",
	})
}
