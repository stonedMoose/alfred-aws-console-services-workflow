package searchers

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/elasticache"
	"github.com/aws/aws-sdk-go-v2/service/elasticache/types"
	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/awspaging"
	"github.com/rkoval/alfred-aws-console-services-workflow/searchers/searchutil"
)

// ElastiCache lists Redis and Memcached clusters through the same API; the
// searcher of each engine fetches them all and presents its own.

func fetchCacheClusters(cfg aws.Config) ([]types.CacheCluster, error) {
	client := elasticache.NewFromConfig(cfg)
	return awspaging.FetchAllPages(func(marker string) ([]types.CacheCluster, *string, error) {
		resp, err := client.DescribeCacheClusters(context.TODO(), &elasticache.DescribeCacheClustersInput{
			MaxRecords: aws.Int32(100),
			Marker:     awspaging.TokenOrNil(marker),
		})
		if err != nil {
			return nil, nil, err
		}
		return resp.CacheClusters, resp.Marker, nil
	})
}

func addCacheClusterToWorkflow(engine string, wf *aw.Workflow, searchArgs searchutil.SearchArgs, cluster types.CacheCluster) {
	if aws.ToString(cluster.Engine) != engine {
		return
	}
	title := aws.ToString(cluster.ARN)
	if cluster.CacheClusterId != nil {
		title = *cluster.CacheClusterId
	}
	searchArgs.AddConsoleResource(wf, searchutil.ConsoleResource{
		Title:       title,
		ConsolePath: "/elasticache/home#" + cacheClusterConsolePage(engine, cluster) + ":id=" + aws.ToString(cluster.CacheClusterId),
		ServiceID:   "elasticache",
		ID:          aws.ToString(cluster.ARN),
		IDPrefix:    searchutil.ARNPrefix,
	}).Subtitle(subtitleFrom(
		joinKnown(" ", aws.ToString(cluster.Engine), aws.ToString(cluster.EngineVersion)),
		aws.ToString(cluster.CacheNodeType),
		aws.ToString(cluster.CacheClusterStatus),
	))
}

// cacheClusterConsolePage is the detail page of a cluster, which differs for
// the members of a replication group.
func cacheClusterConsolePage(engine string, cluster types.CacheCluster) string {
	if cluster.ReplicationGroupId != nil {
		return engine + "-group-detail"
	}
	return engine + "-detail"
}
