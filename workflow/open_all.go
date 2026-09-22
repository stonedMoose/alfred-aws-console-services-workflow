package workflow

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/awsworkflow"
	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

const lastOpenedURLsFile = "last-opened-urls.json"

// consolePage is a console page to open, remembered with the id of the
// service it belongs to.
type consolePage struct {
	serviceID string
	url       string
}

// handleOpenAll opens every console page the query lists when the workflow
// was run to do so, and otherwise offers an item that runs it that way.
func (s *search) handleOpenAll() {
	if s.openAll {
		s.openAllInBrowser()
		return
	}
	s.offerToOpenAll()
}

func (s *search) openAllInBrowser() {
	var openedURLs []string
	for _, page := range s.pagesToOpen() {
		openInBrowser(s.wf, page)
		openedURLs = append(openedURLs, page.url)
	}
	rememberOpenedURLs(s.wf, openedURLs)
}

// pagesToOpen lists the home page of every service, or of every sub-service
// of the service the query names.
func (s *search) pagesToOpen() []consolePage {
	var pages []consolePage
	if s.query.Service == nil {
		for _, service := range s.awsServices {
			pages = append(pages, s.pageOf(service, service))
		}
		return pages
	}
	for _, subService := range s.query.Service.SubServices {
		pages = append(pages, s.pageOf(subService, *s.query.Service))
	}
	return pages
}

// pageOf builds the page of a service, in the region of the service that
// decides it: a sub-service lives in the region of its parent.
func (s *search) pageOf(service, regionOwner awsworkflow.AwsService) consolePage {
	return consolePage{
		serviceID: service.Id,
		url:       util.ConstructAWSConsoleUrl(service.Url, regionOwner.GetRegion(s.args.Cfg)),
	}
}

func openInBrowser(wf *aw.Workflow, page consolePage) {
	log.Printf("opening url: %s ...", page.url)
	if err := wf.RunInBackground("open-service-in-browser-"+page.serviceID, exec.Command("open", page.url)); err != nil {
		panic(err)
	}
	time.Sleep(250 * time.Millisecond) // so that the tabs open more or less in order
}

// rememberOpenedURLs writes the URLs to disk so that an automated check can
// replay them and see which pages redirect or error in AWS.
func rememberOpenedURLs(wf *aw.Workflow, urls []string) {
	urlBytes, err := json.Marshal(urls)
	if err != nil {
		panic(err)
	}
	if err := os.WriteFile(wf.CacheDir()+"/"+lastOpenedURLsFile, urlBytes, 0600); err != nil {
		panic(err)
	}
}

func (s *search) offerToOpenAll() {
	count, described := s.describePagesToOpen()
	s.wf.NewItem(fmt.Sprintf("Open the %d %s in browser", count, described)).
		Subtitle("Fair warning: this may briefly overload your system").
		Icon(aw.IconNote).
		Arg(fmt.Sprintf(`%s -query="%s" -open_all`, os.Args[0], s.query.RawQuery)).
		Var("action", "run-script").
		Valid(true)
}

func (s *search) describePagesToOpen() (count int, described string) {
	if s.query.Service == nil {
		return len(s.awsServices), "services"
	}
	return len(s.query.Service.SubServices), s.query.Service.Id + " sub-services"
}
