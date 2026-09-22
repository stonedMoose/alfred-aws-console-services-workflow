package caching

import (
	"log"
	"os"
	"strconv"
	"time"
)

const (
	defaultMaxCacheAge = 180 * time.Second
	maxCacheAgeEnvVar  = "ALFRED_AWS_CONSOLE_SERVICES_WORKFLOW_MAX_CACHE_AGE_SECONDS"
)

// maxCacheAge is how long cached results stay fresh, overridable through the
// environment.
func maxCacheAge() time.Duration {
	configured := os.Getenv(maxCacheAgeEnvVar)
	if configured == "" {
		return defaultMaxCacheAge
	}
	seconds, err := strconv.Atoi(configured)
	if err != nil {
		panic(err)
	}
	if seconds == 0 {
		return defaultMaxCacheAge
	}
	log.Printf("using custom max cache age of %v seconds", seconds)
	return time.Duration(seconds) * time.Second
}
