package searchers

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/caching"
	"github.com/rkoval/alfred-aws-console-services-workflow/searchers/searchutil"
)

// Searcher lists one kind of AWS resource as Alfred items.
type Searcher interface {
	Search(wf *aw.Workflow, searchArgs searchutil.SearchArgs) error
}

// searchEntities is the template every searcher follows: load the entities of
// its resource kind through the cache, then present each one as an item.
// cacheKey names the resource kind in the cache.
func searchEntities[Entity any](
	wf *aw.Workflow,
	searchArgs searchutil.SearchArgs,
	cacheKey string,
	fetch func(aws.Config) ([]Entity, error),
	present func(*aw.Workflow, searchutil.SearchArgs, Entity),
) error {
	for _, entity := range caching.LoadEntities(wf, searchArgs, cacheKey, fetch) {
		present(wf, searchArgs, entity)
	}
	return nil
}
