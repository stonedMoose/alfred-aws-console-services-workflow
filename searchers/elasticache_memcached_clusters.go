package searchers

import (
	"github.com/aws/aws-sdk-go-v2/service/elasticache/types"
	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/searchers/searchutil"
)

type ElasticacheMemcachedClusterSearcher struct{}

func (s ElasticacheMemcachedClusterSearcher) Search(wf *aw.Workflow, searchArgs searchutil.SearchArgs) error {
	return searchEntities(wf, searchArgs, "elasticache_memcached_clusters", fetchCacheClusters, s.addToWorkflow)
}

func (ElasticacheMemcachedClusterSearcher) addToWorkflow(wf *aw.Workflow, searchArgs searchutil.SearchArgs, cluster types.CacheCluster) {
	addCacheClusterToWorkflow("memcached", wf, searchArgs, cluster)
}
