// Package caching keeps the resources a searcher lists in the workflow cache,
// so that Alfred answers instantly while AWS is only asked in the background.
package caching

import (
	"log"
	"os"
	"os/exec"

	"github.com/aws/aws-sdk-go-v2/aws"
	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/searchers/searchutil"
)

// backgroundFetchJob names the background job that refreshes a stale cache.
const backgroundFetchJob = "fetch"

// LoadEntities returns the entities of one resource kind: freshly fetched from
// AWS when the fetch is forced, otherwise from the cache, with a refresh
// scheduled in the background once the cache has grown stale.
func LoadEntities[Entity any](wf *aw.Workflow, searchArgs searchutil.SearchArgs, cacheKey string, fetch func(aws.Config) ([]Entity, error)) []Entity {
	// TODO optimization: global services (s3 buckets, ...) could share one
	// cache across regions instead of fetching once per region
	cache := entityCache[Entity]{
		wf:         wf,
		searchArgs: searchArgs,
		key:        cacheKey + "_" + searchArgs.Cfg.Region + "_" + searchArgs.Profile,
		fetch:      fetch,
	}
	if searchArgs.ForceFetch {
		return cache.refresh()
	}
	return cache.load()
}

// entityCache is the cache entry of one resource kind for one region and profile.
type entityCache[Entity any] struct {
	wf         *aw.Workflow
	searchArgs searchutil.SearchArgs
	key        string
	fetch      func(aws.Config) ([]Entity, error)
}

func (c entityCache[Entity]) refresh() []Entity {
	log.Printf("fetching from aws ...")
	entities, err := c.fetch(c.searchArgs.Cfg)
	if err != nil {
		recordFetchError(c.wf, err)
		panic(err)
	}
	clearFetchError(c.wf)
	log.Printf("fetched %d results from aws", len(entities))
	c.store(entities)
	return entities
}

func (c entityCache[Entity]) store(entities []Entity) {
	log.Printf("storing %d results with cache key `%s` to %s ...", len(entities), c.key, c.wf.CacheDir())
	if err := c.wf.Cache.StoreJSON(c.key, entities); err != nil {
		panic(err)
	}
}

func (c entityCache[Entity]) load() []Entity {
	if c.isStale() {
		c.refreshInBackground()
		if err := reportLastFetchError(c.wf, c.searchArgs); err != nil {
			return []Entity{}
		}
	}
	if !c.wf.Cache.Exists(c.key) {
		log.Printf("cache with key `%s` did not exist in %s ...", c.key, c.wf.CacheDir())
		c.wf.NewItem("Fetching ...").Icon(aw.IconInfo)
		return []Entity{}
	}
	return c.loadStored()
}

func (c entityCache[Entity]) loadStored() []Entity {
	log.Printf("using cache with key `%s` in %s ...", c.key, c.wf.CacheDir())
	entities := []Entity{}
	if err := c.wf.Cache.LoadJSON(c.key, &entities); err != nil {
		panic(err)
	}
	return entities
}

func (c entityCache[Entity]) isStale() bool {
	maxAge := maxCacheAge()
	if !c.wf.Cache.Expired(c.key, maxAge) {
		return false
	}
	log.Printf("cache with key `%s` was expired (older than %d seconds) in %s", c.key, int(maxAge.Seconds()), c.wf.CacheDir())
	return true
}

// refreshInBackground runs the workflow again with a forced fetch, so that the
// results Alfred shows right away get replaced by fresh ones.
func (c entityCache[Entity]) refreshInBackground() {
	c.wf.Rerun(0.5)
	if c.wf.IsRunning(backgroundFetchJob) {
		log.Printf("background job `%s` already running", backgroundFetchJob)
		return
	}
	cmd := exec.Command(os.Args[0], "-query="+c.searchArgs.FullQuery, "-fetch")
	log.Printf("running `%s` in background as job `%s` ...", cmd, backgroundFetchJob)
	if err := c.wf.RunInBackground(backgroundFetchJob, cmd); err != nil {
		panic(err)
	}
}
