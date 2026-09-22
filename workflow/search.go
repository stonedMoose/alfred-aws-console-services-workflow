package workflow

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/awsworkflow"
	"github.com/rkoval/alfred-aws-console-services-workflow/parsers"
	"github.com/rkoval/alfred-aws-console-services-workflow/searchers"
	"github.com/rkoval/alfred-aws-console-services-workflow/searchers/searchutil"
	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

const contributingGuideURL = "https://github.com/rkoval/alfred-aws-console-services-workflow/blob/master/CONTRIBUTING.md"

// search is one query for services, sub-services or resources, run against
// the AWS configuration the query selects.
type search struct {
	wf          *aw.Workflow
	query       *parsers.Query
	awsServices []awsworkflow.AwsService
	args        searchutil.SearchArgs
	openAll     bool
}

func newSearch(wf *aw.Workflow, query *parsers.Query, awsServices []awsworkflow.AwsService, transport http.RoundTripper, forceFetch, openAll bool) *search {
	cfg := awsworkflow.InitAWS(transport, query.ProfileOverride, query.GetRegionOverride())
	return &search{
		wf:          wf,
		query:       query,
		awsServices: awsServices,
		openAll:     openAll,
		args: searchutil.SearchArgs{
			Cfg:        cfg,
			ForceFetch: forceFetch,
			FullQuery:  query.RawQuery,
			Profile:    awsworkflow.ProfileName(cfg),
		},
	}
}

func (s *search) run() {
	switch {
	case s.query.IsEmpty():
		s.showHome()
	case s.query.HasOpenAll:
		s.handleOpenAll()
	case s.query.Service == nil || s.query.IsBareServiceId():
		s.listServices()
	default:
		s.searchWithinService()
	}
}

// listServices offers every service, filtered by what the user typed so far.
func (s *search) listServices() {
	s.args.Query = s.serviceFilter()
	log.Printf("using searcher associated with services with query %q", s.args.Query)
	addServiceItems(s.wf, s.awsServices, s.args)
	s.filterItems()
}

// serviceFilter is what the service list is filtered by: the words typed, or
// the short name or id of the service the user is still typing.
func (s *search) serviceFilter() string {
	service := s.query.Service
	switch {
	case service == nil:
		return s.query.RemainingQuery
	case service.ShortName != "":
		return strings.ToLower(service.ShortName)
	default:
		return service.Id
	}
}

func (s *search) searchWithinService() {
	service := s.query.Service
	if !s.query.HasDefaultSearchAlias && !service.HasSubServices() {
		s.showUnimplemented(nil, fmt.Sprintf("%s doesn't have sub-services configured (yet)", service.Id))
		return
	}
	s.args.GetRegionFunc = service.GetRegion
	if s.query.TargetsResources() {
		s.searchResources()
		return
	}
	s.listSubServices()
}

// listSubServices offers the sub-services of the service, filtered by what
// the user typed after it.
func (s *search) listSubServices() {
	log.Println("using searcher associated with sub-services")
	s.args.Query = s.subServiceFilter()
	addSubServiceItems(s.wf, *s.query.Service, s.args)
	s.filterItems()
}

func (s *search) subServiceFilter() string {
	if s.query.SubService != nil {
		return s.query.SubService.Id
	}
	return s.query.RemainingQuery
}

// searchResources runs the searcher of the service or sub-service on the rest
// of the query.
func (s *search) searchResources() {
	serviceID, subServiceID := s.query.Service.Id, s.subServiceID()
	log.Println("using searcher associated with " + strings.TrimSpace(serviceID+" "+subServiceID))
	searcher := searchers.SearcherFor(serviceID, subServiceID)
	if searcher == nil {
		s.showUnimplemented(s.query.SubService, fmt.Sprintf("No searcher for `%s` (yet)", strings.TrimSpace(serviceID+" "+subServiceID)))
		return
	}
	s.args.Query = s.query.RemainingQuery
	if err := searcher.Search(s.wf, s.args); err != nil {
		s.wf.FatalError(err)
	}
	s.filterItems()
}

func (s *search) subServiceID() string {
	if s.query.SubService == nil {
		return ""
	}
	return s.query.SubService.Id
}

// showUnimplemented shows the service or sub-service as a plain console link,
// with a pointer to the contributing guide.
func (s *search) showUnimplemented(subService *awsworkflow.AwsService, header string) {
	s.args.IgnoreAutocompleteTerm = true
	if subService == nil {
		addServiceItem(s.wf, *s.query.Service, s.args)
	} else {
		addSubServiceItem(s.wf, *s.query.Service, *subService, s.args)
	}
	util.NewURLItem(s.wf, header).
		Subtitle("Select this result to open the contributing guide to easily add them!").
		Arg(contributingGuideURL).
		Icon(aw.IconNote)
	addUpdateAvailableItem(s.wf)
}

func (s *search) filterItems() {
	filterOnEveryTerm(s.wf, s.args.Query)
}

// multiTermSeparator splits a search into terms that must all match. It is
// deliberately not the configurable search alias: someone who rebinds that
// alias to "." would otherwise find every domain name chopped in half.
const multiTermSeparator = ","

// filterOnEveryTerm keeps the items that match all of the query's terms. awgo's
// Filter both ranks and drops, so running it once per term intersects them. The
// terms go last to first, leaving the surviving ranking to the first term,
// which is the one the person typed most deliberately.
func filterOnEveryTerm(wf *aw.Workflow, query string) {
	terms := searchTerms(query)
	for i := len(terms) - 1; i >= 0; i-- {
		log.Printf("filtering with query %q", terms[i])
		matches := wf.Filter(terms[i])
		log.Printf("%d results match %q", len(matches), terms[i])
	}
}

func searchTerms(query string) []string {
	var terms []string
	for _, term := range strings.Split(query, multiTermSeparator) {
		if term = strings.TrimSpace(term); term != "" {
			terms = append(terms, term)
		}
	}
	return terms
}
