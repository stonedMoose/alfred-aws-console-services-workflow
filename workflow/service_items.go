package workflow

import (
	"strings"

	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/awsworkflow"
	"github.com/rkoval/alfred-aws-console-services-workflow/searchers"
	"github.com/rkoval/alfred-aws-console-services-workflow/searchers/searchutil"
	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

func addServiceItems(wf *aw.Workflow, services []awsworkflow.AwsService, searchArgs searchutil.SearchArgs) {
	for _, service := range services {
		addServiceItem(wf, service, searchArgs)
	}
}

func addSubServiceItems(wf *aw.Workflow, service awsworkflow.AwsService, searchArgs searchutil.SearchArgs) {
	for _, subService := range service.SubServices {
		addSubServiceItem(wf, service, subService, searchArgs)
	}
}

// addServiceItem offers the home page of a service.
func addServiceItem(wf *aw.Workflow, service awsworkflow.AwsService, searchArgs searchutil.SearchArgs) {
	util.NewURLItem(wf, service.Id).
		Subtitle(serviceSubtitle(service)).
		Match(serviceMatchTerms(service)).
		Autocomplete(searchArgs.GetAutocomplete(service.Id)).
		UID(service.Id).
		Arg(util.ConstructAWSConsoleUrl(service.Url, service.GetRegion(searchArgs.Cfg))).
		Icon(util.ServiceIcon(service.Id))
}

// serviceSubtitle reads "Name (ShortName) – Description", flagged when the
// service has sub-services to browse.
func serviceSubtitle(service awsworkflow.AwsService) string {
	var subtitle strings.Builder
	if service.HasSubServices() {
		subtitle.WriteString("🗂 ")
	}
	subtitle.WriteString(service.Name)
	if service.ShortName != "" {
		subtitle.WriteString(" (" + service.ShortName + ")")
	}
	subtitle.WriteString(" – " + service.Description)
	return subtitle.String()
}

// serviceMatchTerms is what a service is found by: its id, its full name
// unless a short name stands for it, and any extra search terms.
func serviceMatchTerms(service awsworkflow.AwsService) string {
	terms := []string{service.Id}
	if service.ShortName == "" {
		terms = append(terms, service.Name)
	}
	terms = append(terms, service.ExtraSearchTerms...)
	return strings.Join(terms, " ")
}

// addSubServiceItem offers the page of a sub-service.
func addSubServiceItem(wf *aw.Workflow, service, subService awsworkflow.AwsService, searchArgs searchutil.SearchArgs) {
	util.NewURLItem(wf, service.Id+" "+subService.Id).
		Subtitle(subServiceSubtitle(service, subService)).
		Match(subService.Id + " " + subService.Name).
		Autocomplete(searchArgs.GetAutocomplete(subService.Id)).
		UID(subService.Id).
		Arg(util.ConstructAWSConsoleUrl(subService.Url, service.GetRegion(searchArgs.Cfg))).
		Icon(util.ServiceIcon(service.Id))
}

func subServiceSubtitle(service, subService awsworkflow.AwsService) string {
	subtitle := searcherMarker(service.Id, subService.Id) + serviceLabel(service) + " – " + subService.Name
	if subService.Description != "" {
		subtitle += " – " + subService.Description
	}
	return subtitle
}

// searcherMarker flags the sub-services whose resources can be searched, with
// a star on the one a bare service id searches.
func searcherMarker(serviceID, subServiceID string) string {
	switch {
	case searchers.IsDefaultSearcher(serviceID, subServiceID):
		return "🔎⭐️ "
	case searchers.SearcherFor(serviceID, subServiceID) != nil:
		return "🔎 "
	default:
		return ""
	}
}

func serviceLabel(service awsworkflow.AwsService) string {
	if service.ShortName != "" {
		return service.ShortName
	}
	return service.Name
}
