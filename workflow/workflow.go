// Package workflow answers Alfred queries: it interprets the query and adds
// the matching services, sub-services or AWS resources as items.
package workflow

import (
	"fmt"
	"log"
	"net/http"

	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/aliases"
	"github.com/rkoval/alfred-aws-console-services-workflow/parsers"
)

// Run answers one query: it adds the matching items to the workflow and sends
// them to Alfred. transport, when set, replaces the HTTP transport of the AWS
// clients; openAll tells that the query's console pages are to be opened
// rather than listed.
func Run(wf *aw.Workflow, rawQuery string, transport http.RoundTripper, forceFetch, openAll bool, ymlPath string) {
	log.Println("using workflow cacheDir: " + wf.CacheDir())
	log.Println("using workflow dataDir: " + wf.DataDir())

	query, awsServices := parsers.NewParser(rawQuery).Parse(ymlPath)
	defer sendFeedback(wf, query)
	log.Printf("using query: %#v", query)

	switch {
	case query.RegionQuery != nil:
		suggestRegions(wf, rawQuery, *query.RegionQuery)
	case query.ProfileQuery != nil:
		suggestProfiles(wf, rawQuery, *query.ProfileQuery)
	default:
		newSearch(wf, query, awsServices, transport, forceFetch, openAll).run()
	}
}

// sendFeedback hands the items to Alfred, with a "nothing found" item when
// there are none. A panic is left for awgo to report.
func sendFeedback(wf *aw.Workflow, query *parsers.Query) {
	if panicValue := recover(); panicValue != nil {
		panic(panicValue)
	}
	if wf.IsEmpty() {
		addNoMatchItem(wf, query)
		addUpdateAvailableItem(wf)
	}
	wf.SendFeedback()
}

func addNoMatchItem(wf *aw.Workflow, query *parsers.Query) {
	title, subtitle := noMatchMessage(query)
	wf.NewItem(title).
		Subtitle(subtitle).
		Icon(aw.IconNote)
}

func noMatchMessage(query *parsers.Query) (title, subtitle string) {
	switch {
	case query.RegionQuery != nil:
		return "No matching regions found", startOverHint(aliases.OverrideAwsRegion)
	case query.ProfileQuery != nil:
		return "No matching profiles found", startOverHint(aliases.OverrideAwsProfile)
	default:
		return "No matching services found", "Try another query (example: `aws ec2 instances`)"
	}
}

func startOverHint(alias string) string {
	return fmt.Sprintf("Try starting over with \"%s\" again to see the full list", alias)
}
