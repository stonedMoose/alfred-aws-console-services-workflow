package searchers

import (
	"context"
	"log"
	"slices"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/searchers/searchutil"
)

// describeClustersBatchSize is the most clusters DescribeClusters accepts per call.
const describeClustersBatchSize = 100

type ECSClusterSearcher struct{}

func (s ECSClusterSearcher) Search(wf *aw.Workflow, searchArgs searchutil.SearchArgs) error {
	return searchEntities(wf, searchArgs, "ecs_clusters", s.fetch, s.addToWorkflow)
}

// fetch lists the cluster ARNs, then describes them in batches to get their
// names and tags.
func (ECSClusterSearcher) fetch(cfg aws.Config) ([]types.Cluster, error) {
	client := ecs.NewFromConfig(cfg)
	clusterARNs, err := listClusterARNs(client)
	if err != nil {
		return nil, err
	}
	if len(clusterARNs) == 0 {
		return []types.Cluster{}, nil
	}
	return describeClusters(client, clusterARNs), nil
}

func listClusterARNs(client *ecs.Client) ([]string, error) {
	var clusterARNs []string
	paginator := ecs.NewListClustersPaginator(client, &ecs.ListClustersInput{})
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(context.TODO())
		if err != nil {
			return nil, err
		}
		clusterARNs = append(clusterARNs, page.ClusterArns...)
	}
	return clusterARNs, nil
}

// describeClusters skips a batch that cannot be described so that the other
// clusters still show.
func describeClusters(client *ecs.Client, clusterARNs []string) []types.Cluster {
	var clusters []types.Cluster
	for batch := range slices.Chunk(clusterARNs, describeClustersBatchSize) {
		resp, err := client.DescribeClusters(context.TODO(), &ecs.DescribeClustersInput{
			Clusters: batch,
			Include:  []types.ClusterField{types.ClusterFieldTags},
		})
		if err != nil {
			log.Printf("skipping %d clusters that could not be described: %v", len(batch), err)
			continue
		}
		clusters = append(clusters, resp.Clusters...)
	}
	return clusters
}

func (ECSClusterSearcher) addToWorkflow(wf *aw.Workflow, searchArgs searchutil.SearchArgs, cluster types.Cluster) {
	arn := aws.ToString(cluster.ClusterArn)
	searchArgs.AddConsoleResource(wf, searchutil.ConsoleResource{
		Title:       aws.ToString(cluster.ClusterName),
		ConsolePath: "/ecs/v2/clusters/" + clusterNameOf(arn) + "/services?region=" + searchArgs.GetRegion(),
		ServiceID:   "ecs",
		ID:          arn,
		IDPrefix:    searchutil.ARNPrefix,
	}).Subtitle(arn)
}

// clusterNameOf is the name the console addresses a cluster by: the part of
// its ARN after the last slash.
func clusterNameOf(arn string) string {
	return arn[strings.LastIndex(arn, "/")+1:]
}
