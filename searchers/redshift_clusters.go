package searchers

import (
	"context"
	"net/url"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/redshift"
	"github.com/aws/aws-sdk-go-v2/service/redshift/types"
	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/awspaging"
	"github.com/rkoval/alfred-aws-console-services-workflow/searchers/searchutil"
)

type RedshiftClusterSearcher struct{}

func (s RedshiftClusterSearcher) Search(wf *aw.Workflow, searchArgs searchutil.SearchArgs) error {
	return searchEntities(wf, searchArgs, "redshift_clusters", s.fetch, s.addToWorkflow)
}

func (RedshiftClusterSearcher) fetch(cfg aws.Config) ([]types.Cluster, error) {
	client := redshift.NewFromConfig(cfg)
	return awspaging.FetchAllPages(func(marker string) ([]types.Cluster, *string, error) {
		resp, err := client.DescribeClusters(context.TODO(), &redshift.DescribeClustersInput{
			MaxRecords: aws.Int32(100), // max allowed by this API
			Marker:     awspaging.TokenOrNil(marker),
		})
		if err != nil {
			return nil, nil, err
		}
		return resp.Clusters, resp.Marker, nil
	})
}

func (RedshiftClusterSearcher) addToWorkflow(wf *aw.Workflow, searchArgs searchutil.SearchArgs, cluster types.Cluster) {
	id := aws.ToString(cluster.ClusterIdentifier)
	searchArgs.AddConsoleResource(wf, searchutil.ConsoleResource{
		Title:       id,
		ConsolePath: "/redshiftv2/home#cluster-details?cluster=" + url.QueryEscape(id),
		ServiceID:   "redshift",
		ID:          aws.ToString(cluster.ClusterNamespaceArn),
		IDPrefix:    searchutil.ARNPrefix,
	}).Subtitle(subtitleFrom(
		aws.ToString(cluster.ClusterStatus),
		aws.ToString(cluster.NodeType),
		formatIfSet("%d nodes", cluster.NumberOfNodes),
	))
}
