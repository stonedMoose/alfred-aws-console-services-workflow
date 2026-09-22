package workflow

import (
	"fmt"
	"log"
	"strings"

	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/aliases"
	"github.com/rkoval/alfred-aws-console-services-workflow/awsconfig"
)

// suggestRegions lists the regions the user may pick while typing a region
// override.
func suggestRegions(wf *aw.Workflow, rawQuery, regionQuery string) {
	for _, region := range awsconfig.AllAWSRegions {
		wf.NewItem(region.Name).
			Subtitle(region.Description).
			Icon(aw.IconWeb).
			Autocomplete(completeOverride(rawQuery, aliases.OverrideAwsRegion, regionQuery, region.Name)).
			UID(region.Name)
	}
	filterSuggestions(wf, "region", regionQuery)
}

// suggestProfiles lists the profiles the user may pick while typing a profile
// override.
func suggestProfiles(wf *aw.Workflow, rawQuery, profileQuery string) {
	for _, profile := range awsconfig.GetAwsProfiles() {
		wf.NewItem(profile.Name).
			Subtitle(profileSubtitle(profile)).
			Icon(aw.IconAccount).
			Autocomplete(completeOverride(rawQuery, aliases.OverrideAwsProfile, profileQuery, profile.Name)).
			UID(profile.Name)
	}
	filterSuggestions(wf, "profile", profileQuery)
}

func profileSubtitle(profile awsconfig.Profile) string {
	if profile.Region == "" {
		return "⚠️ This profile does not specify a region. Functionality will be limited"
	}
	return fmt.Sprintf("🌎 %s", profile.Region)
}

// completeOverride replaces the override being typed with the chosen value.
func completeOverride(rawQuery, alias, typed, chosen string) string {
	return strings.Replace(rawQuery, alias+typed, alias+chosen+" ", 1)
}

func filterSuggestions(wf *aw.Workflow, kind, query string) {
	log.Printf("filtering with %s override %q", kind, query)
	matches := wf.Filter(query)
	log.Printf("%d results match %q", len(matches), query)
}
