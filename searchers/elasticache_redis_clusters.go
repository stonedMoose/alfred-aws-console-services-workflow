package searchers

import (
	"github.com/aws/aws-sdk-go-v2/service/elasticache/types"
	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/searchers/searchutil"
)

type ElasticacheRedisClusterSearcher struct{}

func (s ElasticacheRedisClusterSearcher) Search(wf *aw.Workflow, searchArgs searchutil.SearchArgs) error {
	return searchEntities(wf, searchArgs, "elasticache_redis_clusters", fetchCacheClusters, s.addToWorkflow)
}

func (ElasticacheRedisClusterSearcher) addToWorkflow(wf *aw.Workflow, searchArgs searchutil.SearchArgs, cluster types.CacheCluster) {
	addCacheClusterToWorkflow("redis", wf, searchArgs, cluster)
}
