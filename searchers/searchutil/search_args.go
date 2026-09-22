package searchutil

import (
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

// SearchArgs carries what a searcher needs to know about the current run: the
// AWS configuration to fetch with and the words the user has typed so far.
type SearchArgs struct {
	// Query is the part of the query left for the searcher to filter its results on.
	Query string
	Cfg   aws.Config
	// ForceFetch bypasses the cache and reads from AWS.
	ForceFetch bool
	// FullQuery is the query as the user typed it, the base of every autocompletion.
	FullQuery string
	Profile   string
	// GetRegionFunc tells which region the searched service lives in; global
	// services answer with an empty region.
	GetRegionFunc func(cfg aws.Config) string
	// IgnoreAutocompleteTerm keeps autocompletions equal to the full query, for
	// items shown next to a message rather than as search results.
	IgnoreAutocompleteTerm bool
}

// GetAutocomplete returns the full query with the search term replaced by
// completion, ready for the user to type further.
func (s *SearchArgs) GetAutocomplete(completion string) string {
	if s.IgnoreAutocompleteTerm {
		return s.FullQuery
	}
	if s.Query == "" {
		return s.FullQuery + completion + " "
	}
	return withTrailingSpace(util.ReplaceLast(s.FullQuery, s.Query, completion))
}

func withTrailingSpace(s string) string {
	if strings.HasSuffix(s, " ") {
		return s
	}
	return s + " "
}

// GetRegion returns the region the searched service lives in.
func (s *SearchArgs) GetRegion() string {
	return s.GetRegionFunc(s.Cfg)
}

// matchByIDOrTitle makes the item match and autocomplete on its identifier
// while the user is typing one, recognisable by idPrefix, and on its title
// otherwise.
func (s *SearchArgs) matchByIDOrTitle(item *aw.Item, idPrefix, id, title string) *aw.Item {
	if idPrefix != "" && id != "" && strings.HasPrefix(s.Query, idPrefix) {
		return item.Match(id).Autocomplete(s.GetAutocomplete(id))
	}
	return item.Autocomplete(s.GetAutocomplete(title))
}
